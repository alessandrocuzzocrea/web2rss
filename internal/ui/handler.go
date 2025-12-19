package ui

import (
	"context"
	"html/template"

	"github.com/alessandrocuzzocrea/web2rss/internal/config"
	"github.com/alessandrocuzzocrea/web2rss/internal/db"
	"github.com/alessandrocuzzocrea/web2rss/internal/feed"
)

// Querier defines the interface for database operations needed by the handler
type Querier interface {
	GetFeed(ctx context.Context, id int64) (db.Feed, error)
	ListFeeds(ctx context.Context) ([]db.Feed, error)
	ListFeedsWithItemsCount(ctx context.Context) ([]db.ListFeedsWithItemsCountRow, error)
	CreateFeed(ctx context.Context, arg db.CreateFeedParams) (db.Feed, error)
	UpdateFeed(ctx context.Context, arg db.UpdateFeedParams) error
	DeleteFeed(ctx context.Context, id int64) error
	ListFeedItems(ctx context.Context, feedID int64) ([]db.FeedItem, error)
	DeleteItemsByFeedID(ctx context.Context, feedID int64) error
}

type Handler struct {
	queries     Querier
	templates   *template.Template
	feedService *feed.Service
	config      *config.Config
}

func NewHandler(q Querier, t *template.Template, fs *feed.Service, cfg *config.Config) *Handler {
	return &Handler{
		queries:     q,
		templates:   t,
		feedService: fs,
		config:      cfg,
	}
}
