package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// Login renders the Login page using the shared layout template
//
// @Summary Show login.html
// @Description Render the Login template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /login [get]
func Login(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	pageData := PageData{
		PageTitle: "Login",
		Flashes:   nil,
	}

	templates.LoadAndExecuteTemplate("web/templates/login.html", pageData, responseWriter)

}
