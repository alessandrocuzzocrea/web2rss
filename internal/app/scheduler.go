package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alessandrocuzzocrea/www2rss/internal/db"
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
	// Fetch the webpage
	resp, err := http.Get(feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

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

		// Skip empty items
		// if title == "" || link == "" {
		// 	return
		// }

		// Make link absolute if it's relative
		if strings.HasPrefix(link, "/") {
			baseURL := resp.Request.URL.Scheme + "://" + resp.Request.URL.Host
			link = baseURL + link
		}

		// Upsert the item (will update if exists, insert if new)
		err := a.queries.UpsertFeedItem(ctx, db.UpsertFeedItemParams{
			FeedID:      feed.ID,
			Title:       title,
			Description: db.NewNullString(description),
			Link:        link,
		})

		if err != nil {
			log.Printf("Failed to upsert feed item: %v", err)
		} else {
			newItemsCount++
		}
	})

	// Clean up old items (keep only last 100 items per feed)
	// cutoffTime := time.Now().AddDate(0, 0, -30) // Keep items from last 30 days
	// err = a.queries.DeleteOldFeedItems(ctx, db.DeleteOldFeedItemsParams{
	// 	FeedID:    feed.ID,
	// 	CreatedAt: db.NewNullTime(cutoffTime),
	// })
	// if err != nil {
	// 	log.Printf("Failed to clean up old items for feed %d: %v", feed.ID, err)
	// }

	log.Printf("Feed %d: processed %d items", feed.ID, newItemsCount)
	return nil
}
