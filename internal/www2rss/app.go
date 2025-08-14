package www2rss

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/alessandrocuzzocrea/www2rss/tutorial"
	_ "github.com/mattn/go-sqlite3"
)

// Version represents the current version of the application
const Version = "0.1.0"

// App represents the main application
type App struct {
	// TODO: Add application fields here
}

// New creates a new instance of the application
func New() *App {
	return &App{}
}

// Run starts the application
func (a *App) Run() error {
	ctx := context.Background()
	
	// Open database connection
	db, err := sql.Open("sqlite3", "./blog.db")
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
		Name: "fuckboi",
		Bio:  sql.NullString{String: "A sample author bio", Valid: true},
	}
	
	// Create author in database
	author, err := queries.CreateAuthor(ctx, tutorialParams)
	if err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}
	
	log.Printf("Created author: %+v", author)
	
	return nil
}
