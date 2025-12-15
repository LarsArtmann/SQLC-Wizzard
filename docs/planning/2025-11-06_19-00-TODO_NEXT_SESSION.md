# TODO: Next Session Priorities

**Last Updated:** 2025-11-06
**Current Branch:** `claude/create-sqlc-wizard-011CUrcBQqi6inaWUqF2XMM8`

---

## 🔴 CRITICAL (Must Fix Before v0.1.0)

### 1. Add Comprehensive Ginkgo Tests (4 hours) ⚠️ HIGHEST PRIORITY

**Current:** 0% test coverage
**Target:** 80%+ coverage

**Checklist:**

- [ ] Bootstrap Ginkgo test suites (`ginkgo bootstrap`)
- [ ] Add config package tests
  - [ ] Parser tests (valid + invalid YAML)
  - [ ] Validator tests (all validation rules)
  - [ ] Marshaller round-trip tests
  - [ ] PathOrPaths tests (after implementing)
- [ ] Add template package tests
  - [ ] Microservice template generation (all 3 DBs)
  - [ ] Registry tests
  - [ ] Rule conversion tests
- [ ] Add wizard package tests
  - [ ] Mock huh forms
  - [ ] Test validation logic
- [ ] Add integration tests
  - [ ] Full init command flow
  - [ ] Validate command flow
- [ ] Achieve 80%+ coverage: `go test -cover ./...`

**Why Critical:** No tests = no confidence = not production-ready

---

### 2. Fix Type Safety: Replace `interface{}` (2 hours)

**Current:** Using `interface{}` for Queries/Schema
**Target:** Type-safe union type

**Checklist:**

- [ ] Create `pkg/config/path_or_paths.go`
- [ ] Implement `PathOrPaths` struct
- [ ] Implement `UnmarshalYAML` (handle string and []string)
- [ ] Implement `MarshalYAML`
- [ ] Add `Strings() []string` method
- [ ] Update `SQLConfig` to use `PathOrPaths`
- [ ] Update template generators
- [ ] Add tests for PathOrPaths
- [ ] Verify no `interface{}` remains: `grep -r "interface{}" pkg/config/`

**Why Critical:** Type safety is non-negotiable in Go

---

### 3. Fix Split Brain: SafetyRules (1 hour)

**Current:** Duplication between `SafetyRules` and `RuleConfig`
**Target:** Single source of truth

**Checklist:**

- [ ] Create `internal/domain/rule.go`
- [ ] Move `RuleConfig` to domain package
- [ ] Add `SafetyRules.ToRuleConfigs() []RuleConfig` method
- [ ] Update microservice template to use conversion
- [ ] Remove manual rule creation in template
- [ ] Add unit tests for conversion
- [ ] Verify no duplication: check git diff

**Why Critical:** Split brains cause bugs when definitions drift

---

### 4. Fix Split Brain: Features (1 hour)

**Current:** Duplication between `Features` and `GoGenConfig`
**Target:** Single source of truth

**Checklist:**

- [ ] Create `internal/domain/emit_options.go`
- [ ] Define canonical `EmitOptions` type
- [ ] Add `EmitOptions.ToGoGenConfig() GoGenConfig` method
- [ ] Update `TemplateData` to use `EmitOptions`
- [ ] Update wizard to set `EmitOptions`
- [ ] Update template generation
- [ ] Add unit tests
- [ ] Remove `Features` struct entirely

**Why Critical:** Can cause conflicting configs

---

### 5. Implement Error Package (1.5 hours)

**Current:** Empty `internal/errors/` directory
**Target:** Structured error types with codes

**Checklist:**

- [ ] Create `internal/errors/errors.go` with error codes
- [ ] Define `ErrorCode` enum (CONFIG_INVALID, WIZARD_FAILED, etc.)
- [ ] Create `WizardError` with context
- [ ] Create `ConfigError` with field path
- [ ] Create `TemplateError` with template name
- [ ] Add `Error()` method implementations
- [ ] Add `Unwrap()` for error chains
- [ ] Replace all `fmt.Errorf()` with typed errors
- [ ] Add error tests

**Why Critical:** Better error messages for users

---

### 6. Add Smart Constructors (1 hour)

**Current:** `type ProjectType string` - no validation
**Target:** Validated constructors

**Checklist:**

- [ ] Add `NewProjectType(s string) (ProjectType, error)`
- [ ] Add `NewDatabaseType(s string) (DatabaseType, error)`
- [ ] Add validation logic
- [ ] Update wizard to use constructors
- [ ] Update commands to use constructors
- [ ] Add tests for invalid inputs
- [ ] Make fields private: `type ProjectType struct { value string }`

