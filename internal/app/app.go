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

// Queries defines the interface for database operations
//
//nolint:dupl
type Queries interface {
	GetFeed(ctx context.Context, id int64) (db.Feed, error)
	ListFeeds(ctx context.Context) ([]db.Feed, error)
	ListFeedsWithItemsCount(ctx context.Context) ([]db.ListFeedsWithItemsCountRow, error)
	CreateFeed(ctx context.Context, arg db.CreateFeedParams) (db.Feed, error)
	UpdateFeed(ctx context.Context, arg db.UpdateFeedParams) error
	UpdateFeedLastRefreshedAt(ctx context.Context, arg db.UpdateFeedLastRefreshedAtParams) error
	DeleteFeed(ctx context.Context, id int64) error
	ListFeedItems(ctx context.Context, feedID int64) ([]db.FeedItem, error)
	UpsertFeedItem(ctx context.Context, arg db.UpsertFeedItemParams) ([]int64, error)
	DeleteItemsByFeedID(ctx context.Context, feedID int64) error
}

// App represents the main application
type App struct {
	db        *sql.DB
	queries   Queries
	templates *template.Template
	startTime time.Time
	config    *Config
}

// New creates a new instance of the application
func New(cfg *Config) (*App, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(cfg.DataDir, dataDirPerm); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	database, err := sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA synchronous = FULL;",
		"PRAGMA temp_store = MEMORY;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA strict = ON;",
	}

	for _, p := range pragmas {
		if _, err := database.Exec(p); err != nil {
			_ = database.Close()
			return nil, fmt.Errorf("failed to set %s: %w", p, err)
		}
	}

	// Test database connection
	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	queries := db.New(database)

	templates := template.New("")

	dirs := []string{
		"templates/*.html",
		"templates/partials/*.html",
	}

	// Load templates
	for _, dir := range dirs {
		_, err := templates.ParseGlob(dir)
		if err != nil {
			_ = database.Close()
			return nil, fmt.Errorf("failed to parse templates in %s: %w", dir, err)
		}
	}

	app := &App{
		db:        database,
		queries:   queries,
		templates: templates,
		startTime: time.Now(),
		config:    cfg,
	}

	return app, nil
}

func (a *App) Close() error {
	return a.db.Close()
}

// Run starts the application
func (a *App) Run() error {
	log.Println("Database connection established")

	a.StartFeedScheduler()

	mux := a.Routes()
	port := a.config.Port

	log.Printf("Server starting on :%s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}
