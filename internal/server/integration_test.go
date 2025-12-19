package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alessandrocuzzocrea/web2rss/internal/config"
	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func TestIntegration(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "web2rss-test-*")
	assert.NoError(t, err)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	dbPath := tmpDir + "/test.sqlite3"

	// Initialize database
	database, err := sql.Open("sqlite", dbPath)
	assert.NoError(t, err)

	schema, err := os.ReadFile("../../db/schema.sql")
	assert.NoError(t, err)
	_, err = database.Exec(string(schema))
	assert.NoError(t, err)
	database.Close() // Close it so Server can open it

	cfg := &config.Config{
		Port:        "8080",
		DBPath:      dbPath,
		DataDir:     tmpDir,
		Timezone:    "UTC",
		TemplateDir: "../../templates",
	}

	app, err := New(cfg)
	assert.NoError(t, err)
	defer app.Close()

	// 1. Create a mock website
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
			<html>
				<body>
					<div class="item">
						<h2 class="title">Blog Post 1</h2>
						<a class="link" href="/post1">Read more</a>
					</div>
				</body>
			</html>
		`
		_, _ = fmt.Fprint(w, html)
	}))
	defer ts.Close()

	// Access DB directly for verification since we are in the same package
	queries := db.New(app.db)

	// 2. Add as a feed via queries
	feed, err := queries.CreateFeed(context.Background(), db.CreateFeedParams{
		Name:          "Test Blog",
		Url:           ts.URL,
		ItemSelector:  sql.NullString{String: ".item", Valid: true},
		TitleSelector: sql.NullString{String: ".title", Valid: true},
		LinkSelector:  sql.NullString{String: ".link", Valid: true},
	})
	assert.NoError(t, err)

	// 3. Trigger refresh via FeedService
	err = app.feedService.RefreshFeed(context.Background(), feed)
	assert.NoError(t, err)

	// 4. Check if item was added
	items, err := queries.ListFeedItems(context.Background(), feed.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Blog Post 1", items[0].Title)
	assert.Equal(t, ts.URL+"/post1", items[0].Link)

	// 5. Check RSS output via Handler
	// We get the Mux from the handler
	mux := app.handler.RegisterRoutes()

	req := httptest.NewRequest("GET", fmt.Sprintf("/feed/%d/rss", feed.ID), nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Blog Post 1")
	assert.Contains(t, w.Body.String(), ts.URL+"/post1")
}
