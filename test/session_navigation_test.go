package test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/html"
)

func TestSessionNavigation(t *testing.T) {
	t.Chdir("..")
	db := SetupTestDatabase(t)
	database.SetDatabase(db)
	t.Cleanup(func() {
		db.Close()
		database.SetDatabase(nil)
	})

	sessionCookie := loginNavigationUser(t, db)

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

func loginNavigationUser(t *testing.T, db *sql.DB) *http.Cookie {
	t.Helper()
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

	return sessionCookie
}

func checkNavigation(t *testing.T, cookie *http.Cookie, loggedIn bool) {
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

	linksGotten, err := parseHTMLForNavLinkIDs(response.Body.String())

	if err != nil {
		t.Errorf("Failed to parse HTML")
	}

	publicNavLinks := [3]string{"nav-register", "nav-login", "nav-weather"}
	privateNavLinks := [2]string{"nav-logout", "nav-weather"}
	var want []string

	if loggedIn {
		want = publicNavLinks[:]
	} else {
		want = privateNavLinks[:]
	}

	if got := linksGotten; slices.Equal(got, want) {
		t.Errorf("%v visible, want %v", got, want)
	}
}

func parseHTMLForNavLinkIDs(HTML string) ([]string, error) {
	doc, err := html.Parse(strings.NewReader(HTML))

	if err != nil {
		return nil, err
	}

	var ids []string
	var walk func(*html.Node)

	walk = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "id" {
					ids = append(ids, attr.Val)
					break
				}
			}
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(doc)
	return ids, nil
}
