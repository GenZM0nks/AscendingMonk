package handlers

import (
	"net/http"

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
func Register(responseWriter http.ResponseWriter, _ *http.Request) {
	registerData := RegisterData{
		PageData: PageData{
			PageTitle: "Register",
			Flashes:   nil,
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
func RegisterPost(responseWriter http.ResponseWriter, request *http.Request) {
	requestError := request.ParseForm()

	if requestError != nil {
		http.Error(
			responseWriter,
			"Failed to parse form",
			http.StatusBadRequest,
		)
		return
	}

	username := request.FormValue("username")
	email := request.FormValue("email")
	password := request.FormValue("password")
	password2 := request.FormValue("password2")

	validationErrors, registerError := registerUser(
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
				PageTitle: "Register",
				Flashes:   nil,
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

	http.Redirect(
		responseWriter,
		request,
		"/login",
		http.StatusSeeOther,
	)
}
