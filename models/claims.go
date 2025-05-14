package models

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
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

func GetClaims(r *http.Request) *Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*Claims)
	return claims
}

func WithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}
