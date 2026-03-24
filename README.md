# SQLC-Wizzard

An interactive CLI wizard for generating production-ready `sqlc v2` configurations with smart defaults.

SQLC-Wizzard guides you through creating type-safe SQL code by generating optimized `sqlc.yaml` configurations. Answer a few questions about your project, and get a complete setup with example queries, schema templates, and best-practice safety rules.

## Features

### Interactive Wizard

Start with `sqlc-wizard init` and answer guided questions:

```
$ sqlc-wizard init
? Project Type › Microservice
? Primary Database › PostgreSQL
? Go Package Path › github.com/user/my-api
✓ Created sqlc.yaml
✓ Created schema/ and queries/ directories
✓ Generated example files
```

### Project Templates

| Template     | Description                              |
| ------------ | ---------------------------------------- |
| Hobby        | Simple setup for personal projects       |
| Microservice | Single database, container-optimized     |
| Enterprise   | Multi-database, comprehensive validation |
| API-First    | JSON-focused, REST-friendly              |
| Analytics    | Read-heavy, complex queries              |
| Testing      | In-memory, mock-friendly                 |
| Multi-tenant | Schema-per-tenant patterns               |
| Library      | Embeddable, minimal dependencies         |

### Database Support

- PostgreSQL
- MySQL
- SQLite

### Commands

| Command    | Description                                |
| ---------- | ------------------------------------------ |
| `init`     | Interactive wizard to create configuration |
| `validate` | Validate existing sqlc.yaml                |
| `generate` | Generate example SQL files                 |
| `doctor`   | Check development environment              |
| `migrate`  | Manage configuration migrations            |

## Installation

### Go Install

```bash
go install github.com/LarsArtmann/SQLC-Wizzard@latest
```

### Build from Source

```bash
git clone https://github.com/LarsArtmann/SQLC-Wizzard.git
cd SQLC-Wizzard
go install ./cmd/sqlc-wizard
```

## Usage

### Interactive Mode

```bash
sqlc-wizard init
```

### Non-Interactive Mode

```bash
sqlc-wizard init --non-interactive \
  --project-type=microservice \
  --database=postgresql \
  --package=github.com/user/myapi
```

### Validate Configuration

```bash
sqlc-wizard validate
sqlc-wizard validate --strict
```

### Environment Check

```bash
sqlc-wizard doctor
```

### Generate Example Files

```bash
sqlc-wizard generate --output ./db
```

## Project Structure

```
sqlc-wizzard/
├── cmd/sqlc-wizard/     # CLI entrypoint
├── internal/
│   ├── commands/        # CLI command implementations
│   ├── wizard/          # Interactive TUI wizard
│   ├── templates/       # Project templates
│   ├── generators/      # File generation
│   ├── domain/          # Domain models
│   ├── adapters/        # External interfaces
│   ├── validation/      # Configuration validation
│   └── creators/        # Project creation
├── pkg/
│   └── config/          # sqlc.yaml types
├── generated/           # TypeSpec-generated types
└── templates/           # SQL templates
```

## Configuration

Generated configurations include sensible defaults:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/db/queries"
    schema: "internal/db/schema"
    strict_function_checks: true
    gen:
      go:
        package: "db"
        out: "internal/db"
        emit_json_tags: true
        emit_interface: true
```

## Development

### Prerequisites

- Go 1.21+
- sqlc (for testing generated code)

### Setup

```bash
go mod tidy
go test ./...
just build
```

### Testing

```bash
go test ./...
go test ./internal/wizard -v
go test ./pkg/config -v
```

## Roadmap

### Completed

- Interactive wizard with 8 project templates
- Configuration validation
- Environment health checks
- PostgreSQL, MySQL, SQLite support

### In Progress

- Configuration migration tools
- Enhanced validation rules
- Additional template customization

### Planned

- Web-based configuration generator
- IDE integrations
- Cloud provider templates

## License

MIT License
