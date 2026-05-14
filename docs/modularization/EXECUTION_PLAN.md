# Execution Plan — SQLC-Wizard Modularization

**Date:** 2026-05-13
**Depends on:** `PROPOSAL.md`, `DEPENDENCY_GRAPH.md`
**Branch:** `modularize/split-core-config`

---

## Task Overview

| Tier | Tasks | Impact |
|---|---|---|
| 1% → 51% | 1–3 | Foundational: fix module path, extract core, fix layer violation |
| 4% → 64% | 4–5 | High leverage: extract config, add go.work |
| 20% → 80% | 6–8 | Broad value: update CI, update linters, documentation |
| Remaining | 9–10 | Polish: delete dead code, verify independence |

---

## Task Details

### Task 1: Fix `generated` Module Path Alignment

**Impact:** 1% → 51% (foundational — all subsequent tasks depend on correct module resolution)
**Estimated effort:** 15 min
**Depends on:** Nothing

**What:** Change `generated/go.mod` module path from `sqlc-wizard-types` to `github.com/LarsArtmann/SQLC-Wizzard/generated`.

**Why:** The current mismatch (`module sqlc-wizard-types` vs imports using `github.com/LarsArtmann/SQLC-Wizzard/generated`) relies entirely on the `replace` directive. Aligning the module path makes it explicit and prepares for `go.work`.

**Steps:**
1. Edit `generated/go.mod`: change `module sqlc-wizard-types` to `module github.com/LarsArtmann/SQLC-Wizzard/generated`
2. Run `go mod tidy` in root
3. Run `go build ./...`
4. Run `go test ./...`

**Verification:**
```bash
go build ./... && go test ./... && go vet ./...
```

**Rollback:** `git revert HEAD` — single line change in generated/go.mod

---

### Task 2: Delete Dead Code (`pkg/errors`)

**Impact:** 1% → 51% (removes confusion before extraction)
**Estimated effort:** 10 min
**Depends on:** Nothing (parallel with Task 1)

**What:** Delete `pkg/errors/errors.go` — zero imports across the project.

**Why:** Two error packages is confusing. `pkg/errors` is dead code. Removing it before modularization eliminates the split brain.

**Steps:**
1. Delete `pkg/errors/errors.go`
2. Run `go mod tidy`
3. Run `go build ./...`
4. Run `go test ./...`

**Verification:**
```bash
go build ./... && go test ./... && go vet ./...
```

**Rollback:** `git revert HEAD`

---

### Task 3: Extract `core/` Module

**Impact:** 1% → 51% (foundational — eliminates layer violation, enables config extraction)
**Estimated effort:** 30 min
**Depends on:** Task 1

**What:** Create `core/` directory with its own `go.mod`, move 4 packages from root's `internal/` into it.

**Packages to move:**
- `internal/apperrors` → `core/apperrors`
- `internal/domain` → `core/domain`
- `internal/validation` → `core/validation`
- `internal/schema` → `core/schema`

**Steps:**
1. Create `core/` directory
2. Create `core/go.mod`:
   ```
   module github.com/LarsArtmann/SQLC-Wizzard/core
   
   go 1.26.2
   
   require github.com/LarsArtmann/SQLC-Wizzard/generated v0.0.0-20260504181409-ba8146c868b5
   
   require github.com/samber/lo v1.53.0 // indirect
   ```
3. Move each package directory:
   - `git mv internal/apperrors core/apperrors`
   - `git mv internal/domain core/domain`
   - `git mv internal/validation core/validation`
   - `git mv internal/schema core/schema`
4. Update ALL import paths in moved files:
   - `github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors` → `github.com/LarsArtmann/SQLC-Wizzard/core/apperrors`
   - `github.com/LarsArtmann/SQLC-Wizzard/internal/domain` → `github.com/LarsArtmann/SQLC-Wizzard/core/domain`
   - `github.com/LarsArtmann/SQLC-Wizzard/internal/validation` → `github.com/LarsArtmann/SQLC-Wizzard/core/validation`
   - `github.com/LarsArtmann/SQLC-Wizzard/internal/schema` → `github.com/LarsArtmann/SQLC-Wizzard/core/schema`
5. Update ALL import paths in remaining root files that reference these packages
6. Remove moved packages from root `go.mod` dependencies (go mod tidy handles this)
7. Remove deprecated type aliases in `core/domain/domain.go` (they just re-export `generated` types)
8. Run `go mod tidy` in `core/`
9. Run `go build ./core/...`
10. Run `go test ./core/...`
11. Run `go mod tidy` in root
12. Run `go build ./...`
13. Run `go test ./...`

