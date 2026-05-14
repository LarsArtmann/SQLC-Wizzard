# Library Integration Report — SQLC-Wizard

> How well does SQLC-Wizard leverage the LarsArtmann ecosystem libraries?

**Generated:** 2026-05-13
**Project:** SQLC-Wizard v1.x (Go 1.26.2)
**Reference:** `/home/lars/projects/docs/LIBRARY_GUIDE.md`

---

## Executive Summary

SQLC-Wizard uses **0 out of 19** ecosystem libraries. Every layer — CLI, errors, validation, output, config — is hand-rolled with varying quality. The project carries significant custom code that duplicates the purpose, interfaces, and sometimes the exact API surface of existing ecosystem libraries. The highest-impact integrations would be **cmdguard**, **go-error-family**, **go-output**, and **go-business-rules**.

---

## Integration Matrix

| Library                          | Status   | Should Use?                                                 | Effort | Impact | Priority |
| -------------------------------- | -------- | ----------------------------------------------------------- | ------ | ------ | -------- |
| **cmdguard**                     | Not used | **Yes** — replaces spf13/cobra                              | High   | High   | P1       |
| **go-error-family**              | Not used | **Yes** — replaces `internal/apperrors`                     | Medium | High   | P1       |
| **go-output**                    | Not used | **Yes** — replaces hand-rolled lipgloss output              | Medium | Medium | P2       |
| **go-business-rules**            | Not used | **Yes** — replaces `pkg/config/validator`                   | Medium | Medium | P2       |
| **smart-configs**                | Not used | **Maybe** — for `DATABASE_URL`, sqlc config discovery       | Low    | Medium | P3       |
| **gogenfilter**                  | Not used | **Maybe** — if analyzing existing sqlc-generated code       | Low    | Low    | P3       |
| **go-branded-id**                | Not used | No — no multi-entity ID mixing in this tool                 | —      | —      | —        |
| **go-cqrs-lite**                 | Not used | No — no event sourcing needed                               | —      | —      | —        |
| **go-localfirst**                | Not used | No — no offline/sync requirements                           | —      | —      | —        |
| **go-localsync**                 | Not used | No — no external API sync needed                            | —      | —      | —        |
| **go-finding**                   | Not used | **Maybe** — if building a static analysis pipeline          | Low    | Low    | P4       |
| **cqrs-htmx**                    | Not used | No — no web/HTMX layer                                      | —      | —      | —        |
| **templ-components**             | Not used | No — no web UI                                              | —      | —      | —        |
| **go-commit**                    | Not used | No — no git commit generation in this tool                  | —      | —      | —        |
| **universal-workflow**           | Not used | **Maybe** — for multi-step wizard orchestration             | High   | Low    | P4       |
| **ActaFlow**                     | Not used | No — no actor model needed                                  | —      | —      | —        |
| **go-filewatcher**               | Not used | **Maybe** — for `sqlc-wizard generate --watch`              | Low    | Low    | P4       |
| **project-discovery-sdk**        | Not used | **Maybe** — for auto-detecting project type from filesystem | Medium | Medium | P3       |
| **go-composable-business-types** | Not used | No — no bitemporal/audit trail requirements                 | —      | —      | —        |

---

## Detailed Analysis

### P1: cmdguard — CLI Framework

**Current state:** Uses `spf13/cobra` with manual flag wiring in every command file.

```
cmd/sqlc-wizard/main.go          → cobra.Command with fmt.Fprintf(os.Stderr) error handling
internal/commands/init.go        → 13-flag struct with manual cmd.Flags().StringVar
internal/commands/validate.go    → cobra + lipgloss for output
internal/commands/generate.go    → cobra with manual flag parsing
internal/commands/doctor.go      → cobra with hand-rolled DoctorCheck/DoctorResult types
internal/commands/create.go      → cobra with 7 flags
internal/commands/migrate.go     → cobra with 6 flags + subcommands
```

**What cmdguard provides that cobra doesn't:**

- Type-safe flags from structs with struct tags (no more `cmd.Flags().StringVar`)
- Dependency injection via `Scope` (database, logger, HTTP client managed per-command)
- Environment variable fallbacks with prefix support (12-factor)
- 12 output formats built in via go-output integration
- Graceful shutdown, health checks, and lifecycle hooks
- Eliminates cobra panics and stringly-typed flags

**What SQLC-Wizard would gain:**

