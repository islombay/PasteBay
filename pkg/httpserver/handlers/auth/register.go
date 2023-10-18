package auth

import (
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/auth"
	"PasteBay/pkg/utils/logger/sl"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

func RegisterHandler(log *slog.Logger, db *database.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody models.RequestRegister
		if err := c.BindJSON(&reqBody); err != nil {
			response.ErrorResponse(c, response.ErrorBadRequest)
			return
		}

		isValid := auth.IsUsernameValid(reqBody.Username)
		if !isValid {
			response.ErrorResponse(c, response.ErrorAuthUsername)
			return
		}
		isAvailable, err := db.CheckUsername(reqBody.Username)
		if err != nil || !isAvailable {
			if !isAvailable {
				response.ErrorResponse(c, response.ErrorUserExists)
				return
			}
			log.Error("Error while checking username availability", sl.Err(err))
			response.ErrorResponse(c, response.ErrorCouldNotCreateUser)
			return
		}

		reqBody.Password = strings.TrimSpace(reqBody.Password)
		if reqBody.Password == "" || len(reqBody.Password) < 8 {
			response.ErrorResponse(c, response.ErrorAuthPassword)
			return
		}

		if reqBody.Username == "" || len(reqBody.Username) < 5 {
			response.ErrorResponse(c, response.ErrorAuthUsername)
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
