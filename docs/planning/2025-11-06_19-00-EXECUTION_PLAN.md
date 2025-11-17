# SQLC-Wizard Execution Plan - Micro Tasks (12min each)

## Legend
- **Effort**: Time estimate per task
- **Impact**: How foundational (1-10)
- **Value**: Customer value (1-10)
- **Dep**: Dependencies (task numbers)

## Phase 1: Foundation (Tasks 1-15) - 3 hours

| # | Task | Effort | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 1 | Create go.mod with module path github.com/LarsArtmann/SQLC-Wizzard | 5min | 10 | 10 | - | ⏳ |
| 2 | Add core dependencies: cobra, viper, bubbletea, lipgloss | 10min | 10 | 10 | 1 | ⏳ |
| 3 | Add utility dependencies: samber/lo, samber/mo, samber/do | 8min | 9 | 9 | 1 | ⏳ |
| 4 | Add LarsArtmann/uniflow for user-friendly errors | 5min | 9 | 9 | 1 | ⏳ |
| 5 | Add ginkgo/gomega for testing | 5min | 8 | 8 | 1 | ⏳ |
| 6 | Create directory structure: cmd/sqlc-wizard/ | 3min | 10 | 9 | - | ⏳ |
| 7 | Create directory structure: internal/{wizard,templates,validators,detectors,generators,plugins}/ | 5min | 10 | 9 | 6 | ⏳ |
| 8 | Create directory structure: pkg/{config,database,version}/ | 3min | 10 | 9 | 6 | ⏳ |
| 9 | Create directory structure: templates/{sqlc,queries,workflows,migrations}/ | 3min | 9 | 8 | 6 | ⏳ |
| 10 | Create .gitignore for Go project (bin/, vendor/, *.test, coverage) | 5min | 7 | 7 | - | ⏳ |
| 11 | Create basic Makefile with build, test, install, clean targets | 12min | 6 | 7 | 1 | ⏳ |
| 12 | Setup cmd/sqlc-wizard/main.go with basic cobra/fang integration | 12min | 10 | 10 | 2,6 | ⏳ |
| 13 | Add version command and build-time version injection | 8min | 5 | 6 | 12 | ⏳ |
| 14 | Create internal/wizard/wizard.go interface and base types | 10min | 9 | 9 | 7 | ⏳ |
| 15 | Setup dependency injection container using samber/do | 12min | 8 | 7 | 3,7 | ⏳ |

## Phase 2: Config & Templates (Tasks 16-28) - 2.6 hours

| # | Task | Effort | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 16 | Create pkg/config/types.go with SqlcConfig structs matching sqlc.yaml v2 schema | 12min | 10 | 10 | 8 | ⏳ |
| 17 | Implement pkg/config/parser.go to parse YAML into structs using yaml.v3 | 12min | 10 | 10 | 16 | ⏳ |
| 18 | Add pkg/config/validator.go with basic structure validation | 12min | 9 | 9 | 16 | ⏳ |
| 19 | Create pkg/config/marshaller.go to write structs back to YAML with proper formatting | 12min | 9 | 9 | 16 | ⏳ |
| 20 | Create internal/templates/types.go with TemplateData, ProjectType, DatabaseType enums | 10min | 9 | 9 | 7 | ⏳ |
| 21 | Implement internal/templates/loader.go for embedded template file loading | 12min | 9 | 8 | 20 | ⏳ |
| 22 | Create internal/templates/renderer.go using Go text/template for variable substitution | 12min | 9 | 9 | 20 | ⏳ |
| 23 | Define templates/sqlc/microservice-postgresql.yaml.tmpl with complete sqlc v2 config | 12min | 9 | 10 | 9 | ⏳ |
| 24 | Add template validation to ensure all {{.Variables}} are defined | 10min | 7 | 8 | 22 | ⏳ |
| 25 | Create internal/templates/microservice.go with PostgreSQL defaults | 12min | 9 | 10 | 20,23 | ⏳ |
| 26 | Add emit options defaults (json_tags, interface, prepared_queries, etc.) | 10min | 8 | 9 | 25 | ⏳ |
| 27 | Create templates/queries/postgresql/users-crud.sql with example queries | 12min | 8 | 9 | 9 | ⏳ |
| 28 | Create templates/schema/postgresql/001_users_table.sql with example schema | 12min | 8 | 9 | 9 | ⏳ |

## Phase 3: Interactive Wizard UI (Tasks 29-40) - 2.4 hours

