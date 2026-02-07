# SQLC-Wizzard Troubleshooting Guide

## Quick Diagnosis

```bash
# Run doctor to check system health
wizard doctor

# Verbose mode for debugging
wizard init --verbose

# Check logs
tail -f ~/.local/state/wizard/logs/*.log
```

---

## Common Issues & Solutions

### 1. Installation Issues

#### Issue: "go install: package not found"

**Cause:** Wrong module path or network issues.

**Solutions:**

```bash
# Check GOPATH and GOMODCACHE
echo $GOPATH
echo $GOMODCACHE

# Try with explicit version
go install github.com/LarsArtmann/SQLC-Wizzard/cmd/wizard@latest

# Verify installation
which wizard
wizard version
```

#### Issue: "permission denied" when running wizard

**Cause:** Executable not in PATH or permissions issue.

**Solutions:**

```bash
# Check if wizard is in PATH
which wizard

# Add to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin

# Fix permissions
chmod +x $(go env GOPATH)/bin/wizard
```

---

### 2. SQLC Configuration Issues

#### Issue: "sqlc: no such file or directory"

**Cause:** Incorrect file paths in sqlc.yaml.

**Solutions:**

```yaml
# ❌ BAD: Relative paths without ./sql
version: "2"
sql:
  - name: "db"
    queries: "queries"           # Won't work
    schema: "schema"             # Won't work

# ✅ GOOD: Use relative or absolute paths
version: "2"
sql:
  - name: "db"
    queries: "./sql/queries"      # Works!
    schema: "./sql/schema"        # Works!
```

#### Issue: "sqlc: path must be absolute or start with ./"

**Cause:** Relative path without prefix.

**Solutions:**

```bash
# Check current directory
pwd

# Update sqlc.yaml with correct paths
# Option 1: Absolute paths
queries: "/home/user/project/sql/queries"

# Option 2: Relative paths with ./
queries: "./sql/queries"

# Option 3: Run from project root
cd /path/to/project
sqlc generate  # Uses ./ as current dir
```

#### Issue: "sqlc: package name not found"

**Cause:** Go module path mismatch.

**Solutions:**

```bash
# Check your go.mod
cat go.mod

# Should match package in sqlc.yaml
module github.com/username/project

# sqlc.yaml
gen:
  go:
    package: "project"  # Matches module name
```

#### Issue: "sqlc: no queries found"

**Cause:** Missing `-- name` comments in SQL files.

**Solutions:**

```sql
-- ❌ BAD: No query name
SELECT * FROM users WHERE id = $1;

-- ✅ GOOD: Query name with comment
-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- ✅ GOOD: Alternative syntax
/* name: GetUserByID */
SELECT * FROM users WHERE id = $1;
```

---

### 3. Database Connection Issues

#### Issue: "connection refused" or "timeout"

**Cause:** Database not running or wrong port.

**Solutions:**

```bash
# Check if PostgreSQL is running
ps aux | grep postgres

# Check port
lsof -i :5432

# Start PostgreSQL if needed
brew services start postgresql  # macOS
systemctl start postgresql   # Linux
docker start postgres-container # Docker

# Test connection
psql -U postgres -h localhost -p 5432 -c "SELECT 1"
```

#### Issue: "password authentication failed"

**Cause:** Wrong credentials or missing password.

**Solutions:**

```bash
# Test with psql
psql "postgresql://user:pass@localhost:5432/dbname"

# Set password in environment
export PGPASSWORD="your-password"
psql -U postgres -h localhost -d dbname

# Update connection string
postgresql://user:password@localhost:5432/dbname
```

#### Issue: "database does not exist"

**Cause:** Database not created.

**Solutions:**

```bash
# List databases
psql -U postgres -c "\l"

# Create database
psql -U postgres -c "CREATE DATABASE mydb;"

# Or use template1
psql -U postgres -c "CREATE DATABASE mydb TEMPLATE template1;"
```

#### Issue: "SSL connection error"

**Cause:** SSL configuration mismatch.

**Solutions:**

```yaml
# ❌ BAD: SSL enabled on local DB
database:
  url: "postgresql://user:pass@localhost:5432/db?sslmode=require"

# ✅ GOOD: Disable SSL for local development
database:
  url: "postgresql://user:pass@localhost:5432/db?sslmode=disable"

# ✅ GOOD: Enable SSL for production
database:
  url: "postgresql://user:pass@prod-db.example.com:5432/db?sslmode=require"
```

