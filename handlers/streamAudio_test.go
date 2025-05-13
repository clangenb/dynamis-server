package handlers_test

import (
	"context"
	"dynamis-server/handlers"
	"dynamis-server/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestStreamAudio_ValidAccess(t *testing.T) {
	// Mock claims
	claims := &dynamis_middleware.Claims{
		Subscriptions: []string{"free"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks/1", nil)
	r = r.WithContext(dynamis_middleware.WithClaims(r.Context(), claims))

	// Mock URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("trackID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.StreamAudio(rr, r)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "audio/wav", rr.Header().Get("Content-Type"))
	assert.Equal(t, "test1\n", string(rr.Body.Bytes()))
}

func TestStreamAudio_TrackNotFound(t *testing.T) {
	// Mock claims
	claims := &dynamis_middleware.Claims{
		Subscriptions: []string{"free"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks", nil)
	r = r.WithContext(dynamis_middleware.WithClaims(r.Context(), claims))

	// Mock URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("trackID", "999")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.StreamAudio(rr, r)

	// Assert response
	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "Track not found")
}

func TestStreamAudio_UnauthorizedAccess(t *testing.T) {
	// Mock claims
	claims := &dynamis_middleware.Claims{
		Subscriptions: []string{"premium"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks/1", nil)
	r = r.WithContext(dynamis_middleware.WithClaims(r.Context(), claims))

	// Mock URL parameter
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("trackID", "1")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.StreamAudio(rr, r)

	// Assert response
	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), "Unauthorized access to this track")
}

// Todo: how to mock?
//func TestStreamAudio_FileOpenError(t *testing.T) {
//	// Mock claims
//	claims := &dynamis_middleware.Claims{
//		Subscriptions: []string{"free"},
//	}
//	r := httptest.NewRequest(http.MethodGet, "/tracks/1", nil)
//	r = r.WithContext(dynamis_middleware.WithClaims(r.Context(), claims))
//
//	// Mock URL parameter
//	rctx := chi.NewRouteContext()
//	rctx.URLParams.Add("trackID", "1")
//	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
//
//	// Create response recorder
//	rr := httptest.NewRecorder()
//
//	// Call handler
//	handlers.StreamAudio(rr, r)
//
//	// Assert response
//	assert.Equal(t, http.StatusInternalServerError, rr.Code)
//	assert.Contains(t, rr.Body.String(), "Failed to open audio file")
//}
