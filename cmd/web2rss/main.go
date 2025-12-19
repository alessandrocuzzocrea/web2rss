package main

import (
	"log"

	"github.com/alessandrocuzzocrea/web2rss/internal/config"
	"github.com/alessandrocuzzocrea/web2rss/internal/server"
)

func main() {
	cfg := config.LoadConfig()

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
	defer func() {
		if err := srv.Close(); err != nil {
			log.Printf("Error closing server: %v", err)
		}
	}()

	if err := srv.Run(); err != nil {
		log.Printf("Server error: %v", err)
	}
}
