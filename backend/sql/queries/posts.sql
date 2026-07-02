-- name: CreatePost :one
INSERT INTO posts (id, user_id, title, body, status, created_at, updated_at)
VALUES (gen_random_uuid(), $1, $2, $3, $4, now(), now())
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts
WHERE id = $1;

-- name: ListPosts :many
SELECT 
    p.id,
    p.user_id,
    u.display_name AS author_name,
    p.title,
    p.body,
    p.status,
    p.created_at,
    p.updated_at
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE p.status = 'public'
ORDER BY p.created_at DESC;


-- name: ListUserPosts :many
SELECT id, user_id, title, body, status, created_at, updated_at FROM posts
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE posts
SET title = $2, body = $3, status = $4, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1 AND user_id = $2;