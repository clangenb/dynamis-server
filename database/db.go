package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var DB *sql.DB

const DBPathEnv = "DB_PATH"

func dbFilePath() string {
	path := os.Getenv(DBPathEnv)
	if path == "" {
		path = "data/db.sqlite"
	}
	return path
}

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", dbFilePath())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}

	log.Println("Database connected.")
	createTables()
	return nil
}

func createTables() {
	userSQL := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);`

	subscriptionSQL := `CREATE TABLE IF NOT EXISTS subscriptions (
		user_id TEXT NOT NULL,
		tier TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(userSQL); err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	if _, err := DB.Exec(subscriptionSQL); err != nil {
		log.Fatalf("Failed to create subscriptions table: %v", err)
	}
}