---

### 4. Code Generation Issues

#### Issue: "sqlc: package not found: github.com/lib/pq"

**Cause:** Missing Go dependencies.

**Solutions:**

```bash
# Add dependency
go get github.com/lib/pq

# Or use pgx (recommended)
go get github.com/jackc/pgx/v5

# Update sqlc.yaml
gen:
  go:
    sql_package: "pgx/v5"  # Updated
```

#### Issue: "sqlc: no Go files generated"

**Cause:** No SQL files or no query names.

**Solutions:**

```bash
# Check SQL files exist
ls -la sql/queries/

# Check for query names
grep -n "name:" sql/queries/*.sql

# Regenerate with verbose output
sqlc generate --verbose

# Check output directory
ls -la internal/db/
```

#### Issue: "sqlc: duplicate column name: id"

**Cause:** SQL query returns duplicate columns.

**Solutions:**

```sql
-- ❌ BAD: Both tables have id column
SELECT t.*, u.*
FROM todos t
JOIN users u ON t.user_id = u.id

-- ✅ GOOD: Explicit column selection or alias
SELECT
    t.id AS todo_id,
    t.title,
    u.id AS user_id,
    u.email
FROM todos t
JOIN users u ON t.user_id = u.id

-- ✅ GOOD: Use table alias
SELECT
    t.*,
    u.email
FROM todos t
JOIN users u ON t.user_id = u.id
```

---

### 5. Runtime Issues

#### Issue: "no such table: users"

**Cause:** Tables not created or migrations not run.

**Solutions:**

```bash
# Run migrations
migrate up

# Check database schema
psql -d dbname -c "\d users"

# Verify migration files exist
ls -la sql/schema/*.up.sql

# Manually run schema
psql -d dbname -f sql/schema/000001_create_tables.up.sql
```

#### Issue: "sql: no such function: gen_random_uuid()"

**Cause:** PostgreSQL extension not enabled.

**Solutions:**

```sql
-- Enable uuid-ossp extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Now gen_random_uuid() will work
INSERT INTO users (id, email)
VALUES (gen_random_uuid(), 'test@example.com');
```

#### Issue: "panic: sql: Scan error on column index"

**Cause:** Type mismatch between SQL and Go struct.

**Solutions:**

```sql
-- ❌ BAD: TEXT to INT
SELECT CAST(description AS INTEGER) AS count
FROM todos

-- ✅ GOOD: Compatible types
SELECT COUNT(*) AS count
FROM todos
```

```go
// Check Go struct matches SQL
type Todo struct {
    ID        string    `db:"id"`        // Must match SQL type
    Title     string    `db:"title"`
    CreatedAt time.Time `db:"created_at"`
}
```

---

### 6. Template Issues

#### Issue: "template not found: hobby"

**Cause:** Invalid project type or template not registered.

**Solutions:**

```bash
# List available templates
wizard template list

# Verify correct type
wizard init --type hobby

# Available types:
# - hobby
# - microservice
# - enterprise
# - api-first
# - analytics
# - testing
# - multi-tenant
# - library
```

#### Issue: "template generation failed"

**Cause:** Invalid template data or configuration error.

**Solutions:**

```bash
# Validate template data
wizard validate --file sqlc.yaml

# Check template logs
cat ~/.local/state/wizard/logs/template_*.log

# Try with different output
wizard init --type microservice --output ./custom-output
```

---

### 7. Performance Issues

#### Issue: "Slow query execution"

**Cause:** Missing indexes or inefficient queries.

**Solutions:**

```sql
-- Check query plan
EXPLAIN ANALYZE
SELECT * FROM todos WHERE user_id = $1;

-- Add index if missing
CREATE INDEX idx_todos_user_id ON todos(user_id);
CREATE INDEX idx_todos_status_created ON todos(status, created_at DESC);

-- Use LIMIT to prevent large result sets
SELECT * FROM todos
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 100;
```

#### Issue: "Connection pool exhausted"

**Cause:** Too many connections, not closing properly.

**Solutions:**

```go
// Increase pool size
config.MaxConns = 50

// Always close connections
defer db.Close()

// Use context for timeouts
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := db.QueryContext(ctx, query, args...)
```

---

### 8. Environment-Specific Issues

#### macOS: "command not found: wizard"

