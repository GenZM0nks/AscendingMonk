package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestAbout checks that the about handler renders the expected HTML page
// by looking for parts of it, like 'Our mission' and the image
func TestAbout(test *testing.T) {
	// Below is setup and teardown for htis test
	test.Chdir("..")
	defer test.Chdir("./test")

	request := httptest.NewRequest(http.MethodGet, "/about", nil)
	responseRecorder := httptest.NewRecorder()

	handlers.About(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		test.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			responseRecorder.Code,
		)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		test.Errorf("expected HTML content type, got %q", contentType)
	}

	responseBody := responseRecorder.Body.String()

	if !strings.Contains(responseBody, "Our mission") {
		test.Error("expected the page to contain the mission heading")
	}

	if !strings.Contains(responseBody, `src="/static/monkgroup.png"`) {
		test.Error("expected the page to reference the team image")
	}
}
