package app

import (
	"database/sql"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
)

func (a *App) handleNewFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	if err := a.templates.ExecuteTemplate(w, "new_feed.html", nil); err != nil {
		// print the err
		fmt.Printf("Template error: %v\n", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (a *App) handlePreviewFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	// Escape the HTML so we can show it safely inside <code>
	escapedBody := html.EscapeString(string(bodyBytes))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, `<code id="html-preview" hx-swap-oob="true">%s</code>`, escapedBody)
	fmt.Fprint(w, `
<p id="selected-item-preview" hx-swap-oob="true"><em>Matched 5 items</em></p>
<p id="selected-title-preview" hx-swap-oob="true"><em>First title: "Hello World"</em></p>
<p id="selected-link-preview" hx-swap-oob="true"><em>Link found: https://example.com/post</em></p>
<p id="selected-date-preview" hx-swap-oob="true"><em>Date parsed: 2025-09-13</em></p>`)
}

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
		// print the err
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

func (a *App) handleRefreshFeed(w http.ResponseWriter, r *http.Request) {
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

	feed, err := a.queries.GetFeed(r.Context(), feedID)
	if err != nil {
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	// Trigger feed refresh (this could be a background job in a real application)
	err = a.refreshFeed(r.Context(), feed)
	if err != nil {
		http.Error(w, "Failed to refresh feed", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful refresh
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// func (a *App) fetchFeedPreview(ctx context.Context, url, itemSelector, titleSelector, linkSelector string) ([]FeedItem, error) {
// 	// This is a placeholder implementation.
// 	// In a real application, you would fetch the URL, parse the HTML,
// 	// and extract items based on the provided selectors.

// 	// For demonstration, return some dummy items.
// 	items := []FeedItem{
// 		{Title: "Sample Item 1", Link: "https://example.com/item1"},
// 		{Title: "Sample Item 2", Link: "https://example.com/item2"},
// 		{Title: "Sample Item 3", Link: "https://example.com/item3"},
// 	}

// 	return items, nil
// }

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// func stringToNullString(s string) sql.NullString {
// 	if s == "" {
// 		return sql.NullString{Valid: false}
// 	}
// 	return sql.NullString{String: s, Valid: true}
// }