**Solution:**

```bash
# Add to shell profile
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Verify
which wizard
```

#### Linux: "permission denied"

**Solution:**

```bash
# Install with sudo (not recommended)
sudo go install github.com/LarsArtmann/SQLC-Wizzard/cmd/wizard@latest

# Better: Fix permissions
sudo chown $USER:$(go env GOPATH)/bin/wizard
```

#### Windows: "'wizard' is not recognized"

**Solution:**

```powershell
# Add to PATH
$env:Path = $env:Path + ";C:\Users\$env:USERNAME\go\bin"

# Verify
wizard version

# Or add via System Environment Variables
# Settings > Environment Variables > Path
```

---

## Advanced Troubleshooting

### Enable Debug Logging

```bash
# Set environment variable
export WIZARD_DEBUG=true
export WIZARD_LOG_LEVEL=debug

# Run wizard
wizard init --verbose

# Check logs
tail -f ~/.local/state/wizard/logs/debug.log
```

### Test SQLC Directly

```bash
# Generate code without wizard
sqlc generate

# Validate configuration
sqlc validate

# Check generated files
tree internal/db/

# Test compilation
go build ./internal/db/...
```

### Database Health Checks

```bash
# PostgreSQL
psql -U postgres -c "SELECT version();"
psql -U postgres -c "SELECT datname FROM pg_database WHERE datistemplate = false;"

# MySQL
mysql -u root -e "SELECT VERSION();"
mysql -u root -e "SHOW DATABASES;"

# SQLite
sqlite3 database.db "SELECT sqlite_version();"
sqlite3 database.db ".tables"
```

### Network Diagnostics

```bash
# Check DNS resolution
nslookup database-host.example.com

# Check connectivity
nc -zv database-host.example.com 5432

# Test with telnet
telnet database-host.example.com 5432

# Trace route (if slow)
traceroute database-host.example.com
```

---

## Getting Help

### Community Support

- [GitHub Issues](https://github.com/LarsArtmann/SQLC-Wizzard/issues) - Bug reports and feature requests
- [GitHub Discussions](https://github.com/LarsArtmann/SQLC-Wizzard/discussions) - Q&A and discussions
- [SQLC Documentation](https://docs.sqlc.dev/) - Official SQLC docs

### Report Issues

When reporting issues, include:

1. **Wizard version:** `wizard version`
2. **Go version:** `go version`
3. **OS/Architecture:** `uname -a`
4. **Steps to reproduce:** What you did
5. **Expected behavior:** What should happen
6. **Actual behavior:** What happened
7. **Error message:** Full error output
8. **Configuration:** sqlc.yaml (sanitized)
9. **Logs:** `~/.local/state/wizard/logs/*.log`

### Template Issue

```markdown
## Issue Description

Brief description of the issue.

## Steps to Reproduce

1. Run `wizard init --type hobby`
2. Select PostgreSQL
3. Enter project name
4. ...

## Expected Behavior

What should happen?

## Actual Behavior

What actually happens?

## Environment

- Wizard version: 1.0.0
- Go version: 1.21.0
- OS: macOS 14.0
- Database: PostgreSQL 15

## Logs
```

Paste relevant log output here

```

---

## Summary

Most issues can be resolved by:

1. ✅ **Checking file paths** - Ensure sqlc.yaml paths are correct
2. ✅ **Verifying SQL files** - Check for query name comments
3. ✅ **Testing database connection** - Verify database is running
4. ✅ **Running sqlc directly** - Isolate SQLC issues
5. ✅ **Enabling debug logging** - Get more information
6. ✅ **Checking environment** - Verify PATH and dependencies

### Quick Reference

| Issue | Quick Fix |
| ------ | ---------- |
| Module not found | `go mod tidy` |
| SQLC not found | `go install sqlc.dev/sqlc/cmd/sqlc@latest` |
| Connection refused | Start database service |
| No queries found | Add `-- name` comments to SQL |
| Type mismatch | Check SQL type vs Go struct type |
| Permission denied | `chmod +x wizard` |

### Additional Resources

- [User Guide](./USER_GUIDE.md)
- [Best Practices](./BEST_PRACTICES.md)
- [SQLC Docs](https://docs.sqlc.dev/)
- [PostgreSQL Docs](https://www.postgresql.org/docs/)
- [Go Database/SQL Package](https://pkg.go.dev/database/sql)
```
