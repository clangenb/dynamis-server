package database

import (
	"dynamis-server/models"
	"encoding/json"
	"io/ioutil"
)

func LoadTracks(filePath string) ([]models.Track, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tracks []models.Track
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, err
	}

	return tracks, nil
}
