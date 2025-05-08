package client

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Fetches the decryption key for the given audio file ID
func fetchDecryptionKey(audioID, userID, deviceID, serverURL string) ([]byte, error) {
	url := fmt.Sprintf("%s/key/%s", serverURL, audioID)

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Make the HTTP GET request for the key
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch decryption key: %v", err)
	}
	defer resp.Body.Close()

	// Check if the server response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch decryption key, server returned status: %d", resp.StatusCode)
	}

	// Read the response body, which contains the expiration and encoded key
	var response string
	if _, err := fmt.Fscanf(resp.Body, "%s", &response); err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Split the response into expiration and base64-encoded key
	var expires int64
	var encodedKey string
	if _, err := fmt.Sscanf(response, "%d:%s", &expires, &encodedKey); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Decode the base64-encoded key
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 key: %v", err)
	}

	// Ensure the key is valid for AES-GCM
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key length: expected 32 bytes, got %d", len(key))
	}

	return key, nil
}

// Decrypts the audio content using the given decryption key
func decryptAudio(encryptedData []byte, key []byte) ([]byte, error) {
	// Generate nonce (12 bytes) from the first 12 bytes of the hash
	nonce := key[:12] // You might want to adjust nonce creation as per your server's method

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// Create AES-GCM cipher instance
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %v", err)
	}

	// Decrypt the data
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt audio: %v", err)
	}

	return plaintext, nil
}

// Fetches and saves the encrypted and decrypted audio content
func fetchAndDecryptAudio(audioID, userID, deviceID, serverURL string) error {
	// Step 1: Fetch the encrypted audio
	encryptedAudio, err := fetchAudio(audioID, serverURL)
	if err != nil {
		return fmt.Errorf("failed to fetch audio: %v", err)
	}

	// Step 2: Fetch the decryption key
	key, err := fetchDecryptionKey(audioID, userID, deviceID, serverURL)
	if err != nil {
		return fmt.Errorf("failed to fetch decryption key: %v", err)
	}

	// Step 3: Decrypt the audio
	decryptedAudio, err := decryptAudio(encryptedAudio, key)
	if err != nil {
		return fmt.Errorf("failed to decrypt audio: %v", err)
	}

	// Step 4: Save the decrypted audio
	err = saveAudio(decryptedAudio, audioID)
	if err != nil {
		return fmt.Errorf("failed to save decrypted audio: %v", err)
	}

	fmt.Println("Decrypted audio saved successfully.")
	return nil
}

// Fetches the encrypted audio
func fetchAudio(audioID, serverURL string) ([]byte, error) {
	url := fmt.Sprintf("%s/audio/%s", serverURL, audioID)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch audio: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch audio, server returned status: %d", resp.StatusCode)
	}

	// Read the audio content
	encryptedData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %v", err)
	}

	return encryptedData, nil
}

// Saves the decrypted audio to a file
func saveAudio(data []byte, audioID string) error {
	fileName := fmt.Sprintf("./audio_%s.dec", audioID)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}
