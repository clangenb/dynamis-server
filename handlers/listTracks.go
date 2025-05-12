package handlers

import (
	"dynamis-server/database"
	"dynamis-server/middleware"
	"dynamis-server/models"
	"encoding/json"
	"log"
	"net/http"
)

// ListTracks returns a list of audio tracks filtered by subscription tier.
func ListTracks(w http.ResponseWriter, r *http.Request) {
	// Get user claims from the context
	claims := dynamis_middleware.GetClaims(r)

	// Load tracks from the JSON file
	tracks, err := database.LoadTracks("../data/tracks.json")
	if err != nil {
		http.Error(w, "Failed to load tracks", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Filter tracks based on the user's subscription tier
	var filteredTracks []models.Track
	for _, track := range tracks {
		for _, sub := range claims.Subscriptions {
			if track.Tier == sub {
				filteredTracks = append(filteredTracks, track)
				break
			}
		}
	}

	// Return filtered tracks
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredTracks); err != nil {
		http.Error(w, "Failed to encode tracks", http.StatusInternalServerError)
		log.Println(err)
	}
}
