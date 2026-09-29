# SQLC-Wizzard Agent Guide

## What This Is

Interactive CLI wizard that generates production-ready `sqlc.yaml` (v2) configurations. Go 1.27 CLI app built with cobra + charmbracelet/huh. It scaffolds project files, validates configs, and manages config migrations — it does NOT replace the `sqlc` binary itself. Domain-driven layout: `cmd/` → `internal/commands/` → domain packages (`internal/domain`, `internal/templates`) with infrastructure in `internal/adapters`, `internal/generators`, `internal/creators`.

## Commands

Verified 2026-09-29. The local toolchain must match `go.mod` (go 1.27); if your system Go is older, prefix with `GOTOOLCHAIN=auto`.

```bash
GOTOOLCHAIN=auto go build ./...              # build all packages
GOTOOLCHAIN=auto go build ./cmd/sqlc-wizard  # build the CLI binary
GOTOOLCHAIN=auto go test -race ./...         # full test suite (CI parity)
GOTOOLCHAIN=auto go test -cover ./internal/wizard/  # coverage for one package
golangci-lint run ./...                      # lint (.golangci.yml)
dupl -t 100 -plumbing .                      # duplication check (CI parity)
dprint fmt                                   # format JSON/YAML/Markdown/Dockerfile (dprint.json)
go vet ./...
nix build                                    # Nix build (flake.nix); nix develop for the dev shell
```

CI (`.github/workflows/ci-cd.yml`) runs: build, `go test -v -race -coverprofile`, golangci-lint, dupl. Release (`.github/workflows/release.yml`) runs GoReleaser on version tags (GHCR image, Cosign, SBOM). There is no justfile/Makefile — do not document or create one.

## Type Generation (`generated/`)

- `generated/` is a separate Go module (`module sqlc-wizard-types`) wired into the main module via a `replace` directive in `go.mod` (import path `github.com/LarsArtmann/SQLC-Wizzard/generated`). The import path only resolves because of the replace — do not remove it.
- The module header says it is generated from `api/typespec.tsp`, but no `.tsp` spec is committed. Treat `generated/types.go` as the source of truth for enums (`ProjectType`, `DatabaseType`, emit options, …) and regenerate by editing there.
- NEVER hand-write raw strings for enum values in application code — use the generated constants (`generated.ProjectTypeMicroservice`, not `"microservice"`).

## Architecture

- `cmd/sqlc-wizard/` — entrypoint; `internal/commands/` — cobra commands (`init`, `validate`, `generate`, `doctor`, `migrate` with `create/list/status` subcommands). Commands orchestrate; business logic lives below.
- `internal/wizard/` — huh-based TUI flow; each step is a struct with a `Run()` method; results collect in `WizardResult`.
- `internal/templates/` — 8 project templates implementing the `Template` interface (`Name/Description/DefaultData/Generate/RequiredFeatures`). Two implementation patterns: `ConfiguredTemplate` embedding (preferred, zero-value safe, 7 of 8 templates) and direct `BaseTemplate` embedding (MicroserviceTemplate only — intentional, see gotchas). Shared builders: `config_builder.go`, `build_options.go`, `default_data.go`.
- `internal/domain/` — domain models and enums, including the dual safety-rule system (see gotchas).
- `internal/validation/` — rule transformer converting domain safety policy into sqlc config rules.
- `internal/adapters/`, `internal/generators/`, `internal/creators/` — infrastructure: external-world adapters, file generation, project scaffolding.
- `internal/apperrors/` — structured errors (`apperrors.NewError`, `apperrors.ValidationError`, wrapping with codes). This replaced the old `internal/errors` package (standard-library name collision).
- `pkg/config/` — complete sqlc.yaml v2 schema with parser, marshaller, validator; `pkg/errors/` is a dead leftover package with zero importers (see TODO_LIST T4).
- `internal/testing/` — shared testify helpers/assertions for template tests.

## Conventions

- File policy: keep files under 350 lines; refactor immediately when exceeded (9 files over as of 2026-09-29 — TODO_LIST T1).
- Errors: return `apperrors` types with error codes; wrap with context; never return bare `fmt.Errorf` in domain/application layers.
- Tests: Ginkgo/Gomega suites in most packages (one `RunSpecs` per suite — a second call in the same package fails the suite); plain testify + `internal/testing` helpers for template tests.
- Enums over booleans: new options must use generated enums or typed modes (`domain.EmitMode`, `TypeSafeSafetyRules`), not boolean flags.
- Docs map: FEATURES.md (status), TODO_LIST.md (open work), ROADMAP.md (ideas), CHANGELOG.md (releases). Update the right one, not AGENTS.md.

## Gotchas

- Go 1.27 required by `go.mod`; system toolchains older than that fail with `GOTOOLCHAIN=local` — use `GOTOOLCHAIN=auto`.
- MicroserviceTemplate intentionally stays on `BaseTemplate`: its `Generate()` has 50+ lines of custom logic that does not fit the ConfigBuilder pattern. Do not "finish" this migration without redesigning the builder first.
- Dual safety-rule system: boolean `domain.SafetyRules` and enum `domain.TypeSafeSafetyRules` coexist; the conversion lives in `internal/domain/conversions.go` and `internal/validation/rule_transformer.go`. New code must use the type-safe variant.
- `internal/templates/base.go` still hardcodes `${DATABASE_URL}` and `internal/db` literals despite extracted constants (`constants.go`) — use the constants, don't add more literals.
- `assert.JSONEq` requires valid JSON on BOTH sides. Style names like `"camel"` are plain strings, not JSON — comparing them with JSONEq fails with `invalid character 'c'` (this broke 8 template tests until 2026-09-29). Use `assert.Equal` for scalar comparisons.
- `generateQueryFiles()` in `internal/creators/project_creator.go` exists but is never called — full project scaffolding is not implemented (TODO_LIST T11).
- `pkg/errors` duplicates `internal/apperrors` and is imported by nothing; do not build on it.
- The pnpm-managed TypeSpec devDependencies have no committed `.tsp` spec; `package.json` alone will not regenerate `generated/types.go`.

## Dependencies

- `github.com/spf13/cobra` — CLI framework
- `github.com/charmbracelet/huh` + `lipgloss` — interactive TUI forms and styling
- `github.com/samber/lo`, `github.com/samber/mo` — functional utilities and option types (allowed in domain layer)
- `gopkg.in/yaml.v3` — sqlc.yaml parsing
- `github.com/stretchr/testify`, `github.com/onsi/ginkgo/v2` + `gomega` — testing

## References

- `FEATURES.md` — honest feature inventory with evidence
- `TODO_LIST.md` — ranked open work (harvested from status reports)
- `ROADMAP.md` — long-term ideas and non-goals
- `docs/DOMAIN_LANGUAGE.md` — domain glossary
- `docs/templates/` — template usage, comparison, customization, examples
- `examples/hobby-project/` — complete working example
