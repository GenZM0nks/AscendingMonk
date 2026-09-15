package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestSearchPage checks that the search handler renders the expected HTML page
// by looking for the search label above the input field.
func TestSearchPage(test *testing.T) {
	test.Chdir("..")
	defer test.Chdir("./test")

	request := httptest.NewRequest(http.MethodGet, "/search", nil)
	responseRecorder := httptest.NewRecorder()

	handlers.Search(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		test.Fatalf(
			"expected status %d, got %d\nBody: %s",
			http.StatusOK,
			responseRecorder.Code,
			responseRecorder.Body.String(),
		)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		test.Errorf("expected HTML content type, got %q", contentType)
	}

	responseBody := responseRecorder.Body.String()

	if !strings.Contains(responseBody, "Search") {
		test.Error("expected the page to contain the Search heading")
	}
}

// TestAPISearch checks that the API handler returns JSON with expected fields - the happy path.
func TestAPISearch(test *testing.T) {
	test.Chdir("..")
	database.SetDatabase(SetupTestDatabase(test)) // SetupTestDatabase already opens the DB, so no Connect()

	_, err := database.Execute("INSERT INTO pages (title, url, language, content) VALUES (?, ?, ?, ?)",
		"Fortran",
		"https://en.wikipedia.org/wiki/Fortran",
		"en",
		"Fortran is a general-purpose, compiled imperative programming language.")

	if err != nil {
		test.Fatalf("Failed to insert the Fortran entry into the database\nFrom SQLite: %s.", err.Error())
	}

	test.Cleanup(func() {
		test.Chdir("./test")
		database.Execute("DROP TABLE pages")
		database.Close()
	})

	expectedResults := []handlers.SearchResult{
		{
			Title:       "Fortran",
			URL:         "https://en.wikipedia.org/wiki/Fortran",
			Description: "Fortran is a general-purpose, compiled imperative programming language.",
		},
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/search?q=Fortran&language=en",
		nil,
	)

	response := httptest.NewRecorder()

	handlers.APISearch(response, request)

	if response.Code != http.StatusOK {
		test.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	var actualResults []handlers.SearchResult
	if err := json.Unmarshal(response.Body.Bytes(), &actualResults); err != nil {
		test.Fatalf("response was not valid JSON: %v", err)
	}

	if !reflect.DeepEqual(actualResults, expectedResults) {
		test.Errorf(
			"expected results %+v, got %+v",
			expectedResults,
			actualResults,
		)
	}
}
