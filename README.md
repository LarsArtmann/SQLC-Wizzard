# 🧙‍♂️ SQLC-Wizard

> An interactive CLI wizard that generates perfect sqlc configurations in minutes, not hours.

SQLC-Wizard makes type-safe SQL accessible to everyone by providing an intuitive wizard that guides developers through creating production-ready sqlc setups with smart defaults and comprehensive validation.

## ✨ Features

### 🎯 Interactive Wizard
```
$ sqlc-wizard init
🧙‍♂️ Sqlc Configuration Wizard
Let's create the perfect sqlc setup for your project!

? Project Type › Microservice / CLI / Web API / Library / Enterprise
? Primary Database › PostgreSQL / SQLite / MySQL / Multi-DB
? Go Package Path › github.com/user/my-awesome-api

✓ Created sqlc.yaml
✓ Created .github/workflows/sqlc.yml
✓ Generated example queries
✓ Added database migration template
```

### 🔧 Smart Project Detection
- **Auto-detect** existing database files and schemas
- **Analyze** Go imports to determine SQL package preference  
- **Scan** for existing sqlc configurations and offer upgrades
- **Detect** project structure (microservices, monolith, library)
- **Identify** database-specific features (UUIDs, JSON, FTS)

### 📋 Template Categories
| Template | Description | Best For |
|----------|-------------|-----------|
| 🏠 **Hobby** | Simple SQLite setup | Personal projects, prototypes |
| ⚡ **Microservice** | Single DB, container-optimized | API services, microservices |
| 🏢 **Enterprise** | Multi-DB, comprehensive validation | Large applications, teams |
| 🔧 **API-First** | JSON-focused, REST-friendly | REST/GraphQL backends |
| 📊 **Analytics** | Read-heavy, complex queries | Data platforms, reporting |
| 🧪 **Testing** | In-memory, mock-friendly | Test suites, CI/CD |
| 🌐 **Multi-tenant** | Schema-per-tenant patterns | SaaS applications |
| 📦 **Library** | Embeddable, minimal deps | Go libraries, packages |

### 🚀 Advanced Features

#### Configuration Validation
```bash
$ sqlc-wizard validate
✓ Configuration is valid
⚠️  Warning: Consider enabling emit_prepared_queries for better performance
ℹ️  Suggestion: Add validation rule for no-select-star
🔧 Fix available: sqlc-wizard validate --fix
```

#### Database Migration Assistant
```bash
$ sqlc-wizard migrate --from sqlite --to postgresql
✓ Generated migration scripts
✓ Updated type overrides
✓ Converted SQLite-specific features to PostgreSQL equivalents
```

#### Best Practices Assistant
```bash
$ sqlc-wizard doctor
🩺 Database Health Check
✓ Using prepared queries (performance)
✓ Validation rules enabled (safety)
⚠️  Missing indexes for foreign keys
💡 Suggestion: Add rule require-limit-on-select for large tables
```

## 🚀 Quick Start

### Installation

#### Go Install (Recommended)
```bash
go install github.com/sqlc-wizard/sqlc-wizard@latest
```

#### Build from Source
```bash
git clone https://github.com/sqlc-wizard/sqlc-wizard.git
cd sqlc-wizard
go build -o sqlc-wizard cmd/sqlc-wizard/main.go
```

#### Package Managers
```bash
# Homebrew (coming soon)
brew install sqlc-wizard

# Docker (coming soon)
docker run --rm -v $(pwd):/app sqlc-wizard/sqlc-wizard init
```

### Basic Usage

#### 1. Interactive Wizard (Most Common)
```bash
sqlc-wizard init
```

#### 2. Non-Interactive Mode
```bash
sqlc-wizard generate \
  --project-type=microservice \
  --database=postgresql \
  --package=github.com/user/myapi \
  --output-dir=internal/db
```

#### 3. Validate Existing Configuration
```bash
sqlc-wizard validate
sqlc-wizard validate --fix  # Auto-fix issues
```

#### 4. Health Check
```bash
sqlc-wizard doctor
```

## 📖 Command Reference

### `sqlc-wizard init`
Interactive wizard to create new sqlc configurations.

```bash
sqlc-wizard init [flags]

Flags:
  --project-type string     Project template (hobby, microservice, enterprise, api-first, analytics, testing, multi-tenant, library)
  --database string          Database engine (sqlite, postgresql, mysql, multi)
  --package string          Go package path (e.g., github.com/user/project)
  --output-dir string       Output directory for generated code
  --queries-dir string      SQL queries directory
  --schema-dir string       Database schema directory
  --non-interactive         Skip prompts, use flags only
```

### `sqlc-wizard generate`
Generate configuration without interaction.

