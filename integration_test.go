package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/clangenb/dynamis/client"
	"github.com/clangenb/dynamis/server"
)

// ----- Mock Deriver -----

type MockDeriver struct {
	Key []byte
}

func (m *MockDeriver) DeriveKey(audioID, userID, deviceID string) ([]byte, error) {
	return m.Key, nil
}

// ----- Test -----

func TestFullIntegration(t *testing.T) {
	// Setup constants
	audioID := "testaudio"
	userID := "user123"
	deviceID := "device456"
	key := []byte("mockkeymockkeymockkeymockkeymock") // 32 bytes AES-256
	plaintext := []byte("this is a test audio stream")

	// Encrypt the audio
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatal(err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatal(err)
	}
	nonce := key[:gcm.NonceSize()]
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Write encrypted file to disk
	os.MkdirAll("./audio", 0755)
	audioPath := filepath.Join("audio", audioID+".enc")
	err = os.WriteFile(audioPath, ciphertext, 0644)
	if err != nil {
		t.Fatalf("failed to write encrypted audio: %v", err)
	}
	defer os.Remove(audioPath)

	// Start the real server using httptest
	mockDeriver := &MockDeriver{Key: key}
	srv := server.NewServer("8080", mockDeriver)

	router := srv.Router() // you’ll need to expose this in server.go: `func (s *Server) Router() http.Handler`
	ts := httptest.NewServer(router)
	defer ts.Close()

	// === CLIENT FLOW ===

	// 1. Fetch decryption key
	gotKey, err := client.FetchDecryptionKey(audioID, userID, deviceID, ts.URL)
	if err != nil {
		t.Fatalf("failed to fetch decryption key: %v", err)
	}

	// 2. Fetch encrypted audio
	gotEncrypted, err := client.FetchAudio(audioID, ts.URL)
	if err != nil {
		t.Fatalf("failed to fetch encrypted audio: %v", err)
	}

	// 3. Decrypt it
	decrypted, err := client.DecryptAudio(gotEncrypted, gotKey)
	if err != nil {
		t.Fatalf("failed to decrypt audio: %v", err)
	}

	// 4. Compare
	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("decrypted data mismatch.\nWant: %s\nGot:  %s", plaintext, decrypted)
	}
}