**Why Critical:** Prevent invalid states at compile time

---

### 7. Remove Unused Dependencies (15 minutes)

**Current:** `samber/do` imported but not used
**Target:** Clean `go.mod`

**Checklist:**

- [ ] Remove `samber/do`: `go mod edit -droprequire github.com/samber/do/v2`
- [ ] Run `go mod tidy`
- [ ] Verify removal: `grep "samber/do" go.mod`
- [ ] Commit changes

**Why Critical:** Code bloat, supply chain security

---

## 🟡 HIGH PRIORITY (Quality Improvements)

### 8. Add CI/CD Workflows (1 hour)

**Checklist:**

- [ ] Create `.github/workflows/test.yml`
  - [ ] Run tests on PR
  - [ ] Check coverage
  - [ ] Upload coverage to Codecov
- [ ] Create `.github/workflows/lint.yml`
  - [ ] golangci-lint
  - [ ] go vet
  - [ ] go fmt check
- [ ] Create `.github/workflows/release.yml`
  - [ ] Build binaries
  - [ ] Create GitHub release
  - [ ] Attach artifacts

---

### 9. Add Linting Configuration (30 minutes)

**Checklist:**

- [ ] Create `.golangci.yml` with strict rules
- [ ] Add `go-arch-lint` rules in `.go-arch-lint.yml`
- [ ] Create `.pre-commit-config.yaml`
- [ ] Fix all linter warnings
- [ ] Add to CI/CD

---

### 10. Use Railway-Oriented Programming (2 hours)

**Checklist:**

- [ ] Update `Template.Generate()` → `mo.Result[*SqlcConfig, error]`
- [ ] Update parser functions → `mo.Result`
- [ ] Chain operations with `FlatMap`, `Map`
- [ ] Update commands to match on Results
- [ ] Add Result helper functions

---

## 🟢 MEDIUM PRIORITY (Features)

### 11. Split Large Files (45 minutes)

**Checklist:**

- [ ] Split `wizard.go` into `steps/*.go` files
  - [ ] `steps/project_type.go`
  - [ ] `steps/database.go`
  - [ ] `steps/features.go`
  - [ ] `steps/output.go`
- [ ] Keep each file under 300 lines

---

### 12. Add Hobby Template (1 hour)

**Checklist:**

- [ ] Create `internal/templates/hobby.go`
- [ ] SQLite-focused, minimal setup
- [ ] Single-file queries
- [ ] No complex features
- [ ] Add tests

---

### 13. Add Enterprise Template (2 hours)

**Checklist:**

- [ ] Create `internal/templates/enterprise.go`
- [ ] Multi-DB support
- [ ] Cloud integration
- [ ] Advanced features (all emit options)
- [ ] Add tests

---

### 14. Implement Detectors Package (2 hours)

**Checklist:**

- [ ] Create `internal/detectors/go_mod.go`
  - [ ] Detect package path from go.mod
- [ ] Create `internal/detectors/config.go`
  - [ ] Detect existing sqlc.yaml
- [ ] Create `internal/detectors/database.go`
  - [ ] Detect .sql, .db files
- [ ] Integrate into wizard for smart defaults
- [ ] Add tests

---

## 📊 SUCCESS METRICS

Before releasing v0.1.0:

- [x] MVP works end-to-end ✅
- [ ] Test coverage ≥ 80%
- [ ] Zero `interface{}` types
- [ ] Zero split brains
- [ ] Structured error handling
- [ ] CI/CD passing
- [ ] All linters passing
- [ ] Documentation complete

---

## 🚀 RELEASE CHECKLIST

When ready for v0.1.0:

- [ ] All P0 (Critical) items complete
- [ ] All tests passing
- [ ] Coverage ≥ 80%
- [ ] CI/CD working
- [ ] README updated
- [ ] CHANGELOG.md created
- [ ] Tag release: `git tag v0.1.0`
- [ ] Push tag: `git push origin v0.1.0`
- [ ] GitHub release with binaries

---

## ⏱️ TIME ESTIMATES

**Critical Items (P0):** ~10.5 hours
**High Priority (P1):** ~3.5 hours
**Medium Priority (P2):** ~5.5 hours

**Total to v0.1.0:** ~14 hours (2 work days)

---

## 📝 NOTES

- Start with tests (item #1) - everything else is easier with tests
- Fix type safety before adding features
- Don't add new features until P0 items are complete
- Keep commits small and focused
- Run tests before every commit
- Update this file as you complete items

---

**Good luck! 🚀**
