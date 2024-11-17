-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at,name,url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedByUrl :one
SELECT * FROM feeds WHERE url = $1 LIMIT 1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedAsFetched :one
UPDATE feeds SET last_fetched_at = $2,updated_at = $2 WHERE id = $1 RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds 
WHERE last_fetched_at IS NULL OR last_fetched_at < $1 
ORDER BY last_fetched_at IS NOT NULL,last_fetched_at ASC LIMIT 1;