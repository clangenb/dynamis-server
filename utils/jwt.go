package utils

import (
	"dynamis-server/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWTService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

const JWTSecretEnv = "JWT_SECRET"

func GetJWTSecret() string {
	secret := os.Getenv(JWTSecretEnv)
	if secret == "" {
		panic("JWT_SECRET is not set") // or use a default/fallback for dev
	}
	return secret
}

func NewJWTService() *JWTService {
	secret := GetJWTSecret()
	return &JWTService{
		secret: []byte(secret),
		issuer: issuerFromEnv(),
		ttl:    24 * time.Hour,
	}
}

// Generate a JWT token for the authenticated user
func (s *JWTService) Generate(user *models.User) (string, error) {
	claims := models.Claims{
		UserID:        user.ID,
		Subscriptions: user.Subscriptions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func issuerFromEnv() string {
	env := GetAppEnv()
	switch env {
	case AppEnvDev:
		return "Dynamis Dev"
	case AppEnvProd:
		return "Dynamis Prod"
	default:
		return "Dynamis Dev"
	}
}
