package pastes

import (
	"PasteBay/pkg/database"
	auth2 "PasteBay/pkg/httpserver/handlers/auth"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/blob"
	"PasteBay/pkg/utils/logger/sl"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// DeletePasteHandler godoc
// @Summary Delete Paste
// @Security ApiKeyAuth
// @Tags paste
// @Description Delete paste using ID
// @ID delete_paste
// @Accept json
// @Produce json
// @Param input body models.RequestDeletePaste true "paste info"
// @Success 200 {object} models.ResponseDeletePaste
// @Failure 400,401,404,403,500 {object} response.errorResponse
// @Router /api/paste [delete]
func DeletePasteHandler(log *slog.Logger, db *database.Database, blob *blob.BlobStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody models.RequestDeletePaste
		if err := c.BindJSON(&reqBody); err != nil {
			response.ErrorResponse(c, response.ErrorBadRequest)
			return
		}
		authU, isAuth := c.Get(auth2.MiddlewareAuthUsername)
		if !isAuth {
			response.ErrorResponse(c, response.ErrorUnauthorized)
			return
		}
		user, err := db.GetUserByUsername(authU.(string))
		if err != nil {
			if err.Error() == database.DBNotFound {
				log.Error(fmt.Sprintf("Could not find the user in the db by username, user: %s", authU.(string)), sl.Err(err))
				response.ErrorResponse(c, response.ErrorBadRequest)
				return
			} else {
				log.Error("Could not get the user by username", sl.Err(err))
				response.ErrorResponse(c, response.ErrorServerError)
				return
			}
		}

		pasteObject, err := db.GetPaste(reqBody.Alias)
		if err != nil {
			if err.Error() == database.ErrorNotFound {
				response.ErrorResponse(c, response.ErrorNotFound)
			} else {
				log.Error("Could not get paste from db", sl.Err(err))
				response.ErrorResponse(c, response.ErrorServerError)
			}
			return
		}
		if int64(user.ID) != pasteObject.Author {
			response.ErrorResponse(c, response.ErrorForbidden)
			return
		}
		blob_src, err := db.DeletePaste(pasteObject.ID)
		if err != nil {
			log.Error(fmt.Sprintf("Could not delete paste (id: %s)", pasteObject.ID), sl.Err(err))
			response.ErrorResponse(c, response.ErrorServerError)
			return
		} else {
			err = blob.Delete(blob_src)
			if err != nil {
				log.Error(fmt.Sprintf("Could not delete blob of paste (src: %s)", blob_src), sl.Err(err))
				response.ErrorResponse(c, response.ErrorServerError)
				return
			}
		}
		c.JSON(http.StatusOK, models.ResponseDeletePaste{
			true,
		})
	}
}
