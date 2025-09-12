package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
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
	item_selector := r.FormValue("item_selector")
	title_selector := r.FormValue("title_selector")
	link_selector := r.FormValue("link_selector")

	if name == "" || url == "" {
		http.Error(w, "Name and URL are required", http.StatusBadRequest)
		return
	}

	// Insert the new feed into the database
	_, err := a.queries.CreateFeed(r.Context(), db.CreateFeedParams{
		Name:          name,
		Url:           url,
		ItemSelector:  sql.NullString{String: item_selector, Valid: item_selector != ""},
		TitleSelector: sql.NullString{String: title_selector, Valid: title_selector != ""},
		LinkSelector:  sql.NullString{String: link_selector, Valid: link_selector != ""},
	})

	if err != nil {
		http.Error(w, "Failed to create feed", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) handleEditFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	feedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid feed ID", http.StatusBadRequest)
		return
	}

	// Get feed details
	feed, err := a.queries.GetFeed(r.Context(), feedID)
	if err != nil {
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	var data = struct {
		ID                  int64
		Name                string
		Url                 string
		ItemSelector        string
		TitleSelector       string
		LinkSelector        string
		DescriptionSelector string
		// other fields as needed
	}{
		ID:                  feed.ID,
		Name:                feed.Name,
		Url:                 feed.Url,
		ItemSelector:        nullStringToString(feed.ItemSelector),
		TitleSelector:       nullStringToString(feed.TitleSelector),
		LinkSelector:        nullStringToString(feed.LinkSelector),
		DescriptionSelector: nullStringToString(feed.DescriptionSelector),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	if err := a.templates.ExecuteTemplate(w, "edit_feed.html", data); err != nil {
		//print the err
		fmt.Printf("Template error: %v\n", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (a *App) handleUpdateFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.PathValue("id")

	feedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid feed ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	url := r.FormValue("url")
	item_selector := r.FormValue("item_selector")
	title_selector := r.FormValue("title_selector")
	link_selector := r.FormValue("link_selector")

	if name == "" || url == "" {
		http.Error(w, "Name and URL are required", http.StatusBadRequest)
		return
	}

	// Update the feed in the database
	err = a.queries.UpdateFeed(r.Context(), db.UpdateFeedParams{
		ID:            feedID,
		Name:          name,
		Url:           url,
		ItemSelector:  sql.NullString{String: item_selector, Valid: item_selector != ""},
		TitleSelector: sql.NullString{String: title_selector, Valid: title_selector != ""},
		LinkSelector:  sql.NullString{String: link_selector, Valid: link_selector != ""},
	})

	if err != nil {
		http.Error(w, "Failed to update feed", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful update
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
