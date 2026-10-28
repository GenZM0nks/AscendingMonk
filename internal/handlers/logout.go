package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
)

// APILogout handles logging out users on /api/logout.
//
// @Summary Log users out.
// @Description Log users out and delete their session data.
// @Produce json
// @Success 200 {object} HTTPResponse "Status code and message"
// @Failure 500 {object} HTTPResponse "Error on marshalling 'You were logged out' into JSON"
// @Tags API
// @Router /api/logout [get]
func APILogout(responseWriter http.ResponseWriter, request *http.Request) {
	token, err := request.Cookie("session_token") // Currently ignores if there's no token
	if err == nil {
		database.Execute("DELETE FROM session_tokens WHERE session_value = ?", token.Value)
		http.SetCookie(responseWriter, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   true,
			Path:     "/", // Should ensure that the session token is sent for all endpoints
		})
	}

	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}

// HTTPResponse models an HTTP response with a status code and a message
type HTTPResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
