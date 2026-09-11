package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		test.Error("expected the page to contain a search label")
	}
}
