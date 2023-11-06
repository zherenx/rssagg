-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIME NOT NULL,
    updated_at TIME NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    published_at TIME NOT NULL,
    url TEXT NOT NULL UNIQUE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;