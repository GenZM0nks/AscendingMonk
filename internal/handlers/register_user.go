package handlers

import (
	"database/sql"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"golang.org/x/crypto/bcrypt"
)

// registerUser validates registration data and creates a new user
//
// Validation problems are returned as []ValidationError.
// Unexpected problems, such as database or hashing failures,
// are returned as error
func registerUser(
	username string,
	email string,
	password string,
	password2 string,
) (uint64, []ValidationError, error) {
	validationErrors := []ValidationError{}

	// Check required fields.
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

	// Validate password confirmation when provided
	if password2 != "" && password != password2 {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "password2"},
			Message:  "Passwords do not match",
			Type:     "value_error",
		})
	}

	// Do not query the database if the submitted form is already invalid
	if len(validationErrors) > 0 {
		return 0, validationErrors, nil
	}

	// Check the UNIQUE constraint on username
	var existingUsername string

	usernameRow := database.QueryRow(
		"SELECT username FROM users WHERE username = ?",
		username,
	)

	usernameError := usernameRow.Scan(&existingUsername)

	if usernameError != nil && usernameError != sql.ErrNoRows {
		return 0, nil, usernameError
	}

	if usernameError == nil {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "username"},
			Message:  "The username is already taken",
			Type:     "value_error",
		})
	}

	// Check the UNIQUE constraint on email
	var existingEmail string

	emailRow := database.QueryRow(
		"SELECT email FROM users WHERE email = ?",
		email,
	)

	emailError := emailRow.Scan(&existingEmail)

	if emailError != nil && emailError != sql.ErrNoRows {
		return 0, nil, emailError
	}

	if emailError == nil {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "email"},
			Message:  "The email is already registered",
			Type:     "value_error",
		})
	}

	if len(validationErrors) > 0 {
		return 0, validationErrors, nil
	}

	// Hash the password before storing it
	passwordHash, hashError := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if hashError != nil {
		return 0, nil, hashError
	}

	result, databaseError := database.Execute(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		username,
		email,
		string(passwordHash),
	)

	if databaseError != nil {
		return 0, nil, databaseError
	}

	resultID, lastInsertIDError := result.LastInsertId()

	if lastInsertIDError != nil {
		return 0, nil, lastInsertIDError
	}

	return uint64(resultID), nil, nil
}
