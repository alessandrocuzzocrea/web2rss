package app

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// RSS XML structures
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Language    string `xml:"language,omitempty"`
	PubDate     string `xml:"pubDate,omitempty"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
}

// GET /feed/{id}/ - Generate RSS XML for a feed
func (a *App) handleFeedRSS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")

	feedID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Printf("Invalid feed ID: %v\n", err)
		http.Error(w, "Invalid feed ID", http.StatusBadRequest)
		return
	}

	// Get feed details
	feed, err := a.queries.GetFeed(ctx, feedID)
	if err != nil {
		http.Error(w, "Feed not found", http.StatusNotFound)
		return
	}

	// Get feed items
	items, err := a.queries.ListFeedItems(ctx, feedID)
	if err != nil {
		fmt.Printf("Failed to fetch feed items: %v\n", err)
		http.Error(w, "Failed to fetch feed items", http.StatusInternalServerError)
		return
	}

	// Convert to RSS items
	rssItems := make([]Item, len(items))
	for i, item := range items {
		var pubDate time.Time
		switch {
		case item.Date.Valid:
			pubDate = item.Date.Time
		case item.CreatedAt.Valid:
			pubDate = item.CreatedAt.Time
		default:
			pubDate = time.Now()
		}

		rssItems[i] = Item{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description.String,
			PubDate:     formatRSSDate(pubDate),
		}
	}

	// Create RSS feed
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       feed.Name,
			Link:        feed.Url,
			Description: fmt.Sprintf("RSS feed generated from %s", feed.Url),
			Language:    "en-us",
			PubDate:     formatRSSDate(time.Now()),
			Items:       rssItems,
		},
	}

	// Set headers and encode XML
	w.Header().Set("Content-Type", "application/rss+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Write XML declaration
	_, err = w.Write([]byte(xml.Header))
	if err != nil {
		fmt.Printf("Failed to write XML header: %v\n", err)
		http.Error(w, "Failed to write XML header", http.StatusInternalServerError)
		return
	}

	// Encode RSS to XML
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")
	if err := encoder.Encode(rss); err != nil {
		http.Error(w, "Failed to generate RSS", http.StatusInternalServerError)
		return
	}
}

// formatRSSDate formats time to RFC822 format (RSS standard)
func formatRSSDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC822)
}
