package app

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestHandleHomepage(t *testing.T) {
	mockQ := &mockQueries{
		ListFeedsWithItemsCountFn: func(ctx context.Context) ([]db.ListFeedsWithItemsCountRow, error) {
			return []db.ListFeedsWithItemsCountRow{
				{
					ID:         1,
					Name:       "Test Feed",
					Url:        "http://test.com",
					ItemsCount: 10,
				},
			}, nil
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

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	app.handleHomepage(w, req)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()

	// Check for the 3-dots dropdown
	assert.Contains(t, body, "action-menu")
	assert.Contains(t, body, "⋮")

	// Check for the options inside the dropdown
	assert.Contains(t, body, "/feed/1/edit")
	assert.Contains(t, body, "/feed/1/duplicate")
	assert.Contains(t, body, "/feed/1/reset")
	assert.Contains(t, body, "/feed/1/delete")

	// Check that RSS and Refresh are still present
	assert.Contains(t, body, "/feed/1/rss")
	assert.Contains(t, body, "/feed/1/refresh")
}
