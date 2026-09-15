package test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Empty import since we won't use it directly. The string in sql.Open makes Go runtime select our driver.
)

// SetupTestDatabase starts an empty in-memory sqlite3 database for testing
func SetupTestDatabase(test *testing.T) *sql.DB {
	_database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		test.Fatalf("Failed to open database: %v", err)
	}

	schema := `
	CREATE TABLE pages (
		title TEXT PRIMARY KEY UNIQUE,
		url TEXT NOT NULL UNIQUE,
		language TEXT NOT NULL CHECK(language IN ('en', 'da')) DEFAULT 'en',
		last_updated TIMESTAMP,
		content TEXT NOT NULL
	)`

	if _, err := _database.Exec(schema); err != nil {
		test.Fatalf("Failed to create schema: %v", err)
	}

	return _database
}
