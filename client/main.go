package client

import "fmt"

func main() {
	// Example usage
	serverURL := "http://localhost:8080" // Replace with your actual server URL
	audioID := "test"                    // Example audio ID
	userID := "user123"                  // Example user ID
	deviceID := "device123"              // Example device ID

	// Fetch and decrypt the audio
	if err := fetchAndDecryptAudio(audioID, userID, deviceID, serverURL); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
