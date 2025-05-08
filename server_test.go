package main

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestServeEncryptedAudio(t *testing.T) {
	// Setup: create a dummy audio file
	audioID := "test"
	testContent := []byte("FAKEAUDIO")
	os.MkdirAll("./audio", os.ModePerm)
	err := os.WriteFile("./audio/"+audioID+".enc", testContent, 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	defer os.Remove("./audio/" + audioID + ".enc")

	req := httptest.NewRequest(http.MethodGet, "/audio/"+audioID, nil)
	rr := httptest.NewRecorder()

	NewRouter().ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	if rr.Body.String() != string(testContent) {
		t.Errorf("Expected body %q, got %q", testContent, rr.Body.String())
	}
}

func TestServeDecryptionKey(t *testing.T) {
	audioID := "song123"
	userID := "userA"
	deviceID := "dev123"

	req := httptest.NewRequest(http.MethodGet, "/key/"+audioID, nil)
	req.Header.Set("X-User-ID", userID)
	req.Header.Set("X-Device-ID", deviceID)
	rr := httptest.NewRecorder()

	NewRouter().ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	parts := strings.Split(body, ":")
	if len(parts) != 2 {
		t.Fatalf("Expected key response format 'expiry:base64key', got %s", body)
	}

	// Validate expiry is in the future
	i, err := strconv.ParseInt(parts[0], 10, 64)
	expiryUnix := time.Unix(i, 0)
	if err != nil || expiryUnix.Before(time.Now()) {
		t.Errorf("Invalid expiry: %s", parts[0])
	}

	// Validate key is base64-decodable
	_, err = base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Errorf("Key is not valid base64: %s", parts[1])
	}
}
