package handlers

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

func addPaste(log *slog.Logger, db *database.Database, blob *blob.BlobStorage, aliasPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var addPasteReq models.RequestAddPaste
		if err := c.BindJSON(&addPasteReq); err != nil {
			response.ErrorResponse(c, http.StatusBadRequest, "bad request")
			return
		}

		blobPath, err := blob.Save(addPasteReq.Content)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "could not save content")
			return
		}
		password := ""
		if addPasteReq.AccessPassword != "" {
			password = auth.GenerateHash(addPasteReq.AccessPassword)
		}

		_, alias, err := db.AddPaste(log, database.BodyAddPaste{
			Author:     0,
			Title:      addPasteReq.Title,
			IsPrivate:  addPasteReq.IsPrivate,
			ExpireTime: addPasteReq.ExpireTime,
			ViewsLimit: addPasteReq.ViewsLimit,
			BlobSrc:    blobPath,
			Password:   password,
		})
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "content is not saved")
			return
		}

		c.JSON(http.StatusOK, models.ResponseAddPaste{
			Title:     addPasteReq.Title,
			IsPrivate: addPasteReq.IsPrivate,
			Content:   addPasteReq.Content,
			Alias:     aliasPath + "/" + alias,
		})
	}
}
