package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/alessandrocuzzocrea/web2rss/internal/config"
	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/alessandrocuzzocrea/web2rss/internal/feed"
	"github.com/alessandrocuzzocrea/web2rss/internal/ui"
	_ "modernc.org/sqlite"
)

const (
	dataDirPerm = 0755
)

// Server represents the main application server
type Server struct {
	db          *sql.DB
	handler     *ui.Handler
	feedService *feed.Service
	config      *config.Config
}

// New creates a new instance of the application server
func New(cfg *config.Config) (*Server, error) {
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

	// Initialize Feed Service
	feedService := feed.NewService(queries)

	// Initialize Templates
	templates := template.New("").Funcs(ui.NewTemplateFuncs(cfg))

	dirs := []string{
		cfg.TemplateDir + "/*.html",
		cfg.TemplateDir + "/partials/*.html",
	}

	// Load templates
	for _, dir := range dirs {
		_, err := templates.ParseGlob(dir)
		if err != nil {
			_ = database.Close()
			return nil, fmt.Errorf("failed to parse templates in %s: %w", dir, err)
		}
	}

	// Initialize UI Handler
	handler := ui.NewHandler(queries, templates, feedService, cfg)

	server := &Server{
		db:          database,
		handler:     handler,
		feedService: feedService,
		config:      cfg,
	}

	return server, nil
}

func (s *Server) Close() error {
	return s.db.Close()
}

// Run starts the application
func (s *Server) Run() error {
	log.Println("Database connection established")

	// Start Feed Scheduler
	s.feedService.StartScheduler()

	mux := s.handler.RegisterRoutes()
	port := s.config.Port

	log.Printf("Server starting on :%s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}
