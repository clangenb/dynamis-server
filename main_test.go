package main_test

import (
	"bytes"
	"dynamis-server"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAliceCanAccessTracks(t *testing.T) {
	// Start the server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		main.Main() // Call the main function to start the app
	}))
	defer server.Close()

	// Step 1: Login as Alice
	loginPayload := `{"email": "alice@example.com", "password": "alice"}`
	loginResp, err := http.Post(server.URL+"/login", "application/json", bytes.NewBuffer([]byte(loginPayload)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginResp.StatusCode)

	var loginData map[string]string
	err = json.NewDecoder(loginResp.Body).Decode(&loginData)
	//assert.NoError(t, err)

	log.Println(loginData)

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
	assert.NotEmpty(t, tracks, "Tracks should not be empty")
}
