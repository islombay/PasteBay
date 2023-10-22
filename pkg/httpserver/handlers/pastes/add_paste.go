package pastes

import (
	"PasteBay/pkg/database"
	auth2 "PasteBay/pkg/httpserver/handlers/auth"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/auth"
	"PasteBay/pkg/utils/blob"
	"PasteBay/pkg/utils/logger/sl"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

// AddPaste godoc
// @Summary Add Paste
// @Tags paste
// @Description Add paste and return short alias
// @ID add_paste
// @Accept json
// @Produce json
// @Param input body models.RequestAddPaste true "paste info"
// @Success 200 {object} models.ResponseAddPaste
// @Failure 400,500 {object} response.errorResponse
// @Router /api/paste [post]
func AddPaste(log *slog.Logger, db *database.Database, blob *blob.BlobStorage, aliasPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var addPasteReq models.RequestAddPaste
		if err := c.BindJSON(&addPasteReq); err != nil {
			response.ErrorResponse(c, response.ErrorBadRequest)
			return
		}

		blobPath, err := blob.Save(addPasteReq.Content)
		if err != nil {
			response.ErrorResponse(c, response.ErrorServerError)
			return
		}
		password := ""
		if addPasteReq.AccessPassword != "" {
			password = auth.GenerateHash(addPasteReq.AccessPassword)
		}

		authU, isAuth := c.Get(auth2.MiddlewareAuthUsername)
		author := 0
		if isAuth {
			user, err := db.GetUserByUsername(authU.(string))
			if err != nil {
				if err.Error() == database.DBNotFound {
					log.Error(fmt.Sprintf("Could not find the user in the db by username, user: %s", authU.(string)), sl.Err(err))
				} else {
					log.Error("Could not get the user by username", sl.Err(err))
				}
			} else {
				author = int(user.ID)
			}
		}

		_, alias, err := db.AddPaste(database.BodyAddPaste{
			Author:     author,
			Title:      addPasteReq.Title,
			ExpireTime: addPasteReq.ExpireTime,
			ViewsLimit: addPasteReq.ViewsLimit,
			BlobSrc:    blobPath,
			Password:   password,
		})
		if err != nil {
			blob.Delete(blobPath)
			response.ErrorResponse(c, response.ErrorServerError)
			return
		}

		c.JSON(http.StatusOK, models.ResponseAddPaste{
			Title:   addPasteReq.Title,
			Content: addPasteReq.Content,
			Alias:   aliasPath + "/" + alias,
		})
	}
}
