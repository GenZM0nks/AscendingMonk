package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestregisterPage checks that the register handler renders the expected HTML page
// by looking for parts of it, like 'Sign Up'
func TestRegisterPage(test *testing.T) {
	// Below is setup and teardown for htis test
	test.Chdir("..")
	defer test.Chdir("./test")

	request := httptest.NewRequest(http.MethodGet, "/register", nil)
	responseRecorder := httptest.NewRecorder()

	handlers.Register(responseRecorder, request)

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

	if !strings.Contains(responseBody, "Sign Up") {
		test.Error("expected the page to contain the Sing Up heading")
	}

}