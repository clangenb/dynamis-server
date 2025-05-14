package dynamis_middleware

import (
	"context"
	"dynamis-server/models"
	"dynamis-server/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthorizationHeader = "Authorization"
	BearerPrefix        = "Bearer "
)

// JWTAuth is a middleware that checks for a valid JWT token in the Authorization header.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(AuthorizationHeader)
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, BearerPrefix)
		claims, err := parseToken(tokenStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), models.ClaimsKey, claims))
		next.ServeHTTP(w, r)
	})
}

func parseToken(tokenStr string) (*models.Claims, error) {
	secretKey := utils.GetJWTSecret()

	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", tokenStr)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims from token: %s", tokenStr)
	}

	return claims, nil
}
