package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestLogin checks that the login handler returns expected JSON. The happy path
func TestLogin(test *testing.T) {
	test.Chdir("..")

	password, hashError := bcrypt.GenerateFromPassword(
		[]byte("password"),
		bcrypt.DefaultCost,
	)

	if hashError != nil {
		test.Fatal("Failed to hash password")
	}

	database.SetDatabase(SetupTestDatabase(test)) // SetupTestDatabase already opens the DB, so no Connect()

	_, err := database.Execute("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		"username",
		"email",
		password)

	if err != nil {
		test.Fatalf("Failed to insert the test user entry into the database\nFrom SQLite: %s.", err.Error())
	}

	test.Cleanup(func() {
		test.Chdir("./test")
		database.Execute("DROP TABLE users")
		database.Close()
	})

	formData := url.Values{}
	formData.Set("username", "username")
	formData.Set("password", "password")

	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.Login(responseRecorder, request)

	if responseRecorder.Code != http.StatusSeeOther {
		test.Fatalf("expected status %d, got %d", http.StatusSeeOther, responseRecorder.Code)
	}

	responseBody := responseRecorder.Header().Get("Set-Cookie")

	if !strings.Contains(responseBody, "Path=/;") {
		test.Error("expected the body to contain Path=/;")
	}
	if !strings.Contains(responseBody, "Expires=") {
		test.Error("expected the body to contain Expires=")
	}
	if !strings.Contains(responseBody, "HttpOnly") {
		test.Error("expected the body to contain HttpOnly")
	}
}

// TestLogin_WithWrongUsernameAndPassword checks that the login handler returns expected JSON when wrong user data is passed
func TestLogin_WithWrongUsernameAndPassword(test *testing.T) {
	test.Chdir("..")

	password, hashError := bcrypt.GenerateFromPassword(
		[]byte("password"),
		bcrypt.DefaultCost,
	)

	if hashError != nil {
		test.Fatal("Failed to hash password")
	}

	database.SetDatabase(SetupTestDatabase(test)) // SetupTestDatabase already opens the DB, so no Connect()

	_, err := database.Execute("INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		"username",
		"email",
		password)

	if err != nil {
		test.Fatalf("Failed to insert the test user entry into the database\nFrom SQLite: %s.", err.Error())
	}

	test.Cleanup(func() {
		test.Chdir("./test")
		database.Execute("DROP TABLE users")
		database.Close()
	})

	formData := url.Values{}
	formData.Set("username", "wrongUsername")
	formData.Set("password", "wrongPassword")

	request := httptest.NewRequest(
		http.MethodPost,
		"/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.Login(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		test.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "text/html; charset=utf-8" {
		test.Errorf("expected HTML content type, got %q", contentType)
	}

	responseBody := responseRecorder.Body.String()

	if !strings.Contains(responseBody, "Log In") {
		test.Error("expected the page to contain the Log In heading")
	}
	if !strings.Contains(responseBody, "No user with the provided details exists") {
		test.Error("expected the page to contain the 'No user with the provided details exists error' message")
	}

}
