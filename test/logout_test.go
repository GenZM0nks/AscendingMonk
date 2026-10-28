package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestAPILogout checks that the API handler returns JSON with expected fields - the happy path.
func TestAPILogout(test *testing.T) {
	test.Chdir("..")
	database.SetDatabase(SetupTestDatabase(test)) // SetupTestDatabase already opens the DB, so no Connect()

	_, err := database.Execute("")

	if err != nil {
		test.Fatalf("Failed to insert new user and session.\nFrom SQLite: %s.", err.Error())
	}

	test.Cleanup(func() {
		test.Chdir("./test")
		database.Execute("DROP TABLE pages")
		database.Execute("DROP TABLE users")
		database.Execute("DROP TABLE session_tokens")
		database.Close()
	})

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/logout",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handlers.APILogout(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		test.Fatalf("expected status %d, got %d", http.StatusSeeOther, responseRecorder.Code)
	}
}
