package handlers_test

import (
	"dynamis-server/handlers"
	"dynamis-server/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTracks_ValidFiltering(t *testing.T) {
	setEnv(t)

	// Mock claims
	claims := &models.Claims{
		Subscriptions: []string{"free", "premium"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks", nil)
	r = r.WithContext(models.WithClaims(r.Context(), claims))

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
	claims := &models.Claims{
		Subscriptions: []string{"nonexistent"},
	}
	r := httptest.NewRequest(http.MethodGet, "/tracks", nil)
	r = r.WithContext(models.WithClaims(r.Context(), claims))

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
