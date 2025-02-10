package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret-key")

func GenerateToken(userID string) (string, error) {
	claims := &jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(3 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
