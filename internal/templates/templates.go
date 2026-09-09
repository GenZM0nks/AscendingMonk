package templates

import (
	"fmt"
	"html/template"
	"net/http"
)

// LoadAndExecuteTemplate reads a template file, executes it and writes its HTML to the ResponseWriter
// so the user may see the rendered template in their browser.
func LoadAndExecuteTemplate(pathToTemplate string, templateData any, writer http.ResponseWriter) error {
	template, err := template.ParseFiles("web/templates/layout.html", pathToTemplate)

	if err != nil {
		http.Error(writer, fmt.Sprintf("Unable to load template: %s\nError: %s", pathToTemplate, err.Error()), http.StatusInternalServerError)
		return err
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = template.ExecuteTemplate(writer, "layout", templateData)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
