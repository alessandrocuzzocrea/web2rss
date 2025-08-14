package www2rss

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/alessandrocuzzocrea/www2rss/tutorial"
	_ "github.com/mattn/go-sqlite3"
)

// App represents the main application
type App struct {
    db      *sql.DB
    queries *tutorial.Queries
}

// New creates a new instance of the application
func New() (*App, error) {
    db, err := sql.Open("sqlite3", "./data/www2rss.sqlite3")
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
	
	// Open database connection
	db, err := sql.Open("sqlite3", "./data/www2rss.sqlite3")
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
	queries := tutorial.New(db)
	
	// Create author parameters
	tutorialParams := tutorial.CreateAuthorParams{
		Name: "Test",
		Bio:  sql.NullString{String: "A sample author bio", Valid: true},
		Loller: sql.NullString{String: "A sample loller", Valid: true},
	}
	
	// Create author in database
	author, err := queries.CreateAuthor(ctx, tutorialParams)
	if err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}
	
	log.Printf("Created author: %+v", author)
	
	return nil
}
