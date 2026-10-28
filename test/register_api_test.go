package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// TestRegisterAPIMissingFields verifies that missing register fields
// produce validation errors instead of a handler panic.
func TestRegisterAPIMissingFields(t *testing.T) {
	tests := []struct {
		name string
		form url.Values
	}{
		{
			name: "Username missing",
			form: url.Values{
				"email":     {"mail@example.com"},
				"password":  {"1234"},
				"password2": {"1234"},
			},
		},
		{
			name: "Email missing",
			form: url.Values{
				"username":  {"username"},
				"password":  {"1234"},
				"password2": {"1234"},
			},
		},
		{
			name: "Password missing",
			form: url.Values{
				"username":  {"username"},
				"email":     {"mail@example.com"},
				"password2": {"1234"},
			},
		},
		{
			name: "Username, email missing",
			form: url.Values{
				"password":  {"1234"},
				"password2": {"1234"},
			},
		},
		{
			name: "password != password2",
			form: url.Values{
				"username":  {"username"},
				"email":     {"mail@example.com"},
				"password":  {"1234"},
				"password2": {"1235"},
			},
		},
		{
			name: "Username, email, password missing",
			form: url.Values{
				"password2": {"1234"},
			},
		},
		{
			name: "Username, email, password, password2 missing",
			form: url.Values{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(
				http.MethodPost,
				"/api/register",
				strings.NewReader(tt.form.Encode()),
			)
			request.Header.Set(
				"Content-Type",
				"application/x-www-form-urlencoded",
			)

			recorder := httptest.NewRecorder()

			handlers.RegisterAPI(recorder, request)

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
