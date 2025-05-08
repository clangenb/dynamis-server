package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type MockEncryptor struct {
	MockKey []byte
	MockErr error
}

// DeriveKey mocks the key derivation.
func (m *MockEncryptor) DeriveKey(audioID, userID, deviceID string) ([]byte, error) {
	return m.MockKey, m.MockErr
}

func TestServer_serveEncryptedAudio(t *testing.T) {
	// Setup: create a dummy audio file
	audioID := "test"
	testContent := []byte("FAKEAUDIO")
	os.MkdirAll("./audio", os.ModePerm)
	audioPath := "./audio/" + audioID + ".enc"
	err := os.WriteFile(audioPath, testContent, 0644)
	if err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	fmt.Println("Audio Path: " + audioPath)

	defer os.Remove(audioPath)

	req := httptest.NewRequest(http.MethodGet, "/audio/"+audioID, nil)
	rr := httptest.NewRecorder()

	mockEncryptor := &MockEncryptor{MockKey: []byte("mockkeymockkeymockkeymockkeymo")}
	server := NewServer("8080", mockEncryptor)

	server.Router().ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
	if rr.Body.String() != string(testContent) {
		t.Errorf("Expected body %q, got %q", testContent, rr.Body.String())
	}
}

func TestServer_serveDecryptionKey(t *testing.T) {
	mockKey := []byte("mockkeymockkeymockkeymockkeymo")
	mockEncryptor := &MockEncryptor{MockKey: mockKey}
	server := NewServer("8080", mockEncryptor)

	req := httptest.NewRequest(http.MethodGet, "/key/audio123", nil)
	req.Header.Set("X-User-ID", "user123")
	req.Header.Set("X-Device-ID", "device123")
	rr := httptest.NewRecorder()

	server.Router().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected status 200 OK, got %d", rr.Code)
	}

	// Extract and validate the response format: "timestamp:base64key"
	resp := rr.Body.String()
	var expires int64
	var encodedKey string
	_, err := fmt.Sscanf(resp, "%d:%s", &expires, &encodedKey)
	if err != nil {
		t.Fatalf("Invalid response format: %s", resp)
	}

	decodedKey, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		t.Fatalf("Key is not valid base64: %v", err)
	}
	if string(decodedKey) != string(mockKey) {
		t.Errorf("Expected decoded key %q, got %q", mockKey, decodedKey)
	}
}
