package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// ValidationError describes a single request validation error
type ValidationError struct {
	Location []any  `json:"loc"`
	Message  string `json:"msg"`
	Type     string `json:"type"`
}

// HTTPValidationError contains request validation errors
type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

// AuthResponse describes an authentication API response
type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// RegisterAPI registers a new user
//
// @Summary Register a new user
// @Description Register a new user account
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param password2 formData string false "Password confirmation"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} HTTPValidationError
// @Tags API
// @Router /register [post]
func RegisterAPI(responseWriter http.ResponseWriter, request *http.Request) {
	requestError := request.ParseForm()
	if requestError != nil {
		http.Error(responseWriter, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")
	password2 := request.FormValue("password2")

	validationErrors := []ValidationError{}

	if username == "" {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "username"},
			Message:  "Field required",
			Type:     "missing",
		})
	}

	if email == "" {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "email"},
			Message:  "Field required",
			Type:     "missing",
		})
	}

	if password == "" {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "password"},
			Message:  "Field required",
			Type:     "missing",
		})
	}

	if password2 != "" && password != password2 {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "password2"},
			Message:  "Passwords do not match",
			Type:     "value_error",
		})
	}

	if len(validationErrors) > 0 {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(responseWriter).Encode(HTTPValidationError{
			Detail: validationErrors,
		})

		return
	}

	var existingUsername string
	var existingEmail string

	row := database.QueryRow(
		"SELECT username, email FROM users WHERE username = ? OR email = ? LIMIT 1",
		username,
		email,
	)

	databaseError := row.Scan(&existingUsername, &existingEmail)

	if databaseError != nil && databaseError != sql.ErrNoRows {
		http.Error(
			responseWriter,
			"Failed to check existing user",
			http.StatusInternalServerError,
		)
		return
	}

	if databaseError == nil {
		duplicateErrors := []ValidationError{}

		if existingUsername == username {
			duplicateErrors = append(duplicateErrors, ValidationError{
				Location: []any{"body", "username"},
				Message:  "The username is already taken",
				Type:     "value_error",
			})
		}

		if existingEmail == email {
			duplicateErrors = append(duplicateErrors, ValidationError{
				Location: []any{"body", "email"},
				Message:  "The email is already registered",
				Type:     "value_error",
			})
		}

		if len(duplicateErrors) > 0 {
			responseWriter.Header().Set("Content-Type", "application/json")
			responseWriter.WriteHeader(http.StatusUnprocessableEntity)

			json.NewEncoder(responseWriter).Encode(HTTPValidationError{
				Detail: duplicateErrors,
			})

			return
		}
	}

	passwordHash, hashError := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if hashError != nil {
		http.Error(
			responseWriter,
			"Failed to hash password",
			http.StatusInternalServerError,
		)
		return
	}

	_, databaseError = database.Execute(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username,
		email,
		string(passwordHash),
	)

	if databaseError != nil {
		http.Error(
			responseWriter,
			"Failed to register user",
			http.StatusInternalServerError,
		)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")

	json.NewEncoder(responseWriter).Encode(AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "You were successfully registered and can login now",
	})
}
