package handlers

import (
	"fmt"
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
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

	sqlRowPointer := database.QueryRow("SELECT 1 FROM users WHERE username = ?", formUsername)
	foundUser := &User{}
	//err := sqlRowPointer.Scan(foundUser, sqlRowPointer)
	err := sqlRowPointer.Scan(&foundUser.id, &foundUser.username, &foundUser.email, &foundUser.password)
	if err != nil {
		http.Error(responseWriter, fmt.Sprintf("The provided username (%s) or password not recognized", formUsername), http.StatusNotFound)
	}

}

// This struct with keep user data
type User struct {
	id       uint64
	username string
	email    string
	password string
}
