package handlers

import (
	"fmt"
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// ServeLoginPage renders the Login page using the shared layout template
//
// @Summary Show login.html
// @Description Render the Login template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /login [get]
func ServeLoginPage(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	pageData := PageData{
		PageTitle: "Login",
		Flashes:   nil,
	}

	templates.LoadAndExecuteTemplate("web/templates/login.html", pageData, responseWriter)

}

// Login Handles login requests
//
// @Summary Show login.html
// @Description Accepts a username and password, and compares them to the sqlite database.
// @Produce json
// @Success 200
// @Failure 422 {string} string "error"
// @Tags queries
// @Router /login [post]
func Login(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	requestPointer.ParseForm()
	formUsername := requestPointer.Form["username"][0]
	formPassword := requestPointer.Form["password"][0]

	fmt.Println(formUsername + " " + formPassword)

}

// This struct with keep user data
type User struct {
	username string
	password string
}
