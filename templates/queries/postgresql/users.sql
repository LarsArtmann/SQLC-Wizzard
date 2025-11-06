-- name: GetUser :one
-- Get a single user by ID
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByEmail :one
-- Get a user by their email address
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByUsername :one
-- Get a user by their username
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE username = $1
LIMIT 1;

-- name: ListUsers :many
-- List users with pagination
SELECT id, email, username, full_name, is_active, is_verified, created_at, updated_at
FROM users
WHERE is_active = true
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: CreateUser :one
-- Create a new user
INSERT INTO users (
    email,
    username,
    full_name,
    password_hash
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, email, username, full_name, is_active, is_verified, created_at, updated_at;

-- name: UpdateUser :one
-- Update user details
UPDATE users
SET
    full_name = COALESCE($2, full_name),
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, full_name, is_active, is_verified, created_at, updated_at;

-- name: UpdateUserEmail :one
-- Update user email (requires re-verification)
UPDATE users
SET
    email = $2,
    is_verified = false,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, username, full_name, is_active, is_verified, created_at, updated_at;

-- name: VerifyUser :exec
-- Mark user as verified
UPDATE users
SET
    is_verified = true,
    updated_at = NOW()
WHERE id = $1;

-- name: DeactivateUser :exec
-- Soft delete a user by marking as inactive
UPDATE users
SET
    is_active = false,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
-- Permanently delete a user (use with caution)
DELETE FROM users
WHERE id = $1;

-- name: CountActiveUsers :one
-- Count total active users
SELECT COUNT(*) as count
FROM users
WHERE is_active = true;
