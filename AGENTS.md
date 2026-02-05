# SQLC-Wizard Agent Guide

**Last Updated:** 2025-12-08  
**Project:** Interactive CLI wizard for generating sqlc configurations  
**Primary Language:** Go 1.24.7  
**Architecture:** Domain-Driven Design (DDD) with layered architecture

---

## 🎯 Project Overview

SQLC-Wizard is an interactive CLI tool that generates production-ready `sqlc.yaml` configurations through an intuitive wizard interface. It supports multiple project types (hobby, microservice, enterprise, etc.) and databases (PostgreSQL, MySQL, SQLite) with smart defaults and comprehensive validation.

### Key Features

- Interactive TUI wizard using `charmbracelet/huh`
- Type-safe configuration generation with validation
- Multiple project templates with database-specific optimizations
- CLI commands: `init`, `validate`, `doctor`, `generate`, `migrate`
- Generated Go types from TypeSpec for compile-time safety

---

## 🛠️ Essential Commands

### Build & Development

```bash
# Primary build command (uses justfile)
just build                    # Build to bin/sqlc-wizard with version info

# Development workflow
just dev                      # clean + build + test + find-duplicates

# Installation
just install-local            # Install to GOPATH/bin
```

### Testing

```bash
just test                     # Run all tests with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Integration tests
go test -tags=integration ./...

# Specific test suites
go test ./internal/wizard -v
go test ./pkg/config -v
```

### Code Quality

```bash
just lint                     # Run golangci-lint (if installed)
just fmt                      # Format code with gofmt
just vet                      # Run go vet
just tidy                     # Tidy go modules
just find-duplicates          # Find code duplicates (alias: just fd)
```

### TypeSpec Generation

```bash
just generate-typespec        # Generate Go types from TypeSpec
make types                   # Alternative using Makefile
```

### Verification

```bash
just verify                   # Run build + lint + test (full verification)
```

---

## 🏗️ Project Structure & Architecture

### Core Architecture Principles

- **Domain-Driven Design (DDD)** with clear separation of concerns
- **Layered Architecture**: Domain → Application → Infrastructure
- **Type Safety**: Generated enums prevent invalid states at compile time
- **Template Pattern**: Strategy-based template system

### Directory Structure

```
sqlc-wizard/
├── cmd/sqlc-wizard/          # CLI entrypoint
├── internal/                 # Private application code
│   ├── commands/            # CLI command implementations
│   ├── wizard/              # Interactive TUI wizard logic
│   ├── templates/           # Template system (domain layer)
│   ├── generators/          # File generation (infrastructure)
│   ├── domain/              # Domain models and business logic
│   ├── adapters/            # External interface implementations
│   ├── errors/              # Structured error handling
│   └── testing/             # Test helpers and utilities
├── pkg/                     # Public reusable packages
│   └── config/              # sqlc.yaml configuration types
├── generated/               # TypeSpec-generated Go types
├── templates/              # SQL template files
└── docs/                   # Documentation and status reports
```

---

## 🔧 Development Patterns & Conventions

### Type Safety First

**CRITICAL:** Always use generated types from `./generated` package - never raw strings!

```go
// ✅ CORRECT: Use generated types with validation
projectType, err := templates.NewProjectType("microservice")
if err != nil {
    return err
}

// ❌ WRONG: Never use raw strings for enums
var projectType templates.ProjectType = "microservice" // Unsafe!

// ✅ CORRECT: Use constants for known values
dbType := templates.DatabaseTypePostgreSQL
```

### Error Handling

Use structured error types from `internal/errors`:

```go
// ✅ Use structured errors with context
return errors.ValidationError("project_type", "invalid value")

// ✅ Wrap errors with context
return errors.Wrap(err, errors.ErrorCodeConfigParseFailed, "parser")

// ✅ Use helper functions
return errors.ConfigParseError(path, err)
```

### Testing Patterns

Uses BDD with Ginkgo/Gomega framework:

```go
// Test suite setup
func TestWizardSuite(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Wizard Suite")
}

// BDD-style tests
var _ = Describe("Template Generation", func() {
    Context("with microservice template", func() {
        It("should generate valid sqlc config", func() {
            Expect(config).ToNot(BeNil())
            Expect(config.Version).To(Equal("2"))
        })
    })
})
```

### Template System

All templates implement the `Template` interface:

```go
type Template interface {
    Generate(data generated.TemplateData) (*config.SqlcConfig, error)
    DefaultData() generated.TemplateData
    RequiredFeatures() []string
    Name() string
    Description() string
}
```

---

## 🚨 Critical Gotchas & Non-Obvious Patterns

### TypeSpec Generation Dependency

**CRITICAL:** The project relies on generated Go types from TypeSpec:

1. **First-time setup:** Run `just generate-typespec` to generate types
2. **Generated package:** `./generated` contains all type-safe enums
3. **Module replacement:** `go.mod` has replace directive for local development
4. **Never edit generated files:** They're regenerated from TypeSpec specs

### Domain Layer Constraints

- Domain layer (`internal/templates`, `internal/domain`) can only depend on:
  - `pkg/config`
  - `github.com/samber/lo` and `github.com/samber/mo`
  - Generated types from `./generated`
