package database

import (
	"dynamis-server/models"
	"log"
)

// InitDevDB initializes the development database with test data
func SetupDevEntries() error {

	err := insertIfNotExists(alice())
	if err != nil {
		return err
	}

	err = insertIfNotExists(bob())
	if err != nil {
		return err
	}

	return nil
}

func insertIfNotExists(user *models.User) error {
	maybeUser, _ := GetUserByEmail(user.Email)
	if maybeUser != nil {
		log.Printf("Skip setting up existing dev user %s", user.Email)
		return nil
	}

	return InsertUser(user)
}

const AlicePwd = "alice"
const BobPwd = "bob"

func alice() *models.User {
	hash, _ := models.HashPassword(AlicePwd)
	return &models.User{
		ID:            "1",
		Email:         "alice@example.com",
		PasswordHash:  hash,
		Subscriptions: []string{"sub1", "sub2"},
	}
}

func bob() *models.User {
	hash, _ := models.HashPassword(BobPwd)
	return &models.User{
		ID:            "2",
		Email:         "bob@example.com",
		PasswordHash:  hash,
		Subscriptions: []string{"sub1", "sub2"},
	}
}
