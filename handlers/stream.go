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
	"path/filepath"
)

const AudioRootPathEnv = "AUDIO_ROOT_PATH"

func audioRootPath() string {
	if path, ok := os.LookupEnv(AudioRootPathEnv); ok && path != "" {
		return path
	}
	return "data/audio"
}

func audioFilePath(localPath string) string {
	return filepath.Join(audioRootPath(), localPath)
}

// StreamAudio streams the requested audio file if the user has access.
func StreamAudio(w http.ResponseWriter, r *http.Request) {
	claims := jwtmiddleware.GetClaims(r)
	trackID := chi.URLParam(r, "trackID")
	log.Printf("Requested track ID: %s", trackID)

	tracks, err := database.LoadTracks()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to load tracks", err)
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

	if track == nil {
		respondWithError(w, http.StatusNotFound, "Track not found", nil)
		return
	}

	if !userHasAccess(claims.Subscriptions, track.Tier) {
		respondWithError(w, http.StatusForbidden, "Unauthorized access to this track", nil)
		return
	}

	audioFile, err := os.Open(audioFilePath(track.FilePath))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to open audio file", err)
		return
	}
	defer audioFile.Close()

	w.Header().Set("Content-Type", "audio/wav")
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, audioFile); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to stream audio", err)
	}
}

func userHasAccess(subscriptions []string, requiredTier string) bool {
	for _, sub := range subscriptions {
		if sub == requiredTier {
			return true
		}
	}
	return false
}
