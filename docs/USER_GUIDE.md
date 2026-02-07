# SQLC-Wizzard User Guide

## 📚 Table of Contents

1. [Introduction](#introduction)
2. [Quick Start](#quick-start)
3. [Installation](#installation)
4. [Getting Started](#getting-started)
5. [Configuration](#configuration)
6. [Advanced Features](#advanced-features)
7. [Troubleshooting](#troubleshooting)

## Introduction

SQLC-Wizzard is an interactive CLI tool that simplifies setting up [sqlc](https://sqlc.dev/) projects. It generates production-ready configurations with best practices built-in, saving you hours of setup time.

### Key Features

- 🧙 **Interactive Wizard** - Easy-to-use terminal UI
- 🎨 **Multiple Templates** - Hobby, Microservice, Enterprise, API-First, Analytics, Testing, Multi-Tenant, Library
- 🗄️ **Database Support** - PostgreSQL, MySQL, SQLite
- ✅ **Type-Safe** - Built with TypeSpec for compile-time safety
- 🔧 **Highly Configurable** - Customizable code generation, validation rules, and more

## Quick Start

```bash
# Install SQLC-Wizzard
go install github.com/LarsArtmann/SQLC-Wizzard/cmd/wizard@latest

# Run the wizard
wizard

# Follow the prompts to configure your project
```

## Installation

### Prerequisites

- Go 1.21 or higher
- sqlc (optional, can be installed via wizard)

### Installation Methods

#### Method 1: Go Install (Recommended)

```bash
go install github.com/LarsArtmann/SQLC-Wizzard/cmd/wizard@latest
```

#### Method 2: Build from Source

```bash
git clone https://github.com/LarsArtmann/SQLC-Wizzard.git
cd SQLC-Wizzard
make build
./bin/wizard
```

#### Method 3: Docker

```bash
docker pull ghcr.io/larsartmann/sqlc-wizzard:latest
docker run -it --rm -v $(pwd):/workspace ghcr.io/larsartmann/sqlc-wizzard:latest
```

## Getting Started

### Step 1: Initialize Your Project

Run the wizard and select a template:

```bash
$ wizard
🧙‍♂️  SQLC Configuration Wizard

Let's create a perfect sqlc setup for your project!

📍 Project Type Selection

What type of project are you building?
> 🏠 Hobby - Simple SQLite setup
  ⚡ Microservice - Single DB, container-optimized
  🏢 Enterprise - Multi-DB, comprehensive
  🔧 API-First - JSON-focused, REST-friendly
  📊 Analytics - Read-optimized, warehousing
  🧪 Testing - Isolated, disposable
  🏗️  Multi-Tenant - Shared resources
  📦 Library - Embeddable, minimal deps
```

### Step 2: Choose Database

```bash
📍 Database Selection

Which database will you use?
> 🐘 PostgreSQL - Full-featured, recommended
  🗄️  SQLite - Lightweight, embedded
  🐬 MySQL - Popular, widely supported
```

### Step 3: Configure Features

The wizard will guide you through:

- Project name and package configuration
- Output directories (code, queries, schema)
- Code generation options (interfaces, JSON tags, etc.)
- Safety rules (no SELECT \*, require WHERE, etc.)
- Database features (UUIDs, JSON, arrays, etc.)

### Step 4: Review and Generate

```bash
🎉 Configuration Complete

Project: my-awesome-project
Package: myproject
Type: microservice
Database: postgresql
Output: ./internal/db

Features:
- Interfaces: true
- Prepared Queries: true
- JSON Tags: true

Safety Rules:
- No SELECT *: true
- Require WHERE: true
- Require LIMIT: true

Database Features:
- UUIDs: true
- JSON: true
- Arrays: false
- Full-text: false
```

## Configuration

### sqlc.yaml Structure

The wizard generates a `sqlc.yaml` file in your project root:

```yaml
version: "2"
sql:
  - name: "myproject"
    engine: "postgresql"
    queries: "./sql/queries"
    schema: "./sql/schema"
    gen:
      go:
        package: "myproject"
        out: "./internal/db"
        sql_package: "pgx/v5"
        emit_interface: true
        emit_json_tags: true
```

### Template-Specific Configurations

Each template has optimized defaults:

| Template     | Database   | Code Gen       | Safety   | Use Case               |
| ------------ | ---------- | -------------- | -------- | ---------------------- |
| Hobby        | SQLite     | Basic          | Relaxed  | Side projects, MVPs    |
| Microservice | PostgreSQL | Advanced       | Standard | Containerized services |
| Enterprise   | Multi-DB   | Full           | Strict   | Production systems     |
| API-First    | PostgreSQL | JSON-optimized | Standard | REST APIs              |
| Analytics    | PostgreSQL | Read-optimized | Relaxed  | Data warehouses        |
| Testing      | SQLite     | Isolated       | Relaxed  | Test suites            |
| Multi-Tenant | PostgreSQL | Advanced       | Strict   | SaaS platforms         |
| Library      | Any        | Minimal        | Standard | Shared libraries       |

## Advanced Features

### Custom Code Generation

The wizard supports customizing code generation:

```go
// In your generated code
type DBTX interface {
    Tx(ctx context.Context) *sql.Tx
}

// The wizard can generate:
// - Interfaces for database abstraction
// - Prepared queries for performance
// - JSON tags for API responses
// - Empty slices for consistency
// - Method receivers for OOP patterns
```

### Safety Rules

Enforce SQL best practices:

- **No SELECT \*** - Prevent fetching all columns
- **Require WHERE** - Ensure filtered queries
- **Require LIMIT** - Prevent unlimited result sets
- **Custom CEL Rules** - Add your own validation

### Type-Safe Configuration

All configurations are type-checked at compile time:

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

## Troubleshooting

### Common Issues

#### Issue: "sqlc not found"

**Solution:**

```bash
# Check if sqlc is installed
which sqlc

# Install sqlc
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installation
sqlc version
```

#### Issue: "Package path not found"

**Solution:**
Ensure your Go module path matches your project:

```bash
go mod init github.com/your-username/your-repo
```

#### Issue: "Database connection failed"

**Solution:**

```bash
# Check database is running
psql -U postgres -c "SELECT 1"

# Test connection string
postgresql://user:password@localhost:5432/dbname?sslmode=disable
```

#### Issue: "No queries generated"

**Solution:**
Ensure your query files follow naming conventions:

```sql
-- File: sql/queries/users.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;
```

### Getting Help

```bash
# Check doctor command
wizard doctor

# Verbose mode for debugging
wizard init --verbose

# View logs
tail -f ~/.local/state/wizard/logs/*.log
```

## Next Steps

- See [Tutorial](./tutorial.md) for step-by-step walkthrough
- Check [Advanced Features](./advanced-features.md) for detailed configurations
- Review [Best Practices](./best-practices.md) for production tips
- Visit [GitHub Issues](https://github.com/LarsArtmann/SQLC-Wizzard/issues) for community support

## Contributing

We welcome contributions! See [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

## License

MIT License - See [LICENSE](./LICENSE) for details.
