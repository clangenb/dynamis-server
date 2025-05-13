package database

import (
	"dynamis-server/models"
	"encoding/json"
	"io/ioutil"
	"os"
)

const TracksEnv = "TRACKS_PATH"

func tracksFilePath() string {
	path := os.Getenv(TracksEnv)
	if path == "" {
		path = "data/tracks.json"
	}
	return path
}

// Loads the tracks from the JSON file reading the path from env or from default.
func LoadTracks() ([]models.Track, error) {
	data, err := ioutil.ReadFile(tracksFilePath())
	if err != nil {
		return nil, err
	}

	var tracks []models.Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}
