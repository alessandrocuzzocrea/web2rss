-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = ? LIMIT 1;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY name;

-- name: ListFeedsWithItemsCount :many
SELECT f.*, COUNT(i.id) AS items_count
FROM feeds f
LEFT JOIN feed_items i ON f.id = i.feed_id
GROUP BY f.id
ORDER BY f.name;

-- name: CreateFeed :one
INSERT INTO feeds (name, url, item_selector, title_selector, link_selector, description_selector)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFeed :exec
UPDATE feeds
SET name = ?, url = ?, item_selector = ?, title_selector = ?, link_selector = ?, description_selector = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateFeedLastRefreshedAt :exec
UPDATE feeds
SET last_refreshed_at = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = ?;
