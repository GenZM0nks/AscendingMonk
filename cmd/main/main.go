package main

import (
	"fmt"
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
	err := database.Connect()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer database.Close()

	addr := ":" + os.Getenv("PORT")

	if addr == ":" {
		addr = ":8080"
	}

	files := http.FileServer(http.Dir("./static")) // Request for css won't be filled without this.
	http.Handle("/", files)

	http.Handle("GET /about", handlers.AboutHandler())

	fmt.Printf("Starting server on localhost%s\n", addr)
	http.ListenAndServe(":8080", nil) // nil => DefaultServeMux for routing, request handling.
}
