package middleware

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret-key") // soon can be used in env

func GeneratToken(UserID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id" 	: UserID,
		"exp"		: time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}