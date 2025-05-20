package data

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite" // SQLite driver
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Database connection established.")
	createTables()
}

func createTables() {
	createUserTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"username" TEXT NOT NULL UNIQUE,
		"email" TEXT NOT NULL UNIQUE,
		"password_hash" TEXT NOT NULL,
		"created_at" DATETIME DEFAULT CURRENT_TIMESTAMP,
		"updated_at" DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	statement, err := DB.Prepare(createUserTableSQL)
	if err != nil {
		log.Fatalf("Error preparing create user table statement: %v", err)
	}
	_, err = statement.Exec()
	if err != nil {
		log.Fatalf("Error executing create user table statement: %v", err)
	}
	log.Println("Users table created or already exists.")
}

// CloseDB closes the database connection.
func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed.")
	}
}
