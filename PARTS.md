# SQLC-Wizzard: Component Analysis & Library Extraction Potential

**Last Updated:** 2026-02-28
**Project:** SQLC-Wizzard - Interactive CLI wizard for generating sqlc configurations
**Purpose:** Identify components that could be extracted as reusable libraries/SDKs

---

## Executive Summary

After analyzing the SQLC-Wizzard codebase, **6 components** were identified with varying extraction potential:

| Component                            | Extraction Potential | Priority |
| ------------------------------------ | -------------------- | -------- |
| Type-Safe Enums & Smart Constructors | **HIGH**             | 1        |
| Template Registry & Strategy Pattern | **HIGH**             | 2        |
| Structured Error Package             | **MEDIUM**           | 3        |
| YAML Configuration Parser            | **LOW**              | 4        |
| Step-Based Wizard Framework          | **LOW**              | 5        |
| Adapter Interfaces                   | **SKIP**             | -        |

---

## Component 1: Type-Safe Enums & Smart Constructors

### Location

- `internal/domain/safety_policy.go` (345 lines)
- `internal/domain/emit_modes.go` (309 lines)
- `internal/templates/types.go` (143 lines)
- `generated/types.go` (232 lines)

### Description

A sophisticated type-safe enum system with:

- String-based enums with `IsValid()` validation
- Semantic methods that return booleans (e.g., `IsForbidden()`, `IsRequired()`)
- Smart constructors that prevent invalid states at compile time
- Clear distinction between technical and semantic representations

### Key Types

```go
// Safety Policies
SelectStarPolicy      // allowed, warning, forbidden
WhereClauseRequirement // optional, recommended, required
LimitClauseRequirement // optional, recommended, required
DestructiveOperationPolicy // allowed, warning, forbidden, blocked

// Emit Modes
NullHandlingMode      // nullable, pointer, value
EnumGenerationMode    // string, int, uuid
StructPointerMode     // pointer, value
JSONTagStyle          // snake, camel, pascal
```

### Dependencies

- `github.com/samber/lo` (utility functions)
- `github.com/samber/mo` (functional options)
- Generated types from TypeSpec

### Alternatives

