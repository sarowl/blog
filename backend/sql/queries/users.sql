-- name: CreateUser :one
INSERT INTO users (id, email, display_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;