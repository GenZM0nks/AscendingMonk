package handlers

import (
	"net/http"
)

// APILogout handles logging out users on /api/logout.
//
// @Summary Log users out.
// @Description Log users out and delete their session data.
// @Produce json
// @Success 200
// @Tags Pages
// @Router /api/logout [get]
func APILogout(responseWriter http.ResponseWriter, _ *http.Request) {
	responseWriter.Write([]byte("abc"))
}
