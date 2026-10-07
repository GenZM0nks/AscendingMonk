package handlers

import (
	"encoding/json"
	"fmt"
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
	logoutJSON, err := json.Marshal(HTTPResponse{"You were logged out", 200})

	if err != nil {
		responseWriter.WriteHeader(500)
		errorAsJSON, err := json.Marshal(HTTPResponse{fmt.Sprintf("Error converting logout message to JSON: %s", err), 500})

		if err != nil {
			json.Marshal("Failed to convert error message into an HTTPResponse.")
		}

		responseWriter.Write(errorAsJSON)
		return
	}

	token, err := request.Cookie("session_token") // Currently ignores if there's no token
	if err == nil {
		database.Execute("DELETE FROM session_tokens WHERE session_value = ?", token.String())
	}

	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	responseWriter.Write(logoutJSON)
}

// HTTPResponse models an HTTP response with a status code and a message
type HTTPResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}
