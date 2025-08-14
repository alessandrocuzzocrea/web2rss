package main

import (
	"log"

	"github.com/alessandrocuzzocrea/www2rss/internal/www2rss"
)

func main() {
	log.Printf("Starting www2rss v%s...", www2rss.Version)
	
	app := www2rss.New()
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
	
	log.Println("www2rss application completed successfully")
}
