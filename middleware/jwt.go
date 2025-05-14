package dynamis_middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	ClaimsKey contextKey = "claims"
)

type Claims struct {
	UserID        string   `json:"sub"`
	Subscriptions []string `json:"subscriptions"`
	jwt.RegisteredClaims
}

const jWTSecretEnv = "JWT_SECRET"

func getJwtSecret() string {
	secret := os.Getenv(jWTSecretEnv)
	if secret == "" {
		panic("JWT_SECRET is not set") // or use a default/fallback for dev
	}
	return secret
}

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(getJwtSecret()), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), ClaimsKey, claims))
		next.ServeHTTP(w, r)
	})
}

func GetClaims(r *http.Request) *Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*Claims)
	return claims
}

// WithClaims adds claims to the context
func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}
