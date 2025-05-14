package handlers

import (
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string, err error) {
	http.Error(w, message, statusCode)
	if err != nil {
		log.Printf("Error: %s - %v", message, err)
	} else {
		log.Printf("Error: %s", message)
	}
}
