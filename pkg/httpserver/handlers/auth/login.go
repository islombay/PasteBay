package auth

import (
	"PasteBay/configs"
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/models"
	"PasteBay/pkg/utils/auth"
	"PasteBay/pkg/utils/logger/sl"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

const (
	MiddlewareAuthUsername = "auth_username"
)

// LoginHandler godoc
// @Summary Login
// @Tags account
// @Description Login
// @ID login_user
// @Accept json
// @Produce json
// @Param input body models.RequestLogin true "login info"
// @Success 200 {object} models.ResponseLogin
// @Failure 400,403,404,500 {object} response.errorResponse
// @Router /auth/login [post]
func LoginHandler(log *slog.Logger, db *database.Database, authCfg configs.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody models.RequestLogin
		if err := c.BindJSON(&reqBody); err != nil {
			response.ErrorResponse(c, response.ErrorBadRequest)
			return
		}

		err := CheckRegisBody(c, log, db, map[string]map[string]string{
			"login": {
				"username": reqBody.Username,
				"password": reqBody.Password,
			},
		})
		if err != nil {
			return
		}

		exists, err := db.LoginValid(reqBody.Username, reqBody.Password)
		if err != nil {
			if err.Error() == database.DBNotFound {
				response.ErrorResponse(c, response.ErrorNotFound)
				return
			}
			log.Error("could not check login validness", sl.Err(err))
			response.ErrorResponse(c, response.ErrorServerError)
			return
		}
		if !exists {
			response.ErrorResponse(c, response.ErrorInvalidCredentials)
			return
		}

		token, err := auth.GenerateToken(reqBody.Username, auth.GenerateTokenBody{
			TokenTTL: authCfg.TokenTTL,
			Secret:   authCfg.SecretKey,
		})
		if err != nil {
			response.ErrorResponse(c, response.ErrorServerError)
			return
		}
		c.JSON(http.StatusOK, models.ResponseLogin{
			token,
			true,
		})
	}
}
