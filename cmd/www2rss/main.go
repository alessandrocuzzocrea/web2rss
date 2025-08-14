package main

import (
	"log"

	"github.com/alessandrocuzzocrea/www2rss/internal/app"
)

func main() {
	log.Printf("Starting www2rss...")

	app, err := www2rss.New()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err = app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
	
	log.Println("www2rss application completed successfully")
}
