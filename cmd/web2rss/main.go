package main

import (
	"fmt"
	"log"
	_ "time/tzdata"

	"github.com/alessandrocuzzocrea/web2rss/internal/app"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("web2rss failed: %v", err)
	}
}

func run() error {
	log.Printf("Starting web2rss...")

	cfg := app.LoadConfig()

	application, err := app.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize application: %w", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("Failed to close application: %v", err)
		}
	}()

	if err = application.Run(); err != nil {
		return fmt.Errorf("application failed: %w", err)
	}

	log.Println("web2rss application completed successfully")
	return nil
}