**Files that need import updates (root module):**
- `internal/templates/*.go` — imports `internal/apperrors`, `internal/validation`
- `internal/adapters/*.go` — imports `internal/apperrors`, `internal/schema`, `internal/migration`
- `internal/commands/*.go` — imports `internal/apperrors`
- `internal/creators/*.go` — imports `internal/apperrors`
- `internal/wizard/*.go` — imports `internal/apperrors`, `internal/schema`
- `pkg/config/*.go` — imports `internal/apperrors` → becomes `core/apperrors` (fixes layer violation!)
- `internal/testing/*.go` — imports `internal/domain`

**Verification:**
```bash
# In core/
cd core && go build ./... && go test ./... && go vet ./... && cd ..

# In root
go mod tidy && go build ./... && go test ./... && go vet ./...
```

**Rollback:** `git revert HEAD` — all moves and import updates in one commit

---

### Task 4: Extract `config/` Module

**Impact:** 4% → 64% (high leverage — standalone library potential)
**Estimated effort:** 25 min
**Depends on:** Task 3 (needs core/apperrors to exist as separate module)

**What:** Create `config/` directory with its own `go.mod`, move `pkg/config` into it.

**Steps:**
1. Create `config/` directory
2. Create `config/go.mod`:
   ```
   module github.com/LarsArtmann/SQLC-Wizzard/config
   
   go 1.26.2
   
   require (
       github.com/LarsArtmann/SQLC-Wizzard/core v0.0.0
       github.com/LarsArtmann/SQLC-Wizzard/generated v0.0.0-20260504181409-ba8146c868b5
       gopkg.in/yaml.v3 v3.0.1
   )
   ```
3. Move all files from `pkg/config/` to `config/`:
   - `git mv pkg/config/config_suite_test.go config/`
   - `git mv pkg/config/marshaller.go config/`
   - `git mv pkg/config/parser.go config/`
   - `git mv pkg/config/parser_test.go config/`
   - `git mv pkg/config/path_or_paths.go config/`
   - `git mv pkg/config/path_or_paths_test.go config/`
   - `git mv pkg/config/types.go config/`
   - `git mv pkg/config/validator.go config/`
   - `git mv pkg/config/validator_test.go config/`
4. Update import paths in config files:
   - `github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors` → `github.com/LarsArtmann/SQLC-Wizzard/core/apperrors`
5. Update import paths in ALL root module files that reference `pkg/config`:
   - `github.com/LarsArtmann/SQLC-Wizzard/pkg/config` → `github.com/LarsArtmann/SQLC-Wizzard/config`
6. Update `internal/testing/assertions.go` import: `pkg/config` → `config`
7. Remove `pkg/` directory if empty
8. Run `go mod tidy` in `config/`
9. Run `go build ./config/...`
10. Run `go test ./config/...`
11. Run `go mod tidy` in root
12. Run `go build ./...`
13. Run `go test ./...`

**Verification:**
```bash
# In config/
cd config && go build ./... && go test ./... && go vet ./... && cd ..

# In root
go mod tidy && go build ./... && go test ./... && go vet ./...
```

**Rollback:** `git revert HEAD`

---

### Task 5: Add `go.work` and Remove `replace` Directive

**Impact:** 4% → 64% (high leverage — enables proper multi-module development)
**Estimated effort:** 15 min
**Depends on:** Task 3, Task 4

**What:** Create `go.work` at repo root, remove `replace` directive from root `go.mod`.

**Steps:**
1. Create `go.work`:
   ```go
   go 1.26.2
   
   use (
       ./generated
       ./core
       ./config
       .
   )
   ```
2. Edit root `go.mod`: remove the `replace github.com/LarsArtmann/SQLC-Wizzard/generated => ./generated` line
3. Run `go work sync`
4. Run `go mod tidy` in each module directory
5. Run `go build ./...` at root (workspace mode)
6. Run `go test ./...` at root

**Verification:**
```bash
go work sync && go build ./... && go test ./... && go vet ./...
# Verify each module builds independently
cd core && go build ./... && cd ..
cd config && go build ./... && cd ..
cd generated && go build ./... && cd ..
```

**Rollback:** Delete `go.work`, restore `replace` directive — `git revert HEAD`

---

### Task 6: Update CI/CD Workflow

**Impact:** 20% → 80% (broad value — CI must work with new structure)
**Estimated effort:** 20 min
**Depends on:** Task 5

**What:** Update `.github/workflows/ci-cd.yml` to build/test per module and at workspace level.

**Steps:**
1. Read current CI workflow
2. Add per-module build/test steps:
   ```yaml
   - name: Build & Test generated
     run: cd generated && go build ./... && go test ./...
   - name: Build & Test core
     run: cd core && go build ./... && go test ./...
   - name: Build & Test config
     run: cd config && go build ./... && go test ./...
   - name: Build & Test root (workspace)
     run: go build ./... && go test ./...
   ```
