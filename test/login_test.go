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

// TestAPILogin_WithWrongUsernameAndPassword checks that the login handler reeturns expected JSON when wrong user data is passed
func TestAPILogin_WithWrongUsernameAndPassword(test *testing.T) {
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

	expectedValidationErrors := []handlers.ValidationError{
		handlers.ValidationError{
			Location: []any{"body", "username and password"},
			Message:  "No user with the provided details exists",
			Type:     "value_error",
		},
	}

	expectedResponse := handlers.HTTPValidationError{
		Detail: expectedValidationErrors,
	}

	formData := url.Values{}
	formData.Set("username", "wrongUsername")
	formData.Set("password", "wrongPassword")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.LoginAPI(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnprocessableEntity {
		test.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		test.Fatalf("expected header %s, got %s", "application/json", contentType)
	}

	var actualResponse handlers.HTTPValidationError
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

// TestAPILogin_WithRightUsernameAndWrongPassword checks that the login handler reeturns expected JSON when right username but wrong password is passed
func TestAPILogin_WithRightUsernameAndWrongPassword(test *testing.T) {
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

	expectedValidationErrors := []handlers.ValidationError{
		handlers.ValidationError{
			Location: []any{"body", "username and password"},
			Message:  "No user with the provided details exists",
			Type:     "value_error",
		},
	}

	expectedResponse := handlers.HTTPValidationError{
		Detail: expectedValidationErrors,
	}

	formData := url.Values{}
	formData.Set("username", "username")
	formData.Set("password", "wrongPassword")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.LoginAPI(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnprocessableEntity {
		test.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		test.Fatalf("expected header %s, got %s", "application/json", contentType)
	}

	var actualResponse handlers.HTTPValidationError
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

// TestAPILogin_WithWrongUsernameAndRightPassword checks that the login handler reeturns expected JSON when wrong username but right password is passed
func TestAPILogin_WithWrongUsernameAndRightPassword(test *testing.T) {
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

	expectedValidationErrors := []handlers.ValidationError{
		handlers.ValidationError{
			Location: []any{"body", "username and password"},
			Message:  "No user with the provided details exists",
			Type:     "value_error",
		},
	}

	expectedResponse := handlers.HTTPValidationError{
		Detail: expectedValidationErrors,
	}

	formData := url.Values{}
	formData.Set("username", "wrongUsername")
	formData.Set("password", "password")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.LoginAPI(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnprocessableEntity {
		test.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		test.Fatalf("expected header %s, got %s", "application/json", contentType)
	}

	var actualResponse handlers.HTTPValidationError
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

// TestAPILogin_WithMissingUsernameAndRightPassword checks that the login handler reeturns expected JSON when missing username but right password is passed
func TestAPILogin_WithMissingUsernameAndRightPassword(test *testing.T) {
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

	expectedValidationErrors := []handlers.ValidationError{
		handlers.ValidationError{
			Location: []any{"body", "username"},
			Message:  "Username Field required",
			Type:     "missing",
		},
	}

	expectedResponse := handlers.HTTPValidationError{
		Detail: expectedValidationErrors,
	}

	formData := url.Values{}
	formData.Set("username", "")
	formData.Set("password", "password")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.LoginAPI(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnprocessableEntity {
		test.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		test.Fatalf("expected header %s, got %s", "application/json", contentType)
	}

	var actualResponse handlers.HTTPValidationError
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

// TestAPILogin_WithMissingUsernameAndRightPassword checks that the login handler reeturns expected JSON when right username but missing password is passed
func TestAPILogin_WithRightUsernameAndMissingPassword(test *testing.T) {
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

	expectedValidationErrors := []handlers.ValidationError{
		handlers.ValidationError{
			Location: []any{"body", "password"},
			Message:  "Password Field required",
			Type:     "missing",
		},
	}

	expectedResponse := handlers.HTTPValidationError{
		Detail: expectedValidationErrors,
	}

	formData := url.Values{}
	formData.Set("username", "username")
	formData.Set("password", "")

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/login",
		strings.NewReader(formData.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	responseRecorder := httptest.NewRecorder()

	handlers.LoginAPI(responseRecorder, request)

	if responseRecorder.Code != http.StatusUnprocessableEntity {
		test.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, responseRecorder.Code)
	}

	contentType := responseRecorder.Header().Get("Content-Type")
	if contentType != "application/json" {
		test.Fatalf("expected header %s, got %s", "application/json", contentType)
	}

	var actualResponse handlers.HTTPValidationError
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
