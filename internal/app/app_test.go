package app

// func TestRefreshFeed(t *testing.T) {
// 	a := &App{}
// 	// Note: This test is illustrative. In a real-world scenario, you would set up
// 	// a test database, mock HTTP requests, and verify the results of the refreshFeed function.
// 	// This might involve using a library like "net/http/httptest" to create mock servers
// 	// and responses, as well as setting up and tearing down test data in the database.

// 	// This is a placeholder test function.
// 	// Implement actual tests for the refreshFeed function here.
// 	feeds := []db.Feed{
// 		{ID: 1, Name: "Test Feed 1", Url: "http://example.com/feed1", ItemSelector: sql.NullString{String: ".item", Valid: true}},
// 		{ID: 2, Name: "Test Feed 2", Url: "http://example.com/feed2", ItemSelector: sql.NullString{String: ".entry", Valid: true}},
// 	}

// 	for _, feed := range feeds {
// 		t.Run(feed.Name, func(t *testing.T) {
// 			// Here you would call a.refreshFeed with a mock context and the feed
// 			// and check for expected results.
// 			// For example:
// 			err := a.refreshFeed(context.Background(), feed)
// 			if err != nil {
// 				t.Errorf("Failed to refresh feed %s: %v", feed.Name, err)
// 			}
// 		})
// 	}
// }
