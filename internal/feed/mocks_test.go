package feed

import (
	"context"

	"github.com/alessandrocuzzocrea/web2rss/internal/db"
)

// mockQueries satisfies the Querier interface for testing
//
//nolint:dupl
type mockQueries struct {
	GetFeedFn                   func(ctx context.Context, id int64) (db.Feed, error)
	ListFeedsFn                 func(ctx context.Context) ([]db.Feed, error)
	ListFeedsWithItemsCountFn   func(ctx context.Context) ([]db.ListFeedsWithItemsCountRow, error)
	CreateFeedFn                func(ctx context.Context, arg db.CreateFeedParams) (db.Feed, error)
	UpdateFeedFn                func(ctx context.Context, arg db.UpdateFeedParams) error
	UpdateFeedLastRefreshedAtFn func(ctx context.Context, arg db.UpdateFeedLastRefreshedAtParams) error
	DeleteFeedFn                func(ctx context.Context, id int64) error
	ListFeedItemsFn             func(ctx context.Context, feedID int64) ([]db.FeedItem, error)
	UpsertFeedItemFn            func(ctx context.Context, arg db.UpsertFeedItemParams) ([]int64, error)
	DeleteItemsByFeedIDFn       func(ctx context.Context, feedID int64) error
}

func (m *mockQueries) GetFeed(ctx context.Context, id int64) (db.Feed, error) {
	if m.GetFeedFn != nil {
		return m.GetFeedFn(ctx, id)
	}
	return db.Feed{}, nil
}
func (m *mockQueries) ListFeeds(ctx context.Context) ([]db.Feed, error) {
	if m.ListFeedsFn != nil {
		return m.ListFeedsFn(ctx)
	}
	return nil, nil
}
func (m *mockQueries) ListFeedsWithItemsCount(ctx context.Context) ([]db.ListFeedsWithItemsCountRow, error) {
	if m.ListFeedsWithItemsCountFn != nil {
		return m.ListFeedsWithItemsCountFn(ctx)
	}
	return nil, nil
}
func (m *mockQueries) CreateFeed(ctx context.Context, arg db.CreateFeedParams) (db.Feed, error) {
	if m.CreateFeedFn != nil {
		return m.CreateFeedFn(ctx, arg)
	}
	return db.Feed{}, nil
}
func (m *mockQueries) UpdateFeed(ctx context.Context, arg db.UpdateFeedParams) error {
	if m.UpdateFeedFn != nil {
		return m.UpdateFeedFn(ctx, arg)
	}
	return nil
}
func (m *mockQueries) UpdateFeedLastRefreshedAt(ctx context.Context, arg db.UpdateFeedLastRefreshedAtParams) error {
	if m.UpdateFeedLastRefreshedAtFn != nil {
		return m.UpdateFeedLastRefreshedAtFn(ctx, arg)
	}
	return nil
}
func (m *mockQueries) DeleteFeed(ctx context.Context, id int64) error {
	if m.DeleteFeedFn != nil {
		return m.DeleteFeedFn(ctx, id)
	}
	return nil
}
func (m *mockQueries) ListFeedItems(ctx context.Context, feedID int64) ([]db.FeedItem, error) {
	if m.ListFeedItemsFn != nil {
		return m.ListFeedItemsFn(ctx, feedID)
	}
	return nil, nil
}
func (m *mockQueries) UpsertFeedItem(ctx context.Context, arg db.UpsertFeedItemParams) ([]int64, error) {
	if m.UpsertFeedItemFn != nil {
		return m.UpsertFeedItemFn(ctx, arg)
	}
	return nil, nil
}
func (m *mockQueries) DeleteItemsByFeedID(ctx context.Context, feedID int64) error {
	if m.DeleteItemsByFeedIDFn != nil {
		return m.DeleteItemsByFeedIDFn(ctx, feedID)
	}
	return nil
}