3. Add `go.work sync` step
4. Update cache key to include go.work

**Verification:** Push branch, verify CI passes (or dry-run locally)

**Rollback:** `git revert HEAD`

---

### Task 7: Update Linter Configuration

**Impact:** 20% → 80% (broad value — linting must understand new structure)
**Estimated effort:** 20 min
**Depends on:** Task 5

**What:** Update `.golangci.yml` and `.go-arch-lint.yml` for new module paths.

**Steps:**
1. Update `.golangci.yml` depguard rules:
   - Remove rules about `internal/apperrors` → now `core/apperrors`
   - Update allowed imports for each package
2. Update `.go-arch-lint.yml`:
   - Add `core/` and `config/` as modules
   - Update dependency constraints
3. Run `golangci-lint run ./...` (if installed)
4. Fix any new lint findings

**Verification:**
```bash
go vet ./...
# golangci-lint run ./... (if available)
```

**Rollback:** `git revert HEAD`

---

### Task 8: Update Build Scripts & Documentation

**Impact:** 20% → 80% (broad value — developers must understand new structure)
**Estimated effort:** 20 min
**Depends on:** Task 5

**What:** Update `justfile`, `AGENTS.md`, `README.md`, and `Dockerfile` for new module structure.

**Steps:**
1. Update `justfile`:
   - `just build` should use `go work sync` before build
   - `just test` should test per-module then workspace
   - Add `just test-core`, `just test-config` targets
2. Update `AGENTS.md`:
   - Add module structure section
   - Update build/test commands
   - Document `go.work` usage
3. Update `README.md`:
   - Update project structure diagram
   - Mention multi-module setup
4. Update `Dockerfile`:
   - `COPY go.work go.work`
   - `COPY core/ core/`
   - `COPY config/ config/`
   - Update `go mod download` to workspace-aware

**Verification:**
```bash
just build && just test
```

**Rollback:** `git revert HEAD`

---

### Task 9: Verify Module Independence

**Impact:** Remaining (polish — proof that modularization works)
**Estimated effort:** 15 min
**Depends on:** Task 5

**What:** Prove each module can be built and tested in isolation.

**Steps:**
1. Verify `generated/` builds independently:
   ```bash
   cd generated && go build ./... && go test ./... && cd ..
   ```
2. Verify `core/` builds independently:
   ```bash
   cd core && go build ./... && go test ./... && cd ..
   ```
3. Verify `config/` builds independently:
   ```bash
   cd config && go build ./... && go test ./... && cd ..
   ```
4. Verify root builds with workspace:
   ```bash
   go build ./... && go test ./...
   ```
5. Verify `go mod tidy` is clean in each module (no changes)
6. Verify `go work sync` produces no changes
7. Run `go mod graph` and verify no unexpected cross-module deps

**Verification:** All builds pass. `go mod tidy` changes nothing in any module.

**Rollback:** N/A — read-only verification

---

### Task 10: Final Commit — Update `go.work.sum` and Clean Up

**Impact:** Remaining (polish)
**Estimated effort:** 10 min
**Depends on:** Task 9

**What:** Final cleanup commit ensuring everything is clean.

**Steps:**
1. Run `go work sync` one final time
2. Run `go mod tidy` in every module
3. Remove `pkg/` directory if it's now empty
4. Remove `1763227265_test_migration.*.sql` files from root (test artifacts that shouldn't be committed)
5. Verify no orphaned files
6. Final `go build ./...` and `go test ./...`
7. Commit all cleanup

**Verification:**
```bash
go build ./... && go test ./... && go vet ./... && go work sync
git status  # Should show clean state
```

**Rollback:** `git revert HEAD`

---

## Execution Order Diagram

```
Task 1: Fix generated module path ──────────────────────────┐
Task 2: Delete dead code (pkg/errors) ──────────── parallel  │
                                                              ▼
Task 3: Extract core/ ◄─────────────────────────── depends on Task 1
          │
          ▼
Task 4: Extract config/ ◄──────────────────────── depends on Task 3
          │
          ▼
Task 5: Add go.work + remove replace ◄─────────── depends on Task 3, 4
          │
    ┌─────┼─────────────────────┐
    ▼     ▼                     ▼
Task 6  Task 7               Task 8
 CI/CD   Linters            Docs/Build
    │     │                     │
    └─────┼─────────────────────┘
          ▼
Task 9: Verify independence ◄── depends on Task 5
          │
          ▼
Task 10: Final cleanup
```

## Commit Convention

Each task is one commit with this format:

```
refactor(modularization): <task description>

- <what changed>
- <why>
- <verification>

Part of modularization plan (docs/modularization/PROPOSAL.md)
```
