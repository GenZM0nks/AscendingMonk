package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// Search handles the rendering of the search template on /search.
//
// @Summary Show search.html
// @Description Render the search template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /search [get]
func Search(responseWriter http.ResponseWriter, request *http.Request) {
	pageData := PageData{
		PageTitle: "Search",
		Flashes:   nil,
	}

	searchData := SearchData{
		Data:          pageData,
		Query:         request.URL.Query().Get("q"),
		SearchResults: nil,
	}

	templates.LoadAndExecuteTemplate("web/templates/search.html", searchData, responseWriter)
}
