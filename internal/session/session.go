package session

import (
	"crypto/rand"
	"encoding/base64"

	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
)

func GenerateToken() (string, error) {
	token := make([]byte, 32)
	_, internalError := rand.Read(token)
	if internalError != nil {
		return "", internalError
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

func Authorrize(request *http.Request) error {
	sessionToken, tokenError := request.Cookie("session_token")
	if tokenError != nil {

		fmt.Println(tokenError)
		return tokenError
	}
	if sessionToken.Value == "" {
		fmt.Println("misssing session token value")
		return errors.New("session: session token is missing data")
	}

	fmt.Println("Does have cookie")
	//Get session data
	sessionSqlRow := database.QueryRow("SELECT * FROM session_tokens WHERE session_value = ?", sessionToken.Value)
	var id uint64
	var session_value string
	var csrf_value string
	var created_at time.Time
	missingSessionError := sessionSqlRow.Scan(&session_value, &csrf_value, &created_at, &id)
	if missingSessionError != nil {
		fmt.Println("no session cookie found")
		fmt.Println(missingSessionError)
		return missingSessionError
	}
	fmt.Println(id)
	fmt.Println(session_value)
	fmt.Println(csrf_value)
	fmt.Println(created_at)
	//Check if session is expired
	fmt.Println(time.Now().String() + " is after " + created_at.Add(2*time.Minute).String())
	if time.Now().After(created_at.Add(1 * time.Minute)) {
		_, databaseError := database.Execute("DELETE FROM session_tokens WHERE session_value = ?", session_value)
		if databaseError != nil {
			return databaseError
		} else {
			return errors.New("session: session is expired")
		}
	}

	//Get user data
	userSqlRow := database.QueryRow("SELECT username, password FROM users WHERE id = ?", id)
	var foundUsername string
	var foundPassword string
	missingUserError := userSqlRow.Scan(&foundUsername, &foundPassword)
	if missingUserError != nil {
		return missingUserError
	}
	fmt.Println(id)
	fmt.Println(foundUsername)
	fmt.Println(foundPassword)

	return nil

}
