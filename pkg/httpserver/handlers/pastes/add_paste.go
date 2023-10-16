package pastes

import (
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/auth"
	"PasteBay/pkg/utils/blob"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

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

		_, alias, err := db.AddPaste(database.BodyAddPaste{
			Author:     0,
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