- Replace `InitOptions`, `ValidateOptions`, `GenerateOptions`, `CreateOptions` structs with typed `Command[T, F]` generics
- Replace `RunE: func(cmd *cobra.Command, args []string) error` with typed command handlers
- Replace manual `os.Exit(1)` in `main()` with cmdguard's lifecycle management
- Replace manual error formatting (`fmt.Fprintf(os.Stderr, "Error: %v\n", err)`) with structured error handling

**Recommendation:** Migrate to cmdguard v2. The current cobra usage is straightforward enough for a clean migration. Start with `init` command as proof of concept.

---

### P1: go-error-family — Error Handling

**Current state:** Two parallel error systems:

1. **`internal/apperrors`** (primary, ~500 lines): Custom `Error` struct with `ErrorCode`, `ErrorSeverity`, `ErrorDetails`, `Retryable`, `RequestID`, `UserID`, `Component`, `Timestamp`. Includes `ErrorList` with filtering/grouping. Has `Wrap`, `WrapWithRequestID`, `NewInternal`, `FileNotFoundError`, `ConfigParseError`, etc.

2. **`pkg/errors`** (legacy, ~15 lines): Minimal `BaseError` with `Message` + `Code`. Likely deprecated.

**What go-error-family provides that `apperrors` doesn't:**

- **Behavioral classification** via `Family` (Rejection, Conflict, Transient, Corruption, Infrastructure) — not just severity
- **BSD sysexits.h exit codes** — deterministic `ExitCode(err)` for CLI tools (currently missing)
- **Wix-quality user-facing messages** — auto-generated "What / Why / Fix" presentation at CLI boundary
- **Deterministic diagnostics** — auto-discover PostgreSQL, filesystem, network, git root causes
- **Protocol-based** — share the vocabulary without importing the implementation

**Overlap analysis:**

| Feature            | `internal/apperrors`          | `go-error-family`                              |
| ------------------ | ----------------------------- | ---------------------------------------------- |
| Error codes        | `ErrorCode` string enum       | `Coded` interface                              |
| Severity           | `ErrorSeverity` (4 levels)    | `Family` (5 behavioral levels)                 |
| Retryable          | `Retryable bool`              | `Retryable` interface + family-based inference |
| Wrapping           | `Wrap(err, code, component)`  | `Wrap(err)` with family preservation           |
| JSON serialization | `ToJSON()`                    | Not built-in (use go-output)                   |
| Error list         | `ErrorList` with filter/group | Not provided (different scope)                 |
| Exit codes         | Missing                       | `ExitCode(err)` — critical for CLI             |
| User messages      | Missing                       | `HandleError()` with What/Why/Fix              |
| Diagnostics        | Missing                       | `diagnose.Runner`                              |

**Recommendation:** Replace `internal/apperrors` with `go-error-family`. The `ErrorList` is the only feature not covered, but it's simple enough to keep locally if needed. The biggest wins are behavioral classification, exit codes, and user-facing error messages.

---

### P2: go-output — Multi-Format Output

**Current state:** Hand-rolled output using lipgloss styles:

```go
// internal/ui/styles.go — 90 lines of lipgloss styles
// internal/commands/ui_output.go — PrintSuccess, PrintInfo, PrintError, PrintNextSteps
// internal/commands/validate.go — displayValidationResults with custom lipgloss styles
```

All output is terminal-only. No `--format json`, `--format yaml`, `--format table` support.

**What go-output provides:**

- 12 output formats (JSON, CSV, Markdown, HTML, Mermaid, D2, DOT, YAML, XML, TSV, ASCII tree, styled tables)
- Type-safe format enums with `Parse()`, `IsValid()`, `AllowedValues()`
- `--format` and `--sort-by` flag support
- `TableRenderer` for structured data (perfect for `validate` output, `doctor` results)
- `GraphRenderer` for dependency visualization

**What SQLC-Wizard would gain:**

- `sqlc-wizard validate --format json` for CI integration
- `sqlc-wizard doctor --format table` for machine-readable health checks
- `sqlc-wizard migrate status --format yaml` for automation
- Structured output for all commands instead of fmt.Println

**Recommendation:** Integrate go-output for `validate`, `doctor`, and `migrate status` commands. Keep lipgloss for interactive wizard UI. Start with adding `--format` flag to `validate` command.

---

### P2: go-business-rules — Validation

**Current state:** `pkg/config/validator.go` (~120 lines) with hand-rolled validation:

```go
type ValidationResult struct {
    Errors   []ValidationError
    Warnings []ValidationError
}

type ValidationError struct {
    Field   string
    Message string
}
```

