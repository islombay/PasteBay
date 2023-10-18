package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type tokenClaims struct {
	jwt.StandardClaims
	Username string
}

type GenerateTokenBody struct {
	TokenTTL int
	Secret   string
}

func GenerateToken(username string, body GenerateTokenBody) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(body.TokenTTL) * (time.Hour * 24)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
	})

	return token.SignedString([]byte(body.Secret))
}
