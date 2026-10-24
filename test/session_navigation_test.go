package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
	"golang.org/x/crypto/bcrypt"
)

func TestSessionNavigation(t *testing.T) {
	t.Chdir("..")
	db := SetupTestDatabase(t)
	database.SetDatabase(db)
	t.Cleanup(func() {
		db.Close()
		database.SetDatabase(nil)
	})

	hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("INSERT INTO users (username, email, password) VALUES (?, ?, ?)", "navigation-user", "navigation@example.com", string(hash)); err != nil {
		t.Fatal(err)
	}

	form := url.Values{"username": {"navigation-user"}, "password": {"password"}}
	login := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	login.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	handlers.Login(response, login)
	if response.Code != http.StatusSeeOther {
		t.Fatalf("login returned %d: %s", response.Code, response.Body.String())
	}
	var sessionCookie *http.Cookie
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
		}
	}
	if sessionCookie == nil {
		t.Fatal("login did not set a session cookie")
	}

	checkNavigation := func(t *testing.T, cookie *http.Cookie, loggedIn bool) {
		t.Helper()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		if cookie != nil {
			request.AddCookie(cookie)
		}
		response := httptest.NewRecorder()
		handlers.Search(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("page returned %d: %s", response.Code, response.Body.String())
		}
		for _, id := range []string{"nav-logout", "nav-weather", "nav-login", "nav-register"} {
			want := loggedIn
			if id == "nav-login" || id == "nav-register" {
				want = !loggedIn
			}
			if got := strings.Contains(response.Body.String(), `id="`+id+`"`); got != want {
				t.Errorf("%s visible = %v, want %v", id, got, want)
			}
		}
	}

	t.Run("anonymous", func(t *testing.T) { checkNavigation(t, nil, false) })
	t.Run("after login", func(t *testing.T) { checkNavigation(t, sessionCookie, true) })
	t.Run("unknown token", func(t *testing.T) {
		checkNavigation(t, &http.Cookie{Name: "session_token", Value: "unknown"}, false)
	})
	t.Run("expired session", func(t *testing.T) {
		if _, err := db.Exec("UPDATE session_tokens SET created_at = ? WHERE session_value = ?", time.Now().Add(-3*time.Hour), sessionCookie.Value); err != nil {
			t.Fatal(err)
		}
		checkNavigation(t, sessionCookie, false)
		var count int
		if err := db.QueryRow("SELECT COUNT(*) FROM session_tokens WHERE session_value = ?", sessionCookie.Value).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 0 {
			t.Error("expired session was not deleted")
		}
	})
}
