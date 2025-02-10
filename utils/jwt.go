package utils

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func GenerateToken(userID string) (string, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

	var jwtKey = []byte(os.Getenv("SECRET_KEY"))

	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(3 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func AdminLevel() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := godotenv.Load(); err != nil {
			fmt.Println("Error loading .env file")
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}
		tokenString := strings.Split(authHeader, " ")[1]
		// println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		level := claims["role"].(string)
		// username := claims["username"].(string)

		// Jika level bukan 'admin', kembalikan error
		if level != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have access to this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
