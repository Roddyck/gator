-- name: CreatePost :one
INSERT INTO posts (
    id,
    created_at,
    updated_at,
    title,
    url,
    description,
    published_at,
    feed_id
) VALUES ( $1, $2, $3, $4, $5, $6, $7, $8 )
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name
FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
INNER JOIN users ON feeds.user_id = users.id
WHERE users.id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
