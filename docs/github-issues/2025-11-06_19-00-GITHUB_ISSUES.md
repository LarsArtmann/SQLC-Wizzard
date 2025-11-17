# SQLC-Wizard - Critical Issues and Improvements

**Date:** 2025-11-06
**Branch:** `claude/create-sqlc-wizard-011CUrcBQqi6inaWUqF2XMM8`
**Status:** MVP Complete, Quality Improvements Needed

---

## 🔴 CRITICAL ISSUES (P0 - Must Fix)

### Issue 1: Type Safety Violations - `interface{}` Usage
**Priority:** P0 - CRITICAL
**Impact:** High - Runtime errors, poor type safety
**Effort:** 1 hour

**Problem:**
```go
// pkg/config/types.go:23-24
Queries interface{} `yaml:"queries"` // Can be string or []string
Schema  interface{} `yaml:"schema"`  // Can be string or []string
```

Using `interface{}` destroys compile-time type safety and can cause runtime panics.

**Solution:**
Create a `PathOrPaths` custom type with proper YAML marshaling:
```go
type PathOrPaths struct {
    paths []string
}

func (p *PathOrPaths) UnmarshalYAML(value *yaml.Node) error {
    // Handle string or []string safely
}
```

**Files Affected:**
- `pkg/config/types.go`
- `pkg/config/path_or_paths.go` (new)
- `internal/templates/microservice.go`

---

### Issue 2: Split Brain - SafetyRules Duplication
**Priority:** P0 - CRITICAL
**Impact:** High - Bug-prone, maintenance nightmare
**Effort:** 45 minutes

**Problem:**
Same concept represented in two different ways:
```go
// internal/templates/types.go
type SafetyRules struct {
    NoSelectStar  bool
    RequireWhere  bool
}

// pkg/config/types.go
type RuleConfig struct {
    Name    string
    Rule    string
    Message string
}
```

If we add a rule to one, we must manually add to the other. Easy to forget!

**Solution:**
Single source of truth with conversion methods:
```go
func (sr SafetyRules) ToRuleConfigs() []RuleConfig {
    // Convert boolean flags to CEL rules
}
```

**Files Affected:**
- `internal/templates/types.go`
- `internal/templates/microservice.go`
- `pkg/config/types.go`

---

### Issue 3: Split Brain - Features vs GoGenConfig
**Priority:** P0 - CRITICAL
**Impact:** High - Conflicting configurations possible
**Effort:** 45 minutes

**Problem:**
```go
// internal/templates/types.go
type Features struct {
    EmitInterface   bool
    PreparedQueries bool
}

// pkg/config/types.go
type GoGenConfig struct {
    EmitInterface       bool
    EmitPreparedQueries bool
}
```

Which is the source of truth? Can they conflict? YES!

**Solution:**
Merge into single domain type:
```go
type EmitOptions struct {
    Interface      bool
    PreparedQueries bool
    JSONTags       bool
}

func (e EmitOptions) ToGoGenConfig() GoGenConfig {
    // Convert to config
}
```

---

### Issue 4: No Tests - ZERO Test Coverage
**Priority:** P0 - CRITICAL
**Impact:** High - No quality assurance
**Effort:** 6 hours

**Problem:**
```bash
$ find . -name "*_test.go" | wc -l
0
```

We have **ZERO tests** despite importing Ginkgo and promising BDD tests!

**Solution:**
Add comprehensive Ginkgo test suites:
1. Config package tests (parser, validator, marshaller)
2. Template tests (generation, registry)
3. Wizard tests (UI flows, validation)
4. Integration tests (end-to-end)

**Files to Create:**
- `pkg/config/*_test.go`
- `internal/templates/*_test.go`
- `internal/wizard/*_test.go`
- `test/integration/init_test.go`

---

### Issue 5: No Error Package - Empty Directory
**Priority:** P0 - CRITICAL
**Impact:** Medium - Poor error messages
**Effort:** 1.5 hours

**Problem:**
Created `internal/errors/` directory but never implemented it. All errors are just `fmt.Errorf(...)` strings.

**Solution:**
Implement structured error types:
```go
type ErrorCode string

const (
    ErrConfigInvalid ErrorCode = "CONFIG_INVALID"
    ErrWizardFailed  ErrorCode = "WIZARD_FAILED"
    ErrTemplateFailed ErrorCode = "TEMPLATE_FAILED"
)

type WizardError struct {
    Code    ErrorCode
    Message string
    Context map[string]interface{}
    Cause   error
}
```

