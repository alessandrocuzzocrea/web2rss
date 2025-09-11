package app

import (
	"fmt"
	"net/http"
)

// Routes sets up all HTTP routes for the application.
func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Homepage with HTML template
	mux.HandleFunc("/", a.handleHomepage)

	mux.HandleFunc("POST /feed/", a.handleCreateFeed)
	mux.HandleFunc("GET /feed/{id}/edit", a.handleEditFeed)
	mux.HandleFunc("POST /feed/{id}/edit", a.handleUpdateFeed)

	// Add a health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy"}`)
	})

	// Feed endpoints
	// mux.HandleFunc("/feeds/", a.handleListFeeds)  // List all feeds
	mux.HandleFunc("GET /feed/{id}/rss", a.handleFeedRSS) // Get RSS for specific feed

	return mux
}
