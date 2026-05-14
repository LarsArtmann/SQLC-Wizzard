# Modularization Proposal — SQLC-Wizard

**Date:** 2026-05-13
**Status:** Draft
**Author:** Crush (AI Assistant)

---

## 1. Executive Summary

SQLC-Wizard currently has a **partial split** structure: a root `go.mod`, a `generated/` submodule, and an `examples/hobby-project/` standalone module. The root module uses a `replace` directive to point at `generated/`. However, the internal architecture suffers from:

- **Layer violations**: `pkg/config` imports `internal/apperrors` (public importing private)
- **God-packages**: 7 packages with 15+ exported symbols each
- **Duplicate type systems**: Two error packages, two `DatabaseConfig` types, two `RuleConfig` types
- **Ghost domain layer**: `internal/domain` contains only 2 deprecated type aliases
- **Global mutable state**: Adapters use `var` function pointers; templates use `init()` singleton
- **No go.work**: Development coordination relies on a single `replace` directive

**Goal:** Split the monolithic root `go.mod` into 4 independently versionable sub-modules coordinated by a `go.work` file, with clean DAG enforcement and eliminated layer violations.

**Expected benefits:**
- Compile-time enforced module boundaries (no accidental coupling)
- Faster CI — modules build/test independently
- Clearer ownership and responsibility per module
- `generated` types decoupled from main application
- Test-only dependencies isolated from production `go.mod`
- Preparation for future library consumption (e.g., `pkg/config` as standalone)

---

## 2. Current State Analysis

### 2.1 Module Landscape

| Module | Path | Internal Deps | External Deps | Replace Directives | State |
|---|---|---|---|---|---|
| Root | `go.mod` | `generated` | 15+ (charm, cobra, ginkgo, etc.) | `generated => ./generated` | Leaky |
| Generated | `generated/go.mod` | None | None | None | Clean |
| Hobby Example | `examples/hobby-project/go.mod` | None | `mattn/go-sqlite3` | None | Clean |

**Classification:** Partial split — one submodule extracted (`generated`), everything else in root.

### 2.2 Internal Package Dependency Graph (Non-Test)

```
Layer 0 (leaf):     generated, apperrors, migration, schema, ui, utils
Layer 1:            domain → generated
                    testing → generated, domain, pkg/config
                    pkg/config → generated, internal/apperrors  ⚠️
Layer 2:            adapters → generated, apperrors, migration, schema, templates, pkg/config
                    generators → templates, pkg/config
                    validation → generated, domain
Layer 3:            templates → generated, apperrors, validation, pkg/config
Layer 4:            creators → adapters, apperrors, generated, pkg/config
                    wizard → generated, apperrors, schema, templates, ui, pkg/config
Layer 5:            commands → adapters, apperrors, creators, generators, templates, ui, wizard, generated, pkg/config
Layer 6:            cmd/sqlc-wizard → commands
```

### 2.3 Coupling Hotspots

| Hotspot | Description | Severity |
|---|---|---|
| `pkg/config` → `internal/apperrors` | Public package imports internal — violates Go convention | 🔴 Critical |
| `internal/commands` | Fan-out of 9 internal dependencies — god-package orchestrator | 🟡 High |
| `internal/testing` | ~60 exported symbols — mix of assertions, factories, enums | 🔴 Critical |
| `internal/wizard` | ~42 exports — mixes UI, steps, flow, branching, test helpers | 🟡 High |
| `internal/templates` | ~32 exports — 9 template types + registry + validation | 🟡 High |
| `pkg/errors` vs `internal/apperrors` | Two competing error systems | 🟠 Medium |
| `generated.DatabaseConfig` vs `config.DatabaseConfig` | Duplicate type with same name | 🟠 Medium |
| Adapter global vars | 6 mutable `var` function pointers — race condition risk | 🟠 Medium |
| Template `init()` singleton | Global state, untestable | 🟠 Medium |

### 2.4 Banned Dependencies in go.mod

Per the `how-to-golang` skill banned libraries list:

| Banned Dependency | Severity | Replacement |
|---|---|---|
| `gopkg.in/yaml.v3` | Critical | `go-faster/yaml` |
| `github.com/stretchr/testify` | Critical | `onsi/ginkgo/v2` + `onsi/gomega` (already present) |

