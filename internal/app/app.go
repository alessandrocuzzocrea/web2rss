package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alessandrocuzzocrea/www2rss/tutorial"
	_ "modernc.org/sqlite"
)

// App represents the main application
type App struct {
	db      *sql.DB
	queries *tutorial.Queries
}

// New creates a new instance of the application
func New() (*App, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite", "./data/www2rss.sqlite3")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Optional: set pool settings
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	queries := tutorial.New(db)

	return &App{
		db:      db,
		queries: queries,
	}, nil
}

func (a *App) Close() error {
	return a.db.Close()
}

// Run starts the application
func (a *App) Run() error {
	ctx := context.Background()

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", "./data/www2rss.sqlite3")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	// Create queries instance

	// Create author parameters (just for demo)
	// tutorialParams := tutorial.CreateAuthorParams{
	// 	Name: "Test",
	// 	Bio:  sql.NullString{String: "A sample author bio", Valid: true},
	// 	Loller: sql.NullString{String: "A sample loller", Valid: true},
	// }

	// Create author in database (just for demo)
	// author, err := queries.CreateAuthor(ctx, tutorialParams)
	// if err != nil {
	// 	return fmt.Errorf("failed to create author: %w", err)
	// }

	// log.Printf("Created author: %+v", author)

	mux := a.Routes()
	port := ":8080" // Default port

	if err := http.ListenAndServe(port, mux); err != nil {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	log.Println("Server starting on :8080")

	return nil
}
