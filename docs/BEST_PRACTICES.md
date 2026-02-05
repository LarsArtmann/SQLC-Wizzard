# SQLC Best Practices for Production

## Overview

This guide covers production-ready patterns and best practices for SQLC-generated code.

## Database Schema Design

### 1. Naming Conventions

```sql
-- Table names: plural, snake_case
CREATE TABLE users { ... }
CREATE TABLE user_sessions { ... }

-- Column names: snake_case
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email_address TEXT NOT NULL,  -- Prefer explicit names
    created_at TIMESTAMP,            -- Use _at, _on suffixes for timestamps
    updated_at TIMESTAMP
);

-- Index names: idx_{table}_{column(s)}
CREATE INDEX idx_users_email ON users(email_address);
CREATE INDEX idx_todos_user_status ON todos(user_id, status);
```

### 2. Use Appropriate Types

```sql
-- UUIDs for primary keys
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- TIMESTAMPTZ for all timestamps
created_at TIMESTAMPTZ DEFAULT NOW(),
updated_at TIMESTAMPTZ DEFAULT NOW()

-- TEXT for variable-length strings
email TEXT NOT NULL,              -- Not VARCHAR
description TEXT

-- Use arrays for one-to-many when reasonable
tag_ids UUID[] NOT NULL
```

### 3. Add Constraints

```sql
-- NOT NULL for required fields
email TEXT NOT NULL

-- UNIQUE constraints
email TEXT NOT NULL UNIQUE,
(organization_id, slug) UNIQUE  -- Composite unique

-- CHECK constraints for validation
status TEXT NOT NULL CHECK (status IN ('active', 'inactive', 'suspended'))
priority INTEGER NOT NULL CHECK (priority >= 0 AND priority <= 10)

-- Foreign keys with CASCADE
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
```

## SQL Query Patterns

### 1. Use Parameterized Queries

```sql
-- ✅ GOOD: Parameterized
SELECT * FROM users WHERE id = $1;

-- ❌ BAD: String interpolation
SELECT * FROM users WHERE id = 'some-id';
```

### 2. Limit Result Sets

```sql
-- Always use LIMIT for list queries
SELECT * FROM todos
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;  -- Paginate!

-- Use cursor-based pagination for large datasets
SELECT * FROM todos
WHERE user_id = $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3;
```

### 3. Use Explicit Column Selection

```sql
-- ✅ GOOD: Explicit columns
SELECT id, title, status, created_at
FROM todos
WHERE id = $1;

-- ❌ BAD: SELECT *
SELECT * FROM todos
WHERE id = $1;  -- Fetches unnecessary data
```

### 4. Use Joins Effectively

```sql
-- Simple joins for related data
SELECT
    t.id,
    t.title,
    u.email AS user_email
FROM todos t
JOIN users u ON t.user_id = u.id
WHERE t.id = $1;

-- LEFT JOIN for optional relationships
SELECT
    t.id,
    t.title,
    a.display_name AS assigned_to_name
FROM todos t
LEFT JOIN assignments a ON t.id = a.todo_id
WHERE t.id = $1;
```

## Code Organization

### 1. Package Structure

```
project/
├── cmd/
│   ├── server/main.go         # Application entry
│   └── migrate/main.go       # Migration runner
├── internal/
│   ├── db/
│   │   ├── db.go             # Database utilities
│   │   ├── models.go          # Generated structs
│   │   ├── queries.go         # Generated query functions
│   │   └── handlers.go       # DB helpers
│   ├── handlers/
│   │   ├── todos.go          # HTTP handlers
│   │   └── users.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logger.go
│   └── services/
│       └── todo_service.go   # Business logic
├── sql/
│   ├── schema/
│   │   └── 000001_*.up.sql
│   └── queries/
│       └── todos.sql
├── go.mod
└── sqlc.yaml
```

### 2. Interface Abstraction

```go
// Create database interface for dependency injection
type DBTX interface {
    Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
    Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
    QueryRow(ctx context.Context, sql string, args ...interface{}) *sql.Row
}

// Implement for sql.DB
func (db *sql.DB) DBTX { /* ... */ }

// Implement for sql.Tx
func (tx *sql.Tx) DBTX { /* ... */ }

// Use in queries
func New(ctx context.Context, db DBTX) *Queries {
    return &Queries{db: db}
}
```

### 3. Repository Pattern

```go
// repository/todos.go
type TodoRepository interface {
    Create(ctx context.Context, todo *CreateTodoParams) (*Todo, error)
    GetByID(ctx context.Context, id uuid.UUID) (*Todo, error)
    List(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*Todo, error)
    Update(ctx context.Context, todo *UpdateTodoParams) (*Todo, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type sqlTodoRepository struct {
    q *db.Queries
}

func (r *sqlTodoRepository) Create(ctx context.Context, todo *CreateTodoParams) (*Todo, error) {
    return r.q.CreateTodo(ctx, *todo)
}

// Usage: Dependency injection
func NewTodoService(repo TodoRepository) *TodoService {
    return &TodoService{repo: repo}
}
```

## Error Handling

### 1. Use Contextual Errors

