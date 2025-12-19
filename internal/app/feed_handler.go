package app

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alessandrocuzzocrea/web2rss/internal/db"
)

func (a *App) handleNewFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	a.renderNewFeed(w, nil)
}

func (a *App) handleDuplicateFeed(w http.ResponseWriter, r *http.Request) {
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

	feed, err := a.queries.GetFeed(r.Context(), feedID)
	if err != nil {
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	data := struct {
		ID            int64
		Name          string
		Url           string
		ItemSelector  string
		TitleSelector string
		LinkSelector  string
		DateSelector  string
	}{
		ID:            feed.ID,
		Name:          feed.Name + " (copy)",
		Url:           feed.Url,
		ItemSelector:  nullStringToString(feed.ItemSelector),
		TitleSelector: nullStringToString(feed.TitleSelector),
		LinkSelector:  nullStringToString(feed.LinkSelector),
		DateSelector:  nullStringToString(feed.DateSelector),
	}

	a.renderNewFeed(w, data)
}

func (a *App) renderNewFeed(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Execute the template
	if err := a.templates.ExecuteTemplate(w, "new_feed.html", data); err != nil {
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

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	feedURL := r.FormValue("url")
	if feedURL == "" {
		http.Error(w, "URL required", http.StatusBadRequest)
		return
	}

	itemSelector := r.FormValue("item_selector")
	titleSelector := r.FormValue("title_selector")
	linkSelector := r.FormValue("link_selector")
	dateSelector := r.FormValue("date_selector")

	existingSelectorIDStr := r.FormValue("existing_selector_id")

	if existingSelectorIDStr != "" {
		existingSelectorID, err := strconv.ParseInt(existingSelectorIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid existing selectors ID", http.StatusBadRequest)
			return
		}

		template_feed, err := a.queries.GetFeed(r.Context(), existingSelectorID)
		if err != nil {
			http.Error(w, "Failed to load existing selectors", http.StatusInternalServerError)
			return
		}

		itemSelector = nullStringToString(template_feed.ItemSelector)
		titleSelector = nullStringToString(template_feed.TitleSelector)
		linkSelector = nullStringToString(template_feed.LinkSelector)
		dateSelector = nullStringToString(template_feed.DateSelector)
	}

	// Fetch URL
	resp, err := http.Get(feedURL)
	if err != nil {
		log.Printf("failed to fetch URL %s: %v", feedURL, err)
		http.Error(w, fmt.Sprintf("failed to fetch URL: %v", err), http.StatusBadRequest)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("failed to fetch URL %s: status %d", feedURL, resp.StatusCode)
		http.Error(w, fmt.Sprintf("failed to fetch URL: status %d", resp.StatusCode), http.StatusBadRequest)
		return
	}
	defer func() {
		err = resp.Body.Close()
		_ = fmt.Errorf("failed to close response body: %w", err)
	}()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read URL body: %v", err)
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("failed to parse HTML: %v", err)
		http.Error(w, "failed to parse HTML", http.StatusInternalServerError)
		return
	}

	// First matched item
	first := doc.Find(itemSelector).First()
	firstTitle := strings.TrimSpace(first.Find(titleSelector).Text())
	firstLink, _ := first.Find(linkSelector).Attr("href")
	if firstLink == "" {
		firstLink = strings.TrimSpace(first.Find(linkSelector).Text())
	}

	firstDate := ""
	if dateSelector != "" {
		firstDate = strings.TrimSpace(first.Find(dateSelector).Text())
		if idx := strings.Index(firstDate, " "); idx != -1 {
			firstDate = firstDate[:idx]
		}
	}

	// Absolute link if needed
	if firstLink != "" {
		base, err := url.Parse(feedURL)
		if err == nil {
			rel, err := url.Parse(firstLink)
			if err == nil {
				firstLink = base.ResolveReference(rel).String()
			}
		}
	}

	firstHTML, _ := first.Html()

	// Render Step 2 HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// prepare existing selectors to be used in the template
	type Selector struct {
		ID            string
		Name          string
		ItemSelector  string
		TitleSelector string
		LinkSelector  string
		DateSelector  string
		Url           string
	}

	type PageData struct {
		ExistingSelectors []Selector
		// any other fields, e.g. for preview
		ItemSelector  string
		TitleSelector string
		LinkSelector  string
		DateSelector  string
		FirstHTML     string
		FirstTitle    string
		FirstLink     string
		FirstDate     string
	}

	ExistingSelectors, err := a.queries.ListFeeds(r.Context())
	if err != nil {
		log.Printf("Failed to list feeds: %v", err)
		http.Error(w, "Failed to load feeds", http.StatusInternalServerError)
		return
	}

	// convert to []Selector
	var selectors []Selector
	for _, s := range ExistingSelectors {
		selectors = append(selectors, Selector{
			ID:            strconv.FormatInt(s.ID, 10),
			Name:          s.Name,
			ItemSelector:  nullStringToString(s.ItemSelector),
			TitleSelector: nullStringToString(s.TitleSelector),
			LinkSelector:  nullStringToString(s.LinkSelector),
			DateSelector:  nullStringToString(s.DateSelector),
		})
	}

	data := PageData{
		ExistingSelectors: selectors,
		ItemSelector:      itemSelector,
		TitleSelector:     titleSelector,
		LinkSelector:      linkSelector,
		DateSelector:      dateSelector,
		FirstHTML:         firstHTML,
		FirstTitle:        firstTitle,
		FirstLink:         firstLink,
		FirstDate:         firstDate,
	}

	// lets use feed-selector-partial.html
	if err := a.templates.ExecuteTemplate(w, "feed-selector-partial.html", data); err != nil {
		// print the err
		fmt.Printf("Template error: %v\n", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}
}

func (a *App) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	url := r.FormValue("url")
	item_selector := r.FormValue("item_selector")
	title_selector := r.FormValue("title_selector")
	link_selector := r.FormValue("link_selector")
	date_selector := r.FormValue("date_selector")

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
		DateSelector:  sql.NullString{String: date_selector, Valid: date_selector != ""},
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
		DateSelector        string
	}{
		ID:                  feed.ID,
		Name:                feed.Name,
		Url:                 feed.Url,
		ItemSelector:        nullStringToString(feed.ItemSelector),
		TitleSelector:       nullStringToString(feed.TitleSelector),
		LinkSelector:        nullStringToString(feed.LinkSelector),
		DescriptionSelector: nullStringToString(feed.DescriptionSelector),
		DateSelector:        nullStringToString(feed.DateSelector),
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
	date_selector := r.FormValue("date_selector")

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
		DateSelector:  sql.NullString{String: date_selector, Valid: date_selector != ""},
	})

	if err != nil {
		http.Error(w, "Failed to update feed", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful update
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) handleDeleteFeed(w http.ResponseWriter, r *http.Request) {
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

	// Delete the feed from the database
	err = a.queries.DeleteFeed(r.Context(), feedID)
	if err != nil {
		http.Error(w, "Failed to delete feed", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful deletion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *App) handleResetFeedItems(w http.ResponseWriter, r *http.Request) {
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

	// Delete all items associated with the feed
	err = a.queries.DeleteItemsByFeedID(r.Context(), feedID)
	if err != nil {
		http.Error(w, "Failed to reset feed items", http.StatusInternalServerError)
		return
	}

	// Redirect back to the homepage after successful reset
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

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
