package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var DB *sql.DB
var dbPath string

// Ensure database and table exist
func Init(dbFile string) error {
	dbPath = dbFile
	return initDB()
}

func initDB() error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
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
