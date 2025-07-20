-- +goose Up
CREATE TABLE posts (
    id uuid PRIMARY KEY,
    created_at timestamp NOT NULL,
    updated_at timestamp NOT NULL,
    title text NOT NULL,
    url text UNIQUE NOT NULL,
    description text,
    published_at timestamp NOT NULL,
    feed_id uuid NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    FOREIGN KEY(feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE posts;
