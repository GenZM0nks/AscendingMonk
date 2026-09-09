package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// AboutHandler renders the about template.
//
// @Summary Show about.html
// @Description Render the about template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /about [get]
func AboutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		pageData := PageData{
			PageTitle: "About",
			// User:      *User,
			Flashes: nil,
		}

		templates.LoadAndExecuteTemplate("static/html/about.html", pageData, w)
	}
}
