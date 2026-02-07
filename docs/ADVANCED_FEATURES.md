# SQLC-Wizzard Advanced Features Guide

## Overview

This guide covers advanced features and techniques for using SQLC-Wizzard in production environments.

## Table of Contents

1. [Custom Code Generation](#custom-code-generation)
2. [CEL Validation Rules](#cel-validation-rules)
3. [Type-Safe Configuration](#type-safe-configuration)
4. [Database-Specific Features](#database-specific-features)
5. [Performance Optimization](#performance-optimization)
6. [Multi-Database Projects](#multi-database-projects)

---

## Custom Code Generation

### Understanding Go Code Generation

SQLC generates Go code based on your SQL queries. You can customize this generation through `sqlc.yaml` configuration.

### Generate Options

```yaml
gen:
  go:
    package: "db" # Go package name
    out: "./internal/db" # Output directory
    sql_package: "database/sql" # SQL driver package
    emit_interface: true # Generate interface
    emit_json_tags: true # Add JSON tags
    emit_prepared_queries: true # Use prepared statements
    emit_empty_slices: false # Handle NULL arrays
    emit_enum_valid_method: true # Generate validation method
    emit_all_enum_values: true # Include all enum values
```

### Custom Type Overrides

Customize how SQL types map to Go types:

```yaml
gen:
  go:
    overrides:
      - db_type: "uuid"
        go_type: "UUID"
        nullable: true
        go_import_path: "github.com/google/uuid"
      - db_type: "jsonb"
        go_type: "RawMessage"
        go_import_path: "encoding/json"
      - db_type: "text"
        go_type: "string"
```

### Custom Rename Rules

Change generated struct field names:

```yaml
gen:
  go:
    rename:
      id: "ID"
      uuid: "UUID"
      url: "URL"
      uri: "URI"
      api: "API"
      http: "HTTP"
      json: "JSON"
      db: "DB"
      created_at: "CreatedAt"
      updated_at: "UpdatedAt"
```

### Template-Based Customization

Use SQLC templates to modify generated code:

```go
// Template: internal/templates/models.go.tpl

{{ range .GoStructs }}
type {{ .Name }} struct {
    {{- range .Fields }}
    {{ .Name }} {{ .Type }} `json:"{{ .Tag }}"`
    {{- end }}
}
{{ end }}
```

### Example: Add Custom Methods

```yaml
gen:
  go:
    # Generate interface for database abstraction
    emit_interface: true
    # Add JSON tags for API serialization
    emit_json_tags: true
    # Use pgx driver for PostgreSQL
    sql_package: "github.com/jackc/pgx/v5"
```

---

## CEL Validation Rules

### What is CEL?

CEL (Common Expression Language) is used for SQLC validation rules. You can define custom validation logic.

### Basic Rule Syntax

```yaml
rules:
  - name: "no-select-star"
    rule: "SELECT \\*"
    message: "Do not use SELECT * - select specific columns instead"
```

### Advanced Rule Patterns

#### Rule 1: Require WHERE Clause

```yaml
rules:
  - name: "require-where-on-update"
    rule: "UPDATE .* WHERE"
    message: "All UPDATE queries must have a WHERE clause"
```

#### Rule 2: Require LIMIT on SELECT

```yaml
rules:
  - name: "require-limit-on-select"
    rule: "SELECT .* FROM .* WHERE .* LIMIT"
    message: "SELECT queries with WHERE must have LIMIT"
```

#### Rule 3: Forbid Specific Tables

```yaml
rules:
  - name: "no-direct-deletes"
    rule: "DELETE FROM .* WHERE"
    message: "Use soft deletes instead of direct deletes"
```

### Complex Rules with Logic

```yaml
rules:
  - name: "safe-insert"
    rule: "INSERT INTO .* \\("
    message: "Use explicit column lists for INSERT statements"

  - name: "no-count-star"
    rule: "COUNT\\(\\*\\)"
    message: "Use COUNT(column_name) instead of COUNT(*)"
```

### Template-Specific Rules

Define rules for specific templates:

```yaml
# For analytics template
rules:
  - name: "require-materialized-view"
    rule: "CREATE MATERIALIZED VIEW"
    message: "Use materialized views for analytics queries"

# For API template
rules:
  - name: "no-joins-in-api-queries"
    rule: "SELECT .* FROM .* INNER JOIN"
    message: "API queries should be simple, avoid complex joins"
```

### Testing Rules

Validate your rules with SQLC:

```bash
# Validate SQL files against rules
sqlc vet

# Validate specific directory
sqlc vet --queries ./sql/queries

# Show detailed rule violations
sqlc vet --verbose
```

---

## Type-Safe Configuration

### Using Generated Types

SQLC-Wizzard generates type-safe configuration using TypeSpec:

```go
// generated/types.go
type DatabaseType string

const (
    DatabaseTypePostgreSQL DatabaseType = "postgresql"
    DatabaseTypeMySQL      DatabaseType = "mysql"
    DatabaseTypeSQLite     DatabaseType = "sqlite"
)

func (dt DatabaseType) IsValid() bool {
    switch dt {
    case DatabaseTypePostgreSQL, DatabaseTypeMySQL, DatabaseTypeSQLite:
        return true
    }
    return false
}
```

### Compile-Time Validation

Benefits of type-safe configuration:

```go
// ✅ Type-safe - Compiler catches errors
func CreateConfig(dbType DatabaseType) error {
    if !dbType.IsValid() {
        return fmt.Errorf("invalid database type: %s", dbType)
    }
    return nil
}

// ❌ Type-unsafe - Runtime errors possible
func CreateConfig(dbType string) error {
    if dbType != "postgresql" && dbType != "mysql" {
        return fmt.Errorf("invalid database type: %s", dbType)
    }
    return nil
}
```

### Using Enums for Safety

```go
// generated/types.go
type ProjectType string

const (
    ProjectTypeHobby        ProjectType = "hobby"
    ProjectTypeMicroservice ProjectType = "microservice"
    ProjectTypeEnterprise   ProjectType = "enterprise"
)

func (pt ProjectType) Features() []string {
    switch pt {
    case ProjectTypeMicroservice:
        return []string{"docker", "kubernetes", "monitoring"}
    case ProjectTypeEnterprise:
        return []string{"audit", "multi-tenant", "backup"}
    default:
        return []string{"basic"}
    }
}
```

### Safe Configuration Loading

```go
// Load configuration with validation
func LoadConfig(path string) (*config.Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    // Validate
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("configuration validation failed: %w", err)
    }

    return &config, nil
}
```

---

## Database-Specific Features

### PostgreSQL Features

```yaml
# Use PostgreSQL-specific features
sql:
  - name: "db"
    engine: "postgresql"
    gen:
      go:
        sql_package: "github.com/jackc/pgx/v5"
        overrides:
          - db_type: "uuid"
            go_type: "UUID"
            go_import_path: "github.com/google/uuid"
          - db_type: "jsonb"
            go_type: "RawMessage"
            go_import_path: "encoding/json"
          - db_type: "timestamptz"
            go_type: "time.Time"
          - db_type: "_text"
            go_type: "[]string"
```

### MySQL Features

```yaml
sql:
  - name: "db"
    engine: "mysql"
    gen:
      go:
        sql_package: "github.com/go-sql-driver/mysql"
        overrides:
          - db_type: "json"
            go_type: "RawMessage"
            go_import_path: "encoding/json"
```

### SQLite Features

```yaml
sql:
  - name: "db"
    engine: "sqlite"
    gen:
      go:
        sql_package: "github.com/mattn/go-sqlite3"
        overrides:
          - db_type: "text"
            go_type: "string"
```

### Database-Specific Query Patterns

```sql
-- PostgreSQL-specific: Use WITH for locking
SELECT * FROM users
WHERE id = ?
FOR UPDATE;

-- MySQL-specific: Use specific index hints
SELECT * FROM users FORCE INDEX (idx_email)
WHERE email = ?;

-- SQLite-specific: Use COLLATE for case-insensitive
SELECT * FROM users
WHERE email = ? COLLATE NOCASE;
```

---

## Performance Optimization

### Query Optimization

#### 1. Use Prepared Queries

```yaml
gen:
  go:
    emit_prepared_queries: true
```

```go
// Prepared query is compiled once, executed many times
stmt, err := db.Prepare(ctx, "SELECT * FROM users WHERE id = ?")
if err != nil {
    return nil, err
}

// Execute multiple times
for _, id := range ids {
    user, err := stmt.QueryRow(ctx, id)
    // ...
}
```

#### 2. Use Indexes

```sql
-- Create appropriate indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);
CREATE INDEX idx_todos_user_status ON todos(user_id, status);

-- Composite index for common queries
CREATE INDEX idx_todos_user_status_created ON todos(user_id, status, created_at DESC);
```

#### 3. Optimize SELECT Statements

```sql
-- ❌ BAD: Fetches all columns
SELECT * FROM users WHERE id = ?;

-- ✅ GOOD: Fetches only needed columns
SELECT id, email, full_name FROM users WHERE id = ?;

-- ❌ BAD: No limit
SELECT * FROM logs ORDER BY created_at;

-- ✅ GOOD: Limits result set
SELECT id, message FROM logs ORDER BY created_at DESC LIMIT 100;
```

### Connection Pooling

```go
// Configure connection pool
config, err := pgxpool.ParseConfig(databaseURL)
if err != nil {
    return nil, err
}

config.MaxConns = 25              // Max connections in pool
config.MinConns = 5               // Min connections to keep
config.MaxConnLifetime = time.Hour   // Connection lifetime
config.MaxConnIdleTime = 5 * time.Minute
config.HealthCheckPeriod = time.Minute

pool, err := pgxpool.NewWithConfig(ctx, config)
```

### Batch Operations

```sql
-- Single insert (slow)
INSERT INTO users (id, email) VALUES (?, ?);
INSERT INTO users (id, email) VALUES (?, ?);

-- Batch insert (fast)
INSERT INTO users (id, email) VALUES
    (?, ?),
    (?, ?);
```

```go
// Use transaction for batch
tx, err := db.Begin(ctx)
if err != nil {
    return err
}

_, err = tx.Exec(ctx, "INSERT INTO users (id, email) VALUES (?, ?), (?, ?)", user1ID, user1Email, user2ID, user2Email)
if err != nil {
    tx.Rollback(ctx)
    return err
}

err = tx.Commit(ctx)
```

### Caching Generated Code

```bash
# Cache generated code to avoid regeneration
# Add .go files to .gitignore
echo "**/internal/db/*.go" >> .gitignore

# Generate once and commit
sqlc generate
git add internal/db/
git commit -m "add generated code"
```

---

## Multi-Database Projects

### Supporting Multiple Databases

```yaml
sql:
  # PostgreSQL configuration
  - name: "postgres"
    engine: "postgresql"
    queries: "./sql/postgres/queries"
    schema: "./sql/postgres/schema"
    gen:
      go:
        package: "postgres"
        out: "./internal/postgres"
        sql_package: "github.com/jackc/pgx/v5"

  # MySQL configuration
  - name: "mysql"
    engine: "mysql"
    queries: "./sql/mysql/queries"
    schema: "./sql/mysql/schema"
    gen:
      go:
        package: "mysql"
        out: "./internal/mysql"
        sql_package: "github.com/go-sql-driver/mysql"

  # SQLite configuration
  - name: "sqlite"
    engine: "sqlite"
    queries: "./sql/sqlite/queries"
    schema: "./sql/sqlite/schema"
    gen:
      go:
        package: "sqlite"
        out: "./internal/sqlite"
        sql_package: "github.com/mattn/go-sqlite3"
```

### Database Abstraction Layer

```go
// Create database interface
type Database interface {
    Querier

    // Database operations
    Begin(ctx context.Context) (Transaction, error)
    Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
}

// Implement for each database
type PostgresDB struct {
    db *pgxpool.Pool
    *db.Queries
}

type MySQLDB struct {
    db *sql.DB
    *mysql.Queries
}

type SQLiteDB struct {
    db *sql.DB
    *sqlite.Queries
}

// Use in application
func NewDB(engine string, connString string) (Database, error) {
    switch engine {
    case "postgresql":
        return NewPostgresDB(connString)
    case "mysql":
        return NewMySQLDB(connString)
    case "sqlite":
        return NewSQLiteDB(connString)
    default:
        return nil, fmt.Errorf("unsupported engine: %s", engine)
    }
}
```

### Migration Strategy for Multiple Databases

```bash
# Run migrations for all databases
for db in postgres mysql sqlite; do
    echo "Migrating $db..."
    DATABASE_URL="${DB_URLS[$db]}" migrate up
done

# Or use a single migration file with conditionals
-- Migration file
-- name: 000001_create_tables
-- Run for all databases
{{ if eq .Engine "postgresql" }}
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ...
);
{{ end }}

{{ if eq .Engine "mysql" }}
CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    ...
);
{{ end }}

{{ if eq .Engine "sqlite" }}
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    ...
);
{{ end }}
```

---

## Best Practices Summary

### ✅ Do

1. **Use type-safe configuration** - Leverage generated types
2. **Define validation rules** - Catch SQL errors early
3. **Optimize queries** - Use indexes and limits
4. **Use prepared queries** - Improve performance
5. **Test thoroughly** - Verify generated code
6. **Document decisions** - Keep track of customizations
7. **Version control** - Track configuration changes

### ❌ Don't

1. **Don't use SELECT \*** - Fetch specific columns
2. **Don't hardcode types** - Use type-safe configuration
3. **Don't ignore validation** - Catch errors early
4. **Don't skip tests** - Verify everything works
5. **Don't forget performance** - Optimize hot paths
6. **Don't mix concerns** - Keep configuration clean

## Troubleshooting

### Issue: "Generated code doesn't compile"

**Solutions:**

- Check type overrides match Go syntax
- Verify import paths are correct
- Ensure `go.mod` has required dependencies

### Issue: "Validation rules not being applied"

**Solutions:**

- Check rule syntax is correct
- Verify rules in `sqlc.yaml`
- Run `sqlc vet --verbose` to debug

### Issue: "Performance is poor"

**Solutions:**

- Add indexes to frequently queried columns
- Use prepared queries
- Check query execution plan
- Optimize joins and subqueries

## Next Steps

- [ ] Implement custom validation rules
- [ ] Set up performance monitoring
- [ ] Configure connection pooling
- [ ] Add database-specific optimizations
- [ ] Create abstraction layer for multiple databases

## Resources

- [SQLC Advanced Features](https://docs.sqlc.dev/reference/)
- [CEL Language Guide](https://github.com/google/cel-spec)
- [PostgreSQL Performance](https://wiki.postgresql.org/wiki/Performance_Optimization)
- [MySQL Optimization](https://dev.mysql.com/doc/refman/8.0/en/optimization.html)

## License

MIT License - See [LICENSE](./LICENSE) for details
