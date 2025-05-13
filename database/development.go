package database

import "dynamis-server/models"

// InitDevDb initializes the development database with test data
func InitDevDb(path string) error {
	err := Init(path)
	if err != nil {
		return err
	}

	err = InsertUser(alice())
	if err != nil {
		return err
	}

	err = InsertUser(bob())
	if err != nil {
		return err
	}

	return nil
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
