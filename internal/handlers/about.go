package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// About renders the about page using the shared layout template
//
// @Summary Show about.html
// @Description Render the about template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /about [get]
func About(responseWriter http.ResponseWriter, _ *http.Request) {
	pageData := PageData{
		PageTitle: "About",
		Flashes:   nil,
	}

	templates.LoadAndExecuteTemplate("web/templates/about.html", pageData, responseWriter)
}
