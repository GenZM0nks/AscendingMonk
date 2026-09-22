package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/session"
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

	loginRequest := LoginRequest{
		Username: requestPointer.Form["username"][0],
		Password: requestPointer.Form["password"][0],
	}

	user, validationErrors, loginError := loginUser(loginRequest)

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
	_, _ = database.Execute("DELETE FROM session_tokens") //empty all the test session data
	var databaseError error
	var sessionToken string
	var csrfToken string
	for hasValidToken := true; hasValidToken; hasValidToken = databaseError != nil { //hasValidToken = databaseError != nil could become hasValidToken = someCounter < maxAttempts || databaseError != nil And by extracting this loop into a function, we could return an error, which could warn us that the database session tokens is so filled, that randomly generating multiple strings created duplicates

		var tokenError error
		sessionToken, tokenError = session.GenerateToken()

		if tokenError != nil {
			http.Error(
				responseWriter,
				"Failed to generate token",
				http.StatusInternalServerError,
			)
			return
		}

		tokenError = nil
		csrfToken, tokenError = session.GenerateToken()

		if tokenError != nil {
			http.Error(
				responseWriter,
				"Failed to generate token",
				http.StatusInternalServerError,
			)
			return
		}

		//Tokens are randomly generated, the liklyhood of two identical is slim to none. Nonetheless. If two identical tokens happen to be generated, the database with send back an error (for not upholding the UNIQUE CONSTRAINT), This will cause the program to enter the statement below and crash the system, instead of the intented do-while loop for generating a new token.
		_, databaseError := database.Execute(
			"INSERT INTO session_tokens (session_value, csrf_value, created_at, user_id) VALUES (?, ?, ?, ?)",
			sessionToken,
			csrfToken,
			time.Now(),
			user.id,
		)

		if databaseError != nil {
			fmt.Print(databaseError)
			http.Error(
				responseWriter,
				"Failed to login",
				http.StatusInternalServerError,
			)
			return
		}

	}
	//test
	sqlRowPointer := database.QueryRow("SELECT * FROM session_tokens")
	var session_value string
	var csrf_value string
	var created_at time.Time
	var user_id int64
	sqlRowPointer.Scan(&session_value, &csrf_value, &created_at, &user_id)
	fmt.Println(user_id)
	fmt.Println(session_value)
	fmt.Println(csrf_value)
	fmt.Println(created_at.String())
	//test
	fmt.Print("token generated: ")
	fmt.Println(sessionToken)

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(2 * time.Minute),
		HttpOnly: true,
		Path:     "/", //Should ensure that the session token is sent for all endpoints
	})

	http.SetCookie(responseWriter, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(2 * time.Minute),
		HttpOnly: false,
	})

	http.Redirect(
		responseWriter,
		requestPointer,
		"/",
		http.StatusSeeOther,
	)

}

func loginUser(loginRequest LoginRequest) (*LoginResult, []ValidationError, error) {
	validationErrors := []ValidationError{}

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
		return nil, validationErrors, nil
	}

	sqlRowPointer := database.QueryRow("SELECT id, username, password FROM users WHERE username = ?", loginRequest.Username)
	foundUser := &LoginResult{}
	loginError := sqlRowPointer.Scan(&foundUser.id, &foundUser.Username, &foundUser.Password)
	if loginError != nil && loginError != sql.ErrNoRows {
		return nil, nil, loginError
	}
	if loginError == sql.ErrNoRows {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "username and password"},
			Message:  "No user with the provided details exists",
			Type:     "value_error",
		})
	}

	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	hashError := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(loginRequest.Password))

	if (hashError != nil) && hashError != bcrypt.ErrMismatchedHashAndPassword {
		return nil, nil, hashError
	}

	if hashError == bcrypt.ErrMismatchedHashAndPassword {
		validationErrors = append(validationErrors, ValidationError{
			Location: []any{"body", "username and password"},
			Message:  "No user with the provided details exists",
			Type:     "value_error",
		})
	}

	if len(validationErrors) > 0 {
		return nil, validationErrors, nil
	}

	return foundUser, nil, nil
}

// The method below is only for testing that sessions are set up correctly
func Protected(w http.ResponseWriter, r *http.Request) {
	if err := session.Authorrize(r); err != nil {
		http.Error(
			w,
			"You do not have acces",
			http.StatusForbidden,
		)
		return
	}
	fmt.Println("is authorized")
	http.Error(
		w,
		"secrets, secrets, secrets",
		http.StatusOK,
	)
}

// LoginRequest stores request data
// Allows us pass login data together as one entity
type LoginRequest struct {
	Username string
	Password string
}

// LoginResult stores data from queries.
// Allows us to extract all data when querying
type LoginResult struct {
	id       uint64
	Username string
	Password string
}

// LoginResponse stores data that need to be sent back to the clients browser.
// Allows us send error messages to be displayed and for the username to be automatically set in the form, so the user does not need to repeat the form.
type LoginResponse struct {
	PageData
	Errors   []ValidationError
	Username string
}
