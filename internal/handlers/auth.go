package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/session"
)

// IsLoggedIn returns true if a user is logged in, false if not.
func IsLoggedIn(requestPointer *http.Request) bool {
	err := session.Authorize(requestPointer)
	return err == nil
}
