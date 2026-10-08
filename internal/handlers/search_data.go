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
	URL         string
}

// SearchResultDataWrapper is a struct added only to make the OpenAPI spec. conform with the legacy one.
type SearchResultDataWrapper struct {
	Data []SearchResult `json:"data"`
}
