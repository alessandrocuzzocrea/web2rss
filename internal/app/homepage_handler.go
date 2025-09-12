package app

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
)

// HomePageData contains data for the homepage template
type HomePageData struct {
	Title     string
	Version   string
	GoVersion string
	BuildTime string
	Uptime    string
	Feeds     []db.Feed
}

// handleHomepage renders the homepage using HTML templates
func (a *App) handleHomepage(w http.ResponseWriter, r *http.Request) {
	// Only serve homepage for exact root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	feeds, err := a.queries.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list feeds: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Feeds: %+v\n", feeds) // For debugging; remove in production

	data := HomePageData{
		Title:     "Home",
		Version:   "0.1.0",
		GoVersion: runtime.Version(),
		BuildTime: "2025-09-11", // You can make this dynamic with build flags
		Uptime:    time.Since(a.startTime).Round(time.Second).String(),
		Feeds:     feeds,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	if err := a.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
		return
	}
}
