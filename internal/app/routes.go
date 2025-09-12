package app

import (
	"fmt"
	"net/http"
)

// Routes sets up all HTTP routes for the application.
func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files (CSS, JS, images, etc.)
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/css/", http.StripPrefix("/", fs))
	mux.Handle("/js/", http.StripPrefix("/", fs))
	mux.Handle("/images/", http.StripPrefix("/", fs))

	// Homepage with HTML template
	mux.HandleFunc("/", a.handleHomepage)

	mux.HandleFunc("GET /feed/new", a.handleNewFeed)
	mux.HandleFunc("POST /feed/", a.handleCreateFeed)
	mux.HandleFunc("GET /feed/{id}/edit", a.handleEditFeed)
	mux.HandleFunc("POST /feed/{id}/edit", a.handleUpdateFeed)
	mux.HandleFunc("POST /feed/{id}/refresh", a.handleRefreshFeed)

	// Add a health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{"status": "healthy"}`)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	// Feed endpoints
	// mux.HandleFunc("/feeds/", a.handleListFeeds)  // List all feeds
	mux.HandleFunc("GET /feed/{id}/rss", a.handleFeedRSS) // Get RSS for specific feed

	return mux
}
