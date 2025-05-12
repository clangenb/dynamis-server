package models

import (
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID            string
	Email         string
	PasswordHash  string
	Subscriptions []string
}

// ComparePassword compares the user's stored password hash with the given password
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
