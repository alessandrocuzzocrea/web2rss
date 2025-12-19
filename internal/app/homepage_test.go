package app

import (
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
					CreatedAt: sql.NullTime{
						Time:  time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC),
						Valid: true,
					},
					LastRefreshedAt: sql.NullTime{
						Time:  time.Date(2025, 1, 2, 15, 30, 0, 0, time.UTC),
						Valid: true,
					},
				},
			}, nil
		},
	}

	app := &App{
		queries: mockQ,
	}

	// Load templates
	tmpl := template.New("").Funcs(templateFuncs)
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
	// Check for specific date format (YYYY-MM-DD HH:mm:ss MST)
	// We mocked CreatedAt as "2025-01-01 12:00:00" UTC
	expectedDate := "2025-01-01 12:00:00 UTC"
	assert.Contains(t, w.Body.String(), expectedDate, "Homepage should contain formatted CreatedAt date")

	// We mocked LastRefreshedAt as "2025-01-02 15:30:00" UTC
	expectedRefresh := "2025-01-02 15:30:00 UTC"
	assert.Contains(t, w.Body.String(), expectedRefresh, "Homepage should contain formatted LastRefreshedAt date")
	assert.Contains(t, body, "/feed/1/duplicate")
	assert.Contains(t, body, "/feed/1/reset")
	assert.Contains(t, body, "/feed/1/delete")

	// Check that RSS and Refresh are still present
	assert.Contains(t, body, "/feed/1/rss")
	assert.Contains(t, body, "/feed/1/refresh")
}
