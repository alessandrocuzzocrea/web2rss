package app

import (
	"context"
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
	}

	// Load templates
	tmpl := template.New("")
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
