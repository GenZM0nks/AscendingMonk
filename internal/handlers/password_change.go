package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/session"
	"github.com/GenZM0nks/AscendingMonk/internal/templates"
	"golang.org/x/crypto/bcrypt"
)

// PasswordChange renders the page for changing passwords using the shared layout template
// @Summary Show change_password.html
// @Description Render the password changing template to show in the browser.
// @Produce html
// @Success 200
// @Failure 500 {string} string "error"
// @Failure 400 {string} string "error"
// @Tags Pages
// @Router /users/{username}/change-password [get]
func PasswordChange(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	var username = requestPointer.PathValue("username")

	if username == "" {
		http.Error(
			responseWriter,
			"No value for path variable 'username' found",
			http.StatusBadRequest,
		)
	}

	pageData := ChangePasswordResponse{
		PageData: PageData{
			PageTitle:  "Change password",
			Flashes:    nil,
			IsLoggedIn: IsLoggedIn(requestPointer),
		},
		Errors:   nil,
		Username: username,
	}

	templates.LoadAndExecuteTemplate("web/templates/change_password.html", pageData, responseWriter)
}

// PasswordChangePost handles posted new password
// @Summary Call PasswordChangeHandler
// @Description Redirect after successfully updating a password
// @Produce html
// @Success 303
// @Failure 400 {string} string "error"
// @Failure 500 {string} string "error"
// @Tags Pages
// @Router /users/{username}/change-password [post]
func PasswordChangePost(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	var newPassword = requestPointer.FormValue("password")
	var currentPassword = requestPointer.FormValue("current-password")
	var username = requestPointer.PathValue("username")

	if username == "" {
		http.Error(
			responseWriter,
			"No value for path variable 'username' found",
			http.StatusBadRequest,
		)

		return
	}

	var oldPassword string
	var oldPasswordRow = database.QueryRow("SELECT password FROM users WHERE username = ?", username)
	getOldPasswordError := oldPasswordRow.Scan(&oldPassword)

	if getOldPasswordError != nil ||
		bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(currentPassword)) != nil {
		http.Error(responseWriter, "Invalid username or current password", http.StatusUnauthorized)
		return
	}

	var validator = BasePasswordValidator
	validator.rules = append(validator.rules, DifferentFrom([]byte(oldPassword)))
	validationErrors := validator.Validate(newPassword)
	if len(validationErrors) != 0 {
		changePasswordResponse := ChangePasswordResponse{
			PageData: PageData{
				PageTitle:  "Change password",
				Flashes:    nil,
				IsLoggedIn: IsLoggedIn(requestPointer),
			},
			Errors:   validationErrors,
			Username: username,
		}

		templates.LoadAndExecuteTemplate(
			"web/templates/change_password.html",
			changePasswordResponse,
			responseWriter,
		)

		return
	}

	passwordChangeError := UpdatePassword(username, newPassword)
	if passwordChangeError != nil {
		http.Error(
			responseWriter,
			"Failed to change the password.",
			http.StatusInternalServerError,
		)

		return
	}

	var userID uint64
	userIDRow := database.QueryRow("SELECT id FROM users WHERE username = ?", username)
	scanForIDError := userIDRow.Scan(&userID)

	if scanForIDError != nil {
		http.Error(
			responseWriter,
			"Failed to scan a userID.",
			http.StatusInternalServerError,
		)

		return
	}

	sessionRegistrationError := session.RegisterSession(responseWriter, userID)
	if sessionRegistrationError != nil {
		http.Error(
			responseWriter,
			"Failed to create a session.",
			http.StatusInternalServerError,
		)

		return
	}

	http.Redirect(responseWriter, requestPointer, "/", http.StatusSeeOther)
}

// UpdatePassword updates a valid password and updates CreatedAt to prevent
// the requirement to change password again until the next attack.
func UpdatePassword(username, newPassword string) (err error) {
	passwordHash, hashError := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		bcrypt.DefaultCost,
	)

	if hashError != nil {
		return hashError
	}

	result, err := database.Execute(
		"UPDATE users SET password = ?, created_at = ? WHERE username = ?",
		string(passwordHash), time.Now(), username)

	if err != nil {
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// ChangePasswordResponse passes validation data between the POST and GET
// parts of the handlers for changing passwords
type ChangePasswordResponse struct {
	PageData
	Errors   []error
	Username string
}
