package auth

import (
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func RegisterHandler(log *slog.Logger, db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody models.RequestRegister
		if err := c.BindJSON(&reqBody); err != nil {
			response.ErrorResponse(c, response.ErrorBadRequest)
			return
		}

	}
}
