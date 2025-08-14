package main

import (
	"log"

	www2rss "github.com/alessandrocuzzocrea/www2rss/internal/app"
)

func main() {
	log.Printf("Starting www2rss...")

	app, err := www2rss.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
}
