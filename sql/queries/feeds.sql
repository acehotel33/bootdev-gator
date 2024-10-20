-- name: CreateFeed :one
INSERT INTO feeds ( id, created_at, updated_at, name, url, user_id  ) 
VALUES ( 
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1
RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT f.* 
FROM feeds f
JOIN feed_follows ff ON f.id = ff.feed_id
WHERE ff.user_id = $1
ORDER BY f.last_fetched_at NULLS FIRST
LIMIT 1;
