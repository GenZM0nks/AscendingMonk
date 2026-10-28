package session

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"errors"
	"net/http"
	"time"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
)

// GenerateToken will attempt to create a new random token
// and return the string value of the token
func GenerateToken() (string, error) {
	token := make([]byte, 32)
	_, internalError := rand.Read(token)
	if internalError != nil {
		return "", internalError
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

// Authorize accepts a request, and validates that the request has a session token, that is saved and has not expired
// and can return errors
func Authorize(request *http.Request) error {
	sessionToken, tokenError := request.Cookie("session_token")
	if tokenError != nil {
		return tokenError
	}
	if sessionToken.Value == "" {
		return errors.New("session: session token is missing data")
	}

	// Get session data
	sessionSQLRow := database.QueryRow("SELECT session_value, csrf_value, created_at, user_id FROM session_tokens WHERE session_value = ?", sessionToken.Value)
	var session sessionData
	// The SQLite driver decodes TIMESTAMP columns into time.Time values.
	missingSessionError := sessionSQLRow.Scan(&session.sessionValue, &session.csrfValue, &session.createdAt, &session.id)
	if missingSessionError != nil {
		return missingSessionError
	}

	// Check if session is expired
	if time.Now().After(session.createdAt.Add(2 * time.Hour)) {
		_, databaseError := database.Execute("DELETE FROM session_tokens WHERE session_value = ?", session.sessionValue)
		if databaseError != nil {
			return databaseError
		}
		return errors.New("session: session is expired")
	}

	// Get user data
	userSQLRow := database.QueryRow("SELECT username, password FROM users WHERE id = ?", session.id)
	var foundUsername string
	var foundPassword string
	missingUserError := userSQLRow.Scan(&foundUsername, &foundPassword)
	if missingUserError != nil {
		return missingUserError
	}
	return nil
}

// RegisterSession initializes a session for a user who just logged in.
// It creates a database cell in a sessions table and sets two cookies on the client side,
// session_token and csrf_token. The cookies are valid for two hours.
func RegisterSession(responseWriter http.ResponseWriter, userID uint64) {
	var databaseError error
	var sessionToken string
	var csrfToken string
	for hasValidToken := true; hasValidToken; hasValidToken = databaseError != nil { // hasValidToken = databaseError != nil could become hasValidToken = someCounter < maxAttempts || databaseError != nil And by extracting this loop into a function, we could return an error, which could warn us that the database session tokens is so filled, that randomly generating multiple strings created duplicates

		var tokenError error
		sessionToken, tokenError = GenerateToken()

		if tokenError != nil {
			http.Error(
				responseWriter,
				"Failed to generate token",
				http.StatusInternalServerError,
			)
			return
		}

		tokenError = nil
		csrfToken, tokenError = GenerateToken()

		if tokenError != nil {
			http.Error(
				responseWriter,
				"Failed to generate token",
				http.StatusInternalServerError,
			)
			return
		}

		// Tokens are randomly generated, the liklyhood of two identical is slim to none. Nonetheless. If two identical tokens happen to be generated, the database with send back an error (for not upholding the UNIQUE CONSTRAINT), This will cause the program to enter the statement below and crash the system, instead of the intented do-while loop for generating a new token.
		_, databaseError := database.Execute(
			"INSERT INTO session_tokens (session_value, csrf_value, created_at, user_id) VALUES (?, ?, ?, ?)",
			sessionToken,
			csrfToken,
			time.Now(),
			userID,
		)

		if databaseError != nil {
			fmt.Print(databaseError)
			http.Error(
				responseWriter,
				"Failed to login",
				http.StatusInternalServerError,
			)
			return
		}

	}

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: true,
		Path:     "/", // Should ensure that the session token is sent for all endpoints
	})

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: false,
	})
}

type sessionData struct {
	id           uint64
	sessionValue string
	csrfValue    string
	createdAt    time.Time
}
