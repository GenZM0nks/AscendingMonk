package test

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3" // Empty import since we won't use it directly. The string in sql.Open makes Go runtime select our driver.
)

// SetupTestDatabase starts an empty in-memory sqlite3 database for testing
func SetupTestDatabase(test *testing.T) *sql.DB {
	_database, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		test.Fatalf("Failed to open database: %v", err)
	}

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		test.Fatalf("Failed to read schema: %v", err)
	}

	if _, err := _database.Exec(string(schema)); err != nil {
		test.Fatalf("Failed to create schema: %v", err)
	}

	return _database
}
