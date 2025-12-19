package ui

import (
	"fmt"
	"net/http"
)

// RegisterRoutes sets up all HTTP routes for the application.
func (h *Handler) RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static files (CSS, JS, images, etc.)
	fs := http.FileServer(http.Dir("static/"))
	mux.Handle("/css/", http.StripPrefix("/", fs))
	mux.Handle("/js/", http.StripPrefix("/", fs))
	mux.Handle("/images/", http.StripPrefix("/", fs))

	// Homepage with HTML template
	mux.HandleFunc("/", h.handleHomepage)

	// Feed creation
	mux.HandleFunc("GET /feed/new", h.handleNewFeed)
	mux.HandleFunc("GET /feed/{id}/duplicate", h.handleDuplicateFeed)
	mux.HandleFunc("POST /feed/preview", h.handlePreviewFeed)

	// Create feed
	mux.HandleFunc("POST /feed/", h.handleCreateFeed)

	// Edit feed
	mux.HandleFunc("GET /feed/{id}/edit", h.handleEditFeed)
	mux.HandleFunc("POST /feed/{id}/edit", h.handleUpdateFeed)

	// Delete feed
	mux.HandleFunc("POST /feed/{id}/delete", h.handleDeleteFeed)

	// Reset feed items
	mux.HandleFunc("POST /feed/{id}/reset", h.handleResetFeedItems)

	// Refresh feed
	mux.HandleFunc("POST /feed/{id}/refresh", h.handleRefreshFeed)

	// Feed endpoints
	// mux.HandleFunc("/feeds/", h.handleListFeeds)  // List all feeds
	mux.HandleFunc("GET /feed/{id}/rss", h.handleFeedRSS) // Get RSS for specific feed

	// Add a health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, `{"status": "healthy"}`)
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	return mux
}
