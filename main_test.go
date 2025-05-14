package main_test

import (
	"bytes"
	"dynamis-server"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliceCanAccessTracks(t *testing.T) {
	main.InitializeApp()
	// Start the test server with the application's router
	server := httptest.NewServer(main.SetupRouter())
	defer server.Close()

	// Step 1: Login as Alice
	loginPayload := `{"email": "alice@example.com", "password": "alice"}`
	loginResp, err := http.Post(server.URL+"/login", "application/json", bytes.NewBuffer([]byte(loginPayload)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	var loginData map[string]string
	err = json.NewDecoder(loginResp.Body).Decode(&loginData)
	assert.NoError(t, err)

	token, ok := loginData["token"]
	assert.True(t, ok, "Token not found in login response")

	// Step 2: Access tracks with Alice's token
	req, err := http.NewRequest(http.MethodGet, server.URL+"/tracks", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	tracksResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, tracksResp.StatusCode)

	var tracks []map[string]interface{}
	err = json.NewDecoder(tracksResp.Body).Decode(&tracks)
	assert.NoError(t, err)

	assert.Equal(t, expectedTracks(), tracks, "Response does not match the expected output")
}

func expectedTracks() []map[string]interface{} {
	jsonData := `
	[
		{
			"id": "1",
			"title": "Track 1",
			"file_path": "audio1-test.wav",
			"tier": "free"
		},
		{
			"id": "2",
			"title": "Track 2",
			"file_path": "audio2-test.wav",
			"tier": "premium"
		},
		{
			"id": "3",
			"title": "Track 3",
			"file_path": "audio3-test.wav",
			"tier": "vip"
		}
	]`

	var tracks []map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &tracks)
	if err != nil {
		panic("Error unmarshaling JSON:" + err.Error())
	}
	return tracks
}
