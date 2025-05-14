package utils

import (
	"dynamis-server/middleware"
	"dynamis-server/models"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var signingKey = []byte("your-secret-key")

type JWTService struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

const JWTSecretEnv = "JWT_SECRET"

func NewJWTService() *JWTService {
	secret := os.Getenv(JWTSecretEnv)
	if secret == "" {
		panic("JWT_SECRET is not set") // or use a default/fallback for dev
	}
	return &JWTService{
		secret: []byte(secret),
		issuer: issuerFromEnv(),
		ttl:    24 * time.Hour,
	}
}

// GenerateJWT generates a JWT token for the authenticated user
func (s *JWTService) Generate(user *models.User) (string, error) {
	claims := dynamis_middleware.Claims{
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
