# SQLC-Wizzard Tutorial: Building a REST API

## Overview

This tutorial walks you through building a complete REST API using SQLC-Wizzard. We'll create a **microservice** with **PostgreSQL**, including:

- 🗄️ Database schema
- 🔍 SQL queries
- 🔧 Generated Go code
- 🚀 Simple HTTP handlers

## Prerequisites

- Go 1.21+
- PostgreSQL (or use Docker)
- SQLC-Wizzard installed

```bash
# Install SQLC-Wizzard
go install github.com/LarsArtmann/SQLC-Wizzard/cmd/wizard@latest
```

## Step 1: Initialize Project

```bash
# Create new project directory
mkdir todo-api
cd todo-api

# Initialize Go module
go mod init github.com/username/todo-api

# Run wizard
wizard
```

### Wizard Responses:

```
📍 Project Type Selection
> ⚡ Microservice - Single DB, container-optimized

📍 Database Selection
> 🐘 PostgreSQL - Full-featured, recommended

📍 Project Details
Project name: todo
Package name: todo

📍 Output Configuration
Base directory: ./internal/db
SQL queries directory: ./sql/queries
SQL schema directory: ./sql/schema

📍 Features & Validation
✅ Emit Interface
✅ Emit JSON Tags
✅ No SELECT * Rule
✅ Require WHERE Rule
```

This generates `sqlc.yaml`:

```yaml
version: "2"
sql:
  - name: "todo"
    engine: "postgresql"
    queries: "./sql/queries"
    schema: "./sql/schema"
    gen:
      go:
        package: "todo"
        out: "./internal/db"
        sql_package: "pgx/v5"
        emit_interface: true
        emit_json_tags: true
        emit_prepared_queries: true
```

## Step 2: Create Database Schema

Create `./sql/schema/000001_create_tables.up.sql`:

```sql
-- migrations: 000001_create_tables.up.sql

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Todos table
CREATE TABLE todos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL DEFAULT 'pending',
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_created_at ON todos(created_at DESC);
```

Create rollback `./sql/schema/000001_create_tables.down.sql`:

```sql
-- migrations: 000001_create_tables.down.sql

DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_status;
DROP INDEX IF EXISTS idx_todos_user_id;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS users;
```

## Step 3: Create SQL Queries

Create `./sql/queries/todos.sql`:

```sql
-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todos (user_id, title, description, status, priority)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListTodos :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetTodoByID :one
SELECT * FROM todos
WHERE id = $1;

-- name: UpdateTodo :one
UPDATE todos
SET
    title = COALESCE($2, title),
    description = COALESCE($3, description),
    status = COALESCE($4, status),
    priority = COALESCE($5, priority),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;
```

## Step 4: Generate Go Code

```bash
# Generate code using SQLC
sqlc generate

# Output will be in ./internal/db/
tree internal/db/
```

Output structure:
```
internal/db/
├── db.go              # Database connection utilities
├── models.go           # Generated Go structs
├── todos.sql.go        # Generated query functions
└── sqlc.yaml          # Configuration copy
```

## Step 5: Create HTTP Handlers

Create `./internal/handlers/todos.go`:

```go
package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/LarsArtmann/SQLC-Wizzard/internal/db"
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type TodoHandler struct {
    db *pgxpool.Pool
    queries *db.Queries
}

func NewTodoHandler(db *pgxpool.Pool) *TodoHandler {
    return &TodoHandler{
        db:      db,
        queries: db.New(db),
    }
}

// CreateTodoRequest represents JSON request body
type CreateTodoRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Status      string `json:"status"`
    Priority    int    `json:"priority"`
}

// CreateTodo handles POST /todos
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := getUserIDFromContext(ctx)

    // Parse request
    var req CreateTodoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Create todo
    todo, err := h.queries.CreateTodo(ctx, db.CreateTodoParams{
        UserID:      userID,
        Title:       req.Title,
        Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
        Status:      req.Status,
        Priority:    sql.NullInt32{Int32: int32(req.Priority), Valid: req.Priority != 0},
    })
    if err != nil {
        http.Error(w, "Failed to create todo", http.StatusInternalServerError)
        return
    }

    // Return response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

// ListTodos handles GET /todos
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := getUserIDFromContext(ctx)

    // Parse query params
    limit := int32(10)
    offset := int32(0)

    // Get todos
    todos, err := h.queries.ListTodos(ctx, db.ListTodosParams{
        UserID: userID,
        Limit:  limit,
        Offset: offset,
    })
    if err != nil {
        http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
        return
    }

    // Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todos)
}

// Helper function
func getUserIDFromContext(ctx context.Context) uuid.UUID {
    // In real app, get from JWT/session
    return uuid.MustParse("00000000-0000-0000-0000-000000000000")
}
```

