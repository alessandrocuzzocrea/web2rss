package app

import (
	"bytes"
	"database/sql"
	"fmt"
	"html"
	"io"
	"net/http"
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
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
	}

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
	}

	itemSelector := r.FormValue("item_selector")
	titleSelector := r.FormValue("title_selector")
	linkSelector := r.FormValue("link_selector")
	dateSelector := r.FormValue("date_selector")

	// Fetch HTML
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch URL", http.StatusInternalServerError)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
	}

	// Parse HTML with GoQuery
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, "Failed to parse HTML", http.StatusInternalServerError)
	}

	// Collect preview items
	type previewItem struct {
		Title string
		Link  string
		Date  string
	}

	var previews []previewItem

	doc.Find(itemSelector).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(titleSelector).Text())
		link, exists := s.Find(linkSelector).Attr("href")
		if !exists {
			link = strings.TrimSpace(s.Find(linkSelector).Text())
		}

		var dateStr string
		if dateSelector != "" {
			dateStr = strings.TrimSpace(s.Find(dateSelector).Text())
			if idx := strings.Index(dateStr, " "); idx != -1 {
				dateStr = dateStr[:idx]
			}
		}

		// Make link absolute if relative
		if strings.HasPrefix(link, "/") {
			link = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + link
		}

		previews = append(previews, previewItem{
			Title: title,
			Link:  link,
			Date:  dateStr,
		})
	})

	// Send OOB HTML updates
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// HTML preview (escaped)
	fmt.Fprintf(w, `<code id="html-preview" hx-swap-oob="true">%s</code>`, html.EscapeString(string(bodyBytes)))

	// Item previews
	if len(previews) > 0 {
		fmt.Fprintf(w, `<p id="selected-item-preview" hx-swap-oob="true"><em>Matched %d items</em></p>`, len(previews))
		fmt.Fprintf(w, `<p id="selected-title-preview" hx-swap-oob="true"><em>First title: %q</em></p>`, previews[0].Title)
		fmt.Fprintf(w, `<p id="selected-link-preview" hx-swap-oob="true"><em>Link found: %s</em></p>`, previews[0].Link)
		fmt.Fprintf(w, `<p id="selected-date-preview" hx-swap-oob="true"><em>Date parsed: %s</em></p>`, previews[0].Date)
	} else {
		fmt.Fprint(w, `
<p id="selected-item-preview" hx-swap-oob="true"><em>No items matched</em></p>
<p id="selected-title-preview" hx-swap-oob="true"><em>No title</em></p>
<p id="selected-link-preview" hx-swap-oob="true"><em>No link</em></p>
<p id="selected-date-preview" hx-swap-oob="true"><em>No date</em></p>
`)
	}

	return
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

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
