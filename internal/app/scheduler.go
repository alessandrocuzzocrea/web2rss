package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alessandrocuzzocrea/web2rss/internal/db"
)

// StartFeedScheduler starts a background goroutine that refreshes all feeds every hour
func (a *App) StartFeedScheduler() {
	ticker := time.NewTicker(1 * time.Hour)
	// ticker := time.NewTicker(10 * time.Second)
	go func() {
		defer ticker.Stop()

		// Run immediately on startup
		a.refreshAllFeeds()

		// Then run every hour
		for range ticker.C {
			a.refreshAllFeeds()
		}
	}()

	log.Println("Feed scheduler started - feeds will refresh every hour")
}

// refreshAllFeeds fetches and updates all feeds
func (a *App) refreshAllFeeds() {
	ctx := context.Background()

	feeds, err := a.queries.ListFeeds(ctx)
	if err != nil {
		log.Printf("Failed to fetch feeds for refresh: %v", err)
		return
	}

	log.Printf("Starting refresh of %d feeds", len(feeds))

	for _, feed := range feeds {
		if err := a.refreshFeed(ctx, feed); err != nil {
			log.Printf("Failed to refresh feed %d (%s): %v", feed.ID, feed.Name, err)
		} else {
			log.Printf("Successfully refreshed feed %d (%s)", feed.ID, feed.Name)
		}
	}

	log.Println("Completed feed refresh cycle")
}

// refreshFeed fetches and updates a single feed
func (a *App) refreshFeed(ctx context.Context, feed db.Feed) error {
	if !feed.ItemSelector.Valid {
		return fmt.Errorf("feed %d is missing item selector", feed.ID)
	}

	if !feed.TitleSelector.Valid {
		return fmt.Errorf("feed %d is missing title selector", feed.ID)
	}

	if !feed.LinkSelector.Valid {
		return fmt.Errorf("feed %d is missing link selector", feed.ID)
	}

	// Fetch the webpage
	resp, err := http.Get(feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Extract items using selectors
	var newItemsCount int
	var count []int64
	doc.Find(feed.ItemSelector.String).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(feed.TitleSelector.String).Text())
		link, exists := s.Find(feed.LinkSelector.String).Attr("href")
		if !exists {
			link = strings.TrimSpace(s.Find(feed.LinkSelector.String).Text())
		}

		var description string
		if feed.DescriptionSelector.Valid {
			description, err = s.Find(feed.DescriptionSelector.String).Html()
			if err != nil {
				log.Printf("Failed to get feed item description: %v", err)
			}
		} else {
			description, err = s.Html()
			if err != nil {
				log.Printf("Failed to get feed item description: %v", err)
			}
		}

		var date time.Time
		if feed.DateSelector.Valid && feed.DateSelector.String != "" {
			// Extract the date string from the HTML
			dateStr := strings.TrimSpace(s.Find(feed.DateSelector.String).Text())
			fmt.Println("dateStr:", dateStr)

			// Strip weekday in parentheses if present, e.g., "2025-08-09 (土)" → "2025-08-09"
			if idx := strings.Index(dateStr, " "); idx != -1 {
				dateStr = dateStr[:idx]
			}

			var err error
			// Try common formats
			layouts := []string{
				"2006-01-02", // YYYY-MM-DD
				"2006/01/02", // YYYY/MM/DD
				"02-01-2006", // DD-MM-YYYY
				time.RFC1123, // Mon, 02 Jan 2006 15:04:05 MST
				time.RFC3339, // 2006-01-02T15:04:05Z07:00
			}

			for _, layout := range layouts {
				date, err = time.Parse(layout, dateStr)
				if err == nil {
					break
				}
			}

			if err != nil {
				log.Printf("Failed to parse date '%s': %v", dateStr, err)
			}
		}

		// Skip empty items
		// if title == "" || link == "" {
		// 	return
		// }

		// Make link absolute if it's relative
		if link != "" {
			parsedLink, err := url.Parse(link)
			if err == nil {
				link = resp.Request.URL.ResolveReference(parsedLink).String()
			}
		}

		// Upsert the item (will update if exists, insert if new)
		count, err = a.queries.UpsertFeedItem(ctx, db.UpsertFeedItemParams{
			FeedID:      feed.ID,
			Title:       title,
			Description: db.NewNullString(description),
			Link:        link,
			Date:        db.NewNullTime(date),
		})

		if err != nil {
			log.Printf("Failed to upsert feed item: %v", err)
		} else if len(count) > 0 {
			newItemsCount++
		}
	})

	log.Printf("Feed %d: processed items. Updated %d new items.", feed.ID, newItemsCount)

	// Update the feed's last_refreshed_at timestamp
	if err := a.queries.UpdateFeedLastRefreshedAt(ctx, db.UpdateFeedLastRefreshedAtParams{
		ID:              feed.ID,
		LastRefreshedAt: db.NewNullTime(time.Now().UTC()),
	}); err != nil {
		log.Printf("Failed to update last_refreshed_at for feed %d: %v", feed.ID, err)
	}

	return nil
}
