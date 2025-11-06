-- name: GetUser :one
-- Get a single user by ID
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE id = ?
LIMIT 1;

-- name: GetUserByEmail :one
-- Get a user by their email address
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE email = ?
LIMIT 1;

-- name: GetUserByUsername :one
-- Get a user by their username
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE username = ?
LIMIT 1;

-- name: ListUsers :many
-- List users with pagination
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE is_active = 1
ORDER BY created_at DESC
LIMIT ?
OFFSET ?;

-- name: CreateUser :one
-- Create a new user
INSERT INTO users (
    email,
    username,
    full_name,
    password_hash
) VALUES (
    ?, ?, ?, ?
)
RETURNING id, email, username, full_name, is_active, is_verified, created_at, updated_at;

-- name: UpdateUser :one
-- Update user details
UPDATE users
SET
    full_name = COALESCE(?, full_name),
    updated_at = datetime('now')
WHERE id = ?
RETURNING id, email, username, full_name, is_active, is_verified, created_at, updated_at;

-- name: UpdateUserEmail :one
-- Update user email (requires re-verification)
UPDATE users
SET
    email = ?,
    is_verified = 0,
    updated_at = datetime('now')
WHERE id = ?
RETURNING id, email, username, full_name, is_active, is_verified, created_at, updated_at;

-- name: VerifyUser :exec
-- Mark user as verified
UPDATE users
SET
    is_verified = 1,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeactivateUser :exec
-- Soft delete a user by marking as inactive
UPDATE users
SET
    is_active = 0,
    updated_at = datetime('now')
WHERE id = ?;

-- name: DeleteUser :exec
-- Permanently delete a user (use with caution)
DELETE FROM users
WHERE id = ?;

-- name: CountActiveUsers :one
-- Count total active users
SELECT COUNT(*) as count
FROM users
WHERE is_active = 1;
