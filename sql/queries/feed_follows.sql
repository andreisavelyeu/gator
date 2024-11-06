-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, feed_ID, user_ID)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN users on inserted_feed_follow.user_ID = users.id
INNER JOIN feeds on inserted_feed_follow.feed_ID = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT ff.user_ID, ff.feed_ID, feeds.name from feed_follows ff
INNER JOIN feeds ON feeds.id = ff.feed_ID
INNER JOIN users ON users.id = ff.user_ID
WHERE users.id = $1;

-- name: DeleteFeedFollowByUserAndUrl :exec
DELETE FROM feed_follows 
USING feeds 
WHERE feed_follows.user_ID = $1 
  AND feed_follows.feed_ID = feeds.id
  AND feeds.url = $2;
