-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = ? LIMIT 1;

-- name: ListFeeds :many
SELECT * FROM feeds
ORDER BY name;

-- name: CreateFeed :one
INSERT INTO feeds (name, url, item_selector, title_selector, link_selector, description_selector)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateFeed :exec
UPDATE feeds 
SET name = ?, url = ?, item_selector = ?, title_selector = ?, link_selector = ?, description_selector = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = ?;
