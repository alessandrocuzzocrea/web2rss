package app

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	_ "modernc.org/sqlite"
)

const (
	dataDirPerm = 0755
)

// App represents the main application
type App struct {
	db        *sql.DB
	queries   *db.Queries
	templates *template.Template
	startTime time.Time
}

// New creates a new instance of the application
func New() (*App, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", dataDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	database, err := sql.Open("sqlite", "./data/web2rss.sqlite3")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Optional: set pool settings
	// database.SetMaxOpenConns(10)
	// database.SetMaxIdleConns(5)

	queries := db.New(database)

	// Load templates
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	app := &App{
		db:        database,
		queries:   queries,
		templates: templates,
		startTime: time.Now(),
	}

	return app, nil
}

func (a *App) Close() error {
	return a.db.Close()
}

// Run starts the application
func (a *App) Run() error {
	ctx := context.Background()

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", dataDirPerm); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", "./data/web2rss.sqlite3")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()

	// Test database connection
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	a.StartFeedScheduler()

	mux := a.Routes()
	// host := "localhost"
	port := "8080"

	log.Println("Server starting on :8080")

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}
