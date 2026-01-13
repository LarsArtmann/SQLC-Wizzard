package generators

import "github.com/LarsArtmann/SQLC-Wizzard/internal/templates"

// getQueryTemplate returns the embedded query template for a database type.
func getQueryTemplate(db templates.DatabaseType) string {
	switch db {
	case templates.DatabaseTypePostgreSQL:
		return postgresqlQueriesTemplate
	case templates.DatabaseTypeSQLite:
		return sqliteQueriesTemplate
	case templates.DatabaseTypeMySQL:
		return mysqlQueriesTemplate
	default:
		return ""
	}
}

// getSchemaTemplate returns the embedded schema template for a database type.
func getSchemaTemplate(db templates.DatabaseType) string {
	switch db {
	case templates.DatabaseTypePostgreSQL:
		return postgresqlSchemaTemplate
	case templates.DatabaseTypeSQLite:
		return sqliteSchemaTemplate
	case templates.DatabaseTypeMySQL:
		return mysqlSchemaTemplate
	default:
		return ""
	}
}

const postgresqlQueriesTemplate = `-- name: GetUser :one
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

-- name: DeleteUser :exec
-- Permanently delete a user (use with caution)
DELETE FROM users
WHERE id = $1;
`

const sqliteQueriesTemplate = `-- name: GetUser :one
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

-- name: DeleteUser :exec
-- Permanently delete a user (use with caution)
DELETE FROM users
WHERE id = ?;
`

const mysqlQueriesTemplate = `-- name: GetUser :one
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

-- name: CreateUser :exec
-- Create a new user
INSERT INTO users (
    email,
    username,
    full_name,
    password_hash
) VALUES (
    ?, ?, ?, ?
);

-- name: DeleteUser :exec
-- Permanently delete a user (use with caution)
DELETE FROM users
WHERE id = ?;
`

const postgresqlSchemaTemplate = `-- Example user table schema for PostgreSQL
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
`

const sqliteSchemaTemplate = `-- Example user table schema for SQLite
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    full_name TEXT,
    password_hash TEXT NOT NULL,
    is_active INTEGER NOT NULL DEFAULT 1,
    is_verified INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
`

const mysqlSchemaTemplate = `-- Example user table schema for MySQL
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    full_name VARCHAR(255),
    password_hash VARCHAR(255) NOT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    is_verified TINYINT(1) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_users_email (email),
    INDEX idx_users_username (username)
);
`