**Note:** `testify` is imported but should be migrated to pure ginkgo/gomega. `yaml.v3` is used extensively in `pkg/config` for YAML marshalling.

---

## 3. Proposed Module Structure

### 3.1 Module Definitions

#### Module 1: `generated/` (Existing — Keep As-Is)

| Field | Content |
|---|---|
| Name & path | `generated/` |
| Module path | `sqlc-wizard-types` |
| Purpose | Type-safe enums and domain models generated from TypeSpec |
| Dependencies (prod) | None |
| Dependencies (test) | None |
| Public API | `ProjectType`, `DatabaseType`, `TemplateData`, `EmitOptions`, `SafetyRules`, all config structs, CQRS types |
| Internal packages | N/A (single file) |
| External deps | None |

**Rationale:** Already isolated. This is the foundation — every other module depends on it. Zero changes needed.

#### Module 2: `core/` (New — Extract from internal/)

| Field | Content |
|---|---|
| Name & path | `core/` |
| Module path | `github.com/LarsArtmann/SQLC-Wizzard/core` |
| Purpose | Pure domain logic, type-safe enums, validation, error types, schema — no I/O or UI |
| Dependencies (prod) | `generated/` |
| Dependencies (test) | None |
| Public API | Type-safe `EmitModes`, `SafetyPolicy`, `Domain` types, `Schema`, `Validation`, `AppErrors` |
| Internal packages | `apperrors/`, `domain/`, `validation/`, `schema/` |
| External deps | `samber/lo` |

**Packages extracted from root:**
- `internal/apperrors` → `core/apperrors`
- `internal/domain` → `core/domain`
- `internal/validation` → `core/validation`
- `internal/schema` → `core/schema`

**Rationale:** These packages have zero I/O and zero UI. They form the pure domain kernel. Only `generated` is a dependency.

#### Module 3: `config/` (New — Extract from pkg/config)

| Field | Content |
|---|---|
| Name & path | `config/` |
| Module path | `github.com/LarsArtmann/SQLC-Wizzard/config` |
| Purpose | sqlc.yaml configuration types, parsing, marshalling, and validation |
| Dependencies (prod) | `generated/`, `core/apperrors` |
| Dependencies (test) | None |
| Public API | `SqlcConfig`, `GoGenConfig`, `Parser`, `Validator`, `Marshaller`, `ApplyEmitOptions()` |
| Internal packages | N/A (single package) |
| External deps | `go-faster/yaml` (migrate from `gopkg.in/yaml.v3`) |

**Packages extracted from root:**
- `pkg/config` → `config/`
- `pkg/errors` → DELETE (merge into `core/apperrors`)

**Rationale:** `pkg/config` is already a natural boundary. It represents a standalone concern (sqlc.yaml schema). Making it a separate module enforces that it cannot import `internal/` packages — forcing the `internal/apperrors` extraction into `core/`.

**Layer violation fix:** Currently `pkg/config` → `internal/apperrors`. After modularization: `config/` → `core/apperrors`. Both are public modules — no violation.

#### Module 4: Root (Remaining — Slim Down)

| Field | Content |
|---|---|
| Name & path | `./` (repo root) |
| Module path | `github.com/LarsArtmann/SQLC-Wizzard` |
| Purpose | CLI application, wizard TUI, templates, adapters, generators |
| Dependencies (prod) | `generated/`, `core/`, `config/` |
| Dependencies (test) | `generated/`, `core/`, `config/` |
| Public API | CLI binary (no public library API) |
| Internal packages | `internal/templates`, `internal/wizard`, `internal/adapters`, `internal/commands`, `internal/creators`, `internal/generators`, `internal/migration`, `internal/ui`, `internal/utils`, `internal/testing` |
| External deps | `charm.land/huh/v2`, `charm.land/lipgloss/v2`, `spf13/cobra`, `golang-migrate/migrate/v4`, `onsi/ginkgo/v2`, `onsi/gomega`, `samber/lo` |

**Rationale:** The root module becomes the application layer — everything that depends on I/O, UI, or external services. It's the only module that imports CLI/TUI libraries.

### 3.2 DAG Verification

