-- name: CreateUser :one
INSERT INTO users (id, email, full_name)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE(sqlc.narg('email'), email),
    full_name = COALESCE(sqlc.narg('full_name'), full_name),
    updated_at = (strftime('%s', 'now'))
WHERE id = ?
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: SearchUsersByEmail :many
SELECT * FROM users
WHERE email LIKE ?
ORDER BY created_at DESC;
