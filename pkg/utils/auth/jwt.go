package auth

import (
	"errors"
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

func ParseToken(token string, body GenerateTokenBody) (string, error) {
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid_signing_method")
		}
		return []byte(body.Secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		return claims["Username"].(string), nil
	}
	return "", nil
}
