# TODO List

Open, bounded work. Done items are deleted (they live in CHANGELOG). Ranked by impact/effort.

Harvested from `docs/status/2026-06-14_13-03-File-Size-Refactor-Complete.md`, `2026-03-22_06-05_oxfmt-fix-status.md`, and earlier reports; every item verified against code on 2026-09-29.

## P0 — correctness and hygiene

| ID | Task | Evidence | Effort |
| --- | --- | --- | --- |
| T1 | Refactor 9 files over the 350-line policy limit: `internal/creators/project_creator.go` (434), `internal/wizard/features.go` (392), `internal/adapters/migration_real.go` (378), `internal/templates/template_validation_test.go` (376), `internal/wizard/branching_flow_test.go` (369), `internal/wizard/wizard_step_implementation_test.go` (364), `pkg/config/validator_test.go` (360), `internal/wizard/wizard_run_integration_test.go` (358), `internal/wizard/wizard.go` (355) | `wc -l` 2026-09-29; report §c | Medium |
| T2 | Delete deprecated 21-param `BuildDefaultData`; migrate its 5+ callers to `BuildDefaultDataFromOptions` | `internal/templates/default_data.go:48` (`// Deprecated:`); callers in `configured_template.go:261`, `api_first.go:44`, `enterprise.go:46`, `hobby.go:41`, `multi_tenant.go:94` | Low |
| T3 | Finish `domain.SafetyRules` → `TypeSafeSafetyRules` migration; delete the boolean-flag struct | Old type still referenced in `internal/testing/safety_rules_helpers.go`, `internal/validation/rule_transformer.go`, `rule_transformer_helpers_test.go`, `rule_transformer_unit_test.go` | Medium |
| T4 | Remove or merge `pkg/errors` (dead parallel error package) into `internal/apperrors` | `pkg/errors/errors.go` has zero importers (grep 2026-09-29); split brain with `internal/apperrors` | Low |
| T5 | Delete backup junk: `internal/adapters/adapters_test.go.bak2`, `.bak4`, stray root files `1763227265_test_migration.{up,down}.sql` | Present in repo root and `internal/adapters/` | Trivial |

## P1 — quality gates and coverage

| ID | Task | Evidence | Effort |
| --- | --- | --- | --- |
| T6 | Raise wizard coverage (34.2% on 2026-09-29) toward 60%; it is the core UI component | `go test -cover ./internal/wizard/` 2026-09-29 | High |
| T7 | Raise adapters coverage (22.9%, lowest package) | `go test -cover ./internal/adapters/` 2026-09-29 | Medium |
| T8 | Add CI file-size gate that fails when a `.go` file exceeds 350 lines | CI (`.github/workflows/ci-cd.yml`) runs tests/lint/dupl but no size check; report §f #11 | Low |
| T9 | Add snapshot/golden tests for generated `sqlc.yaml` of all 8 templates | No `snapshot_test.go` in `internal/templates/` (verified 2026-09-29) | Medium |
| T10 | Add missing per-template zero-value tests (APIFirst, Enterprise, Microservice) | `internal/templates/zero_value_test.go` covers only `ConfiguredTemplate` generally | Low |
| T11 | Wire `generateQueryFiles()` into `CreateProject` or delete it | `internal/creators/project_creator.go:62` — "TODO: Full project scaffolding is not yet implemented" | Medium |
| T12 | Replace hardcoded `internal/db`/`${DATABASE_URL}` string repeats with `DefaultPackagePath`/`DefaultDatabaseURL` constants | `internal/templates/base.go:147-167` (goconst: 4× URL, 17× path) | Low |

## P2 — completeness

| ID | Task | Evidence | Effort |
| --- | --- | --- | --- |
| T13 | Add microservice and enterprise example projects (only hobby examples exist) | `examples/` holds `hobby-project/`, `hobby-sqlite/` only | Medium |
| T14 | Ship MySQL starter query/schema templates (wizard supports MySQL engine; starter files are PostgreSQL/SQLite only) | `templates/queries/`, `templates/schema/` | Low |
| T15 | Add CONTRIBUTING.md note about the 350-line file policy | `CONTRIBUTING.md` has no size-policy mention (grep 2026-09-29) | Trivial |
| T16 | Test-suite tidy: consolidate duplicated `BeforeEach` blocks in `internal/creators` tests, add godoc to extracted test helpers, add a table-driven `ValidateAllProjectTypes` covering all 8 types | Harvested from 2026-06-14 report rows 14/15/17 | Low |
| T17 | Clear remaining lint findings: `paralleltest` on template tests (add `t.Parallel()` or justify disabling), investigate `getJSONTagsCaseStyle` warning in `template_validation_test.go` | Harvested from 2026-06-14 report rows 18/23; LSP warnings 2026-09-29 | Low |
