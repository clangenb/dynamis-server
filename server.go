package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var masterKey = []byte("your-32-byte-master-key---------") // 32 bytes!

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/audio/{id}", serveEncryptedAudio)
	r.HandleFunc("/key/{id}", serveDecryptionKey)
	return r
}

func serveEncryptedAudio(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	filePath := fmt.Sprintf("./audio/%s.enc", id)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, file)
}

func serveDecryptionKey(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := r.Header.Get("X-User-ID")
	device := r.Header.Get("X-Device-ID")

	key := deriveKey(id, user, device)
	expires := time.Now().Add(7 * 24 * time.Hour).Unix()
	response := fmt.Sprintf("%d:%s", expires, base64.StdEncoding.EncodeToString(key))

	w.Write([]byte(response))
}

func deriveKey(audioID, userID, deviceID string) []byte {
	// Use a SHA-256 hash as deterministic and fixed-length seed

	input := []byte(audioID + userID + deviceID)
	hash := sha256.Sum256(input)

	block, _ := aes.NewCipher(masterKey)
	gcm, _ := cipher.NewGCM(block)

	nonce := hash[:gcm.NonceSize()] // Always safe (nonce size is 12)
	return gcm.Seal(nil, nonce, hash[:], nil)[:32]
}
