# Hobby SQLite Example

Simple personal project using SQLC-Wizard with SQLite.

## Project Type

- **Type:** Hobby
- **Database:** SQLite
- **Package:** hobby-sqlite
- **Path:** github.com/example/hobby-sqlite

## Features

- ✅ Use UUIDs for primary keys
- ❌ JSON columns (disabled for SQLite)
- ❌ Array columns (disabled for SQLite)
- ❌ Full-text search (disabled for SQLite)

## Project Structure

```
hobby-sqlite/
├── internal/
│   └── db/
│       ├── db.go
│       ├── models.go
│       └── sqlc.yaml
├── sql/
│   ├── schema/
│   │   └── schema.sql
│   └── queries/
│       └── users.sql
├── go.mod
├── go.sum
└── README.md
```

## Setup

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Generate Code

```bash
sqlc generate
```

### 3. Use Generated Code

```go
package main

import (
	"database/sql"
	"log"

	"github.com/example/hobby-sqlite/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	queries := db.New(database)
	// Use generated queries...
}
```

## Example Queries

### Create User

```sql
-- name: CreateUser
INSERT INTO users (id, name, email, created_at)
VALUES (?, ?, ?, ?);
```

### Get User by ID

```sql
-- name: GetUser
SELECT * FROM users
WHERE id = ?;
```

### List Users

```sql
-- name: ListUsers
SELECT * FROM users
ORDER BY created_at DESC;
```

## Database Schema

```sql
CREATE TABLE users (
	id UUID PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Advantages of SQLite for Hobby Projects

- **Zero Configuration:** No database server needed
- **Single File:** Database stored in one file
- **Fast:** Excellent performance for small datasets
- **Portable:** Easy to copy/move database file
- **Cross-Platform:** Works on all operating systems

## When to Upgrade

Consider upgrading to **PostgreSQL** when:
- Need concurrent writes
- Require advanced features (JSON, arrays, full-text search)
- Building production application
- Need better performance at scale

## Next Steps

1. Add your own SQL queries to `sql/queries/`
2. Customize database schema in `sql/schema/`
3. Run `sqlc generate` after changes
4. Import and use generated code in your application

## Resources

- [SQLC Documentation](https://docs.sqlc.dev/)
- [Go SQLite Driver](https://github.com/mattn/go-sqlite3)
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [SQLC-Wizard User Guide](../../docs/user-guide/)
