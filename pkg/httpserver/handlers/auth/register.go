package auth

import (
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/logger/sl"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func RegisterHandler(log *slog.Logger, db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody models.RequestRegister
		if err := c.BindJSON(&reqBody); err != nil {
			response.ErrorResponse(c, response.ErrorBadRequest)
			return
		}

		err := CheckRegisBody(c, log, db, map[string]map[string]string{
			"register": {
				"username": reqBody.Username,
				"password": reqBody.Password,
				"email":    reqBody.Email,
			},
		})
		if err != nil {
			return
		}

		err = db.CreateUser(reqBody.Email, reqBody.Username, reqBody.Password)
		if err != nil {
			log.Error("Could not create user in database", sl.Err(err))
			response.ErrorResponse(c, response.ErrorCouldNotCreateUser)
			return
		}

		c.JSON(http.StatusOK, models.ResponseRegister{
			true,
		})
	}
}