```bash
sqlc-wizard generate [flags]

Flags:
  --template string         Template name or path to custom template
  --config string           Output configuration file (default: sqlc.yaml)
  --database string          Database engine
  --features strings        Database features (fts5, uuid, json, arrays)
  --languages strings       Target languages (go, python, typescript, kotlin)
  --safety strings          Safety features (validation, no-select-star, require-where)
```

### `sqlc-wizard validate`
Validate sqlc configuration files.

```bash
sqlc-wizard validate [file] [flags]

Flags:
  --fix                     Auto-fix common issues
  --strict                  Enable strict validation mode
  --format string           Output format (text, json, yaml)
```

### `sqlc-wizard doctor`
Diagnose common issues and suggest improvements.

```bash
sqlc-wizard doctor [flags]

Flags:
  --check-performance       Check for performance issues
  --check-security          Check for security vulnerabilities
  --check-best-practices    Check for best practices violations
```

### `sqlc-wizard migrate`
Upgrade sqlc configurations between versions.

```bash
sqlc-wizard migrate [flags]

Flags:
  --from string             Source version
  --to string               Target version
  --backup                  Create backup before migration
```

## 🛠️ Template System

### Template Structure
```
templates/
├── hobby/
│   ├── sqlc.yaml.template
│   ├── queries/
│   ├── schema/
│   └── README.md
├── microservice/
│   ├── sqlc.yaml.template
│   ├── docker-compose.yml
│   ├── .github/workflows/
│   └── README.md
└── enterprise/
    ├── sqlc.yaml.template
    ├── monitoring/
    ├── migrations/
    └── README.md
```

### Custom Templates
Create your own templates:

```bash
mkdir ~/.sqlc-wizard/templates/my-template
cat > ~/.sqlc-wizard/templates/my-template/sqlc.yaml.template << 'EOF'
version: "2"
sql:
  - engine: "{{ .Database }}"
    queries: "{{ .QueriesDir }}"
    schema: "{{ .SchemaDir }}"
    gen:
      go:
        package: "{{ .PackageName }}"
        out: "{{ .OutputDir }}"
        emit_json_tags: true
        emit_interface: true
EOF
```

## 🏗️ Project Structure

```
sqlc-wizard/
├── cmd/
│   └── sqlc-wizard/
│       └── main.go              # CLI entrypoint
├── internal/
│   ├── wizard/                  # Interactive wizard logic
│   │   ├── wizard.go           # Main wizard implementation
│   │   ├── steps.go            # Wizard step definitions
│   │   └── ui.go               # TUI components
│   ├── templates/               # Built-in templates
│   │   ├── hobby.go           # Hobby project template
│   │   ├── microservice.go    # Microservice template
│   │   ├── enterprise.go      # Enterprise template
│   │   └── loader.go          # Template loader
│   ├── validators/              # Config validation
│   │   ├── sqlc.go            # sqlc.yaml validator
│   │   ├── database.go        # Database connection validator
│   │   └── best_practices.go  # Best practices checker
│   ├── detectors/               # Project analysis
│   │   ├── project.go         # Project type detection
│   │   ├── database.go        # Database detection
│   │   └── dependencies.go     # Dependency analysis
│   ├── generators/              # Code generation
│   │   ├── sqlc.go            # sqlc.yaml generator
│   │   ├── queries.go         # Example SQL queries
│   │   ├── workflows.go       # GitHub Actions workflows
│   │   └── migrations.go      # Migration templates

├── pkg/
│   ├── config/                  # Config file handling
│   │   ├── sqlc.go            # sqlc.yaml parser
│   │   ├── loader.go          # Config loader
│   │   └── merger.go          # Config merger
│   ├── database/                # DB-specific logic
│   │   ├── sqlite.go          # SQLite-specific features
│   │   ├── postgresql.go      # PostgreSQL-specific features
│   │   ├── mysql.go           # MySQL-specific features
│   │   └── features.go         # Feature detection
│   └── version/                 # Version management
│       ├── parser.go          # Version parser
│       └── migrator.go        # Version migrator
└── templates/
    ├── sqlc/                    # sqlc.yaml templates
    ├── queries/                 # Example SQL queries
    ├── workflows/               # GitHub Actions
    ├── migrations/              # Migration templates
    └── docs/                    # Documentation templates
```

## 🔧 Configuration

### Global Configuration
```bash
# ~/.sqlc-wizard/config.yaml
default_database: postgresql
default_template: microservice
author_name: "Your Name"
author_email: "your.email@example.com"

templates:
  custom_dir: "~/.sqlc-wizard/templates"
  auto_update: true
```

### Environment Variables
```bash
export SQLC_WIZARD_CONFIG_HOME="~/.sqlc-wizard"
export SQLC_WIZARD_TEMPLATE_DIR="~/.sqlc-wizard/templates"
export SQLC_WIZARD_CACHE_DIR="~/.sqlc-wizard/cache"
```

## 🧪 Development

