package handlers

import (
	"dynamis-server/database"
	jwtmiddleware "dynamis-server/middleware"
	"dynamis-server/models"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"os"
)

// StreamAudio streams the requested audio file if the user has access.
func StreamAudio(w http.ResponseWriter, r *http.Request) {
	// Get user claims from the context
	claims := jwtmiddleware.GetClaims(r)

	// Get the track ID from the URL parameters
	trackID := chi.URLParam(r, "trackID")
	log.Println("Requested track ID:", trackID)

	// Load tracks from the JSON file
	tracks, err := database.LoadTracks()
	if err != nil {
		http.Error(w, "Failed to load tracks", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Find the requested track
	var track *models.Track
	for _, t := range tracks {
		if t.ID == trackID {
			track = &t
			break
		}
	}

	// If track not found
	if track == nil {
		http.Error(w, "Track not found", http.StatusNotFound)
		return
	}

	// Check if the user has access to the track based on subscription
	hasAccess := false
	for _, sub := range claims.Subscriptions {
		if track.Tier == sub {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		http.Error(w, "Unauthorized access to this track", http.StatusForbidden)
		return
	}

	// Open the audio file for streaming
	audioFile, err := os.Open(track.FilePath)
	if err != nil {
		http.Error(w, "Failed to open audio file", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer audioFile.Close()

	// Set the correct content type for WAV files
	w.Header().Set("Content-Type", "audio/wav")

	// Stream the audio file
	w.WriteHeader(http.StatusOK)
	_, err = io.Copy(w, audioFile)
	if err != nil {
		http.Error(w, "Failed to stream audio", http.StatusInternalServerError)
		log.Println(err)
	}
}
