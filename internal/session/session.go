package session

import (
	"crypto/rand"
	"encoding/base64"

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
	sessionSQLRow := database.QueryRow("SELECT * FROM session_tokens WHERE session_value = ?", sessionToken.Value)
	var session sessionData
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

type sessionData struct {
	id           uint64
	sessionValue string
	csrfValue    string
	createdAt    time.Time
}
