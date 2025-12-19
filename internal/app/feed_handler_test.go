package app

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestHandlePreviewFeed(t *testing.T) {
	// Create a test server to mock the website
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
			<html>
				<body>
					<div class="item">
						<h2 class="title">Item 1</h2>
						<a class="link" href="/item1">Link 1</a>
					</div>
				</body>
			</html>
		`
		_, _ = fmt.Fprint(w, html)
	}))
	defer ts.Close()

	mockQ := &mockQueries{
		ListFeedsFn: func(ctx context.Context) ([]db.Feed, error) {
			return []db.Feed{}, nil
		},
	}

	app := &App{
		queries: mockQ,
		config: &Config{Timezone: "UTC"},
	}

	// Load templates
	tmpl := template.New("").Funcs(NewTemplateFuncs(app.config))
	_, err := tmpl.ParseGlob("../../templates/*.html")
	assert.NoError(t, err)
	_, err = tmpl.ParseGlob("../../templates/partials/*.html")
	assert.NoError(t, err)
	app.templates = tmpl

	// Create a test request
	form := url.Values{}
	form.Add("url", ts.URL)
	form.Add("item_selector", ".item")
	form.Add("title_selector", ".title")
	form.Add("link_selector", ".link")

	req := httptest.NewRequest("POST", "/preview", nil)
	req.PostForm = form
	w := httptest.NewRecorder()

	app.handlePreviewFeed(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Item 1")
	assert.Contains(t, body, ts.URL+"/item1")
}

func TestHandlePreviewFeedInvalidURL(t *testing.T) {
	mockQ := &mockQueries{
		ListFeedsFn: func(ctx context.Context) ([]db.Feed, error) {
			return []db.Feed{}, nil
		},
	}
	app := &App{
		queries: mockQ,
		config: &Config{Timezone: "UTC"},
	}

	form := url.Values{}
	form.Add("url", "http://nonexistent-website-123.com")
	form.Add("item_selector", ".item")

	req := httptest.NewRequest("POST", "/preview", nil)
	req.PostForm = form
	w := httptest.NewRecorder()

	app.handlePreviewFeed(w, req)

	// Check the response - should be bad request or internal error depending on implementation
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleDuplicateFeed(t *testing.T) {
	mockQ := &mockQueries{
		GetFeedFn: func(ctx context.Context, id int64) (db.Feed, error) {
			return db.Feed{
				ID:            1,
				Name:          "Test Feed",
				Url:           "http://test.com",
				ItemSelector:  sql.NullString{String: ".item", Valid: true},
				TitleSelector: sql.NullString{String: ".title", Valid: true},
				LinkSelector:  sql.NullString{String: ".link", Valid: true},
			}, nil
		},
	}

	app := &App{
		queries: mockQ,
		config: &Config{Timezone: "UTC"},
	}

	// Load templates
	tmpl := template.New("").Funcs(NewTemplateFuncs(app.config))
	_, err := tmpl.ParseGlob("../../templates/*.html")
	assert.NoError(t, err)
	_, err = tmpl.ParseGlob("../../templates/partials/*.html")
	assert.NoError(t, err)
	app.templates = tmpl

	req := httptest.NewRequest("GET", "/feed/1/duplicate", nil)
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	app.handleDuplicateFeed(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	assert.Contains(t, body, "Test Feed (copy)")
	assert.Contains(t, body, "http://test.com")
	assert.Contains(t, body, "hx-trigger=\"change delay:500ms, load\"")
	assert.Contains(t, body, "name=\"item_selector\" value=\".item\"")
}
