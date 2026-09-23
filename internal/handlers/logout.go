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
// @Success 200
// @Tags Pages
// @Router /api/logout [get]
func APILogout(responseWriter http.ResponseWriter, request *http.Request) {
	searchResultsJSON, err := json.Marshal("You were logged out")

	if err != nil {
		responseWriter.Write(fmt.Appendln(nil, "Error converting search results to JSON: %w\n", err.Error()))
	}

	token, err := request.Cookie("session_token") // Currently ignores if there's no token
	if err == nil {
		database.Execute("DELETE FROM session_tokens WHERE session_value = ?", token.String())
	}

	responseWriter.Write(searchResultsJSON)
}