---

### Issue 6: Weak Type Constructors - No Validation
**Priority:** P0 - CRITICAL
**Impact:** Medium - Invalid values possible
**Effort:** 1 hour

**Problem:**
```go
type ProjectType string
type DatabaseType string
```

These are just type aliases! Can pass any string:
```go
data.ProjectType = ProjectType("invalid-garbage") // Compiles! 💥
```

**Solution:**
Smart constructors with validation:
```go
func NewProjectType(s string) (ProjectType, error) {
    if !isValidProjectType(s) {
        return ProjectType{}, errors.ErrInvalidProjectType
    }
    return ProjectType{value: s}, nil
}
```

---

### Issue 7: Unused Dependencies - `samber/do`
**Priority:** P0 - Code Smell
**Impact:** Low - Bloat
**Effort:** 5 minutes

**Problem:**
```bash
$ grep "samber/do" go.mod
github.com/samber/do/v2 v2.0.0

$ grep -r "do\." . --include="*.go"
# NO RESULTS!
```

Imported but never used!

**Solution:**
Remove from `go.mod`:
```bash
go mod edit -droprequire github.com/samber/do/v2
go mod tidy
```

---

## ⚠️ HIGH PRIORITY (P1 - Should Fix)

### Issue 8: No Railway-Oriented Programming
**Priority:** P1 - HIGH
**Effort:** 2 hours

We imported `samber/mo` but don't use `Result[T, error]` anywhere!

**Solution:**
```go
func (t *Template) Generate(data TemplateData) mo.Result[*SqlcConfig, error] {
    return mo.TupleToResult(generateConfig(data))
}
```

---

### Issue 9: No CI/CD Workflows
**Priority:** P1 - HIGH
**Effort:** 1 hour

**Files to Create:**
- `.github/workflows/test.yml`
- `.github/workflows/lint.yml`
- `.github/workflows/release.yml`

---

### Issue 10: Large Files - wizard.go is 290 lines
**Priority:** P1 - HIGH
**Effort:** 45 minutes

Split into:
- `internal/wizard/steps/project_type.go`
- `internal/wizard/steps/database.go`
- `internal/wizard/steps/features.go`
- `internal/wizard/steps/output.go`

---

## 📝 MEDIUM PRIORITY (P2 - Nice to Have)

### Issue 11: Missing Templates
Only microservice template implemented. Need:
- Hobby (SQLite-focused)
- Enterprise (multi-DB)
- API-First (JSON-focused)

### Issue 12: No Detectors Package
Directory exists but empty. Should auto-detect:
- `go.mod` for package path
- Existing `sqlc.yaml`
- Database files

### Issue 13: No Architecture Linting
Should use `go-arch-lint` to enforce package dependencies.

---

## ✅ COMPLETED WORK

### ✅ CLI Foundation
- Cobra-based CLI with commands
- Version command with build-time injection

### ✅ Init Command
- Interactive mode with beautiful TUI
- Non-interactive mode with flags
- Works end-to-end

### ✅ Validate Command
- Configuration validation
- Best practice warnings
- Colored output

### ✅ Config Package
- YAML parsing/marshalling
- Validation with errors/warnings
- Supports sqlc v2 schema

### ✅ Template System
- Microservice template (PostgreSQL, MySQL, SQLite)
- Smart defaults per database
- CEL-based safety rules

### ✅ Generators
- sqlc.yaml generation
- Example SQL queries (CRUD)
- Example schema (users table)

---

## 📊 METRICS

**Lines of Code:** ~2,500
**Test Coverage:** 0% ❌
**Type Safety Score:** 6/10 ⚠️
**Documentation:** 5/10 ⚠️
**Architecture Quality:** 7/10 ⚠️

---

## 🎯 RECOMMENDED NEXT STEPS

1. **Fix type safety** (interface{} → proper types)
2. **Add tests** (Ginkgo BDD suites)
3. **Implement error package**
4. **Fix split brains**
5. **Add CI/CD**
6. **Release v0.1.0**

---

## 📅 TIMELINE ESTIMATE

**Critical Fixes (P0):** 8 hours
**Quality Improvements (P1):** 6 hours
**Feature Additions (P2):** 8 hours
**Total:** ~22 hours

---

## 🔗 BRANCH

All work on: `claude/create-sqlc-wizard-011CUrcBQqi6inaWUqF2XMM8`

**Last Push:** 2025-11-06
**Commits:** 11 commits
**Status:** All pushed ✅
