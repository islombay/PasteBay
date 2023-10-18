package auth

import (
	"PasteBay/pkg/database"
	"PasteBay/pkg/httpserver/response"
	"PasteBay/pkg/utils/auth"
	"PasteBay/pkg/utils/logger/sl"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strings"
)

func CheckRegisBody(c *gin.Context, log *slog.Logger, db *database.Database, data map[string]map[string]string) error {
	regisBody, exists := data["register"]
	if exists {
		value, exists := regisBody["username"]
		if exists {
			isValid := auth.IsUsernameValid(value)
			if !isValid {
				response.ErrorResponse(c, response.ErrorAuthUsername)
				return errors.New("username not valid")
			}

			isAvailable, err := db.CheckUsername(value)
			if err != nil || !isAvailable {
				if !isAvailable {
					response.ErrorResponse(c, response.ErrorUserExists)
					return errors.New("username taken")
				}
				log.Error("Error while checking username availability", sl.Err(err))
				response.ErrorResponse(c, response.ErrorCouldNotCreateUser)
				return errors.New("username could not be gotten")
			}

			if value == "" || len(value) < 5 {
				response.ErrorResponse(c, response.ErrorAuthUsername)
				return errors.New("username could not be gotten")
			}
		}
		value, exists = regisBody["password"]
		if exists {
			value = strings.TrimSpace(value)
			if value == "" || len(value) < 8 {
				response.ErrorResponse(c, response.ErrorAuthPassword)
				return errors.New("password")
			}
		}
		value, exists = regisBody["email"]
		if exists {
			if value == "" {
				response.ErrorResponse(c, response.ErrorBadRequest)
				return errors.New("email")
			}
		}
	} else if loginBody, exists := data["login"]; exists {
		value, exists := loginBody["password"]
		if exists {
			value = strings.TrimSpace(value)
			if value == "" || len(value) < 8 {
				response.ErrorResponse(c, response.ErrorLoginPassword)
				return errors.New("password")
			}
		}

		value, exists = regisBody["username"]
		if exists {
			isValid := auth.IsUsernameValid(value)
			if !isValid {
				response.ErrorResponse(c, response.ErrorLoginUsername)
				return errors.New("username not valid")
			}

			if value == "" || len(value) < 5 {
				response.ErrorResponse(c, response.ErrorLoginUsername)
				return errors.New("username could not be gotten")
			}
		}
	}
	return nil
}
