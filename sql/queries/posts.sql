-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, url, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsForUser :many
SELECT feeds.name, feeds.url AS feed_url, posts.id, posts.title, posts.description, posts.published_at, posts.url AS post_url 
FROM posts
JOIN feeds ON feeds.id = posts.feed_id
WHERE posts.feed_id IN (
    SELECT feed_id FROM feed_follows
    WHERE feed_follows.user_id = $1
)
ORDER BY published_at DESC
LIMIT $2;