| Library                                                                                 | Stars | Description                 | Limitation                                      |
| --------------------------------------------------------------------------------------- | ----- | --------------------------- | ----------------------------------------------- |
| [abice/go-enum](https://github.com/abice/go-enum)                                       | 900+  | Code generator for enums    | No semantic methods, just string/int conversion |
| [alvaroloes/enumer](https://github.com/alvaroloes/enumer)                               | 800+  | Enum code generator         | Basic validation only, no semantic helpers      |
| [golang/go](https://go.googlesource.com/proposal/+/master/design/43651/type-parameters) | N/A   | Go 1.18+ generics           | No built-in enum support                        |
| [hitzgg/go-enums](https://github.com/hitzgg/go-enums)                                   | 100+  | Generic enum implementation | Basic, no semantic layer                        |

### Our Unique Value

1. **Semantic Methods**: Unlike alternatives that only provide `String()` and `Parse()`, our enums have:

   ```go
   // Instead of: if policy == "forbidden" { ... }
   // We provide: if policy.IsForbidden() { ... }
   ```

2. **Domain-Specific Validation**: Enums are validated against allowed values with clear error messages

3. **Migration-Safe**: Both old boolean flags and new type-safe enums coexist during migration

4. **Self-Documenting**: Method names clearly express intent (`IsForbidden()` vs `== "forbidden"`)

5. **Pattern**: String-based internal representation with type-safe external API

### Recommendation: **EXTRACT (High Priority)**

Create a reusable library: `github.com/larsartmann/go-semantic-enum`

**Scope:**

- Generic enum generator with semantic method support
- Configuration-driven semantic methods
- Validation framework
- Migration helpers for boolean → enum transitions

---

## Component 2: Template Registry & Strategy Pattern

### Location

- `internal/templates/registry.go` (75 lines)
- `internal/templates/types.go` (143 lines)
- `internal/templates/*.go` (8 template implementations)

### Description

A clean implementation of:

- **Strategy Pattern**: Each template implements a common interface
- **Registry Pattern**: Central registration and discovery
- **Factory Pattern**: Create templates by name/type

### Key Interfaces

```go
type Template interface {
    Generate(data generated.TemplateData) (*config.SqlcConfig, error)
    DefaultData() generated.TemplateData
    RequiredFeatures() []string
    Name() string
    Description() string
}

type Registry interface {
    Register(template Template)
    Get(name string) (Template, error)
    List() []Template
    HasTemplate(name string) bool
}
```

### Dependencies

- Domain types only (no external dependencies)
- `github.com/samber/lo` (optional, for utilities)

### Alternatives

| Library                                                       | Stars | Description          | Limitation                           |
| ------------------------------------------------------------- | ----- | -------------------- | ------------------------------------ |
| [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) | 5000+ | Plugin system        | Overkill for in-process strategies   |
| [uber-go/dig](https://github.com/uber-go/dig)                 | 4000+ | Dependency injection | Focused on DI, not strategy pattern  |
| [samber/do](https://github.com/samber/do)                     | 2000+ | DI framework         | Focused on DI, not registry/strategy |
| Manual implementation                                         | N/A   | Interface + map      | No discovery, no metadata            |

### Our Unique Value

1. **Self-Describing Templates**: Templates provide their own metadata (`Name()`, `Description()`)

2. **Feature Requirements**: Templates declare what features they need (`RequiredFeatures()`)

3. **Default Data Generation**: Each template provides sensible defaults (`DefaultData()`)

4. **Type-Safe Registration**: Compile-time interface enforcement

5. **Discovery-Friendly**: `List()` enables dynamic UI generation

### Recommendation: **EXTRACT (High Priority)**

Create a reusable library: `github.com/larsartmann/go-strategy-registry`

**Scope:**

- Generic strategy interface definitions
- Thread-safe registry implementation
- Metadata support (name, description, features)
- Discovery and listing utilities
- Factory pattern helpers

---

## Component 3: Structured Error Package

### Location

- `internal/apperrors/error_types.go` (215 lines)

### Description

A comprehensive structured error system with:

- Typed error codes (string-based enum)
- Error severity levels (info, warning, error, critical)
- Rich error details (field, value, expected, actual, message)
- Context support (component, request_id, user_id, timestamp)
- Retryable flag for transient errors
- Smart constructors with format string validation

### Key Types

```go
type ErrorCode string        // Typed error identifier
type ErrorSeverity string    // info, warning, error, critical

type ErrorDetails struct {
    Field    string
    Value    any
    Expected string
    Actual   string
    Message  string
}

type Error struct {
    Code       ErrorCode
    Message    string
    Severity   ErrorSeverity
    Details    []ErrorDetails
    Component  string
    RequestID  string
    UserID     string
    Retryable  bool
    Timestamp  time.Time
    Cause      error
}
```

### Dependencies

- Standard library only
- No external dependencies

### Alternatives

| Library                                                       | Stars   | Description                 | Limitation                          |
| ------------------------------------------------------------- | ------- | --------------------------- | ----------------------------------- |
| [cockroachdb/errors](https://github.com/cockroachdb/errors)   | 1500+   | Error wrapping with context | No structured details, no severity  |
| [pkg/errors](https://github.com/pkg/errors)                   | 8000+   | Stack traces                | Deprecated, no structure            |
| [larsartmann/uniflow](https://github.com/larsartmann/uniflow) | Unknown | Error handling              | Check if available                  |
| [rotisserie/eris](https://github.com/rotisserie/eris)         | 1500+   | Error handling with context | Different approach, less structured |
| [go-errors/errors](https://github.com/go-errors/errors)       | 800+    | Stack traces                | Basic wrapping only                 |

### Our Unique Value

1. **Structured Details**: Field-level error information for validation feedback

2. **Severity Levels**: Built-in classification (info/warning/error/critical)

3. **Rich Context**: Request ID, user ID, component tracking for observability

4. **Retryable Flag**: Explicit marking for transient vs permanent errors

5. **Format String Validation**: Smart constructors catch format string errors

6. **Error Code System**: Typed codes enable programmatic error handling

### Recommendation: **EXTRACT (Medium Priority)**

Create a reusable library: `github.com/larsartmann/go-structured-errors`

**Scope:**

- Core error types and interfaces
- Smart constructors
- Severity classification
- Error details builder
- JSON/YAML serialization for APIs
- Integration with logging frameworks

**Note:** Check `larsartmann/uniflow` first - if it covers this, extend it instead of creating a new library.

---

## Component 4: YAML Configuration Parser

### Location

- `pkg/config/types.go` (161 lines)

### Description

A type-safe YAML configuration system with:

- Complete `sqlc.yaml` v2 schema representation
- Custom YAML unmarshaling (`PathOrPaths` for flexible input)
- Parser, Validator, Marshaller pattern
- Deep validation of nested structures

### Key Types

```go
type SqlcConfig struct {
    Version string         `yaml:"version"`
    SQL     []SQLConfig    `yaml:"sql"`
}

type SQLConfig struct {
    Engine  string         `yaml:"engine"`
    Schema  PathOrPaths    `yaml:"schema"`  // Flexible: string or []string
    Queries PathOrPaths    `yaml:"queries"`
    Gen     GenConfig      `yaml:"gen"`
    Rules   []RuleConfig   `yaml:"rules"`
}
```

### Dependencies

- `gopkg.in/yaml.v3` (YAML parsing)
- Standard library

### Alternatives

| Library                                         | Stars  | Description                  | Limitation                         |
| ----------------------------------------------- | ------ | ---------------------------- | ---------------------------------- |
| [knadh/koanf](https://github.com/knadh/koanf)   | 3000+  | Configuration management     | More complex, different philosophy |
| [spf13/viper](https://github.com/spf13/viper)   | 27000+ | Configuration with env/flags | BANNED per HOW_TO_GOLANG.md        |
| [go-yaml/yaml](https://github.com/go-yaml/yaml) | 2000+  | YAML parsing                 | We already use this                |

### Our Unique Value

1. **PathOrPaths**: Flexible YAML unmarshaling for single or multiple values

2. **Domain-Specific**: Tailored for sqlc.yaml schema

3. **Validation-Ready**: Types designed for validation, not just parsing

### Recommendation: **KEEP INTERNAL (Low Priority)**

**Reason:** This is highly domain-specific to sqlc.yaml. The only reusable pattern is `PathOrPaths`, which is a small utility. The value of extracting it is minimal compared to the maintenance overhead.

**Alternative:** Extract `PathOrPaths` as a small utility if needed elsewhere.

---

## Component 5: Step-Based Wizard Framework

### Location

- `internal/wizard/wizard.go` (272 lines)

### Description

A step-based wizard framework with:

- StepInterface for composable wizard steps
- Dependency injection for testing
- Built-in validation at each step
- State management across steps

### Key Interfaces

```go
type StepInterface interface {
    Execute() error
    Validate() error
    Name() string
}

type WizardDependencies struct {
    Input  io.Reader
    Output io.Writer
    // ... testable dependencies
}
```

### Dependencies

- `github.com/charmbracelet/huh` (TUI library)
- `github.com/charmbracelet/bubbletea` (TUI framework)

### Alternatives

| Library                                                               | Stars  | Description       | Limitation                 |
| --------------------------------------------------------------------- | ------ | ----------------- | -------------------------- |
| [charmbracelet/huh](https://github.com/charmbracelet/huh)             | 4000+  | Interactive forms | We already use this        |
| [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) | 25000+ | TUI framework     | We already use this        |
| [AlecAivazis/survey](https://github.com/AlecAivazis/survey)           | 4000+  | CLI surveys       | Less modern, huh is better |

### Our Unique Value

1. **Step Pattern**: Composable, reusable steps

2. **Testable Dependencies**: DI for unit testing

3. **Validation Integration**: Built-in step validation

### Recommendation: **KEEP INTERNAL (Low Priority)**

**Reason:** This is tightly coupled to `charmbracelet/huh`. The step pattern is useful but not substantial enough for a separate library. The value is in the specific wizard implementation, not the framework.

**Consider:** If we build more wizards, extract a minimal step framework.

---

## Component 6: Adapter Interfaces (Hexagonal Architecture)

### Location

- `internal/adapters/interfaces.go` (117 lines)

### Description

Adapter interfaces for hexagonal architecture:

- `SQLCAdapter` - sqlc CLI integration
- `DatabaseAdapter` - Database connection/inspection
- `CLIAdapter` - CLI output/formatting
- `TemplateAdapter` - Template rendering
- `FileSystemAdapter` - File operations

### Dependencies

- Context from standard library
- Domain types

### Alternatives

| Library                                     | Stars | Description           | Limitation                  |
| ------------------------------------------- | ----- | --------------------- | --------------------------- |
| [samber/do](https://github.com/samber/do)   | 2000+ | DI framework          | Different purpose           |
| [uber-go/fx](https://github.com/uber-go/fx) | 6000+ | Application framework | BANNED per HOW_TO_GOLANG.md |

### Our Unique Value

1. **Domain-Specific**: Interfaces designed for our specific use case

2. **Context-Aware**: All methods accept context

3. **Testable**: Easy mocking for tests

### Recommendation: **SKIP (Do Not Extract)**

**Reason:** These interfaces are highly domain-specific. Hexagonal architecture ports should be defined per-application, not shared. Extracting them would create artificial constraints.

---

## Implementation Priorities

### Phase 1: High-Value Extractions

1. **go-semantic-enum** (Component 1)
   - Highest value: Unique pattern not available elsewhere
   - Effort: Medium (needs generalization)
   - Impact: Reusable across all Go projects

2. **go-strategy-registry** (Component 2)
   - High value: Clean, reusable pattern
   - Effort: Low (well-isolated already)
   - Impact: Useful for any plugin/strategy system

### Phase 2: Medium-Value Extractions

3. **go-structured-errors** (Component 3)
   - Check `larsartmann/uniflow` first
   - If new library needed, coordinate with uniflow design
   - Effort: Medium
   - Impact: Better error handling across projects

### Phase 3: Keep Internal

4. **YAML Configuration Parser** (Component 4)
   - Keep internal, too domain-specific
   - Consider extracting `PathOrPaths` utility

5. **Step-Based Wizard Framework** (Component 5)
   - Keep internal, tightly coupled to huh
   - Revisit if we build more wizards

6. **Adapter Interfaces** (Component 6)
   - Do not extract, define per-application

---

## Next Steps

1. [ ] Create `go-semantic-enum` repository
2. [ ] Create `go-strategy-registry` repository
3. [ ] Evaluate `larsartmann/uniflow` for error handling
4. [ ] Update HOW_TO_GOLANG.md with new library references
5. [ ] Migrate SQLC-Wizzard to use extracted libraries

---

## References

- [HOW_TO_GOLANG.md](../library-policy/HOW_TO_GOLANG.md) - Project standards
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Project architecture
- [generated/types.go](./generated/types.go) - TypeSpec-generated types
