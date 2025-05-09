package server

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Server holds the configuration and dependencies for the HTTP server.
type Server struct {
	Port    string
	Deriver KeyDeriver
	router  *mux.Router
}

// NewServer initializes a new server with the provided port and encryptor.
func NewServer(port string, encryptor KeyDeriver) *Server {
	s := &Server{
		Port:    port,
		Deriver: encryptor,
		router:  mux.NewRouter(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.HandleFunc("/audio/{id}", s.serveEncryptedAudio)
	s.router.HandleFunc("/key/{id}", s.serveDerivedCEK)
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", s.Port), s.router)
}

func (s *Server) Router() *mux.Router {
	return s.router
}

// serveEncryptedAudio serves encrypted audio files from the server.
func (s *Server) serveEncryptedAudio(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	filePath := fmt.Sprintf("./audio/%s.enc", id)

	fmt.Println("File Path: " + filePath)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", 404)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, file)
}

// serveDerivedCEK serves a decryption key for the requested audio.
func (s *Server) serveDerivedCEK(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user := r.Header.Get("X-User-ID")
	device := r.Header.Get("X-Device-ID")

	// Derive the key using the Deriver.
	key, err := s.Deriver.DeriveCEK(id, user, device)
	if err != nil {
		http.Error(w, "Failed to derive key", http.StatusInternalServerError)
		return
	}

	// Set an expiration for the key (7 days from now)
	expires := time.Now().Add(7 * 24 * time.Hour).Unix()

	// Return the key and expiration time in the response.
	response := fmt.Sprintf("%d:%s", expires, base64.StdEncoding.EncodeToString(key))
	w.Write([]byte(response))
}
