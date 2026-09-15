package handlers

import (
	"encoding/json"
	"net/http"
)

// ValidationError describes a single request validation error
type ValidationError struct {
	Location []any  `json:"loc"`
	Message  string `json:"msg"`
	Type     string `json:"type"`
}

// HTTPValidationError contains request validation errors
type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

// AuthResponse describes an authentication API response
type AuthResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// RegisterAPI registers a new user
//
// @Summary Register a new user
// @Description Register a new user account
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Username"
// @Param email formData string true "Email"
// @Param password formData string true "Password"
// @Param password2 formData string false "Password confirmation"
// @Success 200 {object} AuthResponse
// @Failure 422 {object} HTTPValidationError
// @Tags API
// @Router /register [post]
func RegisterAPI(responseWriter http.ResponseWriter, request *http.Request) {
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
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusUnprocessableEntity)

		json.NewEncoder(responseWriter).Encode(HTTPValidationError{
			Detail: validationErrors,
		})

		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")

	json.NewEncoder(responseWriter).Encode(AuthResponse{
		StatusCode: http.StatusOK,
		Message:    "You were successfully registered and can login now",
	})
}