| # | Task | Effort | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 29 | Create internal/wizard/ui/styles.go with lipgloss theme definitions | 10min | 8 | 8 | 7 | ⏳ |
| 30 | Implement internal/wizard/ui/prompts.go with bubbletea select prompt | 12min | 9 | 9 | 29 | ⏳ |
| 31 | Add text input prompt with validation using bubbletea textinput | 12min | 9 | 9 | 30 | ⏳ |
| 32 | Create multi-select prompt for features (UUIDs, JSON, etc.) | 12min | 8 | 8 | 30 | ⏳ |
| 33 | Implement internal/wizard/steps/project_type.go step (Microservice/Hobby/Enterprise) | 12min | 9 | 10 | 30 | ⏳ |
| 34 | Create internal/wizard/steps/database.go step (PostgreSQL/SQLite/MySQL) | 12min | 9 | 10 | 30 | ⏳ |
| 35 | Implement internal/wizard/steps/package_path.go with go.mod detection | 12min | 9 | 9 | 31 | ⏳ |
| 36 | Add internal/wizard/steps/features.go for database feature selection | 10min | 7 | 8 | 32 | ⏳ |
| 37 | Create internal/wizard/steps/output_dirs.go for path configuration | 10min | 8 | 8 | 31 | ⏳ |
| 38 | Implement internal/wizard/orchestrator.go to run steps sequentially | 12min | 10 | 10 | 33-37 | ⏳ |
| 39 | Add step navigation (back, forward, quit) with Railway pattern | 12min | 8 | 7 | 38 | ⏳ |
| 40 | Create internal/wizard/state.go to maintain wizard state across steps | 10min | 9 | 8 | 38 | ⏳ |

## Phase 4: Init Command (Tasks 41-48) - 1.6 hours

| # | Task | Effort | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 41 | Create cmd/sqlc-wizard/commands/init.go with cobra command definition | 10min | 10 | 10 | 12 | ⏳ |
| 42 | Add flags to init command (--non-interactive, --project-type, --database, etc.) | 12min | 9 | 9 | 41 | ⏳ |
| 43 | Implement init command interactive mode using wizard orchestrator | 12min | 10 | 10 | 38,41 | ⏳ |
| 44 | Add init command non-interactive mode using flags only | 10min | 8 | 8 | 42 | ⏳ |
| 45 | Create internal/generators/sqlc_generator.go to write sqlc.yaml | 12min | 10 | 10 | 19,22 | ⏳ |
| 46 | Implement internal/generators/query_generator.go for example SQL files | 10min | 8 | 9 | 27 | ⏳ |
| 47 | Add internal/generators/schema_generator.go for database schema | 10min | 8 | 9 | 28 | ⏳ |
| 48 | Create success message with next steps using lipgloss formatting | 8min | 7 | 8 | 29 | ⏳ |

## Phase 5: Validation (Tasks 49-53) - 1 hour

| # | Task | Errors | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 49 | Create internal/validators/sqlc_validator.go with YAML schema validation | 12min | 9 | 9 | 18 | ⏳ |
| 50 | Add best practices validation (emit_interface, prepared_queries checks) | 12min | 8 | 8 | 49 | ⏳ |
| 51 | Implement security validation (no unsafe patterns, require-where checks) | 12min | 8 | 9 | 49 | ⏳ |
| 52 | Create cmd/sqlc-wizard/commands/validate.go command | 10min | 9 | 9 | 12 | ⏳ |
| 53 | Add validate command with --fix flag for auto-corrections | 12min | 8 | 8 | 49,52 | ⏳ |

## Phase 6: Detection & Smart Defaults (Tasks 54-57) - 48min

| # | Task | Effort | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 54 | Create internal/detectors/project_detector.go to read go.mod | 12min | 8 | 8 | 7 | ⏳ |
| 55 | Add internal/detectors/database_detector.go to find existing DB files | 12min | 7 | 7 | 7 | ⏳ |
| 56 | Implement internal/detectors/config_detector.go to find existing sqlc.yaml | 12min | 7 | 7 | 17 | ⏳ |
| 57 | Integrate detectors into wizard for smart defaults | 12min | 8 | 9 | 38,54-56 | ⏳ |

## Phase 7: Testing (Tasks 58-60) - 36min

| # | Task | Effort | Impact | Value | Dep | Status |
|---|------|--------|--------|-------|-----|--------|
| 58 | Create internal/wizard/wizard_test.go with ginkgo suite setup | 12min | 8 | 7 | 5,14 | ⏳ |
| 59 | Write BDD tests for wizard orchestrator flow | 12min | 8 | 7 | 38,58 | ⏳ |
| 60 | Add integration test: init command end-to-end | 12min | 8 | 8 | 43,58 | ⏳ |

---

## Summary

**Total Micro Tasks**: 60
**Total Estimated Time**: ~12 hours
**MVP (Tasks 1-48)**: ~9.6 hours - Ships working `init` command
**V2 (Tasks 49-57)**: ~1.8 hours - Adds validation & smart defaults
**Testing (Tasks 58-60)**: ~36min - Quality assurance

**Critical Path**: 1→2→6→12→14→20→22→25→30→33→38→41→43→45
**Parallel Opportunities**: Templates (23-28), Detectors (54-57), UI components (29-32)

## Customer Value Milestones

1. **After Task 15**: Project structure exists, can build
2. **After Task 28**: All templates ready
3. **After Task 40**: Wizard UI complete
4. **After Task 48**: 🎉 **MVP SHIPPED** - User can run `sqlc-wizard init`
5. **After Task 53**: Validation adds safety
6. **After Task 57**: Smart defaults improve UX
7. **After Task 60**: Tests ensure quality

## Next Steps

Execute tasks sequentially:
1. Run task
2. Verify it works
3. `git add . && git commit -m "detailed message"`
4. Continue to next task
5. When milestone reached: `git push`
