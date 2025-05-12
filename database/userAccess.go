package database

import (
	"database/sql"
	"dynamis-server/models"
)

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserSubscriptions(userID string) ([]string, error)
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(email string) (*models.User, error) {
	row := DB.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	// Fetch subscriptions
	subscriptions, err := GetUserSubscriptions(user.ID)
	if err != nil {
		return nil, err
	}
	user.Subscriptions = subscriptions
	return &user, nil
}

// InsertUser inserts a new user into the database
func InsertUser(user *models.User) error {
	query := `INSERT INTO users (id, email, password_hash) VALUES (?, ?, ?)`
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Email, user.PasswordHash)
	return err
}

// GetUserSubscriptions fetches the user's subscription tiers from the database
func GetUserSubscriptions(userID string) ([]string, error) {
	rows, err := DB.Query("SELECT tier FROM subscriptions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []string
	for rows.Next() {
		var tier string
		if err := rows.Scan(&tier); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, tier)
	}
	return subscriptions, nil
}
