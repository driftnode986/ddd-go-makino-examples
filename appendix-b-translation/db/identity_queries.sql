-- name: InsertUser :exec
INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4);

-- name: UpdateUser :execrows
UPDATE users SET email = $2, name = $3 WHERE id = $1;

-- name: GetUserByID :one
SELECT id, email, name, created_at FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, name, created_at FROM users WHERE email = $1;
