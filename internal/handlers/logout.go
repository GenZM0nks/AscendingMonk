package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/session"
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
	session.DropSession(request, responseWriter)
	http.Redirect(responseWriter, request, "/", http.StatusSeeOther)
}

// HTTPResponse models an HTTP response with a status code and a message
type HTTPResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
