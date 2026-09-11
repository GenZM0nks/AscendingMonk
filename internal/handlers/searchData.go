package handlers

// SearchData combines PageData and a query for the search template.
type SearchData struct {
	Data          PageData
	Query         string
	SearchResults []SearchResult
}

// SearchResult is the result of a query using our search engine.
type SearchResult struct {
	Title       string
	Description string
}
