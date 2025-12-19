package app

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestFormatRSSDate(t *testing.T) {
	// Test with a valid time
	testTime := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	expected := "Mon, 25 Dec 2023 10:30:00 +0000"
	result := formatRSSDate(testTime)
	assert.Equal(t, expected, result)

	// Test with zero time
	zeroTime := time.Time{}
	result = formatRSSDate(zeroTime)
	assert.Equal(t, "", result)
}

func TestRSSGeneration(t *testing.T) {
	// Test the RSS structure generation directly
	items := []Item{
		{
			Title:       "Test Item 1",
			Link:        "https://example.com/item1",
			Description: "This is a test item",
			PubDate:     "Mon, 25 Dec 2023 10:30:00 +0000",
		},
	}

	// Create RSS feed
	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "Test Feed",
			Link:        "https://example.com",
			Description: "RSS feed generated from https://example.com",
			Language:    "en-us",
			PubDate:     "Mon, 25 Dec 2023 10:30:00 +0000",
			Items:       items,
		},
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(rss, "", "  ")
	assert.NoError(t, err)

	// Check that the XML contains expected elements
	xmlStr := xml.Header + string(output)
	assert.Contains(t, xmlStr, "<rss version=\"2.0\">")
	assert.Contains(t, xmlStr, "<title>Test Feed</title>")
	assert.Contains(t, xmlStr, "<link>https://example.com</link>")
	assert.Contains(t, xmlStr, "<description>RSS feed generated from https://example.com</description>")
	assert.Contains(t, xmlStr, "<item>")
	assert.Contains(t, xmlStr, "<title>Test Item 1</title>")
	assert.Contains(t, xmlStr, "<link>https://example.com/item1</link>")
	assert.Contains(t, xmlStr, "<description>This is a test item</description>")
	assert.Contains(t, xmlStr, "<pubDate>Mon, 25 Dec 2023 10:30:00 +0000</pubDate>")

	// Verify it's valid XML by unmarshaling it back
	var parsedRSS RSS
	err = xml.Unmarshal(output, &parsedRSS)
	assert.NoError(t, err)

	// Verify the parsed data
	assert.Equal(t, "2.0", parsedRSS.Version)
	assert.Equal(t, "Test Feed", parsedRSS.Channel.Title)
	assert.Equal(t, "https://example.com", parsedRSS.Channel.Link)
	assert.Equal(t, "RSS feed generated from https://example.com", parsedRSS.Channel.Description)
	assert.Equal(t, "en-us", parsedRSS.Channel.Language)
	assert.Len(t, parsedRSS.Channel.Items, 1)
	assert.Equal(t, "Test Item 1", parsedRSS.Channel.Items[0].Title)
	assert.Equal(t, "https://example.com/item1", parsedRSS.Channel.Items[0].Link)
	assert.Equal(t, "This is a test item", parsedRSS.Channel.Items[0].Description)
	assert.Equal(t, "Mon, 25 Dec 2023 10:30:00 +0000", parsedRSS.Channel.Items[0].PubDate)
}

