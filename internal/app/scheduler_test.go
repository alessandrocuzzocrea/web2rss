package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/stretchr/testify/assert"
)



func TestRefreshFeed(t *testing.T) {
	// Create a test server to mock the website
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
			<html>
				<body>
					<div class="item">
						<h2 class="title">Item 1</h2>
						<a class="link" href="/item1">Link 1</a>
						<span class="date">2023-12-25</span>
					</div>
					<div class="item">
						<h2 class="title">Item 2</h2>
						<a class="link" href="https://other.com/item2">Link 2</a>
						<span class="date">2023/12/26</span>
					</div>
				</body>
			</html>
		`
		fmt.Fprint(w, html)
	}))
	defer ts.Close()

	// Mock queries
	var upsertedItems []db.UpsertFeedItemParams
	mockQ := &mockQueries{
		UpsertFeedItemFn: func(ctx context.Context, params db.UpsertFeedItemParams) ([]int64, error) {
			upsertedItems = append(upsertedItems, params)
			return []int64{int64(len(upsertedItems))}, nil
		},
		UpdateFeedLastRefreshedAtFn: func(ctx context.Context, params db.UpdateFeedLastRefreshedAtParams) error {
			return nil
		},
	}

	app := &App{
		queries: mockQ,
	}

	feed := db.Feed{
		ID:            1,
		Url:           ts.URL,
		ItemSelector:  db.NewNullString(".item"),
		TitleSelector: db.NewNullString(".title"),
		LinkSelector:  db.NewNullString(".link"),
		DateSelector:  db.NewNullString(".date"),
	}

	err := app.refreshFeed(context.Background(), feed)
	assert.NoError(t, err)

	assert.Len(t, upsertedItems, 2)

	// Verify Item 1 (relative link resolved)
	assert.Equal(t, "Item 1", upsertedItems[0].Title)
	assert.Equal(t, ts.URL+"/item1", upsertedItems[0].Link)
	assert.Equal(t, 2023, upsertedItems[0].Date.Time.Year())
	assert.Equal(t, time.Month(12), upsertedItems[0].Date.Time.Month())
	assert.Equal(t, 25, upsertedItems[0].Date.Time.Day())

	// Verify Item 2 (absolute link kept)
	assert.Equal(t, "Item 2", upsertedItems[1].Title)
	assert.Equal(t, "https://other.com/item2", upsertedItems[1].Link)
	assert.Equal(t, 2023, upsertedItems[1].Date.Time.Year())
	assert.Equal(t, time.Month(12), upsertedItems[1].Date.Time.Month())
	assert.Equal(t, 26, upsertedItems[1].Date.Time.Day())
}
