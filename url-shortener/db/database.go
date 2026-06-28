package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() *sql.DB {
	var err error

	// Ensure data directory exists
	os.MkdirAll("data", os.ModePerm)

	DB, err = sql.Open("sqlite3", "./data/url_shortener.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	createTables()

	return DB
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	urlTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		original_url TEXT NOT NULL,
		short_code TEXT NOT NULL UNIQUE,
		expiry DATETIME,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(urlTable)
	if err != nil {
		log.Fatal("Failed to create urls table:", err)
	}
}
