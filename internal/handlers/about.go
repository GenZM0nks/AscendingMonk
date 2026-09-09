package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// About renders the about page using the shared layout template
func About(responseWriter http.ResponseWriter, request *http.Request) {
	pageTemplates, templateError := template.ParseFiles(
		"web/templates/layout.html",
		"web/templates/about.html",
	)
	if templateError != nil {
		log.Println(templateError)
		http.Error(
			responseWriter,
			"Unable to load the about page",
			http.StatusInternalServerError,
		)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")

	renderError := pageTemplates.ExecuteTemplate(responseWriter, "layout", nil)
	if renderError != nil {
		log.Println(renderError)
	}
}