Also: `internal/wizard/output.go` has `ValidateConfiguration()`, `internal/adapters/template_real.go` has `ValidateTemplateData()`, and `internal/domain/safety_policy.go` has `IsValid()` methods scattered across enums.

**What go-business-rules provides:**

- Severity-aware validation (Info → Warning → Error → Critical) — maps exactly to the current ErrorSeverity
- JSON-serializable validation results for API responses
- Mergeable results from multiple validation sources
- Composite combinators — `All()`, `Any()`, `When()` for conditional rules
- 20+ pre-built rules with fluent builder API

**Overlap:**

- Current `ValidationResult` is a subset of go-business-rules' `ValidationResultError`
- Current `ErrorSeverity` is almost identical to go-business-rules' `Severity`
- Current field-by-field validation could use the builder pattern

**Recommendation:** Integrate go-business-rules for `pkg/config/validator` and wizard validation. The severity model maps directly. Use the builder pattern for config validation rules.

---

### P3: smart-configs — Configuration

**Current state:** Configuration is hardcoded YAML parsing with `os.ReadFile`:

```go
// pkg/config/parser.go
func ParseFile(path string) (*SqlcConfig, error) {
    data, err := os.ReadFile(path)
    // ... yaml.Unmarshal
}
```

No environment variable support, no `.env` fallback, no CI/CD-aware suggestions, no tool detection.

**What smart-configs provides:**

- Multi-source resolution: CLI args → env → `.env` → CLI tools → config files → cache → defaults
- CI/CD-aware suggestions (GitHub Actions, Docker, K8s detection)
- Service-specific resolvers (Turso DB URLs, GitHub tokens via `gh` CLI)
- Type-safe generics — `Get[int]("PORT")`, `Get[bool]("DEBUG")`

**What SQLC-Wizard would gain:**

- `DATABASE_URL` discovery from environment, `.env`, or `gh` CLI
- `sqlc` binary discovery with installation guidance
- CI/CD-aware error messages when config is missing
- Tool detection for `sqlc-wizard doctor` command

**Recommendation:** Integrate for `doctor` command's tool detection and `DATABASE_URL` resolution in migration commands. Start with `smart-configs.Get()` for environment variable fallback.

---

### P3: project-discovery-sdk — Project Discovery

**Current state:** Project type is selected manually via wizard or `--project-type` flag. No auto-detection.

**What project-discovery-sdk provides:**

- Language detection via go-enry
- File filtering (skip generated files via gogenfilter)
- Project snapshotting — capture structure metadata

**What SQLC-Wizard would gain:**

- Auto-detect project type from filesystem (if `go.mod` exists → Go project, if `Dockerfile` exists → containerized, etc.)
- Auto-detect database from existing SQL files or migration directories
- Smart defaults based on existing project structure

**Recommendation:** Consider for `init` command auto-detection. Medium effort, but significantly improves UX.

---

### P3: gogenfilter — Generated Code Detection

**Current state:** Not used.

