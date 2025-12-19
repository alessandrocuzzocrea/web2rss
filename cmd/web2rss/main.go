package main

import (
	"log"

	"github.com/alessandrocuzzocrea/web2rss/internal/app"
)

func main() {
	log.Printf("Starting web2rss...")

	cfg := app.LoadConfig()

	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Close()

	if err = app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}

	log.Println("web2rss application completed successfully")
}
