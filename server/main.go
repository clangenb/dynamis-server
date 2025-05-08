package server

import (
	"log"
	"os"
)

func main() {
	// Load the master key from an environment variable (or hardcoded for testing)
	masterKey := []byte(os.Getenv("AUDIO_ENCRYPTION_KEY"))
	if len(masterKey) != 32 {
		log.Fatal("Master key must be 32 bytes long.")
	}

	// Initialize the Deriver with the master key
	encryptor := NewEncryptor(masterKey)

	// Create the server and start it
	server := NewServer("8080", encryptor)
	if err := server.Start(); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
