-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  published_at TIMESTAMP NOT NULL,
  url VARCHAR(255) UNIQUE NOT NULL,
  title VARCHAR(255) NOT NULL,
  description VARCHAR(255),
  feed_ID UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;

