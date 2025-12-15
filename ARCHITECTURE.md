# SQLC-Wizard Architecture

**Last Updated:** 2025-11-06
**Version:** 0.1.0-dev (MVP)

---

## 📐 Architecture Overview

SQLC-Wizard follows **Domain-Driven Design (DDD)** principles with clear separation between:

- **Domain logic** (business rules, templates, validation)
- **Application logic** (commands, workflows)
- **Infrastructure** (file I/O, YAML parsing, TUI)

### Architectural Patterns Used

1. ✅ **Domain-Driven Design (DDD)** - Clear domain model
2. ✅ **Layered Architecture** - Separation of concerns
3. ✅ **Strategy Pattern** - Templates as strategies
4. ✅ **Registry Pattern** - Template discovery
5. ⚠️ **Railway-Oriented Programming** - Partial (needs `mo.Result`)
6. ⚠️ **Hexagonal Architecture** - Partial (needs ports/adapters)
7. ❌ **CQRS** - Not yet (commands exist, queries don't)
8. ❌ **Event-Driven** - Not yet (no events)

---

## 📦 Package Structure

```
sqlc-wizard/
├── cmd/
│   └── sqlc-wizard/
│       └── main.go                 # CLI entrypoint
│
├── internal/                       # Private application code
│   ├── commands/                   # CLI commands (application layer)
│   │   ├── init.go                # Init command
│   │   └── validate.go            # Validate command
│   │
│   ├── wizard/                     # Interactive wizard (UI layer)
│   │   └── wizard.go              # Wizard orchestration
│   │
│   ├── templates/                  # Template system (domain layer)
│   │   ├── types.go               # Domain types (ProjectType, Features, etc.)
│   │   ├── registry.go            # Template registry
│   │   └── microservice.go        # Microservice template strategy
│   │
│   ├── generators/                 # File generators (infrastructure)
│   │   ├── generator.go           # File generation
│   │   └── embedded_templates.go  # Embedded SQL templates
│   │
│   ├── errors/                     # Error handling (⚠️ EMPTY - TODO)
│   └── detectors/                  # Project detection (⚠️ EMPTY - TODO)
│
├── pkg/                            # Public reusable packages
│   ├── config/                    # sqlc.yaml config (domain model)
│   │   ├── types.go              # Config types
│   │   ├── parser.go             # YAML parsing
│   │   ├── validator.go          # Validation
│   │   └── marshaller.go         # YAML writing
│   │
│   ├── database/                  # Database-specific (⚠️ EMPTY - TODO)
│   └── version/                   # Version management (⚠️ EMPTY - TODO)
│
└── templates/                      # SQL template files
    ├── queries/
    │   ├── postgresql/
    │   ├── sqlite/
    │   └── mysql/
    └── schema/
        ├── postgresql/
        ├── sqlite/
        └── mysql/
```

---

## 🎯 Dependency Rules (go-arch-lint)

```yaml
# .go-arch-lint.yml (TODO: Implement)
version: 1

allow:
  # Domain layer (most restrictive)
  - from: internal/templates
    to:
      - pkg/config
      - github.com/samber/lo
      - github.com/samber/mo

  # Application layer
  - from: internal/commands
    to:
      - internal/wizard
      - internal/templates
      - internal/generators
      - pkg/config

  # Infrastructure layer (least restrictive)
  - from: internal/generators
    to:
      - internal/templates
      - pkg/config
      - os
      - io

  # Wizard can depend on everything
  - from: internal/wizard
    to:
      - internal/templates
      - pkg/config
      - github.com/charmbracelet/huh

deny:
  # Domain cannot depend on infrastructure
  - from: internal/templates
    to:
      - internal/generators
      - internal/wizard

  # Config cannot depend on application
  - from: pkg/config
    to:
      - internal/*
```

---

## 🔄 Data Flow

### Init Command Flow

```
User Input (CLI/TUI)
  ↓
Commands Layer (init.go)
  ↓
Wizard Layer (wizard.go) ← Interactive prompts
  ↓
Templates Layer (microservice.go) ← Business logic
  ↓
Config Layer (types.go) ← Domain model
  ↓
Generator Layer (generator.go) ← File I/O
  ↓
File System (sqlc.yaml, queries, schema)
```

### Validate Command Flow

```
File System (sqlc.yaml)
  ↓
Config Layer (parser.go) ← Parse YAML
  ↓
Config Layer (validator.go) ← Validate
  ↓
Commands Layer (validate.go) ← Format output
  ↓
User Output (colored terminal)
```

---

## 🏛️ Design Principles

### 1. Make Illegal States Unrepresentable

**Bad (Current):**

```go
type ProjectType string // Any string is valid!
```

**Good (TODO):**

```go
type ProjectType struct {
    value string
}

func NewProjectType(s string) (ProjectType, error) {
    if !isValid(s) {
        return ProjectType{}, ErrInvalidProjectType
    }
    return ProjectType{value: s}, nil
}
```

### 2. Single Source of Truth

**Bad (Current - Split Brain):**

```go
// Two representations of same concept!
type SafetyRules struct { NoSelectStar bool }
type RuleConfig struct { Name string, Rule string }
```

**Good (TODO):**

```go
type SafetyRule interface {
    ToRuleConfig() RuleConfig
}
```

### 3. Composition Over Inheritance

**Good (Current):**

```go
type Template interface {
    Generate(data TemplateData) (*SqlcConfig, error)
}

type MicroserviceTemplate struct{}
func (t *MicroserviceTemplate) Generate(...) { }
```

### 4. Dependency Inversion

**Current:** Direct dependencies everywhere

**TODO:** Use interfaces

```go
type ConfigWriter interface {
    Write(cfg *SqlcConfig, path string) error
}

type FileConfigWriter struct{}
func (w *FileConfigWriter) Write(...) error { }
```

---

## 🧪 Testing Strategy

### Test Pyramid (TODO)

```
      /\
     /  \    E2E Tests (10%)
    /____\
   /      \   Integration Tests (20%)
  /________\
 /          \  Unit Tests (70%)
/____________\
```

**Current:** 0% coverage ❌

### Test Categories

1. **Unit Tests** (70%)
   - Config parsing/validation
   - Template generation
   - Type conversions
   - Error handling

2. **Integration Tests** (20%)
   - Full wizard flow
   - File generation
   - Command execution

3. **E2E Tests** (10%)
   - CLI invocation
   - Real file system
   - Actual sqlc validation

### BDD with Ginkgo (TODO)

```go
var _ = Describe("Microservice Template", func() {
    Context("when generating PostgreSQL config", func() {
        It("should include UUID support", func() {
            data := TemplateData{Database: DatabaseTypePostgreSQL}
            cfg, err := template.Generate(data)

            Expect(err).ToNot(HaveOccurred())
            Expect(cfg.SQL[0].Gen.Go.Overrides).To(ContainElement(
                Override{DBType: "uuid", GoType: "UUID"},
            ))
        })
    })
})
```

---

## 🔐 Type Safety Guidelines

### ✅ DO

1. Use strong types for domain concepts
2. Validate at construction time (smart constructors)
3. Use custom marshaling for complex types
4. Make invalid states unrepresentable

### ❌ DON'T

1. Use `interface{}` (never!)
2. Use string/int for enums (use types!)
3. Allow default zero values to be valid
4. Skip validation

---

## 🎨 Code Style

### File Size Limits

- **Max 300 lines per file**
- **Max 50 lines per function**
- **Max 5 parameters per function**

**Current violations:**

- `internal/generators/embedded_templates.go` - 270 lines ⚠️
- `internal/wizard/wizard.go` - 290 lines ⚠️

### Naming Conventions

```go
// Types: PascalCase
type ProjectType struct { }

// Functions: camelCase (exported: PascalCase)
func parseConfig() { }
func NewTemplate() Template { }

// Constants: PascalCase or SCREAMING_SNAKE_CASE
const DefaultOutputDir = "internal/db"
const ERR_CONFIG_INVALID = "CONFIG_INVALID"

// Variables: camelCase
var projectName string
```

---

## 🚀 Future Architecture

### Planned Improvements

1. **Hexagonal Architecture**
   - Define ports (interfaces)
   - Implement adapters (file, HTTP, etc.)

2. **Event-Driven Architecture**
   - Emit events for extensibility
   - Configuration via events

3. **CQRS**
   - Separate read/write models
   - Query layer for config inspection

4. **Railway-Oriented Programming**
   - Use `mo.Result[T, error]` everywhere
   - Chain operations with `FlatMap`

---

## 📚 References

- [Domain-Driven Design](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Railway-Oriented Programming](https://fsharpforfunandprofit.com/rop/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

## 🤔 Architectural Decision Records (ADRs)

### ADR-001: Why Cobra instead of urfave/cli?

**Decision:** Use `spf13/cobra` for CLI framework

**Rationale:**

- Industry standard (used by kubectl, gh, docker)
- Excellent documentation
- Automatic completion generation
- Nested command support

**Status:** ✅ Accepted

### ADR-002: Why charmbracelet/huh instead of survey?

**Decision:** Use `charmbracelet/huh` for TUI

**Rationale:**

- Beautiful out-of-the-box styling
- Modern, actively maintained
- Consistent with charmbracelet ecosystem (lipgloss, bubbletea)
- Better form validation support

**Status:** ✅ Accepted

### ADR-003: Why embed templates as Go constants? (TODO: Revisit)

**Decision:** Embed SQL templates as Go string constants

**Rationale:**

- No runtime file dependencies
- Single binary distribution
- Fast access (no I/O)

**Concerns:**

- Should use `//go:embed` instead for real files
- Easier to edit as actual .sql files
- Better syntax highlighting

**Status:** ⚠️ Needs Review

---

**End of Architecture Document**
