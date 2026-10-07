package handlers

import (
	"encoding/json"
	"errors"
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

	if errTitle != nil {
		return nil, errTitle
	}

	defer rowsTitle.Close()

	dbQueryContent := `SELECT title, url, content FROM pages WHERE language = ? AND content LIKE ? AND title not LIKE ?`
	rowsContent, err := database.Query(dbQueryContent, language, "%"+query+"%", "%"+query+"%")

	if err != nil {
		return nil, err
	}

	defer rowsContent.Close()

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

	if rowsContent.Err() != nil {
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
// @Summary Fetch query data.
// @Description Fetch query data from the database.
// @Produce json
// @Param q query string true "Search query"
// @Param language query string false "Two-letter language string like 'en' for English, the default option."
// @Success 200 {object} SearchResultDataWrapper "A collection of search results relevant to the query"
// @Failure 422 {object} ErrorWithMessageAndStatusCode "Unprocessable Entity"
// @Failure 400 {object} ErrorWithMessageAndStatusCode "Client Error - missing query parameter 'q'"
// @Tags API
// @Router /api/search [get]
func APISearch(responseWriter http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("q")
	responseWriter.Header().Set("Content-Type", "application/json")

	if query == "" {
		writeJSONError(responseWriter, "%s.\n", errors.New("required query parameter 'q' not given"), 400)
		return
	}

	language := request.URL.Query().Get("language")
	searchResults, err := fetchSearchResults(query, language)

	if err != nil {
		writeJSONError(responseWriter, "Error reading search results from database: %s\n", err, 422)
		return
	}

	data := SearchResultDataWrapper{Data: searchResults}
	searchResultsJSON, err := json.Marshal(data)

	if err != nil {
		writeJSONError(responseWriter, "Error converting search results to JSON: %s\n", err, 422)
		return
	}

	responseWriter.Write(searchResultsJSON)
}

func writeJSONError(responseWriter http.ResponseWriter, errorFormatString string, err error, statusCode int) {
	responseWriter.WriteHeader(statusCode)
	_error := ErrorWithMessageAndStatusCode{statusCode, fmt.Sprintf(errorFormatString, err.Error())}
	_errorAsJSON, err := json.Marshal(_error)

	if err != nil {
		responseWriter.WriteHeader(500)
		responseWriter.Write(fmt.Appendf(nil, "Got error, then failed to convert error to JSON: %s\n", err.Error()))
		return
	}

	responseWriter.Write(_errorAsJSON)
}

// ErrorWithMessageAndStatusCode is a base struct for HTTP errors.
type ErrorWithMessageAndStatusCode struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
