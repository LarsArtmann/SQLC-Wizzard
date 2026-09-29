# Features

Honest feature inventory. Status answers one question: **does working code exist, and if not, why not?**

Verified against code on 2026-09-29 (build passes, all 14 packages pass `go test ./...`).

| Feature | Status | Evidence |
| --- | --- | --- |
| Interactive wizard (`init`) | FULLY_FUNCTIONAL | `internal/wizard/` step flow; package tests pass (34.2% coverage) |
| 8 project templates (hobby, microservice, enterprise, api-first, analytics, testing, multi-tenant, library) | FULLY_FUNCTIONAL | `internal/templates/registry.go`; per-template tests pass |
| ConfiguredTemplate pattern (zero-value safe templates) | FULLY_FUNCTIONAL | `internal/templates/configured_template.go`; `zero_value_test.go` |
| Config validation (`validate`, `--strict`) | FULLY_FUNCTIONAL | `internal/commands/validate.go`, `internal/validation/rule_transformer.go`; pkg/config coverage 64.5% |
| Example file generation (`generate`) | PARTIALLY_FUNCTIONAL | `internal/creators/project_creator.go:62` — TODO: full project scaffolding not wired into `CreateProject` |
| Environment doctor (`doctor`) | FULLY_FUNCTIONAL | `internal/commands/doctor.go` |
| Config migration (`migrate create/list/status`) | FULLY_FUNCTIONAL | `internal/migration/status.go`, `internal/commands/migrate_*.go` |
| Non-interactive init flags | FULLY_FUNCTIONAL | `internal/commands/init.go:43-59` (`--non-interactive`, `--project-type`, `--database`, `--package-name`, …) |
| TypeSpec-generated type-safe enums | FULLY_FUNCTIONAL | `generated/types.go` (separate module wired via `replace` in `go.mod`) |
| SQL starter templates | PARTIALLY_FUNCTIONAL | `templates/queries/`, `templates/schema/` ship PostgreSQL + SQLite only; no MySQL starter files |
| Safety rules (no SELECT \*, require WHERE, …) | PARTIALLY_FUNCTIONAL | Dual system: `domain.TypeSafeSafetyRules` (new) coexists with boolean `domain.SafetyRules` (4 files still on old type) |
| CI quality gate (test + race + lint + dupl) | FULLY_FUNCTIONAL | `.github/workflows/ci-cd.yml` (test, race, coverprofile, golangci-lint, dupl) |
| Release automation (GoReleaser on tags, GHCR image) | FULLY_FUNCTIONAL | `.github/workflows/release.yml`, `.goreleaser.yml` (tags v0.1.0, v0.2.0 cut) |
| Documentation suite (user guide, tutorial, best practices, troubleshooting, migration, advanced) | FULLY_FUNCTIONAL | `docs/USER_GUIDE.md`, `docs/TUTORIAL.md`, `docs/BEST_PRACTICES.md`, `docs/TROUBLESHOOTING.md`, `docs/MIGRATION_GUIDE.md`, `docs/ADVANCED_FEATURES.md` |
| Example projects | PARTIALLY_FUNCTIONAL | `examples/hobby-project/`, `examples/hobby-sqlite/` only; microservice and enterprise examples missing |
| Nix flake dev environment + build | FULLY_FUNCTIONAL | `flake.nix` (`packages.default`, devShell, `apps.default`) |

## Known gaps (not features yet)

| Item | Status | Why |
| --- | --- | --- |
| Wizard test coverage ≥ 60% | PLANNED | 34.2% on 2026-09-29; tracked in TODO_LIST T6 |
| Snapshot/golden tests for generated configs | PLANNED | No `snapshot_test.go`; tracked in TODO_LIST T12 |
| File-size policy gate in CI | PLANNED | 9 files > 350 lines on 2026-09-29; tracked in TODO_LIST T7 |

## Counts

Computed from repo on 2026-09-29: 16 feature rows (11 FULLY_FUNCTIONAL, 5 PARTIALLY_FUNCTIONAL, 0 BROKEN, 0 DISABLED), 3 PLANNED gaps.
