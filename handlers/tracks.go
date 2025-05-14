package handlers

import (
	"dynamis-server/database"
	"dynamis-server/models"
	"encoding/json"
	"net/http"
)

// ListTracks returns a list of audio tracks filtered by subscription tier.
func ListTracks(w http.ResponseWriter, r *http.Request) {
	claims := models.GetClaims(r)

	tracks, err := database.LoadTracks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load tracks", err)
		return
	}

	// Create a map for quick subscription lookup
	subscriptionSet := make(map[string]struct{})
	for _, sub := range claims.Subscriptions {
		subscriptionSet[sub] = struct{}{}
	}

	// Filter tracks based on the user's subscriptions
	var filteredTracks []models.Track
	for _, track := range tracks {
		if _, ok := subscriptionSet[track.Tier]; ok {
			filteredTracks = append(filteredTracks, track)
		}
	}

	// Respond with the filtered tracks
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredTracks); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to encode tracks", err)
	}
}
