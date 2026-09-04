# SQLC-Wizard User Guide

**Last Updated:** January 13, 2026\
**Version:** 1.0.0

---

## 📚 Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Project Types](#project-types)
4. [Configuration Options](#configuration-options)
5. [Common Use Cases](#common-use-cases)
6. [Troubleshooting](#troubleshooting)

---

## 📦 Installation

SQLC-Wizard can be installed using several methods. Choose the one that best fits your workflow.

### Method 1: Go Install (Recommended)

This is the simplest method if you have Go installed.

```bash
go install github.com/LarsArtmann/SQLC-Wizzard/cmd/sqlc-wizard@latest
```

After installation, verify it works:

```bash
sqlc-wizard version
```

**Output:**

```
sqlc-wizard 1.0.0
  commit: abc123def456
  built:  2026-01-13T12:00:00Z
```

### Method 2: Homebrew (macOS)

If you're on macOS, you can install via Homebrew:

```bash
brew install sqlc-wizard
```

After installation, verify it works:

```bash
sqlc-wizard version
```

### Method 3: Binary Download

Download the appropriate binary for your platform from the [releases page](https://github.com/LarsArtmann/SQLC-Wizzard/releases).

**Supported Platforms:**

- Linux (AMD64, ARM64)
- macOS (AMD64, ARM64)
- Windows (AMD64)

**Download Example (Linux AMD64):**

```bash
curl -L -o sqlc-wizard https://github.com/LarsArtmann/SQLC-Wizzard/releases/download/v1.0.0/sqlc-wizard-linux-amd64
chmod +x sqlc-wizard
./sqlc-wizard version
```

### Method 4: Docker

If you prefer using Docker:

```bash
docker pull ghcr.io/larsartmann/sqlc-wizard:latest
docker run --rm -it -v $(pwd):/workspace ghcr.io/larsartmann/sqlc-wizard:latest
```

**Note:** The `-v $(pwd):/workspace` flag mounts your current directory into the container so the wizard can write files.

### Method 5: Build from Source

If you want to build from source:

```bash
git clone https://github.com/LarsArtmann/SQLC-Wizzard.git
cd SQLC-Wizzard
go build -o sqlc-wizard cmd/sqlc-wizard/main.go
./sqlc-wizard version
```

## Prerequisites

- **Go:** 1.24 or higher (for Go install or building from source)
- **sqlc:** Any version (wizard generates compatible configurations)
- **Database:** PostgreSQL, MySQL, or SQLite (depending on your project)
- **Terminal:** TUI requires terminal with color support

## Verification

After installation, verify that sqlc-wizard is in your PATH:

```bash
which sqlc-wizard
```

**Expected Output (Go install):**

```
/home/yourname/go/bin/sqlc-wizard
```

**Expected Output (Homebrew):**

```
/usr/local/bin/sqlc-wizard
```

If the command is not found, add the appropriate directory to your PATH:

**For Go install:**

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

**For Homebrew (macOS):**

```bash
export PATH=$PATH:/usr/local/bin
```

**To make this permanent**, add the above line to your `~/.bashrc`, `~/.zshrc`, or `~/.profile`.

## Next Steps

Once installed, proceed to the [Quick Start](#quick-start) guide to create your first project.

---

**Need Help?** Check the [Troubleshooting](#troubleshooting) section or [open an issue](https://github.com/LarsArtmann/SQLC-Wizzard/issues).

---

## ⚡ Quick Start

This guide will walk you through creating your first SQLC project using the wizard.

### Step 1: Initialize New Project

Run the wizard to start creating your project:

```bash
sqlc-wizard
```

**Expected Output:**

```
╔═══════════════════════════════════════════════════════════════╗
║                                                                    ║
║   🧙 SQLC-Wizard 1.0.0 - Interactive Configuration Generator   ║
║                                                                    ║
║   Create perfect sqlc configurations in minutes, not hours            ║
║                                                                    ║
╚═════════════════════════════════════════════════════════════════╝

Press Enter to continue...
```

### Step 2: Choose Project Type

The wizard will ask you to select a project type:

```
━━━ Project Type ━━━

What type of project are you creating?

> hobby           Simple personal project with SQLite
  microservice     Microservice with PostgreSQL
  enterprise       Enterprise application with advanced features
  api-first       API-focused project with multi-database support
  analytics       Data analytics pipeline
  testing         Testing/verification framework

↑/↓ navigate  enter to select  ? for help
```

**For Quick Start:** Select `hobby` (press Enter)

### Step 3: Choose Database

Select your database engine:

```
━━━ Database Selection ━━━

Which database are you using?

> postgresql      PostgreSQL (recommended for production)
  mysql           MySQL
  sqlite          SQLite (recommended for hobby/testing)

↑/↓ navigate  enter to select  ? for help
```

**For Quick Start:** Select `sqlite` (press Enter)

### Step 4: Configure Project Details

Enter your project details:

```
━━━ Project Details ━━━

Package name: myproject
Package path: github.com/myorg/myproject
```

**Quick Start Values:**

- **Package name:** `myproject`
- **Package path:** `github.com/myorg/myproject`

### Step 5: Configure Features

Select optional features:

```
━━━ Features ━━━

Enable features (use arrow keys + space to toggle):

>[x] Use UUIDs for primary keys
 [ ] Use JSON columns
 [ ] Use array columns
 [ ] Enable full-text search
 [ ] Use generated queries
 [ ] Use generated schema

↑/↓ navigate  space to toggle  enter to continue
```

**For Quick Start:**

- Keep "Use UUIDs" checked (press space to toggle)
- Leave everything else unchecked
- Press Enter to continue

### Step 6: Configure Output

Configure where files are generated:

```
━━━ Output Configuration ━━━

Where should files be generated?

Base directory: ./internal/db
Queries directory: ./sql/queries
Schema directory: ./sql/schema

↑/↓ navigate  enter to confirm  ? for help
```

**Quick Start Values:**

- Press Enter to accept defaults

### Step 7: Completion

The wizard will generate your configuration:

```
━━━ Configuration Generated ━━━

✅ Generated sqlc.yaml
✅ Generated ./internal/db directory structure
✅ Created sample SQL files
✅ Created Go package skeleton

Configuration saved to: sqlc.yaml

Next steps:
  1. Edit SQL queries in ./sql/queries/
  2. Run: sqlc generate
  3. Import generated code in your Go code
  4. Start building your application!

Press Enter to exit...
```

### Step 8: Verify Generated Configuration

Check the generated `sqlc.yaml`:

```bash
cat sqlc.yaml
```

**Expected Output:**

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "sqlite"
    gen:
      go:
        out: "internal/db"
        sql_package: "db"
        emit_json_tags: true
        emit_prepared_queries: true
```

### Step 9: Generate Code

Use sqlc to generate your Go code:

```bash
sqlc generate
```

**Expected Output:**

```
# package db
...
```

### Step 10: Use Generated Code

Import the generated code in your Go application:

```go
package main

import (
    "database/sql"
    "log"

    "github.com/myorg/myproject/internal/db"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    queries := db.New(db)
    // Use generated queries here...
}
```

## Quick Start Summary

You've now:

1. ✅ Installed SQLC-Wizard
2. ✅ Created a new hobby project with SQLite
3. ✅ Generated sqlc.yaml configuration
4. ✅ Generated Go code using sqlc

## Next Steps

- Add your own SQL queries to `./sql/queries/`
- Customize database schema in `./sql/schema/`
- Regenerate code: `sqlc generate`
- Check [Project Types](#project-types) for more advanced options

**Congratulations!** You're ready to use SQLC-Wizard in your project.

---

**Need more help?** Check out:

- [Project Types](#project-types) - Learn about other templates
- [Configuration Options](#configuration-options) - All available settings
- [Troubleshooting](#troubleshooting) - Common issues and solutions

---

## 🏗 Project Types

SQLC-Wizard provides several project templates optimized for different use cases.

### Hobby Project

**Best For:** Personal projects, prototypes, small applications

**Database:** SQLite (default) or PostgreSQL

**Features:**

- Simple directory structure
- Minimal configuration options
- Fast setup (≤ 2 minutes)
- Low memory footprint

**Generated Structure:**

```
myproject/
├── internal/
│   └── db/
│       ├── db.go           (database wrapper)
│       ├── models.go       (generated models)
│       └── sqlc.yaml      (sqlc configuration)
├── sql/
│   ├── schema/          (database schemas)
│   └── queries/         (SQL queries)
└── go.mod
```

**Use When:**

- Building a personal blog or app
- Creating a prototype or MVP
- Learning SQLC
- Building small tools or utilities

---

### Microservice Project

**Best For:** Microservices, APIs, backend services

**Database:** PostgreSQL (recommended) or MySQL

**Features:**

- API token authentication queries
- Service health check queries
- Transaction support
- Connection pooling configuration
- Prepared query optimization

**Generated Structure:**

```
microservice/
├── internal/
│   ├── db/
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── auth_tokens.sql    (pre-generated)
│   │   └── health_checks.sql (pre-generated)
│   └── api/
│       └── main.go            (HTTP server skeleton)
├── sql/
│   ├── schema/
│   │   ├── auth.sql          (authentication schema)
│   │   └── users.sql         (user management schema)
│   └── queries/
│       └── users.sql         (sample queries)
├── Dockerfile
└── docker-compose.yml       (with PostgreSQL)
```

**Use When:**

- Building microservices architecture
- Creating REST APIs or GraphQL services
- Need service-to-service authentication
- Building cloud-native applications

---

### Enterprise Project

**Best For:** Large-scale applications, multi-tenant systems

**Database:** PostgreSQL or MySQL

**Features:**

- Audit logging tables
- Soft delete queries
- Row-level security support
- Multi-database configurations
- Migration tooling support
- Performance monitoring queries

**Generated Structure:**

```
enterprise/
├── internal/
│   ├── db/
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── audit.sql              (pre-generated)
│   │   ├── soft_delete.sql        (pre-generated)
│   │   └── row_security.sql       (pre-generated)
│   └── config/
│       └── config.go            (configuration loader)
├── sql/
│   ├── schema/
│   │   ├── audit.sql             (audit trail schema)
│   │   ├── users.sql             (user management)
│   │   ├── permissions.sql       (role-based access)
│   │   └── migrations.sql        (migration support)
│   └── queries/
│       └── users.sql             (sample queries)
├── migrations/
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── Dockerfile
├── docker-compose.yml           (multi-database)
└── config.yaml                 (application config)
```

**Use When:**

- Building SaaS applications
- Need audit trails and compliance
- Multi-tenant architecture
- Enterprise data governance
- Complex permission systems

---

### API-First Project

**Best For:** API-focused applications, mobile backends, web services

**Database:** Multi-database support (PostgreSQL, MySQL, SQLite)

**Features:**

- Rate limiting queries
- API key management
- Request/response logging
- Pagination queries
- Response caching support

**Generated Structure:**

```
api-first/
├── internal/
│   ├── db/
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── rate_limits.sql     (pre-generated)
│   │   ├── api_keys.sql       (pre-generated)
│   │   └── cache_support.sql   (pre-generated)
│   └── http/
│       ├── handlers.go        (HTTP request handlers)
│       ├── middleware.go     (auth, rate limiting)
│       └── router.go         (route setup)
├── sql/
│   ├── schema/
│   │   ├── api.sql              (API management)
│   │   ├── users.sql            (user data)
│   │   └── sessions.sql         (session management)
│   └── queries/
│       └── users.sql             (sample queries)
├── Dockerfile
└── swagger.yaml                (API documentation stub)
```

**Use When:**

- Building REST/GraphQL APIs
- Mobile application backends
- Rate limiting required
- API versioning
- Response caching optimization

---

### Analytics Project

**Best For:** Data pipelines, BI tools, reporting systems

**Database:** PostgreSQL (recommended) or MySQL

**Features:**

- Time-series queries
- Aggregation functions
- Window functions support
- Materialized view queries
- Data import/export queries

**Generated Structure:**

```
analytics/
├── internal/
│   ├── db/
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── timeseries.sql     (pre-generated)
│   │   └── aggregations.sql    (pre-generated)
│   └── pipeline/
│       ├── extract.go         (data extraction)
│       ├── transform.go       (data transformation)
│       └── load.go            (data loading)
├── sql/
│   ├── schema/
│   │   ├── events.sql           (event tracking)
│   │   ├── metrics.sql          (metrics storage)
│   │   └── reports.sql          (report queries)
│   └── queries/
│       └── aggregations.sql    (sample queries)
├── Dockerfile
└── airflow_dag.py             (Airflow DAG stub)
```

**Use When:**

- Building data warehouses
- Creating BI dashboards
- Real-time analytics
- ETL/ELT pipelines
- Report generation systems

---

### Testing Project

**Best For:** Test frameworks, QA tools, verification systems

**Database:** SQLite (default) or PostgreSQL

**Features:**

- Test data seeding queries
- Test cleanup queries
- Test assertion queries
- Mock data generation
- Test result storage

**Generated Structure:**

```
testing/
├── internal/
│   ├── db/
│   │   ├── db.go
│   │   ├── models.go
│   │   ├── seed_data.sql      (test data)
│   │   ├── cleanup.sql        (test cleanup)
│   │   └── assertions.sql     (test helpers)
│   └── runner/
│       ├── suite.go           (test suite)
│       ├── setup.go           (test setup)
│       └── teardown.go        (test teardown)
├── sql/
│   ├── schema/
│   │   ├── test_data.sql       (test fixtures)
│   │   └── expected.sql       (expected results)
│   └── queries/
│       └── test_helpers.sql    (sample queries)
└── Dockerfile
```

**Use When:**

- Building integration test frameworks
- Creating QA tools
- Database testing utilities
- Test data management
- Automated testing pipelines

## Choosing the Right Project Type

| Project Type | Complexity           | Features | Database          | Use Case |
| ------------ | -------------------- | -------- | ----------------- | -------- |
| Hobby        | ⭐ Simple            | Basic    | Personal projects |          |
| Microservice | ⭐⭐⭐ Moderate      | Advanced | APIs, Services    |          |
| Enterprise   | ⭐⭐⭐⭐⭐⭐ Complex | Multi-DB | SaaS, Large apps  |          |
| API-First    | ⭐⭐⭐⭐ Advanced    | Multi-DB | REST/GraphQL APIs |          |
| Analytics    | ⭐⭐⭐⭐ Moderate    | Advanced | Data pipelines    |          |
| Testing      | ⭐⭐ Simple          | Basic    | QA, Testing tools |          |

**Recommendation:** Start with **Hobby** or **Microservice** type, then upgrade as needed.

---

**Next:** Learn about [Configuration Options](#configuration-options).

---

## ⚙️ Configuration Options

SQLC-Wizard provides many configuration options to customize your project.

### Output Configuration

Configure where generated files are placed.

**Options:**

- **Base directory:** Root directory for generated code (default: `./internal/db`)
- **Queries directory:** SQL query files (default: `./sql/queries`)
- **Schema directory:** Database schema files (default: `./sql/schema`)

**Best Practices:**

- Use `internal/db` for private database code
- Keep SQL files in `sql/` directory
- Separate queries from schema in subdirectories

**Example:**

```
Output Configuration

Base directory: ./internal/db
Queries directory: ./sql/queries
Schema directory: ./sql/schema
```

### Database Features

Configure database-specific options.

**Available Features:**

- **Use UUIDs:** Generate UUID columns for primary keys (recommended)
- **Use JSON columns:** Support JSONB data types (PostgreSQL)
- **Use arrays:** Support array data types (PostgreSQL)
- **Full-text search:** Enable full-text search queries
- **Generated queries:** Include helper queries in generated code
- **Generated schema:** Include sample schema files

**Feature Dependencies:**

| Feature           | PostgreSQL | MySQL | SQLite |
| ----------------- | ---------- | ----- | ------ |
| UUIDs             | ✅         | ✅    | ✅     |
| JSON columns      | ✅         | ✅    | ❌     |
| Array columns     | ✅         | ❌    | ❌     |
| Full-text search  | ✅         | ✅    | ❌     |
| Generated queries | ✅         | ✅    | ✅     |
| Generated schema  | ✅         | ✅    | ✅     |

**Recommendations:**

- ✅ Always enable "Use UUIDs" (better security and distribution)
- ✅ Enable "JSON columns" for flexible data storage (PostgreSQL/MySQL)
- ✅ Enable "Full-text search" for content-heavy applications
- ⚠️ Avoid arrays unless needed (harder to query)
- ✅ Enable "Generated queries" for faster development
- ✅ Enable "Generated schema" for quick start

### Project Package Configuration

Configure Go package details.

**Options:**

- **Package name:** Name of generated Go package (default: `myproject`)
- **Package path:** Full Go module path (default: `github.com/myorg/myproject`)

**Best Practices:**

- Use lowercase package names (e.g., `db`, `models`)
- Use full module paths for packages (e.g., `github.com/org/project/internal/db`)
- Match package name to directory structure

**Example:**

```
Project Details

Package name: db
Package path: github.com/myorg/myproject/internal/db
```

### sqlc Configuration Options

Wizard generates optimized sqlc.yaml with these options:

#### Go Options

```yaml
gen:
  go:
    out: "internal/db" # Output directory
    sql_package: "db" # Package name
    emit_json_tags: true # JSON struct tags
    emit_prepared_queries: true # Prepared queries
    emit_interface: true # Generate interfaces
```

**Explanation:**

- `emit_json_tags`: Add `json:` struct tags for API responses
- `emit_prepared_queries`: Use prepared statements (better performance)
- `emit_interface`: Generate DB interface (easier mocking in tests)

#### Database Options

```yaml
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "sqlite" # or postgresql, mysql
```

**Supported Engines:**

- `postgresql`: PostgreSQL 12+ (recommended for production)
- `mysql`: MySQL 8.0+ (good compatibility)
- `sqlite`: SQLite 3.35+ (good for development/testing)

### Advanced Options

For advanced users, you can manually edit `sqlc.yaml` after generation.

**Common Manual Adjustments:**

```yaml
# Add additional output languages
gen:
  go:
    out: "internal/db"
  typescript:
    out: "web/src/db"
    package: "db"

# Override sqlc rules
rules:
  - engine: "postgresql"
    schema: "sql/schema"
    queries: "sql/queries"

# Add strict mode
strict_generate: true
```

**When to Manually Edit:**

- Need TypeScript output (in addition to Go)
- Require custom sqlc rules
- Want strict type checking
- Need specific sqlc overrides

---

## 🐛 Troubleshooting

This section covers common issues and their solutions.

### Installation Issues

#### "command not found: sqlc-wizard"

**Problem:** `sqlc-wizard` command not found after installation.

**Solutions:**

1. **Verify Installation**

   ```bash
   which sqlc-wizard
   ```

   **Expected Output:**
   - Go install: `/home/yourname/go/bin/sqlc-wizard`
   - Homebrew: `/usr/local/bin/sqlc-wizard`

2. **Add to PATH**

   For Go install:

   ```bash
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
   ```

   For Homebrew (macOS):

   ```bash
   echo 'export PATH=$PATH:/usr/local/bin' >> ~/.zshrc
   source ~/.zshrc
   ```

3. **Reopen Terminal**
   - Close and reopen your terminal for PATH changes to take effect

---

#### "permission denied" when running binary

**Problem:** Binary doesn't have execute permission.

**Solution:**

```bash
chmod +x sqlc-wizard
```

---

#### "invalid checksum" when downloading binary

**Problem:** Downloaded file is corrupted or incomplete.

**Solutions:**

1. **Redownload binary**

   ```bash
   rm sqlc-wizard
   curl -L -o sqlc-wizard https://github.com/LarsArtmann/SQLC-Wizzard/releases/download/v1.0.0/sqlc-wizard-linux-amd64
   ```

2. **Verify checksum** (if provided)
   ```bash
   sha256sum sqlc-wizard
   ```
   Compare with checksum from release notes.

---

### Wizard Issues

#### "TUI: terminal not supported"

**Problem:** Terminal doesn't support TUI features.

**Solutions:**

1. **Use SSH with proper terminal support**

   ```bash
   ssh -t user@host sqlc-wizard
   ```

   Note the `-t` flag for pseudo-terminal allocation.

2. **Use CI/CD mode** (non-interactive)
   - For automation, use pre-generated configurations
   - See [CI/CD Examples](../guides/ci-cd.md) for details

3. **Use Docker with proper TTY**
   ```bash
   docker run --rm -it -v $(pwd):/workspace ghcr.io/larsartmann/sqlc-wizard:latest
   ```
   Note the `-it` flags for interactive TTY.

---

#### "connection refused" when connecting to database

**Problem:** Database not running or wrong port.

**Solutions:**

1. **Check if database is running**

   ```bash
   # PostgreSQL
   pg_isready -h localhost -p 5432

   # MySQL
   mysqladmin -h localhost -p 3306 ping

   # SQLite (no daemon needed)
   ls -la *.db
   ```

2. **Check database connection settings**

   ```bash
   # PostgreSQL
   cat ~/.pgpass
   host:localhost:5432:dbname:user:password

   # MySQL
   cat ~/.my.cnf
   [client]
   host = localhost
   port = 3306
   ```

3. **Start database** (if not running)

   ```bash
   # PostgreSQL (Docker)
   docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres

   # MySQL (Docker)
   docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql

   # SQLite (no daemon needed)
   # Just ensure .db file exists
   ```

---

#### "syntax error" in sqlc.yaml

**Problem:** Generated configuration has syntax error.

**Solutions:**

1. **Validate YAML syntax**

   ```bash
   # Install yamllint (if not installed)
   pip install yamllint

   # Validate sqlc.yaml
   yamllint sqlc.yaml
   ```

2. **Check indentation** (YAML is indentation-sensitive)

   ```yaml
   # Correct (2 spaces)
   version: "2"
   sql:
     - schema: "sql/schema"

   # Incorrect (tabs)
   version: "2"
   sql:
     - schema: "sql/schema"
   ```

3. **Regenerate configuration**
   ```bash
   rm sqlc.yaml
   sqlc-wizard
   ```

---

### Code Generation Issues

#### "no queries found" when running `sqlc generate`

**Problem:** SQL query files not in expected location.

**Solutions:**

1. **Check query directory structure**

   ```bash
   # Should match sqlc.yaml
   ls -la sql/queries/
   ```

2. **Verify sqlc.yaml configuration**

   ```yaml
   sql:
     - schema: "sql/schema"
       queries: "sql/queries" # Check this path
   ```

3. **Create sample query file**

   ```bash
   mkdir -p sql/queries
   cat > sql/queries/users.sql <<EOF
   -- name: GetUser
   SELECT * FROM users WHERE id = ?;
   EOF
   ```

4. **Run sqlc generate again**
   ```bash
   sqlc generate
   ```

---

#### "type mismatch" in generated code

**Problem:** SQL types don't match Go types expected by sqlc.

**Solutions:**

1. **Check SQL column types**

   ```sql
   -- Use integer for id
   CREATE TABLE users (
     id INTEGER PRIMARY KEY,
     name TEXT NOT NULL
   );

   -- Or use uuid extension (PostgreSQL)
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

   CREATE TABLE users (
     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
     name TEXT NOT NULL
   );
   ```

2. **Check sqlc overrides**

   ```yaml
   # Override column types
   override:
     column_type:
       "id": "uuid.UUID"
   ```

3. **Regenerate with wizard** (if using wrong database engine)
   ```bash
   rm sqlc.yaml
   sqlc-wizard  # Select correct database
   sqlc generate
   ```

---

### Performance Issues

#### "slow code generation" with many tables

**Problem:** Large schemas take time to generate code.

**Solutions:**

1. **Use caching** (sqlc feature)

   ```bash
   sqlc generate --cache
   ```

2. **Split schema files** (if possible)

   ```bash
   # Instead of one huge schema.sql
   mkdir -p sql/schema
   # Split into users.sql, orders.sql, products.sql
   ```

3. **Increase resources** (if in CI/CD)
   ```yaml
   # CI configuration (example)
   resources:
     limits:
       cpus: "2"
       memory: "4Gi"
   ```

---

### Migration Issues

#### "migration failed" after upgrading sqlc

**Problem:** sqlc version incompatibility or breaking changes.

**Solutions:**

1. **Check sqlc version**

   ```bash
   sqlc version
   ```

2. **Review breaking changes** in sqlc release notes
   - Visit: https://docs.sqlc.dev/
   - Check release notes for your version

3. **Regenerate with wizard**

   ```bash
   # Backup current configuration
   cp sqlc.yaml sqlc.yaml.backup

   # Regenerate
   rm sqlc.yaml
   sqlc-wizard

   # Review changes
   diff sqlc.yaml.backup sqlc.yaml
   ```

4. **Adjust for new features** (if needed)
   - Check if new sqlc features available
   - Update configuration manually if needed

---

### Getting More Help

If you're still stuck:

1. **Check Documentation**
   - [User Guide](.) - This guide
   - [Project Types](#project-types) - Template details
   - [Configuration Options](#configuration-options) - All settings

2. **Check Examples**
   - [Hobby Example](../examples/hobby-sqlite/)
   - [Microservice Example](../examples/microservice-pg/)
   - [Enterprise Example](../examples/enterprise-multi/)

3. **Open an Issue**
   - [Report Bug](https://github.com/LarsArtmann/SQLC-Wizzard/issues/new?template=bug_report.md)
   - [Request Feature](https://github.com/LarsArtmann/SQLC-Wizzard/issues/new?template=feature_request.md)
   - [Ask Question](https://github.com/LarsArtmann/SQLC-Wizzard/discussions)

4. **Check Community**
   - [GitHub Discussions](https://github.com/LarsArtmann/SQLC-Wizzard/discussions)
   - [Stack Overflow](https://stackoverflow.com/questions/tagged/sqlc)
   - [Reddit r/golang](https://reddit.com/r/golang)

---

**Need more advanced troubleshooting?** Check [Migration Guide](../guides/migration.md) for common migration scenarios.
