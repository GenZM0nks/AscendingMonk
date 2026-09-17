package handlers

import (
	"database/sql"
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
func ServeLoginPage(responseWriter http.ResponseWriter, _ *http.Request) {
	pageData := LoginResponse{
		PageData: PageData{
			PageTitle: "Login",
			Flashes:   nil,
		},
		Errors:   nil,
		Username: "",
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
	requestError := requestPointer.ParseForm()

	if requestError != nil {
		http.Error(
			responseWriter,
			"Failed to parse form",
			http.StatusBadRequest,
		)
		return
	}

	loginRequest := LoginData{
		Username: requestPointer.Form["username"][0],
		Password: requestPointer.Form["password"][0],
	}

	validationErrors, loginError := loginUser(loginRequest)

	if loginError != nil {
		http.Error(
			responseWriter,
			"Failed to login user",
			http.StatusInternalServerError,
		)
		return
	}

	if len(validationErrors) > 0 {
		loginResponse := LoginResponse{
			PageData: PageData{
				PageTitle: "Login",
				Flashes:   nil,
			},
			Errors:   validationErrors,
			Username: loginRequest.Username,
		}

		templates.LoadAndExecuteTemplate(
			"web/templates/login.html",
			loginResponse,
			responseWriter,
		)

		return
	}

	http.Redirect(
		responseWriter,
		requestPointer,
		"/",
		http.StatusSeeOther,
	)

}

func loginUser(loginRequest LoginData) ([]ValidationError, error) {
	validationErrors := []ValidationError{}

	fmt.Println(loginRequest.Username + " " + loginRequest.Password)

	if loginRequest.Username == "" {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "username"},
			Message:  "Username Field required",
			Type:     "missing",
		})
	}

	if loginRequest.Password == "" {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "password"},
			Message:  "Password Field required",
			Type:     "missing",
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	fmt.Println("Passed missing field check")

	sqlRowPointer := database.QueryRow("SELECT username, password FROM users WHERE username = ? AND password = ?", loginRequest.Username, loginRequest.Password) //for testing b69f71a0ab8bec2aa3055dc8745cce81
	foundUser := &LoginData{}
	loginError := sqlRowPointer.Scan(&foundUser.Username, &foundUser.Password)
	if loginError != nil && loginError != sql.ErrNoRows {
		return nil, loginError
	}

	if loginError == sql.ErrNoRows {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "username and password"},
			Message:  "No user with the provided details exists",
			Type:     "value_error",
		})
	}

	if len(validationErrors) > 0 {
		return validationErrors, nil
	}

	return nil, nil
}

// This struct with keep user data
type User struct {
	Id       uint64
	Username string
	Email    string
	Password string
}

type LoginData struct {
	Username string
	Password string
}

type LoginResponse struct {
	PageData
	Errors   []ValidationError
	Username string
}
