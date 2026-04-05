package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		//check if the header exist
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header missing",
			})
			c.Abort()
			return
		}

		// format must : bearer <token>
		tokenString := strings.Split(authHeader, " ")
		if len(tokenString) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "invalid token format",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString[1], func (token *jwt.Token) (interface{}, error)  {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error" : "invalid token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}