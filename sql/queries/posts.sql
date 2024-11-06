-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, published_at, url, title, description, feed_ID)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPosts :many
SELECT * from posts order by published_at DESC limit $1;