```go
import "errors"

var (
    ErrNotFound   = errors.New("resource not found")
    ErrConflict   = errors.New("resource conflict")
    ErrValidation = errors.New("validation failed")
)

// Wrap with context
if err := repo.GetByID(ctx, id); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
        return nil, fmt.Errorf("todo %s not found: %w", id, ErrNotFound)
    }
    return nil, fmt.Errorf("failed to get todo %s: %w", id, err)
}
```

### 2. Handle Database Errors

```go
import "github.com/jackc/pgx/v5/pgconn"

switch err {
case pgconn.ErrNoRows:
    // Return 404 for GET requests
    return nil, ErrNotFound
case pgconn.ErrConstraintViolation:
    // Return 409 for conflicts
    return nil, ErrConflict
case pgconn.ErrConnectionDoesNotExist:
    // Return 503 for connection issues
    return nil, fmt.Errorf("database unavailable: %w", err)
default:
    // Return 500 for other errors
    return nil, fmt.Errorf("internal error: %w", err)
}
```

## Performance Optimization

### 1. Use Indexes

```sql
-- Index foreign keys
CREATE INDEX idx_todos_user_id ON todos(user_id);

-- Index query columns
CREATE INDEX idx_todos_status_created ON todos(status, created_at DESC);

-- Composite indexes for common query patterns
CREATE INDEX idx_todos_user_status ON todos(user_id, status);
```

### 2. Use Prepared Queries

```go
// SQLC generates prepared queries automatically
// Just use them correctly

// ✅ GOOD: Single query execution
todo, err := queries.CreateTodo(ctx, params)

// ❌ BAD: N+1 queries in loop
for _, title := range titles {
    _, _ := queries.CreateTodo(ctx, CreateTodoParams{Title: title})
}
```

### 3. Connection Pooling

```go
// Configure connection pool
config, _ := pgxpool.ParseConfig(databaseURL)
config.MaxConns = 25              // Max connections
config.MinConns = 5               // Min connections
config.MaxConnLifetime = time.Hour   // Connection lifetime
config.MaxConnIdleTime = 5 * time.Minute
config.HealthCheckPeriod = time.Minute

pool, _ := pgxpool.NewWithConfig(ctx, config)
```

### 4. Batch Operations

```sql
-- Use INSERT ... SELECT for bulk inserts
INSERT INTO todos (user_id, title)
SELECT user_id, title
FROM temp_import_data
WHERE processed = false;
```

## Security

### 1. Use Parameterized Queries

```sql
-- Prevents SQL injection
SELECT * FROM users WHERE email = $1;
-- Never concatenate: 'WHERE email = ''' || $1 || ''''
```

### 2. Least Privilege

```sql
-- Application user shouldn't be superuser
CREATE USER todo_app WITH PASSWORD 'secure-password';
GRANT SELECT, INSERT, UPDATE, DELETE ON todos TO todo_app;
GRANT SELECT, UPDATE ON users TO todo_app;
```

### 3. Row-Level Security

```sql
-- Add RLS policies
CREATE POLICY user_todos_policy ON todos
FOR ALL
USING (user_id = current_user_id())
WITH CHECK (user_id = current_user_id());

ALTER TABLE todos ENABLE ROW LEVEL SECURITY;
```

## Testing

### 1. Test Generated Code

```go
func TestTodoRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    queries := db.NewQueries(db)

    // Test insertion
    params := CreateTodoParams{
        UserID: testUserID,
        Title:  "Test Todo",
    }
    todo, err := queries.CreateTodo(context.Background(), params)

    assert.NoError(t, err)
    assert.Equal(t, params.Title, todo.Title)
}
```

### 2. Test Database Interactions

```go
func TestTodoRepository_Concurrency(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    // Test concurrent operations
    done := make(chan bool)
    for i := 0; i < 10; i++ {
        go func(id int) {
            _, _ := queries.CreateTodo(ctx, params)
            done <- true
        }(i)
    }

    // Wait for all operations
    for i := 0; i < 10; i++ {
        <-done
    }
}
```

## Deployment

### 1. Environment Configuration

```yaml
# config/production.yaml
database:
  url: ${DATABASE_URL}
  pool:
    max_conns: 25
    min_conns: 5

server:
  port: ${PORT:-8080}
  read_timeout: 5s
  write_timeout: 10s
```

### 2. Migration Strategy

```bash
# Run migrations on startup
migrate up

# Rollback if needed
migrate down 1

# Get status
migrate status
```

### 3. Monitoring

```go
// Add metrics to queries
func (r *sqlTodoRepository) Create(ctx context.Context, params *CreateTodoParams) (*Todo, error) {
    start := time.Now()

    todo, err := r.q.CreateTodo(ctx, *params)

    // Record metrics
    duration := time.Since(start)
    metrics.QueryDuration.Observe(duration.Seconds())
    metrics.QueryErrors.Inc(float64(boolToUint64(err != nil)))

    return todo, err
}
```

## Summary

Follow these best practices for:

✅ Type-safe database interactions
✅ Performant queries
✅ Secure against SQL injection
✅ Maintainable codebase
✅ Testable architecture
✅ Production-ready deployment

### Additional Resources:

- [SQLC Documentation](https://docs.sqlc.dev/)
- [PostgreSQL Best Practices](https://wiki.postgresql.org/wiki/Don't_Do_This)
- [Go Database Best Practices](https://go.dev/doc/database/sql-package)
