package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestAPILogin checks that the login handler reeturns expected JSON. The happy path
func TestAPILogin(test *testing.T) {
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

	expectedResponse := handlers.AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "You were successfully logged in",
	}

	formData := url.Values{}
	formData.Set("username", "username")
	formData.Set("password", "password")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.LoginAPI(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		test.Fatalf("expected status %d, got %d", http.StatusOK, responseRecorder.Code)
	}

	var actualResponse handlers.AuthResponse
	parsingError := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
	if parsingError != nil {
		test.Fatalf("response was not valid JSON: %v", err)
	}

	if !reflect.DeepEqual(actualResponse, expectedResponse) {
		test.Errorf(
			"expected results %+v, got %+v",
			expectedResponse,
			actualResponse,
		)
	}

}