func TestHandleFeedRSS(t *testing.T) {
	// Create a mock app with test queries
	now := time.Now()

	mockQ := &mockQueries{
		GetFeedFn: func(ctx context.Context, id int64) (db.Feed, error) {
			if id == 1 {
				return db.Feed{
					ID:   1,
					Name: "Test Feed",
					Url:  "https://example.com",
				}, nil
			}
			return db.Feed{}, fmt.Errorf("feed not found")
		},
		ListFeedItemsFn: func(ctx context.Context, feedID int64) ([]db.FeedItem, error) {
			if feedID == 1 {
				return []db.FeedItem{
					{
						ID:    1,
						Title: "Test Item 1",
						Link:  "https://example.com/item1",
						Description: sql.NullString{
							String: "This is a test item",
							Valid:  true,
						},
						CreatedAt: sql.NullTime{
							Time:  now,
							Valid: true,
						},
					},
				}, nil
			}
			return []db.FeedItem{}, nil
		},
	}

	app := &App{
		queries: mockQ,
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/feed/1/", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	// Call the handler function directly
	app.handleFeedRSS(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/rss+xml; charset=utf-8", w.Header().Get("Content-Type"))

	// Parse the response body to verify it's valid XML
	body := w.Body.String()
	assert.True(t, strings.HasPrefix(body, xml.Header))

	// Parse the RSS XML to verify structure
	var rss RSS
	err := xml.Unmarshal([]byte(body), &rss)
	assert.NoError(t, err)

	// Verify RSS structure
	assert.Equal(t, "2.0", rss.Version)
	assert.Equal(t, "Test Feed", rss.Channel.Title)
	assert.Equal(t, "https://example.com", rss.Channel.Link)
	assert.Equal(t, "RSS feed generated from https://example.com", rss.Channel.Description)
	assert.Equal(t, "en-us", rss.Channel.Language)
	assert.Len(t, rss.Channel.Items, 1)

	// Verify item
	assert.Equal(t, "Test Item 1", rss.Channel.Items[0].Title)
	assert.Equal(t, "https://example.com/item1", rss.Channel.Items[0].Link)
	assert.Equal(t, "This is a test item", rss.Channel.Items[0].Description)
}

func TestHandleFeedRSSInvalidID(t *testing.T) {
	// Create an app instance
	app := &App{}

	// Create a test request with invalid ID
	req := httptest.NewRequest("GET", "/feed/invalid/", nil)
	req.SetPathValue("id", "invalid")
	w := httptest.NewRecorder()

	// Call the handler
	app.handleFeedRSS(w, req)

	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid feed ID")
}

func TestHandleFeedRSSFeedNotFound(t *testing.T) {
	// Create a mock app with test queries
	mockQ := &mockQueries{
		GetFeedFn: func(ctx context.Context, id int64) (db.Feed, error) {
			return db.Feed{}, fmt.Errorf("feed not found")
		},
	}

	app := &App{
		queries: mockQ,
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/feed/999/", nil)
	req.SetPathValue("id", "999")
	w := httptest.NewRecorder()

	// Call the handler function directly
	app.handleFeedRSS(w, req)

	// Check the response
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Feed not found")
}

func TestHandleFeedRSSItemsError(t *testing.T) {
	// Create a mock app with test queries that return an error when fetching items
	mockQ := &mockQueries{
		GetFeedFn: func(ctx context.Context, id int64) (db.Feed, error) {
			return db.Feed{
				ID:   1,
				Name: "Test Feed",
				Url:  "https://example.com",
			}, nil
		},
		ListFeedItemsFn: func(ctx context.Context, feedID int64) ([]db.FeedItem, error) {
			return []db.FeedItem{}, fmt.Errorf("failed to fetch items")
		},
	}

	app := &App{
		queries: mockQ,
	}

	// Create a test request
	req := httptest.NewRequest("GET", "/feed/1/", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	// Call the handler function directly
	app.handleFeedRSS(w, req)

	// Check the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to fetch feed items")
}

func TestRSSXMLConformance(t *testing.T) {
	// Test that the generated RSS XML conforms to RSS 2.0 specification
	now := time.Now()

	items := []Item{
		{
			Title:       "Test Item",
			Link:        "https://example.com/item",
			Description: "This is a test item description",
			PubDate:     formatRSSDate(now),
		},
	}

	rss := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "Test Feed",
			Link:        "https://example.com",
			Description: "RSS feed generated from https://example.com",
			Language:    "en-us",
			PubDate:     formatRSSDate(now),
			Items:       items,
		},
	}

	// Generate XML
	output, err := xml.MarshalIndent(rss, "", "  ")
	assert.NoError(t, err)

	xmlStr := xml.Header + string(output)

	// Check RSS 2.0 required elements
	assert.Contains(t, xmlStr, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	assert.Contains(t, xmlStr, "<rss version=\"2.0\">")

	// Check channel required elements
	assert.Contains(t, xmlStr, "<title>Test Feed</title>")
	assert.Contains(t, xmlStr, "<link>https://example.com</link>")
	assert.Contains(t, xmlStr, "<description>RSS feed generated from https://example.com</description>")

	// Check item required elements
	assert.Contains(t, xmlStr, "<item>")
	assert.Contains(t, xmlStr, "<title>Test Item</title>")
	assert.Contains(t, xmlStr, "<link>https://example.com/item</link>")

	// Check that the XML is well-formed by parsing it
	var parsedRSS RSS
	err = xml.Unmarshal([]byte(xmlStr), &parsedRSS)
	assert.NoError(t, err)

	// Verify structure
	assert.Equal(t, "2.0", parsedRSS.Version)
	assert.Equal(t, "Test Feed", parsedRSS.Channel.Title)
	assert.Equal(t, "https://example.com", parsedRSS.Channel.Link)
	assert.Equal(t, "RSS feed generated from https://example.com", parsedRSS.Channel.Description)
	assert.Len(t, parsedRSS.Channel.Items, 1)
	assert.Equal(t, "Test Item", parsedRSS.Channel.Items[0].Title)
	assert.Equal(t, "https://example.com/item", parsedRSS.Channel.Items[0].Link)
}
