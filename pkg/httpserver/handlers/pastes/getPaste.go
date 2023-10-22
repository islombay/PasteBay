package pastes

import (
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/auth"
	"PasteBay/pkg/utils/blob"
	"PasteBay/pkg/utils/logger/sl"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

// GetPaste godoc
// @Summary Get paste
// @Tags paste
// @Description Using alias of the paste, get paste
// @ID get_paste
// @Accept json
// @Produce json
// @Param input body models.RequestGetPaste false "paste info"
// @Success 200 {object} models.ResponseGetPaste
// @Failure 400,403,404,500 {object} response.errorResponse
// @Router /api/paste [get]
// @Router /api/paste [post]
func GetPaste(log *slog.Logger, db *database.Database, blob *blob.BlobStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqbody models.RequestGetPaste
		if c.Request.Method == "POST" {
			if err := c.BindJSON(&reqbody); err != nil {
				response.ErrorResponse(c, response.ErrorBadRequest)
				return
			}
		}

		alias := c.Param("alias")

		pasteObject, err := db.GetPaste(alias)
		if err != nil {
			if err.Error() == database.ErrorNotFound {
				response.ErrorResponse(c, response.ErrorNotFound)
			} else {
				log.Error("Could not get paste from db", sl.Err(err))
				response.ErrorResponse(c, response.ErrorServerError)
			}
			return
		}

		// TODO: Expire_time implementation
		if !(pasteObject.ExpireTime.Equal(pasteObject.CreatedAt)) {
			now := time.Now()

			fmt.Println(pasteObject.ExpireTime, now)
			fmt.Printf("-------\nTime After: %v\n--------\n", now.After(pasteObject.ExpireTime))
			fmt.Printf("-------\nTime Before: %v\n--------\n", now.Before(pasteObject.ExpireTime))
			if !pasteObject.ExpireTime.After(now) {
				log.Debug("Paste time up")
				blob_src, err := db.DeletePaste(pasteObject.ID)
				if err != nil {
					log.Error(fmt.Sprintf("Could not delete paste (id: %s)", pasteObject.ID), sl.Err(err))
				} else {
					err = blob.Delete(blob_src)
					if err != nil {
						log.Error(fmt.Sprintf("Could not delete blob of paste (src: %s)", blob_src), sl.Err(err))
					}
				}

				response.ErrorResponse(c, response.ErrorNotFound)
				return
			}
		}

		if pasteObject.AccessPassword != "" {
			if !(reqbody.Password != "" && auth.GenerateHash(reqbody.Password) == pasteObject.AccessPassword) {
				response.ErrorResponse(c, response.ErrorInvalidCredentials)
				return
			}
		}

		if (pasteObject.ViewsLimit <= int64(pasteObject.ViewsCount)) && pasteObject.ViewsLimit != -1 {
			log.Debug("Paste view limit reached")
			blob_src, err := db.DeletePaste(pasteObject.ID)
			if err != nil {
				log.Error(fmt.Sprintf("Could not delete paste (id: %s)", pasteObject.ID), sl.Err(err))
			} else {
				err = blob.Delete(blob_src)
				if err != nil {
					log.Error(fmt.Sprintf("Could not delete blob of paste (src: %s)", blob_src), sl.Err(err))
				}
			}

			response.ErrorResponse(c, response.ErrorNotFound)
			return
		}

		db.ViewIncreasePaste(pasteObject.ID)

		blobContent, err := blob.GetContent(pasteObject.BlobSrc)
		if err != nil {
			log.Error(fmt.Sprintf("Could not load the content of the blob (src: %s)", pasteObject.BlobSrc), sl.Err(err))
			response.ErrorResponse(c, response.ErrorNotFound)
			return
		}

		response := models.ResponseGetPaste{
			CreatedAt:  pasteObject.CreatedAt,
			UpdatedAt:  pasteObject.UpdatedAt,
			Title:      pasteObject.Title,
			ViewsCount: pasteObject.ViewsCount,
			Content:    blobContent,
		}
		if int(pasteObject.Author) != -1 {
			log.Debug("Author exists")
			author_obj, err := db.GetUserByID(int(pasteObject.Author))
			if err != nil {
				log.Error("Could not get user by id for get paste function", sl.Err(err))
			} else {
				author_res := models.ResponseGetPasteAuthor{
					LastLogin: author_obj.LastLogin,
					Username:  author_obj.Username,
				}
				response.Author = author_res
			}
		}
		c.JSON(http.StatusOK, response)
	}
}
