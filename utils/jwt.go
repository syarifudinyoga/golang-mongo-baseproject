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

func AuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Token tidak ada, kirim error dan hentikan proses
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			c.Abort() // Hentikan eksekusi lebih lanjut
			return
		}

		// Pisahkan Bearer dari token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // Artinya tidak ada "Bearer" dalam header
			// Format token salah, kirim error dan hentikan proses
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort() // Hentikan eksekusi lebih lanjut
			return
		}

		// Verifikasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			// Token tidak valid, kirim error dan hentikan proses
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort() // Hentikan eksekusi lebih lanjut
			return
		}

		c.Next() // Jika token valid, lanjutkan ke handler berikutnya
	}
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
