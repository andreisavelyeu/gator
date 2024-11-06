-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_ID UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_ID UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE (user_ID, feed_ID)
);

-- +goose Down
DROP TABLE feed_follows;
