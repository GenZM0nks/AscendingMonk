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


	dbQueryTitle := `SELECT title, url, content FROM pages WHERE language = ? AND title LIKE ?`
	rowsTitle, errTitle := database.Query(dbQueryTitle, language, "%"+query+"%")

	dbQueryContent := `SELECT title, url, content FROM pages WHERE language = ? AND content LIKE ? AND title not LIKE ?`
	rowsContent, err := database.Query(dbQueryContent, language, "%"+query+"%", "%"+query+"%")

	if err != nil || errTitle != nil{
		return nil, err
	}

	for rowsTitle.Next() {
		var searchResult SearchResult

		err = rowsTitle.Scan(
			&searchResult.Title,
			&searchResult.URL,
			&searchResult.Description,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, searchResult)
	}

	for rowsContent.Next() {
		var searchResult SearchResult

		err = rowsContent.Scan(
			&searchResult.Title,
			&searchResult.URL,
			&searchResult.Description,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, searchResult)
	}

	if rowsContent.Err() != nil{
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
// @Router / [get]
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