## Step 6: Create Main Application

Create `./cmd/server/main.go`:

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    _ "github.com/LarsArtmann/SQLC-Wizzard/internal/db"
    "github.com/LarsArtmann/SQLC-Wizzard/internal/handlers"
)

func main() {
    // Connect to database
    ctx := context.Background()
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://localhost:5432/todo?sslmode=disable"
    }

    config, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        log.Fatal(err)
    }

    pool, err := pgxpool.NewWithConfig(ctx, config)
    if err != nil {
        log.Fatal(err)
    }
    defer pool.Close()

    // Run migrations (optional)
    // ... run your migration tool here ...

    // Setup handlers
    todoHandler := handlers.NewTodoHandler(pool)

    // Setup routes
    mux := http.NewServeMux()
    mux.HandleFunc("/todos", todoHandler.ListTodos).Methods("GET")
    mux.HandleFunc("/todos", todoHandler.CreateTodo).Methods("POST")

    // Start server
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    log.Println("Server starting on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
```

## Step 7: Run and Test

```bash
# Start PostgreSQL (Docker)
docker run --name postgres-todo -e POSTGRES_PASSWORD=todo -p 5432:5432 -d postgres:15

# Set environment variable
export DATABASE_URL="postgres://postgres:todo@localhost:5432/todo?sslmode=disable"

# Run migrations
go run cmd/migrate/main.go

# Start server
go run cmd/server/main.go
```

### Test with curl:

```bash
# Create todo
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Learn SQLC",
    "description": "Master SQL code generation",
    "status": "pending",
    "priority": 1
  }'

# List todos
curl http://localhost:8080/todos
```

## Step 8: Production Tips

### 1. Use Prepared Queries

Generated code already uses prepared queries for security and performance:

```go
// This is automatically generated
func (q *Queries) CreateTodo(ctx context.Context, db DBTX, params CreateTodoParams) (*Todo, error) {
    // ... uses prepared statement internally
}
```

### 2. Type Safety

All queries are type-safe:

```go
// Compile-time error if wrong type
params := db.CreateTodoParams{
    UserID: userID,
    Title: 123, // ❌ Error: expected string, got int
}
```

### 3. Database Transactions

```go
// Use transactions for multiple operations
tx, err := pool.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback()

todo, err := queries.CreateTodo(ctx, tx, params)
if err != nil {
    return err
}

if err := tx.Commit(); err != nil {
    return err
}
```

### 4. Error Handling

```go
// Handle specific errors
switch err {
case pgx.ErrNoRows:
    // Not found - return 404
case pgx.ErrConnDone:
    // Connection closed - return 503
default:
    // Internal server error - return 500
}
```

## Summary

You've built a complete REST API with:

✅ Type-safe database queries
✅ Generated Go code (no manual SQL)
✅ Production-ready configuration
✅ Best practices enforced
✅ Easy to maintain and extend

### Next Steps:

- Add authentication (JWT/OAuth)
- Add input validation
- Add pagination
- Add caching layer
- Deploy to production

### Additional Resources:

- [User Guide](./USER_GUIDE.md)
- [Advanced Features](./advanced-features.md)
- [Best Practices](./best-practices.md)
- [SQLC Documentation](https://docs.sqlc.dev/)
