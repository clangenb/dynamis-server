package client

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestFetchDecryptionKey(t *testing.T) {
	// Mock the server response for decryption key
	mockKey := "mockkeymockkeymockkeymockkeymock"
	mockServer := http.NewServeMux()
	mockServer.HandleFunc("/key/test", func(w http.ResponseWriter, r *http.Request) {
		// Simulate the server's response body: expires:key
		expiration := time.Now().Add(7 * 24 * time.Hour).Unix()
		encodedKey := base64.StdEncoding.EncodeToString([]byte(mockKey))
		response := fmt.Sprintf("%d:%s", expiration, encodedKey)
		w.Write([]byte(response))
	})

	// Run a test server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mockServer,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()
	defer server.Close()

	// Wait a bit for the server to be ready
	time.Sleep(time.Second)

	// Call the function to test
	serverURL := "http://localhost:8080"
	key, err := fetchDecryptionKey("test", "user123", "device123", serverURL)
	if err != nil {
		t.Fatalf("Error fetching decryption key: %v", err)
	}

	// Ensure the key is correct
	expectedKey := mockKey
	if string(key) != expectedKey {
		t.Errorf("Expected key %q, but got %q", expectedKey, string(key))
	}
}

func TestDecryptAudio(t *testing.T) {
	// Sample key and encrypted data (simulated)
	mockKey := []byte("mockkeymockkeymockkeymockkeymock")
	encryptedData := []byte("encrypted-audio-data")

	// Mock the AES decryption (you would normally use actual encryption for real tests)
	decryptedData, err := decryptAudio(encryptedData, mockKey)
	if err != nil {
		t.Fatalf("Error decrypting audio: %v", err)
	}

	// Verify that decrypted data is as expected (for the test, we are simulating)
	expectedDecryptedData := []byte("decrypted-audio-data") // Example expected result
	if string(decryptedData) != string(expectedDecryptedData) {
		t.Errorf("Expected decrypted data %q, but got %q", expectedDecryptedData, decryptedData)
	}
}

func TestSaveAudio(t *testing.T) {
	// Prepare sample data
	audioID := "test"
	data := []byte("some-decrypted-audio")

	// Mock the file system using ioutil.TempFile
	tempFile := "./" + audioID + ".dec"

	// Attempt to save the audio
	err := saveAudio(data, audioID)
	if err != nil {
		t.Fatalf("Error saving audio: %v", err)
	}

	// Verify that the file is saved
	_, err = os.Stat(tempFile)
	if os.IsNotExist(err) {
		t.Errorf("Expected file %s to be created, but it was not found", tempFile)
	}

	// Cleanup
	os.Remove(tempFile)
}

func TestFetchAudio(t *testing.T) {
	// Mock the server response for encrypted audio
	mockServer := http.NewServeMux()
	mockServer.HandleFunc("/audio/test", func(w http.ResponseWriter, r *http.Request) {
		// Simulate the server's response with encrypted audio data
		w.Write([]byte("mock-encrypted-audio"))
	})

	// Run a test server
	server := &http.Server{
		Addr:    ":8081",
		Handler: mockServer,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()
	defer server.Close()

	// Wait a bit for the server to be ready
	time.Sleep(time.Second)

	// Call the function to test
	serverURL := "http://localhost:8081"
	audioData, err := fetchAudio("test", serverURL)
	if err != nil {
		t.Fatalf("Error fetching audio: %v", err)
	}

	// Check if the audio content is as expected
	expectedData := []byte("mock-encrypted-audio")
	if string(audioData) != string(expectedData) {
		t.Errorf("Expected audio data %q, but got %q", expectedData, audioData)
	}
}

func TestFetchAndDecryptAudio(t *testing.T) {
	// Prepare mock server to simulate the full fetch-decrypt process
	mockServer := http.NewServeMux()

	// Mock /audio endpoint
	mockServer.HandleFunc("/audio/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("encrypted-audio"))
	})

	// Mock /key endpoint
	mockServer.HandleFunc("/key/test", func(w http.ResponseWriter, r *http.Request) {
		expiration := time.Now().Add(7 * 24 * time.Hour).Unix()
		mockKey := "mockkeymockkeymockkeymockkeymo"
		encodedKey := base64.StdEncoding.EncodeToString([]byte(mockKey))
		response := fmt.Sprintf("%d:%s", expiration, encodedKey)
		w.Write([]byte(response))
	})

	// Run the mock server
	server := &http.Server{
		Addr:    ":8082",
		Handler: mockServer,
	}
	go server.ListenAndServe()
	defer server.Close()

	// Run the full test to fetch, decrypt, and save audio
	serverURL := "http://localhost:8082"
	if err := fetchAndDecryptAudio("test", "user123", "device123", serverURL); err != nil {
		t.Fatalf("Error during full audio fetch and decrypt: %v", err)
	}
}
