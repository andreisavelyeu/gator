-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, url, name, user_ID)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;


-- name: GetFeeds :many
SELECT f.name, f.url, u.name as user_name from feeds f INNER JOIN users u ON u.id = f.user_ID ORDER BY f.updated_at;

-- name: GetFeed :one
SELECT * from feeds where url = $1;