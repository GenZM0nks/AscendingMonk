package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

	expectedResult := "You were logged out"

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/logout",
		nil,
	)

	responseRecorder := httptest.NewRecorder()

	handlers.APILogout(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		test.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	var actualResults string
	if err := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResults); err != nil {
		test.Fatalf("response was not valid JSON: %v", err)
	}

	if !reflect.DeepEqual(actualResults, expectedResult) {
		test.Errorf(
			"expected results %+v, got %+v",
			expectedResult,
			actualResults,
		)
	}
}