### Prerequisites
- Go 1.21+
- sqlc (for testing)
- Docker (optional, for database testing)

### Setup
```bash
git clone https://github.com/sqlc-wizard/sqlc-wizard.git
cd sqlc-wizard

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Build
go build -o bin/sqlc-wizard cmd/sqlc-wizard/main.go

# Install locally
go install ./cmd/sqlc-wizard
```

### Testing
```bash
# Run all tests
go test ./...

# Run integration tests
go test -tags=integration ./...

# Test with different databases
docker-compose up -d postgresql mysql sqlite
go test -tags=integration ./internal/detectors/...
```

### Contributing
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-wizard-step`
3. Make your changes and add tests
4. Run tests: `go test ./...`
5. Submit a pull request

## 📚 Examples

### Example 1: New Microservice
```bash
$ sqlc-wizard init
🧙‍♂️ Sqlc Configuration Wizard

? Project Type › Microservice
? Primary Database › PostgreSQL  
? Project Name › user-service
? Go Package Path › github.com/company/user-service
? Database Features › ✓ UUIDs ✓ JSON columns
? Safety Features › ✓ Validation rules ✓ No SELECT *
? Output Directory › internal/db

✓ Created sqlc.yaml
✓ Created .github/workflows/sqlc.yml
✓ Generated example queries
✓ Added migration template
```

Generated `sqlc.yaml`:
```yaml
version: "2"
sql:
  - name: "user_service"
    engine: "postgresql"
    queries: "internal/db/queries"
    schema: "internal/db/schema"
    strict_function_checks: true
    strict_order_by: true
    database:
      uri: "${DATABASE_URL}"
      managed: true
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        build_tags: "postgres,pgx"
        emit_json_tags: true
        emit_db_tags: true
        emit_prepared_queries: true
        emit_interface: true
        emit_exact_table_names: true
        emit_empty_slices: true
        json_tags_case_style: "camel"
        omit_unused_structs: true
```

### Example 2: Configuration Upgrade
```bash
$ sqlc-wizard migrate --from v1 --to v2
✓ Backing up sqlc.yaml to sqlc.yaml.backup
✓ Migrating configuration from v1 to v2
✓ Migration completed successfully
```

### Example 3: Health Check
```bash
$ sqlc-wizard doctor
🩺 Database Health Check

✓ Configuration is valid sqlc v2
✓ Using prepared queries for performance
✓ Interface enabled for testability
⚠️  No validation rules configured
💡 Suggestion: Add no-select-star rule for security
⚠️  Missing indexes on foreign key columns
💡 Suggestion: Add performance monitoring
✓ Using appropriate Go types for database columns
```

## 🔌 Integrations

### IDE Extensions
- **VS Code**: Auto-completion, validation, and wizard UI
- **GoLand**: Integration with database tools
- **Vim/Neovim**: LSP integration for sqlc.yaml

### CI/CD Integration
```yaml
# .github/workflows/sqlc-wizard.yml
name: SQLC Validation

on: [push, pull_request]

jobs:
  sqlc-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: sqlc-wizard/setup-action@v1
      - name: Validate sqlc configuration
        run: sqlc-wizard validate --strict
      - name: Run sqlc checks
        run: sqlc-wizard doctor --check-security --check-performance
```

### Framework Integrations
- **Gin**: Optimized templates for REST APIs
- **Echo**: Echo-specific query patterns  
- **Chi**: Chi router integration
- **Fiber**: Fiber-optimized configurations

## 🤝 Community

- **GitHub**: https://github.com/sqlc-wizard/sqlc-wizard
- **Discord**: Join our Discord community
- **Discussions**: GitHub Discussions for questions and ideas
- **Twitter**: Follow @sqlcwizard for updates

## 📈 Roadmap

### Phase 1: Core Wizard ✅
- [x] Basic interactive wizard
- [x] Essential templates (Hobby, Microservice, Enterprise)
- [x] SQLite and PostgreSQL support
- [x] Configuration validation

### Phase 2: Advanced Features (In Progress)
- [ ] MySQL support
- [x] Multi-database configurations
- [ ] Configuration upgrade/migration
- [x] Doctor/diagnostics system

### Phase 3: Ecosystem Integration (Planned)
- [ ] IDE extensions (VS Code, GoLand)
- [ ] Web-based configuration generator
- [ ] Framework-specific templates (Gin, Echo, Chi)
- [ ] Cloud provider templates (AWS RDS, GCP CloudSQL, Azure)

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **sqlc Team**: For building the amazing sqlc tool
- **GoReleaser**: Inspiration for the wizard CLI pattern
- **Charm**: For the excellent TUI components (bubbletea, lipgloss)
- **Community**: For feedback, contributions, and feature requests

---

## 🧙‍♂️ Make sqlc configuration magical!

*Generated with ❤️ by the SQLC-Wizard team*