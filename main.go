package main

import (
	"log"
	"net/http"

	www2rss "github.com/alessandrocuzzocrea/www2rss/internal/app"
)

func main() {
	log.Printf("Starting www2rss...")

	app, err := www2rss.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	// Set up routes
	app.Routes()

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
