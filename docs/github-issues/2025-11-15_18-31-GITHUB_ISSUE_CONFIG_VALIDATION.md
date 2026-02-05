# 🔧 SQLC Config Validation, Extension & Improvement

**Priority:** HIGH (#2)  
**Complexity:** MEDIUM-HIGH  
**Estimated Time:** 2-3 days  
**Impact:** HIGH - Core configuration management

---

## 🎯 **Problem Statement**

SQLC configuration (`sqlc.yaml`) is complex and users struggle with:

- Configuration validation errors (hard to debug)
- Missing optimal settings
- Version compatibility issues
- Database-specific optimizations
- Advanced feature discovery

**This creates friction and prevents users from getting optimal SQLC performance.**

---

## 🎪 **Solution Vision**

Create intelligent SQLC configuration assistant:

```bash
# Validate existing configuration
sqlc-wizard config validate --file sqlc.yaml

# Fix common issues automatically
sqlc-wizard config fix --file sqlc.yaml --auto

# Extend with optimizations
sqlc-wizard config optimize --file sqlc.yaml --performance

# Generate configuration wizard
sqlc-wizard config init --database postgresql --type microservice

# Show missing features/optimizations
sqlc-wizard config suggest --file sqlc.yaml
```

---

## 🏗️ **Technical Implementation**

### **Phase 1: Validation Engine**

```go
type ConfigValidator struct {
    rules     []ValidationRule
    schema    *ConfigSchema
    version   string
}

type ValidationIssue struct {
    Severity   string  // ERROR, WARNING, INFO
    Message    string
    Field      string
    Suggestion string
    AutoFix    bool
}

type ValidationRule interface {
    Validate(config *SqlcConfig) []ValidationIssue
    Fix(config *SqlcConfig) error
}
```

### **Phase 2: Optimization Engine**

```go
type ConfigOptimizer struct {
    database  DatabaseType
    project   ProjectType
    rules     []OptimizationRule
}

type OptimizationRule interface {
    Analyze(config *SqlcConfig) []OptimizationSuggestion
    Apply(config *SqlcConfig) error
}

type OptimizationSuggestion struct {
    Impact     string  // HIGH, MEDIUM, LOW
    Description string
    Change     ConfigChange
    Reason     string
}
```

### **Phase 3: Extension System**

```go
type ConfigExtension interface {
    Name() string
    AppliesTo(config *SqlcConfig) bool
    Extend(config *SqlcConfig) error
}

// Database-specific extensions
type PostgreSQLExtension struct{}
type MySQLExtension struct{}
type SQLiteExtension struct{}
```

---

## 🧪 **Example Usage & Output**

### **Validation Example**

```bash
sqlc-wizard config validate --file sqlc.yaml

🔍 SQLC Configuration Validation Report
========================================

❌ ERROR: Missing 'version' field
   Impact: sqlc cannot determine compatibility
   Suggestion: Add 'version: "2"' or 'version: "1.x"'
   Auto-fix: ✅ Available

⚠️  WARNING: Missing 'output' configuration
   Impact: Generated files will go to unexpected locations
   Suggestion: Add 'output: "internal/db"'
   Auto-fix: ✅ Available

💡 INFO: No performance optimizations detected
   Impact: May miss SQLC performance features
   Suggestion: Consider adding emit_json_tags, emit_prepared_queries
   Auto-fix: ✅ Available
```

### **Optimization Example**

```bash
sqlc-wizard config optimize --file sqlc.yaml --performance

🚀 Configuration Optimization Report
===================================

HIGH: PostgreSQL performance settings not optimized
   Current: Basic sqlc.yaml
   Suggested: Add PostgreSQL-specific settings
   Impact: 15-25% query performance improvement
   Changes:
     - add emit_json_tags: true
     - add emit_prepared_queries: true
     - add emit_exact_table_names: false
     - add emit_empty_slices: true
     - add overrides:
         - db: public
         - go_type: "types.UUID"
           column: "id"

MEDIUM: Missing development optimizations
   Impact: Slower development experience
   Suggested: Add development-friendly settings
   Changes:
     - add emit_method_with_db_argument: true
     - add emit_pointers_for_null_types: true
```

### **Auto-Fix Example**

```bash
sqlc-wizard config fix --file sqlc.yaml --auto

🔧 Auto-Fixing Configuration Issues
==================================

✅ Fixed missing 'version' field (set to "2")
✅ Added 'output' directory (set to "internal/db")
✅ Added performance optimizations for PostgreSQL
✅ Added development-friendly settings
✅ Validated configuration passes sqlc check

Configuration fixed and validated successfully! ✅
```

---

## 🎯 **Core Features**

### **1. Comprehensive Validation**

```go
type ValidationRules struct {
    SyntaxValidation    // YAML syntax checking
    RequiredFields      // Required sqlc.yaml fields
    TypeValidation      // Field type checking
    VersionCompat       // SQLC version compatibility
    DatabaseValidation  // Database-specific rules
    PathValidation     // File/directory path validation
}
```

### **2. Performance Optimization**

```go
type OptimizationRules struct {
    DatabaseOptimizations // PostgreSQL, MySQL, SQLite specific
    QueryOptimizations   // Query generation optimizations
    OutputOptimizations   // Generated code optimizations
    DevelopmentOptimizations // Developer experience improvements
}
```

### **3. Configuration Extensions**

```go
type Extensions struct {
    DatabaseExtensions   // Database-specific settings
    FrameworkExtensions  // Framework integrations
    PerformanceExtensions // Performance tuning
    DevelopmentExtensions // Developer tools
}
```

---

## 📁 **Configuration Templates**

### **By Database Type**

```yaml
# PostgreSQL optimized
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries/"
    schema: "db/schema/"
    gen:
      go:
        package: "db"
        out: "internal/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: true
        emit_exact_table_names: false
```

### **By Project Type**

```yaml
# Microservice optimized
sql:
  - gen:
      go:
        emit_pointers_for_null_types: true
        emit_method_with_db_argument: true
        emit_empty_slices: true
        overrides:
          - column: "created_at"
            go_type: "time.Time"
```

### **By Use Case**

```yaml
# API development optimized
sql:
  - gen:
      go:
        emit_json_tags: true
        emit_db_tags: true
        emit_interface: true
```

---

## 🛠️ **Implementation Plan**

### **Day 1: Validation Engine**

- [ ] Create `ConfigValidator` core structure
- [ ] Implement basic validation rules
- [ ] Add YAML parsing and error reporting
- [ ] Create CLI validation command

### **Day 2: Optimization Engine**

- [ ] Implement `ConfigOptimizer` structure
- [ ] Add database-specific optimization rules
- [ ] Create performance analysis logic
- [ ] Add auto-fix capabilities

### **Day 3: Extension System**

- [ ] Create extension interface and registry
- [ ] Implement database-specific extensions
- [ ] Add project-type specific extensions
- [ ] Create template system

---

## 🔧 **Technical Architecture**

### **New Components**

```go
// internal/commands/config.go
func NewConfigCommand() *cobra.Command

// internal/config/
type Validator interface {}
type Optimizer interface {}
type Extension interface {}

// internal/rules/
type ValidationRule interface {}
type OptimizationRule interface {}
```

### **Integration Points**

- Extend `pkg/config/` with validation logic
- Use `generated/` types for optimization rules
- Leverage `internal/adapters/` for file operations
- Add to main CLI command registry

---

## 🎯 **Acceptance Criteria**

### **Validation Features**

- [ ] Detects all common sqlc.yaml configuration errors
- [ ] Provides clear, actionable error messages
- [ ] Auto-fixes common configuration issues
- [ ] Validates SQLC version compatibility

### **Optimization Features**

- [ ] Suggests database-specific optimizations
- [ ] Provides performance impact estimates
- [ ] Applies optimizations automatically with user consent
- [ ] Validates optimized configurations

### **Extension Features**

- [ ] Supports PostgreSQL, MySQL, SQLite specific settings
- [ ] Provides project-type specific customizations
- [ ] Creates configuration templates from existing projects
- [ ] Allows custom rule/extension definitions

### **User Experience**

- [ ] Commands work without requiring SQLC expertise
- [ ] Provides clear explanations for all suggestions
- [ ] Allows selective application of optimizations
- [ ] Creates backup before making changes

---

## 📊 **Success Metrics**

### **Configuration Quality**

- [ ] Reduction in configuration errors (> 80%)
- [ ] Improvement in user-generated config quality
- [ ] Adoption of optimization suggestions (> 70%)

### **User Experience**

- [ ] Time to valid configuration < 1 minute
- [ ] User satisfaction rate > 4.5/5
- [ ] Support ticket reduction for config issues

### **Performance Impact**

- [ ] Measurable query performance improvements
- [ ] Development experience improvements
- [ ] Production deployment success rate

---

## 🎯 **Why HIGH PRIORITY (#2)**

This feature is **HIGHLY VALUABLE** because:

1. **Solves Core Pain Point:** Configuration is biggest SQLC friction point
2. **Improves Performance:** Optimized configs significantly improve performance
3. **Reduces Support Burden:** Auto-fixes common configuration issues
4. **Enables Advanced Features:** Helps users discover SQLC capabilities
5. **Improves Success Rate:** Better configs lead to more successful projects

**This feature will make SQLC dramatically easier and more performant for all users.**

---

## 📋 **Definition of Done**

- [ ] `sqlc-wizard config validate` command implemented
- [ ] `sqlc-wizard config fix` command with auto-fix capabilities
- [ ] `sqlc-wizard config optimize` command with suggestions
- [ ] `sqlc-wizard config suggest` command with missing features
- [ ] Support for PostgreSQL, MySQL, SQLite optimizations
- [ ] Configuration templates for all project types
- [ ] Comprehensive test coverage (> 90%)
- [ ] Documentation with examples and troubleshooting

---

**This issue will transform SQLC configuration from a complex chore into an intelligent, automated experience!** 🔧✨

---

_Created: 2025-11-15_  
_Priority: HIGH (#2)_  
_Ready for implementation_ 🎯
