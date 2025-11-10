-- name: GetUser :one
SELECT id, username, created_at FROM users WHERE id = $1;