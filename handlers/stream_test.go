package handlers_test

import (
	"context"
	"dynamis-server/database"
	"dynamis-server/handlers"
	"dynamis-server/models"
	"dynamis-server/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setEnv(t *testing.T) {
	set(t, database.TracksEnv, "../data/dev/tracks.json")
	set(t, handlers.AudioRootPathEnv, "../data/dev/audio")
	set(t, utils.JWTSecretEnv, "test-secret")
}

func set(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("Failed to set environment variable %v: %v", handlers.AudioRootPathEnv, err)
	}
}

func TestAudioRootPath(t *testing.T) {
	// Test with environment variable set
	os.Setenv(handlers.AudioRootPathEnv, "/custom/audio/path")

	assert.Equal(t, "/custom/audio/path", handlers.AudioRootPath())

	// Test with environment variable unset
	os.Unsetenv(handlers.AudioRootPathEnv)
	assert.Equal(t, "data/audio", handlers.AudioRootPath())
}

func TestAudioFilePath(t *testing.T) {
	// Mock the root path
	os.Setenv(handlers.AudioRootPathEnv, "/custom/audio/path")

	// Test file path joining
	expectedPath := filepath.Join("/custom/audio/path", "test/file.wav")
	assert.Equal(t, expectedPath, handlers.AudioFilePath("test/file.wav"))
}

func TestStreamAudio_ValidAccess(t *testing.T) {
	setEnv(t)

	// Mock claims
	claims := &models.Claims{
		Subscriptions: []string{"free"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks/1", nil)
	r = r.WithContext(models.WithClaims(r.Context(), claims))

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
	setEnv(t)

	// Mock claims
	claims := &models.Claims{
		Subscriptions: []string{"free"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks", nil)
	r = r.WithContext(models.WithClaims(r.Context(), claims))

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
	setEnv(t)

	// Mock claims
	claims := &models.Claims{
		Subscriptions: []string{"premium"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks/1", nil)
	r = r.WithContext(models.WithClaims(r.Context(), claims))

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
//	claims := &models.Claims{
//		Subscriptions: []string{"free"},
//	}
//	r := httptest.NewRequest(http.MethodGet, "/tracks/1", nil)
//	r = r.WithContext(models.WithClaims(r.Context(), claims))
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