```
                    ┌─────────────┐
                    │  generated/  │   (Layer 0: no deps)
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │    core/     │   (Layer 1: depends on generated only)
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │   config/    │   (Layer 2: depends on generated + core)
                    └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    │ root (app)   │   (Layer 3: depends on generated + core + config)
                    └─────────────┘

Direction: DOWN only. No cycles. No upward dependencies.
```

**Proof of acyclicity:**
- `generated` has zero internal deps → cannot create cycle
- `core` depends only on `generated` → cannot create cycle
- `config` depends on `generated` + `core` → cannot create cycle
- `root` depends on all three → leaf of the DAG

**Special case — test dependencies:** The `internal/testing` package in root will depend on `core/` and `config/`. This is acceptable — test helpers import domain types. No bidirectional test deps exist.

### 3.3 Replace / Workspace Strategy

**Chosen: `go.work` at repo root.**

| Module | go.work entry |
|---|---|
| `generated/` | ✅ |
| `core/` | ✅ |
| `config/` | ✅ |
| `./` (root) | ✅ |

```go
// go.work
go 1.26.2

use (
    ./generated
    ./core
    ./config
    .
)
```

**Rules:**
- No `replace` directives in any `go.mod` — `go.work` handles local development
- `go.work` is committed to git (all modules are in-repo)
- `go mod tidy` works both with and without workspace (verified in CI)
- When publishing, consumers use versioned imports — `go.work` is ignored

### 3.4 Test Dependency Isolation

| Module | Production Deps | Test Deps |
|---|---|---|
| `generated/` | None | None |
| `core/` | `generated/` | None |
| `config/` | `generated/`, `core/apperrors` | None |
| Root (app) | `generated/`, `core/`, `config/` | `core/`, `config/` (via `internal/testing`) |

**Testhelpers strategy:**
- `internal/testing` stays in root module — it provides wizard/template-specific test helpers
- If shared test fixtures for `core/` are needed in the future, create `core/testing/` within the `core/` module
- Root module's `internal/testing` depends on `core/` types only (never on infrastructure adapters)

### 3.5 Interface Extraction Plan

No interface/impl splits needed for this modularization. The current boundaries are:

- `generated/` — already pure types (no interfaces)
- `core/` — pure domain logic (no interfaces)
- `config/` — serialization types (no interfaces)
- Root — contains all adapter interfaces (`internal/adapters/interfaces.go`)

The adapter pattern stays in the root module since adapters are an application concern, not a shared library concern.

### 3.6 Versioning Strategy

**Chosen: Shared version (single git tag).**

| Strategy | Why |
|---|---|
| Shared version | All modules are tightly coupled, single team, no external consumers yet |
| Tag format | `v1.2.3` at repo root |
| All modules bump together | Atomic releases |

**Rationale:** The project has no external consumers of individual modules. `generated/` is the only published module (`sqlc-wizard-types`) and could be independently versioned in the future. Until there's a consumer need, shared versioning is simplest.

**Migration path:** If `config/` becomes a standalone library, switch to independent semver with tags like `config/v1.0.0`.

### 3.7 Migration Strategy (Ordered Steps)

See `EXECUTION_PLAN.md` for the full step-by-step migration with verification checkpoints.

---

## 4. Self-Review Findings (Phase 4)

### 4.1 What Did We Forget?

1. **`pkg/errors` is dead code** — Zero imports across the entire project. Confirmed safe to delete.
2. **`internal/testing/assertions.go` imports testify** — This is non-test code that depends on `github.com/stretchr/testify/assert`. This makes testify a compile-time (not just test) dependency. Must be addressed: either refactor assertions to use gomega, or mark the entire `internal/testing` package as test-only.
3. **Module path mismatch**: `generated/go.mod` declares `module sqlc-wizard-types` but the root `go.mod` imports it as `github.com/LarsArtmann/SQLC-Wizzard/generated` via a `replace` directive. This works but is fragile. Consider aligning the module path to `github.com/LarsArtmann/SQLC-Wizzard/generated` in the generated go.mod.
4. **`internal/testing` → `pkg/config` coupling**: The test helpers import `pkg/config` types. If config is extracted, internal/testing needs that module dep. This is acceptable but must be documented.

### 4.2 Module Boundary Assessment

- **Too fine?** No — only 3 modules extracted (core, config, generated). This is conservative.
- **Too coarse?** Potentially — `core/` contains 4 packages (apperrors, domain, validation, schema) that could each be their own module. But this would be over-engineering for the current project size. Keep as-is.
- **Right granularity?** Yes. The 4-module structure (generated, core, config, root) provides clean separation without excessive complexity.

