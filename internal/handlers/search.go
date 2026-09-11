package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// Queries the database for search results on /search and /api/search
func fetchSearchResults(query string, language string) (results []SearchResult, err error) {
	results = make([]SearchResult, 0)

	if language == "" {
		language = "en"
	}

	if query == "" {
		return results, nil
	}

	dbQuery := `SELECT title, url, content FROM pages WHERE language = ? AND content LIKE ?`
	rows, err := database.Query(dbQuery, language, "%"+query+"%")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var searchResult SearchResult

		err = rows.Scan(
			&searchResult.Title,
			&searchResult.URL,
			&searchResult.Description,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, searchResult)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return results, nil
}

// Search handles the rendering of the search template on /search.
//
// @Summary Show search.html
// @Description Render the search template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /search [post]
func Search(responseWriter http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("q")
	language := request.URL.Query().Get("language")
	searchResults, err := fetchSearchResults(query, language)

	if err != nil {
		responseWriter.Write(fmt.Appendln(nil, "Error reading search results from database: %w\n", err.Error()))
	}

	pageData := PageData{
		PageTitle: "Search",
		Flashes:   nil,
	}

	searchData := SearchData{
		Data:          pageData,
		Query:         query,
		SearchResults: searchResults,
	}

	templates.LoadAndExecuteTemplate("web/templates/search.html", searchData, responseWriter)
}

// APISearch handles the rendering of the search template on /search.
//
// @Summary Fetch query data.
// @Description Fetch query data from the database.
// @Produce json
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /api/search [get]
func APISearch(responseWriter http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("q")
	language := request.URL.Query().Get("language")
	searchResults, err := fetchSearchResults(query, language)

	if err != nil {
		responseWriter.Write(fmt.Appendln(nil, "Error reading search results from database: %w\n", err.Error()))
	}

	searchResultsJSON, err := json.Marshal(searchResults)

	if err != nil {
		responseWriter.Write(fmt.Appendln(nil, "Error converting search results to JSON: %w\n", err.Error()))
	}

	responseWriter.Write(searchResultsJSON)
}