- **NEVER** import infrastructure or UI packages in domain layer

### Boolean Flag Deprecation

The codebase is migrating from boolean flags to type-safe enums:

```go
// ❌ DEPRECATED: Boolean flags
type SafetyRules struct {
    NoSelectStar bool  // Being replaced with enum
    RequireWhere bool  // Being replaced with enum
}

// ✅ NEW: Type-safe enums (in progress)
type TypeSafeSafetyRules struct {
    SelectStarPolicy SelectStarPolicy    // Enum with validation
    WhereRequirement WhereRequirement   // Enum with validation
}
```

### Dual Configuration System

There are currently two configuration systems (migration in progress):

1. **OLD:** Boolean-heavy `EmitOptions` and `SafetyRules` (deprecated)
2. **NEW:** Type-safe `TypeSafeEmitOptions` and `TypeSafeSafetyRules`

**Always prefer the new type-safe versions for new code.**

---

## 📋 Package-Specific Guidelines

### internal/wizard

- Uses `charmbracelet/huh` for TUI components
- Implements step-by-step wizard flow with validation
- Each wizard step is a separate struct with `Run()` method
- Results collected in `WizardResult` struct

### internal/templates

- Pure domain logic - no I/O or UI dependencies
- Uses strategy pattern for different project types
- Template registry manages template discovery
- All templates implement `Template` interface

### pkg/config

- Represents complete `sqlc.yaml` v2 schema
- Includes all sqlc configuration options
- Supports Go, Kotlin, Python, TypeScript generation
- Provides `ApplyEmitOptions()` helper for type-safe configuration

### internal/commands

- Implements CLI commands using `spf13/cobra`
- Each command in separate file (`init.go`, `validate.go`, etc.)
- Uses dependency injection for adapters and services
- Commands orchestrate domain layer, don't contain business logic

### generated/

- Contains TypeSpec-generated Go types
- Provides compile-time type safety for enums
- **DO NOT EDIT** - regenerate with `just generate-typespec`
- Types include `ProjectType`, `DatabaseType`, `EmitOptions`, etc.

---

## 🧪 Testing Approach

### Test Organization

- Unit tests: `*_test.go` files alongside source
- Integration tests: `integration_test.go` files
- BDD tests: Using Ginkgo framework with `Describe/Context/It`
- Test utilities: `internal/testing` package

### Test Data Management

- Use table-driven tests for multiple scenarios
- Leverage `internal/testing/helpers.go` for common patterns
- Property-based testing for complex validation logic

### Key Test Suites

- `internal/wizard`: Wizard flow and UI interactions
- `internal/templates`: Template generation logic
- `pkg/config`: Configuration parsing and validation
- `internal/generators`: File generation functionality

---

## 🔍 Code Quality Standards

### File Size Limits

- **Max 300 lines per file**
- **Max 50 lines per function**
- **Max 5 parameters per function**
- Violations should be refactored immediately

### Naming Conventions

```go
// Types: PascalCase
type ProjectType struct {}

// Functions: camelCase (exported: PascalCase)
func parseConfig() {}
func NewTemplate() Template {}

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const DefaultOutputDir = "internal/db"
const ERR_CONFIG_INVALID = "CONFIG_INVALID"

// Generated types: Use constants, NO RAW STRINGS!
const ProjectTypeMicroservice = ProjectType(generated.ProjectTypeMicroservice)
```

### Documentation Standards

- Public functions must have Go doc comments
- Complex business logic requires explanatory comments
- Use `DEPRECATED:` comments for legacy code
- Include migration paths for deprecated APIs

---

## 🚀 Build & Deployment

### Version Management

- Version info injected via ldflags during build
- Git-based versioning: `git describe --tags --always --dirty`
- Development builds use "dev" version

### Dependencies

- Go modules managed with `go mod tidy`
- TypeSpec dependencies in `package.json`
- Use `just deps` to download all Go dependencies

### Release Process

```bash
# Build release binary
go build -ldflags "-X main.Version=$(git describe --tags)" -o sqlc-wizard cmd/sqlc-wizard/main.go

# Install to system
go install ./cmd/sqlc-wizard
```

---

## 🐛 Common Issues & Solutions

### TypeSpec Generation Issues

```bash
# If generated types are missing
just generate-typespec

# If module replacement doesn't work
go mod tidy
go clean -modcache
go mod download
```

### Test Failures

- Integration tests require database setup
- Some tests may need specific environment variables
- Use `go test -v` for verbose output to debug issues

### Linter Issues

- Install golangci-lint: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- Run `just lint` to check code quality
- Fix issues before submitting PRs

---

## 🎯 Development Priorities

### Current Migration Focus

1. **Type Safety:** Complete migration from boolean flags to type-safe enums
2. **Test Coverage:** Increase coverage, especially for wizard components
3. **Documentation:** Improve inline documentation and examples
4. **Architecture:** Implement proper dependency injection

### Quality Gates

- All tests must pass (`just test`)
- Code must lint cleanly (`just lint`)
- No new boolean flag usage - use generated enums
- File size limits enforced (< 300 lines)

---

**Generated by Crush Assistant**  
_For agents working on SQLC-Wizard codebase_
