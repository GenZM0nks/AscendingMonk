package session

import (
	"crypto/rand"
	"encoding/base64"

	"net/http"
)

func GenerateToken() (string, error) {
	token := make([]byte, 32)
	_, internalError := rand.Read(token)
	if internalError != nil {
		return "", internalError
	}
	return base64.URLEncoding.EncodeToString(token), nil
}

var AuthError error

func Authorrize(request *http.Request) error {
	sessionToken, tokenError := request.("session_token");
	if tokenError != nill || sessionToken.Value == ""{
		return AuthError
	}
}
