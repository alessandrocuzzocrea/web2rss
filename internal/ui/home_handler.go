package ui

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/alessandrocuzzocrea/web2rss/internal/config"
	"github.com/alessandrocuzzocrea/web2rss/internal/db"
)

// HomePageData contains data for the homepage template
type HomePageData struct {
	Title      string
	Version    string
	CommitHash string
	GoVersion  string
	BuildTime  string
	Uptime     string
	Feeds      []db.ListFeedsWithItemsCountRow
}

// handleHomepage renders the homepage using HTML templates
func (h *Handler) handleHomepage(w http.ResponseWriter, r *http.Request) {
	// Only serve homepage for exact root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	feeds, err := h.queries.ListFeedsWithItemsCount(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list feeds: %v", err), http.StatusInternalServerError)
		return
	}
	// fmt.Printf("Feeds: %+v\n", feeds) // For debugging; remove in production

	data := HomePageData{
		Title:      "Home",
		Version:    "0.1.0",
		CommitHash: config.CommitHash,
		GoVersion:  runtime.Version(),
		BuildTime:  "2025-09-11", // You can make this dynamic with build flags
		Uptime:     "N/A",        // time.Since(h.startTime) -> h.startTime is not available in Handler yet.
		// We can add StartTime to Handler or calculate it differently.
		// For now, let's just put "N/A" or pass it from Config if we tracked it there.
		Feeds: feeds,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}
}
