package utils

import (
	"dynamis-server/middleware"
	"dynamis-server/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var signingKey = []byte("your-secret-key")

// GenerateJWT generates a JWT token for the authenticated user
func GenerateJWT(user *models.User) (string, error) {
	claims := dynamis_middleware.Claims{
		UserID:        user.ID,
		Subscriptions: user.Subscriptions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "your_project",
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}
