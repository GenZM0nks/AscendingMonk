package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestLoginAPIMissingFields verifies that missing login fields
// produce validation errors instead of a handler panic.
func TestLoginAPIMissingFields(t *testing.T) {
	tests := []struct {
		name string
		form url.Values
	}{
		{
			name: "both fields missing",
			form: url.Values{},
		},
		{
			name: "username missing",
			form: url.Values{
				"password": {"example-password"},
			},
		},
		{
			name: "password missing",
			form: url.Values{
				"username": {"example-user"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/login",
				strings.NewReader(tt.form.Encode()),
			)
			request.Header.Set(
				"Content-Type",
				"application/x-www-form-urlencoded",
			)

			recorder := httptest.NewRecorder()

			handlers.LoginAPI(recorder, request)

			if recorder.Code != http.StatusUnprocessableEntity {
				t.Fatalf(
					"expected status %d, got %d\nBody: %s",
					http.StatusUnprocessableEntity,
					recorder.Code,
					recorder.Body.String(),
				)
			}

			if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
				t.Errorf("expected JSON content type, got %q", contentType)
			}
		})
	}
}
