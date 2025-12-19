package app

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

	// Initialize base database schema (we need migrations for this to be really clean,
	// but for now we'll just open it and sqlc will work if schema exists)
	// Actually, let's just use the schema.sql file if we can find it.

	database, err := sql.Open("sqlite", dbPath)
	assert.NoError(t, err)

	schema, err := os.ReadFile("../../db/schema.sql")
	assert.NoError(t, err)
	_, err = database.Exec(string(schema))
	assert.NoError(t, err)

	cfg := &Config{
		Port:     "8080",
		DBPath:   dbPath,
		DataDir:  tmpDir,
		Timezone: "UTC",
	}

	app := &App{
		db:      database,
		queries: db.New(database),
		config:  cfg,
	}

	// Load templates
	tmpl := template.New("").Funcs(NewTemplateFuncs(cfg))
	_, err = tmpl.ParseGlob("../../templates/*.html")
	assert.NoError(t, err)
	_, err = tmpl.ParseGlob("../../templates/partials/*.html")
	assert.NoError(t, err)
	app.templates = tmpl

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

	// 2. Add as a feed via queries
	feed, err := app.queries.CreateFeed(context.Background(), db.CreateFeedParams{
		Name:          "Test Blog",
		Url:           ts.URL,
		ItemSelector:  sql.NullString{String: ".item", Valid: true},
		TitleSelector: sql.NullString{String: ".title", Valid: true},
		LinkSelector:  sql.NullString{String: ".link", Valid: true},
	})
	assert.NoError(t, err)

	// 3. Trigger refresh
	err = app.refreshFeed(context.Background(), feed)
	assert.NoError(t, err)

	// 4. Check if item was added
	items, err := app.queries.ListFeedItems(context.Background(), feed.ID)
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Blog Post 1", items[0].Title)
	assert.Equal(t, ts.URL+"/post1", items[0].Link)

	// 5. Check RSS output
	req := httptest.NewRequest("GET", fmt.Sprintf("/feed/%d/", feed.ID), nil)
	req.SetPathValue("id", fmt.Sprintf("%d", feed.ID))
	w := httptest.NewRecorder()
	app.handleFeedRSS(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Blog Post 1")
	assert.Contains(t, w.Body.String(), ts.URL+"/post1")
}
