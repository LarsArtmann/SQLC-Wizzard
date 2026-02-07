# Hobby Project Example

A complete example of a hobby project using SQLC-Wizzard.

## Project Structure

```
examples/hobby-project/
├── README.md              # This file
├── go.mod                 # Go module
├── sqlc.yaml              # SQLC configuration
└── sql/
    ├── schema/             # Database schema
    │   └── 000001_create_tables.up.sql
    └── queries/            # SQL queries
        └── users.sql
```

## Getting Started

### 1. Initialize the Project

```bash
# From the project root
cd examples/hobby-project

# Install dependencies
go mod tidy

# Generate code from SQL files
sqlc generate
```

### 2. Database Setup

#### Option A: SQLite (Easiest - No setup required)

```bash
# SQLite is embedded, no database setup needed
# The project will create a 'database.db' file automatically
```

#### Option B: PostgreSQL

```bash
# Start PostgreSQL
brew services start postgresql  # macOS
# or
sudo systemctl start postgresql  # Linux

# Create database
psql -U postgres -c "CREATE DATABASE hobby_project;"

# Set environment variable
export DATABASE_URL="postgresql://postgres@localhost:5432/hobby_project?sslmode=disable"
```

### 3. Run Migrations

```bash
# Apply schema
psql $DATABASE_URL < sql/schema/000001_create_tables.up.sql

# Or use a migration tool
# migrate up
```

### 4. Generate Go Code

```bash
# Generate code from SQL queries
sqlc generate

# Output will be in internal/db/
tree internal/db/
```

Expected output:

```
internal/db/
├── db.go              # Database connection utilities
├── models.go          # Generated Go structs
├── sqlc.yaml         # Configuration copy
└── users.sql.go       # Generated query functions
```

### 5. Use the Generated Code

```go
package main

import (
    "context"
    "fmt"
    "database/sql"

    "github.com/LarsArtmann/SQLC-Wizzard/examples/hobby-project/internal/db"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Open database connection
    db, err := sql.Open("sqlite3", "./database.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Use generated queries
    queries := db.New(db)

    // Create a user
    ctx := context.Background()
    user, err := queries.CreateUser(ctx, db.CreateUserParams{
        Email:    "user@example.com",
        FullName: "John Doe",
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Created user: %s (%s)\n", user.FullName, user.Email)
}
```

## Files Description

### sqlc.yaml

Configuration for SQLC code generation:

```yaml
version: "2"
sql:
  - name: "hobby"
    engine: "postgresql"
    queries: "./sql/queries"
    schema: "./sql/schema"
    gen:
      go:
        package: "hobby"
        out: "./internal/db"
        sql_package: "database/sql"
```

### sql/schema/000001_create_tables.up.sql

Database migration file:

```sql
-- Create users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    full_name TEXT,
    created_at INTEGER DEFAULT (strftime('%s', 'now'))
);

-- Create indexes
CREATE INDEX idx_users_email ON users(email);
```

### sql/queries/users.sql

SQL query definitions with `-- name` comments:

```sql
-- name: CreateUser :one
INSERT INTO users (id, email, full_name)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET email = COALESCE(sqlc.narg('email'), email),
    full_name = COALESCE(sqlc.narg('full_name'), full_name)
WHERE id = ?1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;
```

## Common Operations

### Add a New User

```go
user, err := queries.CreateUser(ctx, db.CreateUserParams{
    Email:    "newuser@example.com",
    FullName: "New User",
})
```

### Get User by ID

```go
user, err := queries.GetUserByID(ctx, "user-123")
```

### List All Users

```go
users, err := queries.ListUsers(ctx)
```

### Update User

```go
updatedUser, err := queries.UpdateUser(ctx, db.UpdateUserParams{
    ID:       "user-123",
    FullName: sql.NullString{String: "Updated Name", Valid: true},
})
```

### Delete User

```go
err := queries.DeleteUser(ctx, "user-123")
```

## Database Support

This example works with multiple databases:

| Database   | SQL Package                    | Notes                         |
| ---------- | ------------------------------ | ----------------------------- |
| SQLite     | database/sql                   | No external database required |
| PostgreSQL | github.com/lib/pq              | Requires PostgreSQL server    |
| MySQL      | github.com/go-sql-driver/mysql | Requires MySQL server         |

## Testing

```bash
# Run the example
go run main.go

# Or run tests
go test ./...
```

## Customization

To customize this example for your project:

1. **Update sqlc.yaml**
   - Change engine (postgresql, mysql, sqlite)
   - Update package name
   - Modify output directory

2. **Add more tables**
   - Create new migration file: sql/schema/000002_create_posts.up.sql
   - Define your table schema

3. **Add more queries**
   - Create new query file: sql/queries/posts.sql
   - Add queries with `-- name` comments

4. **Regenerate code**
   - Run: `sqlc generate`
   - All Go code will be updated automatically

## Next Steps

1. **Build a REST API**
   - Add HTTP handlers
   - Serve the generated code over HTTP

2. **Add validation**
   - Validate input data
   - Handle errors gracefully

3. **Add tests**
   - Unit tests for your application logic
   - Integration tests for database operations

4. **Deploy**
   - Package for distribution
   - Deploy to cloud or self-host

## Resources

- [SQLC Documentation](https://docs.sqlc.dev/)
- [User Guide](../../docs/USER_GUIDE.md)
- [Best Practices](../../docs/BEST_PRACTICES.md)
- [Troubleshooting](../../docs/TROUBLESHOOTING.md)

## License

Same as parent project (MIT)