**Potential use case:** If `sqlc-wizard` ever needs to analyze existing sqlc-generated code (e.g., `sqlc-wizard migrate` detecting what's already been generated), gogenfilter could detect and skip generated files.

**Recommendation:** Low priority. Only relevant if adding code analysis features.

---

### Not Recommended

| Library                          | Why Not                                                                             |
| -------------------------------- | ----------------------------------------------------------------------------------- |
| **go-branded-id**                | SQLC-Wizard has a single entity type (sqlc config). No multi-entity ID mixing risk. |
| **go-cqrs-lite**                 | No event sourcing, no aggregates, no command/query separation needed.               |
| **go-localfirst**                | No offline/sync requirements.                                                       |
| **go-localsync**                 | No external API sync requirements.                                                  |
| **cqrs-htmx**                    | No web layer — this is a CLI tool.                                                  |
| **templ-components**             | No web UI — TUI uses charmbracelet/huh.                                             |
| **go-commit**                    | No git commit generation in this tool's scope.                                      |
| **ActaFlow**                     | No actor model needed.                                                              |
| **go-composable-business-types** | No bitemporal data, audit trails, or money arithmetic.                              |

---

## Improvement Roadmap

### Phase 1: Foundation (P1)

| Step | Action                                                                   | Files Changed                          | Effort |
| ---- | ------------------------------------------------------------------------ | -------------------------------------- | ------ |
| 1.1  | Add `go-error-family` dependency                                         | `go.mod`                               | 1h     |
| 1.2  | Replace `internal/apperrors` error types with `go-error-family` protocol | `internal/apperrors/*.go`, all callers | 4h     |
| 1.3  | Add `HandleError()` to `main()` for Wix-quality error messages           | `cmd/sqlc-wizard/main.go`              | 1h     |
| 1.4  | Remove `pkg/errors` (already superseded by `apperrors`)                  | `pkg/errors/errors.go`                 | 15m    |
| 1.5  | Evaluate cmdguard migration proof-of-concept on `version` command        | `cmd/sqlc-wizard/main.go`              | 4h     |

### Phase 2: Output & Validation (P2)

| Step | Action                                                       | Files Changed                   | Effort |
| ---- | ------------------------------------------------------------ | ------------------------------- | ------ |
| 2.1  | Add `go-output` dependency                                   | `go.mod`                        | 30m    |
| 2.2  | Add `--format` flag to `validate` command                    | `internal/commands/validate.go` | 2h     |
| 2.3  | Add `--format` flag to `doctor` command                      | `internal/commands/doctor.go`   | 2h     |
| 2.4  | Add `go-business-rules` dependency                           | `go.mod`                        | 30m    |
| 2.5  | Refactor `pkg/config/validator.go` to use `ValidatorBuilder` | `pkg/config/validator.go`       | 3h     |

### Phase 3: Smart Features (P3)

| Step | Action                                                        | Files Changed                    | Effort |
| ---- | ------------------------------------------------------------- | -------------------------------- | ------ |
| 3.1  | Add `smart-configs` for `DATABASE_URL` resolution             | `internal/commands/migrate_*.go` | 2h     |
| 3.2  | Add `smart-configs` for tool detection in `doctor`            | `internal/commands/doctor.go`    | 2h     |
| 3.3  | Evaluate `project-discovery-sdk` for auto-detection in `init` | `internal/commands/init.go`      | 4h     |

---

## Current Dependency Audit

| Dependency                             | Version | Ecosystem Replacement                              | Notes                                                 |
| -------------------------------------- | ------- | -------------------------------------------------- | ----------------------------------------------------- |
| `spf13/cobra`                          | v1.10.2 | **cmdguard**                                       | Cobra footguns: panics, stringly-typed flags, no DI   |
| `charm.land/huh/v2`                    | v2.0.3  | Keep                                               | Good TUI library, no replacement needed               |
| `charm.land/lipgloss/v2`               | v2.0.3  | Keep for TUI + **go-output** for structured output | Dual-use: lipgloss for wizard, go-output for commands |
| `charm.land/log/v2`                    | v2.0.0  | Keep                                               | Structured logging                                    |
| `gopkg.in/yaml.v3`                     | v3.0.1  | Keep                                               | Needed for sqlc.yaml parsing                          |
| `github.com/samber/lo`                 | v1.53.0 | Keep                                               | Generic helpers                                       |
| `github.com/google/uuid`               | v1.6.0  | Keep                                               | UUID generation                                       |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | Keep                                               | Migration engine                                      |
| `github.com/onsi/ginkgo/v2`            | v2.28.3 | Keep                                               | BDD testing framework                                 |
| `github.com/onsi/gomega`               | v1.40.0 | Keep                                               | BDD matchers                                          |
| `github.com/stretchr/testify`          | v1.11.1 | Consider removing                                  | Conflicts with ginkgo/gomega — pick one               |

---

## Key Findings

### Strengths

- Clean DDD layer separation — adapters, domain, commands, wizard are well-separated
- Type-safe enums in `internal/domain/` (safety policies, emit modes) show awareness of "impossible states unrepresentable"
- Good test coverage with Ginkgo BDD framework
- Adapter pattern for external dependencies enables clean testing

### Weaknesses

- **Two error systems** (`internal/apperrors` + `pkg/errors`) — neither uses `go-error-family`
- **No structured output** — all output is `fmt.Println` or lipgloss, no `--format json` support
- **No exit codes** — CLI exits with `os.Exit(1)` for all errors regardless of type
- **Cobra without DI** — each command manually wires its dependencies
- **Dual validation** — `pkg/config/validator` + scattered `IsValid()` methods + wizard validation
- **`pkg/errors` is dead code** — `BaseError` has 2 fields, used nowhere meaningfully
- **Conflicting test frameworks** — both `testify` and `ginkgo/gomega` imported

---

_"Zero ecosystem integrations means maximum reinvention. The codebase is functional but carries ~800 lines of custom error/validation code that duplicates existing, battle-tested libraries."_
