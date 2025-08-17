-- name: GetFeedItem :one
SELECT * FROM feed_items
WHERE id = ? LIMIT 1;

-- name: ListFeedItems :many
SELECT * FROM feed_items
WHERE feed_id = ?
ORDER BY created_at DESC;

-- name: UpsertFeedItem :exec
INSERT OR REPLACE INTO feed_items (feed_id, title, description, link, updated_at)
VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP);

-- name: DeleteFeedItem :exec
DELETE FROM feed_items
WHERE id = ?;

-- name: DeleteOldFeedItems :exec
DELETE FROM feed_items 
WHERE feed_id = ? AND created_at < ?;
