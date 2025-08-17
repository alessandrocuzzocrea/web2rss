-- name: GetFeed :one
SELECT * FROM feeds
WHERE id = ? LIMIT 1;

-- name: GetFirstFeed :one
select * from feeds limit 1;

-- -- name: ListAuthors :many
-- SELECT * FROM authors
-- ORDER BY name;

-- -- name: CreateAuthor :one
-- INSERT INTO authors (
--   name, bio, loller
-- ) VALUES (
--   ?, ?, ?
-- )
-- RETURNING *;

-- -- name: UpdateAuthor :exec
-- UPDATE authors
-- set name = ?,
-- bio = ?
-- WHERE id = ?;

-- -- name: DeleteAuthor :exec
-- DELETE FROM authors
-- WHERE id = ?;
