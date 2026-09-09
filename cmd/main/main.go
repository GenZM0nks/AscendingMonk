package main

import (
	"log"
	"net/http"
	"os"

	"github.com/GenZM0nks/AscendingMonk/internal/database"
	"github.com/GenZM0nks/AscendingMonk/internal/handlers"
)

// @title           AscendingMonkAPI
// @version         1.0
// @description     This is a Go rewrite of whoknows.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://github.com/GenZM0nks/AscendingMonk
// @contact.email  apisupport@genzmonks.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
func main() {
	router := http.NewServeMux()

	err := database.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer database.Close()

	addr := ":" + os.Getenv("PORT")

	if addr == ":" {
		addr = ":8080"
	}

	router.HandleFunc("GET /about", handlers.About)

	staticFiles := http.FileServer(http.Dir("web/static"))
	router.Handle("GET /static/", http.StripPrefix("/static/", staticFiles))

	log.Printf("Server listening on http://localhost:%s\n", addr)

	serverError := http.ListenAndServe(addr, router)
	if serverError != nil {
		log.Fatal(serverError)
	}
}
