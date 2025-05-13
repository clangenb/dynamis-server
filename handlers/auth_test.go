package handlers_test

import (
	"bytes"
	"dynamis-server/database"
	"dynamis-server/handlers"
	"dynamis-server/models"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginHandler_ValidCredentials(t *testing.T) {
	setupTestDB(t)

	password := "validpassword"
	hash, _ := models.HashPassword(password)
	user := &models.User{
		ID:            "123",
		Email:         "test@example.com",
		PasswordHash:  hash,
		Subscriptions: []string{"sub1", "sub2"},
	}
	database.InsertUser(user)

	// Create request
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "validpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.LoginHandler(rr, req)

	// Assert response
	assert.Equal(t, http.StatusOK, rr.Code)
	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	assert.NotEmpty(t, response["token"])
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	setupTestDB(t)

	password := "validpassword"
	hash, _ := models.HashPassword(password)
	user := &models.User{
		ID:            "123",
		Email:         "test@example.com",
		PasswordHash:  hash,
		Subscriptions: []string{"sub1", "sub2"},
	}
	database.InsertUser(user)

	// Create request
	requestBody, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "wrongpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.LoginHandler(rr, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestLoginHandler_InvalidRequest(t *testing.T) {
	// Create request with invalid JSON
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handlers.LoginHandler(rr, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func setupTestDB(t *testing.T) {
	set(t, database.DBPathEnv, ":memory:")

	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init test db: %v", err)
	}
}
