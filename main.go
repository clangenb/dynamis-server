package main

import (
	"dynamis-server/database"
	"dynamis-server/handlers"
	"dynamis-server/middleware"
	"dynamis-server/utils"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

// Curl the dev db with:
//
// 1. Get the JWT token for Alice
// curl -X POST http://localhost:8080/login \                                                                                                                         ok | 10:36:30
//   -H "Content-Type: application/json" \
//   -d '{"email": "alice@example.com", "password": "alice"}'
//2. Use the token to access the tracks
// curl -X GET http://localhost:8080/tracks -H "Authorization: Bearer <JWT_TOKEN>"
// 3. Use the token to access a specific track
//  curl -X GET http://localhost:8080/tracks/<trackID> -H "Authorization: Bearer <JWT_TOKEN>"

// SetupRouter initializes the router with all routes and middleware.
func SetupRouter() http.Handler {
	// Create a new router
	r := chi.NewRouter()

	// Routes
	r.Post("/login", handlers.LoginHandler) // Login endpoint

	// Secure routes (need JWT auth)
	r.With(dynamis_middleware.JWTAuth).Get("/tracks", handlers.ListTracks)
	r.With(dynamis_middleware.JWTAuth).Get("/tracks/{trackID}", handlers.StreamAudio)

	return r
}

func InitializeApp() {
	// Initialize database
	database.InitDB()

	switch utils.GetAppEnv() {
	case utils.AppEnvDev:
		log.Println("Running in development mode")
		err := database.SetupDevEntries()
		if err != nil {
			log.Fatalf("Failed to set up dev entries: %v", err)
		}
	case utils.AppEnvProd:
		log.Println("Running in production mode")
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (continuing anyway)")
	}

	InitializeApp()

	// Start the server
	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", SetupRouter())
}
