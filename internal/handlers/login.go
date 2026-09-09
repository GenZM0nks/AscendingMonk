package handlers

import (
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/templates"
)

func Login(responseWriter http.ResponseWriter, requestPointer *http.Request) {
	pageData := PageData{
		PageTitle: "Login",
		Flashes:   []string{"Hello"},
	}

	templates.LoadAndExecuteTemplate("web/templates/login.html", pageData, responseWriter)

}
