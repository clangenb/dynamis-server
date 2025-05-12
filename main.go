package main

import (
	"dynamis-server/database"
	"dynamis-server/handlers"
	"dynamis-server/middleware"
	"github.com/go-chi/chi/v5"

	"log"
	"net/http"
)

func main() {
	// Initialize database
	database.Init()

	// Create a new router
	r := chi.NewRouter()

	// Routes
	r.Post("/login", handlers.LoginHandler) // Login endpoint

	// Secure routes (need JWT auth)
	r.With(dynamis_middleware.JWTAuth).Get("/audio", handlers.ListTracks)
	r.With(dynamis_middleware.JWTAuth).Get("/audio/{trackID}", handlers.StreamAudio)

	// Start the server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
