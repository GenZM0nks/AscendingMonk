package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3" // Empty import since we won't use it directly. The string in sql.Open makes Go runtime select our driver.
)

var database *sql.DB // Remove this global and refactor if / when we build a service layer

// Connect initializes a database connection and sets `database` equal to its handle on success.
// On failure, print errors with details.
// This method expects the environment variable `DB_PATH` to be set at runtime in a .env file.
func Connect() error {
	database, err := sql.Open("sqlite3", os.Getenv("DB_PATH"))

	if err != nil {
		return fmt.Errorf("Failed to open database: %w", err)
	}

	err = database.Ping()
	if err != nil {
		return fmt.Errorf("Failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection if `database` is initialized,
// blocking new requests and shutting down.
func Close() error {
	if database != nil {
		return database.Close()
	}

	return nil
}

// Query executes a parameterized query (like a SELECT), `query`,  using `args`
// and returns as many rows as match it.
func Query(query string, args ...any) (*sql.Rows, error) {
	if database == nil {
		return nil, fmt.Errorf("Database not initialized. Call Connect()")
	}

	return database.Query(query, args...)
}

// QueryRow executes a parameterized query (like a SELECT), `query`, using `args“
// which returns a single row
func QueryRow(query string, args ...any) *sql.Row {
	return database.QueryRow(query, args...)
}

// Execute executes a parameterized (INSERT, UPDATE, DELETE) query, `query`, using `args`
// without returning any rows
func Execute(query string, args ...any) (sql.Result, error) {
	if database == nil {
		return nil, fmt.Errorf("Database not initialized")
	}

	return database.Exec(query, args...)
}
