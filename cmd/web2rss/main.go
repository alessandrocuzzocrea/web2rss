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
	defer srv.Close()

	if err := srv.Run(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
