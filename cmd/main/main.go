package main

import (
	"log"
	"net/http"

	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// main registers the routes and starts the web server
func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /about", handlers.About)

	staticFiles := http.FileServer(http.Dir("web/static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFiles))

	log.Println("Server listening on http://localhost:8080")

	serverError := http.ListenAndServe(":8080", router)
	if serverError != nil {
		log.Fatal(serverError)
	}
}
