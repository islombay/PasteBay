package middlewares

import (
	"PasteBay/configs"
	auth2 "PasteBay/pkg/httpserver/handlers/auth"
	"PasteBay/pkg/utils/auth"
	"github.com/gin-gonic/gin"
	"strings"
)

func MiddlewareAuth(config configs.AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(tokenHeader) == 2 {
			jwtToken := tokenHeader[1]
			if jwtToken != "" {
				username, err := auth.ParseToken(jwtToken, auth.GenerateTokenBody{
					TokenTTL: config.TokenTTL,
					Secret:   config.SecretKey,
				})
				if err != nil {
					c.Next()
					return
				}
				c.Set(auth2.MiddlewareAuthUsername, username)
				c.Next()
			}
		}
		//c.Abort()
		c.Next()
	}
}
