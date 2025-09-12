package main

import (
	"log"

	web2rss "github.com/alessandrocuzzocrea/web2rss/internal/app"
)

func main() {
	log.Printf("Starting web2rss...")

	app, err := web2rss.New()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()
}
