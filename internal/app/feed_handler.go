package app

import (
	"net/http"

	"github.com/alessandrocuzzocrea/www2rss/internal/db"
)

func (a *App) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	url := r.FormValue("url")

	if name == "" || url == "" {
		http.Error(w, "Name and URL are required", http.StatusBadRequest)
		return
	}

	// Insert the new feed into the database
	_, err := a.queries.CreateFeed(r.Context(), db.CreateFeedParams{
		Name: name,
		Url:  url,
	})
	if err != nil {
		http.Error(w, "Failed to create feed", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
