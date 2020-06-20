package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"example.com/gin_server/handlers"
	"example.com/gin_server/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// IsAuthorized : checks auth with JWT and Bearer Authentication
func IsAuthorized(d *storage.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")

		// Splitting auth header by spaces
		authHeader := strings.Split(c.GetHeader("Authorization"), " ")

		if len(authHeader) != 2 {
			handlers.HandleError(http.StatusBadRequest, errors.New("Invalid auth format"), c)
			return
		}

		if authHeader[0] == "Bearer" {
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(authHeader[1], claims, func(token *jwt.Token) (interface{}, error) {
				// Check signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				// I'm not doing anymore validation stuff here. I'm only using one signing key
				return []byte(os.Getenv("SECRET")), nil
			})

			if token.Valid {
				fmt.Println("claims", claims)
				log.Println("authorized")
				return
			}

			if err != nil {
				log.Println(err.Error())
			}
		}

		handlers.HandleError(http.StatusUnauthorized, errors.New("Unauthorized"), c)
		return
	}
}
