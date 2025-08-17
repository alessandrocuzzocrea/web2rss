package app

import (
	"fmt"
	"net/http"
)

// Routes sets up all HTTP routes for the application.
func (a *App) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Add a simple health check endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "www2rss is running!", "all": "lol"}`)
	})

	// Add a health endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy"}`)
	})

	return mux
}
