package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// Register renders the Register page using the shared layout template
//
// @Summary Show register.html
// @Description Render the Register template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /register [get]
func Register(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	pageData := PageData{
		PageTitle: "Register",
		Flashes:   nil,
	}

	templates.LoadAndExecuteTemplate("web/templates/register.html", pageData, responseWriter)

}