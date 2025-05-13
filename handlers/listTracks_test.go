package handlers_test

import (
	"dynamis-server/database"
	"dynamis-server/handlers"
	"dynamis-server/middleware"
	"dynamis-server/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTracks_ValidFiltering(t *testing.T) {
	setEnv(t)

	// Mock claims
	claims := &dynamis_middleware.Claims{
		Subscriptions: []string{"free", "premium"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks", nil)
	r = r.WithContext(dynamis_middleware.WithClaims(r.Context(), claims))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.ListTracks(rr, r)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	var response []models.Track
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "Track 1", response[0].Title)
	assert.Equal(t, "Track 2", response[1].Title)
}

func TestListTracks_NoMatchingTracks(t *testing.T) {
	setEnv(t)

	// Mock claims
	claims := &dynamis_middleware.Claims{
		Subscriptions: []string{"nonexistent"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks", nil)
	r = r.WithContext(dynamis_middleware.WithClaims(r.Context(), claims))

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.ListTracks(rr, r)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	var response []models.Track
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response)
}

func setEnv(t *testing.T) {
	err := os.Setenv(database.TracksEnv, "../data/tracks.json")
	if err != nil {
		t.Fatalf("Failed to set environment variable: %v", err)
	}
}
