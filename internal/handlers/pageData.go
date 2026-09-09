package handlers

// PageData stores data for templates nested in layout.html.
// Allows us to set titles, flashes (pop up messages) and a user for isLoggedIn conditional rendering.
type PageData struct {
	PageTitle string
	User      string // Edit when we have an actual user
	Flashes   []string
}
