package main

import (
	"log"
	"net/http"

	"github.com/alessandrocuzzocrea/www2rss/internal/www2rss"
)

func main() {
	log.Printf("Starting www2rss v%s...")

	app, err := www2rss.New()
    if err != nil {
        log.Fatal(err)
    }
    defer app.Close()

    app.Routes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
