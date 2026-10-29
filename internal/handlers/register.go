package handlers

import (
	"fmt"
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/session"
	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

// RegisterData contains the data used to render the registration page
type RegisterData struct {
	PageData
	Errors   []ValidationError
	Username string
	Email    string
}

// Register renders the Register page using the shared layout template
//
// @Summary Show register.html
// @Description Render the Register template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /register [get]
func Register(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	registerData := RegisterData{
		PageData: PageData{
			PageTitle:  "Register",
			Flashes:    nil,
			IsLoggedIn: IsLoggedIn(requestPointer),
		},
		Errors:   nil,
		Username: "",
		Email:    "",
	}

	templates.LoadAndExecuteTemplate(
		"web/templates/register.html",
		registerData,
		responseWriter,
	)
}

// RegisterPost handles registration submitted from the browser.
func RegisterPost(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	requestError := requestPointer.ParseForm()

	if requestError != nil {
		http.Error(
			responseWriter,
			"Failed to parse form",
			http.StatusBadRequest,
		)
		return
	}

	username := requestPointer.FormValue("username")
	email := requestPointer.FormValue("email")
	password := requestPointer.FormValue("password")
	password2 := requestPointer.FormValue("password2")

	userID, validationErrors, registerError := registerUser(
		username,
		email,
		password,
		password2,
	)

	if registerError != nil {
		http.Error(
			responseWriter,
			"Failed to register user",
			http.StatusInternalServerError,
		)
		return
	}

	if len(validationErrors) > 0 {
		registerData := RegisterData{
			PageData: PageData{
				PageTitle:  "Register",
				Flashes:    nil,
				IsLoggedIn: IsLoggedIn(requestPointer),
			},
			Errors:   validationErrors,
			Username: username,
			Email:    email,
		}

		templates.LoadAndExecuteTemplate(
			"web/templates/register.html",
			registerData,
			responseWriter,
		)

		return
	}

	tokenGenerationError := session.RegisterSession(responseWriter, userID)
	if tokenGenerationError != nil {
		http.Error(
			responseWriter,
			fmt.Sprintf("Failed to generate session- and/or CSRF tokens for a session for user \"%s\"", username),
			http.StatusInternalServerError,
		)
		return
	}

	http.Redirect(
		responseWriter,
		requestPointer,
		"/search",
		http.StatusSeeOther,
	)
}