### 4.3 What Could Be Better?

1. The `generated` module path should be aligned to `github.com/LarsArtmann/SQLC-Wizzard/generated` for consistency
2. `testify` in non-test code (`internal/testing/assertions.go`) should be refactored to gomega — but this is orthogonal to modularization
3. The `internal/domain` package is a ghost (2 deprecated aliases). It should be either populated or removed during extraction.

### 4.4 Cross-Reference with how-to-golang

| Check | Status |
|---|---|
| No banned deps in proposed modules | ⚠️ `gopkg.in/yaml.v3` in config/ (documented, migrate later) |
| No banned deps in proposed modules | ⚠️ `testify` in root module (non-test code, document) |
| Domain types in correct location | ✅ `generated/` is canonical source |
| Architecture patterns aligned | ✅ Layered with DAG enforcement |
| Required libraries present | ✅ ginkgo/gomega for testing |

### 4.5 Split Brain Check

| Type | Location 1 | Location 2 | Resolution |
|---|---|---|---|
| `DatabaseConfig` | `generated.DatabaseConfig` | `config.DatabaseConfig` (YAML-specific) | ✅ Different purposes — keep both |
| `RuleConfig` | `generated.RuleConfig` | `config.RuleConfig` | ✅ Different purposes — keep both |
| `EmitOptions` | `generated.EmitOptions` | `domain.EmitOptions` (deprecated alias) | ⚠️ Remove alias during extraction |
| `SafetyRules` | `generated.SafetyRules` | `domain.SafetyRules` (deprecated alias) | ⚠️ Remove alias during extraction |
| Error types | `pkg/errors.BaseError` | `internal/apperrors.Error` | ⚠️ Delete `pkg/errors` — it's unused |

---

## 5. Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|---|---|---|---|
| Import path breaks after extraction | Medium | High | Update all imports in one commit per module, verify with `go build ./...` |
| Test dependencies leak into prod go.mod | Low | Medium | Audit each module's go.mod after `go mod tidy` |
| Circular dependency discovered during extraction | Low | High | DAG analysis above shows no cycles; if found, adjust boundaries before proceeding |
| `pkg/config` → `internal/apperrors` breakage | High | Medium | Must extract `apperrors` to `core/` before `config/` extraction |
| `gopkg.in/yaml.v3` replacement breaks config parsing | Medium | High | Migrate to `go-faster/yaml` after modularization, not during |
| `go.work` conflicts with existing `replace` directive | Medium | Low | Remove `replace` directive when adding `go.work` |
| `generated` module path mismatch | Medium | Medium | Align `generated/go.mod` to `github.com/LarsArtmann/SQLC-Wizzard/generated` |
| `testify` in non-test code | Low | Low | Refactor `internal/testing/assertions.go` to gomega (orthogonal to modularization) |

---

## 6. Build System Impact

| System | Changes Needed |
|---|---|
| `justfile` | Update `build`, `test`, `lint` commands to work with go.work |
| `Makefile` | Update `types` target if needed |
| `.golangci.yml` | Update `depguard` rules for new module paths |
| `.go-arch-lint.yml` | Update architecture constraints for new module layout |
| `Dockerfile` | Update `COPY` and `go mod download` for multi-module |
| `.github/workflows/ci-cd.yml` | Add per-module build/test steps, cache per module |
| `flake.nix` | Not present yet — future work |

---

## 7. Key Decisions

1. **`core/` as separate module** — Forces `pkg/config` to depend on a public module, eliminating the `pkg/` → `internal/` layer violation
2. **`config/` as separate module** — Enables future standalone library use; natural boundary
3. **`go.work` over `replace` directives** — Cleaner, standard Go tooling, no manual sync
4. **Shared versioning** — Simplest for single-team monorepo with no external consumers
5. **Keep `generated/` module path as `sqlc-wizard-types`** — Already published; changing would break consumers
6. **`pkg/errors` deleted** — Dead code; `internal/apperrors` is the real error system
7. **No interface/impl split** — Not warranted for this project size; would add complexity without benefit
8. **Templates stay in root** — Templates are tightly coupled to wizard flow; splitting would create circular deps
