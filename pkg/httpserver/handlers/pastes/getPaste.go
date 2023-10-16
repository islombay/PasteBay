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
	"strconv"
	"time"
)

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
				response.ErrorResponse(c, response.ErrorServerError)
			}
			return
		}

		// TODO: Expire_time implementation
		if !(pasteObject.ExpireTime.Equal(pasteObject.CreatedAt)) {
			now := time.Now()

			fmt.Println(pasteObject.ExpireTime.Compare(now))

			if !pasteObject.ExpireTime.After(now) {
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
		}

		response := models.ResponseGetPaste{
			CreatedAt:  pasteObject.CreatedAt,
			UpdatedAt:  pasteObject.UpdatedAt,
			Author:     strconv.FormatInt(pasteObject.Author, 10),
			Title:      pasteObject.Title,
			ViewsCount: pasteObject.ViewsCount,
			Content:    blobContent,
		}
		c.JSON(http.StatusOK, response)
	}
}
