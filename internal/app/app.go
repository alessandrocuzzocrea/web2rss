package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alessandrocuzzocrea/www2rss/internal/db"
	_ "modernc.org/sqlite"
)

// App represents the main application
type App struct {
	db      *sql.DB
	queries *db.Queries
}

// New creates a new instance of the application
func New() (*App, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	database, err := sql.Open("sqlite", "./data/www2rss.sqlite3")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Optional: set pool settings
	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)

	queries := db.New(database)

	return &App{
		db:      database,
		queries: queries,
	}, nil
}

func (a *App) Close() error {
	return a.db.Close()
}

// Run starts the application
func (a *App) Run() error {
	ctx := context.Background()

	// Create data directory if it doesn't exist
	if err := os.MkdirAll("./data", 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite", "./data/www2rss.sqlite3")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connection established")

	// Create queries instance

	// Create author parameters (just for demo)
	// dbstoreParams := dbstore.CreateAuthorParams{
	// 	Name: "Test",
	// 	Bio:  sql.NullString{String: "A sample author bio", Valid: true},
	// 	Loller: sql.NullString{String: "A sample loller", Valid: true},
	// }

	// Create author in database (just for demo)
	// author, err := queries.CreateAuthor(ctx, tutorialParams)
	// if err != nil {
	// 	return fmt.Errorf("failed to create author: %w", err)
	// }

	// log.Printf("Created author: %+v", author)

	// lets grab the first row of feeds table
	// feed, err := a.queries.GetFeed(ctx, 1)
	// if err != nil {
	// 	log.Printf("No feeds found: %v", err)
	// } else {
	// 	log.Printf("First feed: %+v", feed)
	// }

	// log.Printf("Feed URL: %s", feed.Url)

	// // get the url body
	// resp, err := http.Get(feed.Url)
	// if err != nil {
	// 	log.Printf("Failed to get feed URL: %v", err)
	// 	return err
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	log.Printf("Failed to get feed URL: %s", resp.Status)
	// 	return fmt.Errorf("failed to get feed URL: %s", resp.Status)
	// }

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Printf("Failed to read feed body: %v", err)
	// 	return err
	// }

	// log.Printf("Feed body: %s", body)

	// extract this html with this select: <div class="event">

	// Create a new HTML document
	// doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	// if err != nil {
	// 	log.Printf("Failed to parse feed body as HTML: %v", err)
	// 	return err
	// }

	// feedItemSelector := ""
	// if feed.ItemSelector.Valid {
	// 	feedItemSelector = feed.ItemSelector.String
	// }

	// feedItemTitleSelector := ""
	// if feed.TitleSelector.Valid {
	// 	feedItemTitleSelector = feed.TitleSelector.String
	// }

	// feedItemLinkSelector := ""
	// if feed.LinkSelector.Valid {
	// 	feedItemLinkSelector = feed.LinkSelector.String
	// }

	// Find the event elements
	// doc.Find(feedItemSelector).Each(func(i int, s *goquery.Selection) {
	// 	entryTitle := s.Find(feedItemTitleSelector).Text()
	// 	entryLink := s.Find(feedItemLinkSelector).AttrOr("href", "")
	// 	entryDescription := s.Text()

	// 	// log.Printf("Entry %d: %s %s %s", i+1, entryTitle, entryLink, entryDescription)
	// 	a.queries.UpsertFeedItem(ctx, dbstore.UpsertFeedItemParams{
	// 		FeedID:      feed.ID,
	// 		Title:       entryTitle,
	// 		Description: sql.NullString{String: entryDescription, Valid: true},
	// 		Link:        entryLink,
	// 	})
	// })

	// Start the feed scheduler (refreshes feeds every hour)
	a.StartFeedScheduler()

	mux := a.Routes()
	host := "localhost"
	port := "8080"

	log.Println("Server starting on :8080")

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), mux); err != nil {
		return fmt.Errorf("HTTP server failed: %w", err)
	}

	return nil
}
