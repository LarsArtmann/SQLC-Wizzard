# Template Customization Guide

This guide explains how to customize templates for your specific needs.

---

## Overview

SQLC-Wizard templates provide sensible defaults, but you often need to customize them for your project requirements. This guide covers all customization options with examples.

---

## Quick Reference

| Customization | Files Affected | Complexity | Impact |
|----------------|-----------------|------------|--------|
| **Database Engine** | TemplateData | Low | High (SQLite → PostgreSQL → MySQL) |
| **Database URL** | TemplateData | Low | Medium (dev → staging → production) |
| **Package Name** | TemplateData | Low | Medium (db → mypackage) |
| **Output Paths** | TemplateData | Low | Low (directory structure) |
| **Features** | TemplateData | Medium | High (add UUID, JSON, arrays, etc.) |
| **Code Generation Options** | TemplateData | Medium | High (prepared queries, interfaces, etc.) |
| **Validation Rules** | TemplateData | High | High (strict vs flexible) |
| **Type Overrides** | TemplateData | Medium | High (custom Go types for DB types) |
| **Rename Rules** | TemplateData | Low | Medium (snake_case → camelCase) |
| **Safety Rules** | TemplateData | High | High (allow/disallow queries) |
| **Build Tags** | TemplateData | Low | Medium (postgres, pgx, etc.) |

---

## Customization Options

### 1. Database Engine

Change the database engine used by template.

**Supported Engines:**
- `postgres` (PostgreSQL)
- `mysql` (MySQL)
- `sqlite` (SQLite)

**When to Customize:**
- Switching from SQLite (Hobby/Testing) to PostgreSQL (Enterprise/API First)
- Switching between PostgreSQL and MySQL
- Changing database vendor

**Example:**

```go
// Start with Hobby template (SQLite)
hobbyTemplate := templates.NewHobbyTemplate()
data := hobbyTemplate.DefaultData()

// Customize: Change to PostgreSQL
data.Database.Engine = templates.DatabaseTypePostgreSQL

// Customize: Change to MySQL
data.Database.Engine = templates.DatabaseTypeMySQL
```

**Generated Configuration:**

```yaml
# SQLite (Hobby template default)
engine: sqlite

# PostgreSQL (after customization)
engine: postgresql

# MySQL (after customization)
engine: mysql
```

**Impact:**
- ✅ Switches database driver
- ✅ Updates type overrides (different per database)
- ✅ Updates build tags (different per database)
- ✅ Requires appropriate database connection string
- ❌ May require code changes (different SQL dialects)

**Considerations:**
1. **SQL Dialect:** PostgreSQL, MySQL, and SQLite have different SQL syntax
2. **Type Support:** Not all features available in all databases (e.g., PostgreSQL arrays vs MySQL JSON)
3. **Driver Compatibility:** Ensure your Go driver matches database
4. **Testing:** Test with both databases if you might switch

---

### 2. Database URL

Change the database connection string.

**When to Customize:**
- Development vs. staging vs. production databases
- Different environments (local, Docker, Kubernetes)
- Using connection pooling (custom URL parameters)
- Using managed database services (Heroku, AWS RDS, Google Cloud SQL)

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Local development (SQLite)
data.Database.URL = "file:dev.db"

// Customize: Local development (PostgreSQL)
data.Database.URL = "postgres://user:password@localhost:5432/mydb?sslmode=disable"

// Customize: Staging environment
data.Database.URL = "postgres://user:password@staging-db:5432/mydb?sslmode=require"

// Customize: Production environment (from environment variable)
data.Database.URL = os.Getenv("DATABASE_URL")

// Customize: Production (AWS RDS PostgreSQL)
data.Database.URL = "postgres://user:password@mydb.cxxxxx.us-west-2.rds.amazonaws.com:5432/mydb?sslmode=require"

// Customize: Production (Google Cloud SQL)
data.Database.URL = "postgres://user:password@/cloudsql/project-id:region:mydb"

// Customize: Production (Heroku Postgres)
data.Database.URL = os.Getenv("DATABASE_URL")  // Heroku sets this

// Customize: Production with connection pooling
data.Database.URL = "postgres://user:password@prod-db:5432/mydb?pool_max_conns=20&pool_min_conns=5&pool_max_conn_lifetime=1h"
```

**Generated Configuration:**

```yaml
# Simple SQLite URL
database:
  uri: file:dev.db

# PostgreSQL URL with SSL
database:
  uri: ${DATABASE_URL}

# PostgreSQL URL with connection pooling
database:
  uri: postgres://user:password@prod-db:5432/mydb?pool_max_conns=20&pool_min_conns=5&pool_max_conn_lifetime=1h

# AWS RDS PostgreSQL URL
database:
  uri: postgres://user:password@mydb.cxxxxx.us-west-2.rds.amazonaws.com:5432/mydb?sslmode=require
```

**Impact:**
- ✅ Switches database environment
- ✅ Supports environment variables (recommended)
- ✅ Supports connection pooling parameters
- ✅ Works with managed database services
- ❌ Requires valid connection string
- ❌ Database must be accessible

**Best Practices:**
1. **Use Environment Variables** → Never hardcode database URLs in code
2. **Use Connection Pooling** → Add pool parameters to URL
3. **Use SSL in Production** → Always use `sslmode=require` or `sslmode=verify-full`
4. **Use Managed Services** → Let Heroku/AWS/Google manage database complexity
5. **Keep Dev URLs Simple** → Use `localhost` and `sslmode=disable` for local development

**Environment-Specific URLs:**

```go
// Development
devURL := "postgres://user:password@localhost:5432/devdb?sslmode=disable"

// Staging
stagingURL := "postgres://user:password@staging-db:5432/stagedb?sslmode=require"

// Production
prodURL := os.Getenv("DATABASE_URL")  // From environment variable

// Switch based on environment
env := os.Getenv("APP_ENV")
var dbURL string
switch env {
case "development":
    dbURL = devURL
case "staging":
    dbURL = stagingURL
case "production":
    dbURL = prodURL
default:
    dbURL = devURL
}
data.Database.URL = dbURL
```

---

### 3. Package Name

Change the Go package name for generated code.

**When to Customize:**
- Changing default `db` package to something more specific
- Multiple databases in same project (avoid naming conflicts)
- Library development (package name is public API)

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Change package name
data.Package.Name = "mymodels"

// Customize: Use project name as package
data.Package.Name = "myproject/db"
```

**Generated Configuration:**

```yaml
# Default (db)
gen:
  go:
    package: db

# Customized (mymodels)
gen:
  go:
    package: mymodels
```

**Impact:**
- ✅ Changes import path in generated code
- ✅ Allows multiple databases in same project
- ✅ Provides public API name for libraries
- ❌ Requires updating import paths in application code

**Best Practices:**
1. **Use Descriptive Names** → `users`, `products`, `orders` not `models`, `db`, `data`
2. **Avoid Conflicts** → Use unique package names if multiple databases
3. **Use Nested Packages** → `myproject/db` is better than `mydb`
4. **Keep It Simple** → Don't make package names too long or complex

**Example with Multiple Databases:**

```go
// Primary database
primaryData := templates.NewHobbyTemplate().DefaultData()
primaryData.Package.Name = "myproject/primary"
primaryData.Database.URL = os.Getenv("PRIMARY_DB_URL")

// Analytics database
analyticsData := templates.NewAnalyticsTemplate().DefaultData()
analyticsData.Package.Name = "myproject/analytics"
analyticsData.Database.URL = os.Getenv("ANALYTICS_DB_URL")
```

---

### 4. Output Paths

Change where generated code is written.

**When to Customize:**
- Different directory structure (internal/db → db)
- Shared output directory for multiple templates
- Monorepo setup (packages/)
- Testing output (testdata/)

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Base directory
data.Output.BaseDir = "internal/db"

// Customize: Queries directory
data.Output.QueriesDir = "internal/db/queries"

// Customize: Schema directory
data.Output.SchemaDir = "internal/db/schema"

// Customize: Monorepo structure
data.Output.BaseDir = "packages/db"
data.Output.QueriesDir = "packages/db/queries"
data.Output.SchemaDir = "packages/db/schema"

// Customize: Testing paths
data.Output.BaseDir = "testdata/db"
data.Output.QueriesDir = "testdata/db/queries"
data.Output.SchemaDir = "testdata/db/schema"
```

**Generated Configuration:**

```yaml
# Default (db/structure)
output:
  base_dir: db
  queries_dir: db/queries
  schema_dir: db/schema

# Internal structure (recommended)
output:
  base_dir: internal/db
  queries_dir: internal/db/queries
  schema_dir: internal/db/schema

# Monorepo structure
output:
  base_dir: packages/db
  queries_dir: packages/db/queries
  schema_dir: packages/db/schema

# Testing structure
output:
  base_dir: testdata/db
  queries_dir: testdata/db/queries
  schema_dir: testdata/db/schema
```

**Impact:**
- ✅ Changes where generated code is written
- ✅ Allows custom directory structures
- ✅ Supports monorepos and multi-package projects
- ❌ Requires updating import paths in application code
- ❌ May require adjusting Go module paths

**Best Practices:**
1. **Use Internal Structure** → `internal/db` is Go best practice
2. **Separate Directories** → `internal/db` is better than flat `db`
3. **Testing Paths** → Use `testdata/` for test fixtures
4. **Monorepo Structure** → Use `packages/` for shared packages
5. **Keep It Consistent** → Don't mix structures across project

**Directory Structure Examples:**

```bash
# Go standard structure
internal/
  db/
    models.go      # Generated
    querier.go    # Generated
    queries/       # Generated
    schema/        # Generated

# Monorepo structure
packages/
  db/
    models.go      # Generated
    queries/       # Generated
  analytics/
    models.go      # Generated
    queries/       # Generated

# Testing structure
testdata/
  db/
    schema/        # Test schemas
    fixtures/      # Test data
  db/
    models.go      # Generated
    queries/       # Generated
```

---

### 5. Features (Database Features)

Enable or disable database features (UUID, JSON, arrays, full-text search, etc.).

**Available Features:**
- `UseUUIDs` - Use UUID types for primary keys
- `UseJSON` - Use JSONB (PostgreSQL) or JSON (MySQL) types
- `UseArrays` - Use array types (PostgreSQL)
- `UseFullText` - Use full-text search types (PostgreSQL)
- `UseManaged` - Use managed database mode

**When to Customize:**
- Adding UUID support for distributed systems
- Adding JSON document storage
- Adding array support for many-to-many relationships
- Adding full-text search for analytics
- Combining features for advanced use cases

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Enable UUID support
data.Database.UseUUIDs = true

// Customize: Enable JSON support
data.Database.UseJSON = true

// Customize: Enable array support
data.Database.UseArrays = true

// Customize: Enable full-text search
data.Database.UseFullText = true

// Customize: Enable ALL features
data.Database.UseUUIDs = true
data.Database.UseJSON = true
data.Database.UseArrays = true
data.Database.UseFullText = true
```

**Generated Configuration:**

```yaml
# No features (Hobby template default)
database:
  use_uuids: false
  use_json: false
  use_arrays: false
  use_full_text: false

# UUID enabled
database:
  use_uuids: true

# JSON enabled
database:
  use_json: true

# ALL features enabled (Enterprise template default)
database:
  use_uuids: true
  use_json: true
  use_arrays: true
  use_full_text: true
```

**Impact:**
- ✅ Enables advanced database features
- ✅ Adds appropriate type overrides (UUID → uuid.UUID, etc.)
- ✅ Enables use of complex data types
- ❌ Requires adding Go imports (uuid, encoding/json)
- ❌ Increases generated code complexity
- ❌ May require schema changes (add UUID columns, etc.)

**Best Practices:**
1. **Enable Features in Pairs** → UUID + JSON, Arrays + Full-Text
2. **Test Feature Combinations** → Not all combinations may work
3. **Document Feature Usage** → Keep track of which features are enabled
4. **Consider Migration Path** → Adding features to existing tables is harder

**Feature Combinations:**

```go
// Simple (Hobby template)
data.Database.UseUUIDs = false
data.Database.UseJSON = false
data.Database.UseArrays = false
data.Database.UseFullText = false

// API-focused (Microservice template)
data.Database.UseUUIDs = true
data.Database.UseJSON = true
data.Database.UseArrays = false
data.Database.UseFullText = false

// Enterprise (all features)
data.Database.UseUUIDs = true
data.Database.UseJSON = true
data.Database.UseArrays = true
data.Database.UseFullText = true

// Analytics (data-heavy)
data.Database.UseUUIDs = false
data.Database.UseJSON = true
data.Database.UseArrays = true
data.Database.UseFullText = true
```

**Schema Changes for Features:**

```sql
-- Adding UUID support to existing table
ALTER TABLE users ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE users DROP COLUMN id;

-- Adding JSON support to existing table
ALTER TABLE users ADD COLUMN metadata JSONB;

-- Adding array support to existing table
ALTER TABLE users ADD COLUMN tags TEXT[];
```

---

### 6. Code Generation Options

Customize how sqlc generates Go code.

**Available Options:**
- `EmitJSONTags` - Generate JSON tags on structs
- `EmitPreparedQueries` - Generate prepared query methods
- `EmitInterface` - Generate Querier interface
- `EmitEmptySlices` - Emit `omitempty` for empty slices
- `EmitResultStructPointers` - Use pointers for result structs
- `EmitParamsStructPointers` - Use pointers for params structs
- `EmitEnumValidMethod` - Generate `Valid()` method for enums
- `EmitAllEnumValues` - Generate all enum value constants
- `JSONTagsCaseStyle` - JSON tag naming (snake, camel, pascal)
- `BuildTags` - Build tags for conditional compilation

**When to Customize:**
- Enabling/disabling JSON tags (API vs internal use)
- Enabling/disabling prepared queries (performance vs flexibility)
- Enabling/disabling interface generation (mocking vs simplicity)
- Changing JSON tag style (snake_case vs camelCase)
- Adding enum validation (type safety)
- Adding build tags (conditional compilation)

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Enable JSON tags (for APIs)
data.Validation.EmitOptions.EmitJSONTags = true

// Customize: Disable JSON tags (for internal use)
data.Validation.EmitOptions.EmitJSONTags = false

// Customize: Enable prepared queries (performance)
data.Validation.EmitOptions.EmitPreparedQueries = true

// Customize: Disable prepared queries (flexibility)
data.Validation.EmitOptions.EmitPreparedQueries = false

// Customize: Enable interface generation (mocking)
data.Validation.EmitOptions.EmitInterface = true

// Customize: Disable interface generation (simplicity)
data.Validation.EmitOptions.EmitInterface = false

// Customize: Enable pointer result structs (memory efficiency)
data.Validation.EmitOptions.EmitResultStructPointers = true

// Customize: Enable pointer params structs (nullable fields)
data.Validation.EmitOptions.EmitParamsStructPointers = true

// Customize: Enable enum validation (type safety)
data.Validation.EmitOptions.EmitEnumValidMethod = true

// Customize: Enable all enum values (flexibility)
data.Validation.EmitOptions.EmitAllEnumValues = true

// Customize: Set JSON tag style
data.Validation.EmitOptions.JSONTagsCaseStyle = "camel"  // camelCase
data.Validation.EmitOptions.JSONTagsCaseStyle = "snake"  // snake_case
data.Validation.EmitOptions.JSONTagsCaseStyle = "pascal"  // PascalCase
```

**Generated Configuration:**

```yaml
# Minimal (Hobby template default)
gen:
  go:
    emit_json_tags: false
    emit_prepared_queries: false
    emit_interface: false
    emit_empty_slices: true
    emit_result_struct_pointers: false
    emit_params_struct_pointers: false

# API-optimized (API First template default)
gen:
  go:
    emit_json_tags: true
    emit_prepared_queries: true
    emit_interface: true
    emit_empty_slices: true
    json_tags_case_style: camel

# Enterprise (all features, Enterprise template default)
gen:
  go:
    emit_json_tags: true
    emit_prepared_queries: true
    emit_interface: true
    emit_empty_slices: true
    emit_result_struct_pointers: true
    emit_params_struct_pointers: true
    emit_enum_valid_method: true
    emit_all_enum_values: true
    json_tags_case_style: camel
```

**Impact:**
- ✅ Changes generated code structure
- ✅ Enables/disables advanced features
- ✅ Changes JSON serialization behavior
- ✅ Changes memory usage (pointers vs values)
- ✅ Changes testing patterns (interfaces vs no interfaces)
- ❌ May require application code changes
- ❌ More complex generated code with all features enabled

**Best Practices:**
1. **APIs Use JSON Tags** → Enable `emit_json_tags` for public APIs
2. **Internal Code May Not** → Disable if JSON not needed
3. **Prepared Queries are Safer** → Always enable in production
4. **Interfaces Enable Mocking** → Enable `emit_interface` for easier testing
5. **Pointer Types Save Memory** → Enable for large structs
6. **Enum Validation Prevents Bugs** → Enable `emit_enum_valid_method`

**Option Recommendations:**

```go
// For Production APIs
data.Validation.EmitOptions.EmitJSONTags = true
data.Validation.EmitOptions.EmitPreparedQueries = true
data.Validation.EmitOptions.EmitInterface = true
data.Validation.EmitOptions.JSONTagsCaseStyle = "camel"

// For Internal Services
data.Validation.EmitOptions.EmitJSONTags = false
data.Validation.EmitOptions.EmitPreparedQueries = false
data.Validation.EmitOptions.EmitInterface = false

// For Type Safety
data.Validation.EmitOptions.EmitResultStructPointers = true
data.Validation.EmitOptions.EmitParamsStructPointers = true
data.Validation.EmitOptions.EmitEnumValidMethod = true
data.Validation.EmitOptions.EmitAllEnumValues = true

// For Simplicity
data.Validation.EmitOptions.EmitInterface = false
data.Validation.EmitOptions.EmitResultStructPointers = false
data.Validation.EmitOptions.EmitParamsStructPointers = false
```

---

### 7. Validation Rules

Customize database validation settings.

**Available Rules:**
- `StrictFunctions` - Enable strict function validation
- `StrictOrderBy` - Enable strict ORDER BY validation
- `NoSelectStar` - Disallow `SELECT *`
- `RequireWhere` - Require WHERE clause
- `NoDropTable` - Disallow DROP TABLE
- `NoTruncate` - Disallow TRUNCATE TABLE
- `RequireLimit` - Require LIMIT clause

**When to Customize:**
- Enabling production-level validation (Enterprise template)
- Enabling strict mode for data integrity
- Enabling safety rules to prevent dangerous queries
- Disabling validation for development flexibility

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Enable strict function checks
data.Validation.StrictFunctions = true

// Customize: Enable strict ORDER BY checks
data.Validation.StrictOrderBy = true

// Customize: Enable all safety rules (production)
data.Validation.SafetyRules.NoSelectStar = true
data.Validation.SafetyRules.RequireWhere = true
data.Validation.SafetyRules.NoDropTable = true
data.Validation.SafetyRules.NoTruncate = true
data.Validation.SafetyRules.RequireLimit = true
```

**Generated Configuration:**

```yaml
# Relaxed (Hobby template default)
strict_function_checks: false
strict_order_by: false
rules: []

# Strict (Enterprise template default)
strict_function_checks: true
strict_order_by: true
rules:
  - name: no-select-star
    rule: SELECT * is not allowed
  - name: require-where
    rule: WHERE clause is required
  - name: no-drop-table
    rule: DROP TABLE is not allowed
  - name: no-truncate
    rule: TRUNCATE TABLE is not allowed
  - name: require-limit
    rule: LIMIT clause is required
```

**Impact:**
- ✅ Enforces data integrity
- ✅ Prevents dangerous queries
- ✅ Improves code quality
- ✅ Fail-fast on bad queries
- ❌ More restrictive (harder to write queries)
- ❌ May break existing queries (development vs production)

**Best Practices:**
1. **Enable Strict Mode in Production** → All safety rules enabled
2. **Disable Strict Mode in Development** → Allow any query for flexibility
3. **Use Environment-Based Rules** → Different rules for dev vs prod
4. **Review Failing Queries** → Understand why queries fail validation
5. **Gradual Migration** → Start relaxed, add rules as code matures

**Environment-Based Validation:**

```go
// Development - relaxed rules
if os.Getenv("APP_ENV") == "development" {
    data.Validation.StrictFunctions = false
    data.Validation.StrictOrderBy = false
    data.Validation.SafetyRules.NoSelectStar = false
} else {
    // Production - strict rules
    data.Validation.StrictFunctions = true
    data.Validation.StrictOrderBy = true
    data.Validation.SafetyRules.NoSelectStar = true
    data.Validation.SafetyRules.RequireWhere = true
}
```

---

### 8. Type Overrides

Customize how database types map to Go types.

**When to Customize:**
- Adding custom type for specific database type
- Using different Go library for a type
- Mapping JSON types to custom structs
- Adding nullable behavior for types

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Add custom type override for UUID (using v5)
data.Database.TypeOverrides = []config.Override{
    {
        DBType:       "uuid",
        GoType:       "uuid.UUID",
        GoImportPath: "github.com/google/uuid/v5",
    },
}

// Customize: Add custom type for JSON (using custom struct)
data.Database.TypeOverrides = []config.Override{
    {
        DBType:       "jsonb",
        GoType:       "CustomJSON",
        GoImportPath: "myproject/types",
    },
}

// Customize: Add nullable behavior
data.Database.TypeOverrides = []config.Override{
    {
        DBType:       "int",
        GoType:       "sql.NullInt64",
        GoImportPath: "database/sql",
        Nullable:      true,
    },
}
```

**Generated Configuration:**

```yaml
# Default UUID override
overrides:
  - db_type: uuid
    go_type: UUID
    go_import_path: github.com/google/uuid

# Custom UUID override (v5)
overrides:
  - db_type: uuid
    go_type: uuid.UUID
    go_import_path: github.com/google/uuid/v5

# Custom JSON override (using custom struct)
overrides:
  - db_type: jsonb
    go_type: CustomJSON
    go_import_path: myproject/types

# Nullable int override
overrides:
  - db_type: int
    go_type: sql.NullInt64
    nullable: true
```

**Impact:**
- ✅ Allows custom Go types for database columns
- ✅ Supports nullable behavior
- ✅ Can use alternative Go libraries
- ✅ Enables custom JSON structures
- ❌ Requires importing custom types
- ❌ May break existing code
- ❌ Less common patterns (harder to understand)

**Best Practices:**
1. **Use Standard Types** → Prefer `uuid.UUID`, `json.RawMessage` over custom
2. **Document Custom Types** → Clearly document why you're overriding defaults
3. **Test Custom Types** → Verify they work with sqlc and Go
4. **Consider Backward Compatibility** → Ensure changes don't break existing code
5. **Use Nullable for Optional Fields** → Better than zero values

---

### 9. Rename Rules

Customize how database columns are renamed to Go fields.

**When to Customize:**
- Changing field naming convention (snake_case → camelCase)
- Adding custom rename rules for specific columns
- Adding common rename rules (id → ID, url → URL, etc.)

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Add custom rename rule
data.Database.RenameRules = map[string]string{
    "first_name": "FirstName",
    "last_name": "LastName",
}

// Customize: Remove common rename rule
data.Database.RenameRules = map[string]string{}

// Customize: Override default rename rules
data.Database.RenameRules = map[string]string{
    "id": "id",  // Don't rename to ID
    "url": "url", // Don't rename to URL
}
```

**Generated Configuration:**

```yaml
# Default rename rules (BaseTemplate)
rename:
  id: ID
  uuid: UUID
  url: URL
  uri: URI
  json: JSON
  api: API
  http: HTTP
  db: DB
  otp: OTP

# Custom rename rules
rename:
  id: ID
  uuid: UUID
  url: URL
  first_name: FirstName
  last_name: LastName

# No rename rules (use database names)
rename: {}
```

**Impact:**
- ✅ Changes Go field names
- ✅ Enforces naming conventions
- ✅ Prevents common naming issues (id → ID)
- ✅ Supports API naming (camelCase vs snake_case)
- ❌ Changes generated code field names
- ❌ May break existing application code
- ❌ Requires consistency across project

**Best Practices:**
1. **Follow Go Naming** → Use `ID`, `URL`, `UUID` for common fields
2. **API Naming** → Use camelCase for JSON fields
3. **Internal Naming** → Use snake_case for internal code
4. **Document Your Naming** → Keep rename rules documented
5. **Be Consistent** → Use same naming across all templates

---

### 10. Build Tags

Customize Go build tags for conditional compilation.

**When to Customize:**
- Different builds for different databases
- Conditional compilation of features
- Testing with different build tags
- Production optimizations

**Examples:**

```go
// Start with any template
data := template.DefaultData()

// Customize: Add build tags for PostgreSQL
data.Database.BuildTags = "postgres,pgx"

// Customize: Add build tags for MySQL
data.Database.BuildTags = "mysql"

// Customize: Add build tags for SQLite
data.Database.BuildTags = "sqlite"

// Customize: Add build tags for testing
data.Database.BuildTags = "test,mysql"

// Customize: Add build tags for production optimizations
data.Database.BuildTags = "prod,postgres,pgx,optimizations"
```

**Generated Configuration:**

```yaml
# Default build tags (PostgreSQL)
build_tags: postgres,pgx

# MySQL build tags
build_tags: mysql

# SQLite build tags
build_tags: sqlite

# Testing build tags
build_tags: test,mysql
```

**Impact:**
- ✅ Enables conditional compilation
- ✅ Allows different builds for different databases
- ✅ Supports testing with build tags
- ✅ Can include custom build tags for features
- ❌ Requires Go build tag support
- ❌ More complex build process
- ❌ May require Makefile or build script

**Best Practices:**
1. **Use Database-Specific Tags** → `postgres`, `mysql`, `sqlite`
2. **Use Standard Tags** → `test`, `prod`, `debug`
3. **Document Build Tags** → Clearly document what each tag does
4. **Test All Build Combinations** → Ensure code compiles with different tags
5. **Use Makefile** → Manage multiple build configurations

**Makefile Example:**

```makefile
# Build all database drivers
.PHONY: all postgres mysql sqlite all

postgres:
    go build -tags postgres,pgx

mysql:
    go build -tags mysql

sqlite:
    go build -tags sqlite

all:
    go build -tags "postgres,pgx mysql sqlite"

# Build with optimizations
.PHONY: optimized
optimized:
    go build -tags "prod,postgres,pgx,optimizations"

# Run tests with different tags
.PHONY: test test-all
test:
    go test -tags postgres

test-all:
    go test -tags "postgres mysql sqlite"
```

---

## Advanced Customization Patterns

### Pattern 1: Starting from One Template and Customizing

Start with the template closest to your needs, then customize:

```go
// Step 1: Choose base template
template := templates.NewHobbyTemplate()

// Step 2: Get defaults
data := template.DefaultData()

// Step 3: Customize for your needs
data.Database.Engine = templates.DatabaseTypePostgreSQL
data.Database.UseUUIDs = true
data.Database.UseJSON = true
data.Validation.EmitOptions.EmitJSONTags = true
data.Validation.EmitOptions.EmitPreparedQueries = true
data.Validation.EmitOptions.EmitInterface = true

// Step 4: Generate config
config, err := template.Generate(data)
```

**Pros:**
- ✅ Faster than starting from scratch
- ✅ Uses proven defaults
- ✅ Incremental customization
- ✅ Easier to understand full configuration

**Cons:**
- ❌ May carry unnecessary features from base template
- ❌ May need to disable some defaults
- ❌ Less flexibility than building from scratch

### Pattern 2: Combining Features from Multiple Templates

Start with one template, then add features from others:

```go
// Step 1: Start with Hobby template (simple)
data := templates.NewHobbyTemplate().DefaultData()

// Step 2: Add Enterprise features
data.Database.UseUUIDs = true
data.Database.UseJSON = true
data.Database.UseArrays = true
data.Database.UseFullText = true

// Step 3: Add Microservice code gen options
data.Validation.EmitOptions.EmitPreparedQueries = true
data.Validation.EmitOptions.EmitInterface = true

// Step 4: Add API First JSON options
data.Validation.EmitOptions.JSONTagsCaseStyle = "camel"

// Step 5: Generate config
config, err := templates.NewHobbyTemplate().Generate(data)
```

**Pros:**
- ✅ Mix and match features from different templates
- ✅ Customize exactly what you need
- ✅ Leverage multiple template patterns
- ✅ Create hybrid solutions

**Cons:**
- ❌ Requires understanding of all templates
- ❌ May create inconsistent combinations
- ❌ Harder to debug if issues arise

### Pattern 3: Environment-Specific Customization

Different configurations for different environments:

```go
// Development configuration
func GetDevConfig() generated.TemplateData {
    data := templates.NewHobbyTemplate().DefaultData()
    data.Database.URL = "file:dev.db"
    data.Validation.StrictFunctions = false
    return data
}

// Staging configuration
func GetStagingConfig() generated.TemplateData {
    data := templates.NewMicroserviceTemplate().DefaultData()
    data.Database.URL = "postgres://user:password@staging-db:5432/mydb?sslmode=require"
    data.Validation.StrictFunctions = true
    return data
}

// Production configuration
func GetProdConfig() generated.TemplateData {
    data := templates.NewEnterpriseTemplate().DefaultData()
    data.Database.URL = os.Getenv("DATABASE_URL")
    data.Validation.StrictFunctions = true
    data.Validation.StrictOrderBy = true
    data.Validation.SafetyRules.NoSelectStar = true
    data.Validation.SafetyRules.RequireWhere = true
    return data
}

// Choose config based on environment
env := os.Getenv("APP_ENV")
var config generated.TemplateData
switch env {
case "development":
    config = GetDevConfig()
case "staging":
    config = GetStagingConfig()
case "production":
    config = GetProdConfig()
default:
    config = GetDevConfig()  // Default to dev
}

// Generate config
sqlcConfig, err := config.Generate(config)
```

**Pros:**
- ✅ Separate configs for each environment
- ✅ Production uses environment variables
- ✅ Development is flexible
- ✅ Easy to switch between environments
- ✅ Prevents production config in code

**Cons:**
- ❌ More complex code
- ❌ Requires managing multiple configs
- ❌ Easy to misconfigure wrong environment

### Pattern 4: Progressive Customization

Start simple, add features incrementally:

```go
// Iteration 1: Start with Hobby template (simple)
data := templates.NewHobbyTemplate().DefaultData()

// Iteration 2: Add PostgreSQL support
data.Database.Engine = templates.DatabaseTypePostgreSQL

// Iteration 3: Add UUID support
data.Database.UseUUIDs = true

// Iteration 4: Add JSON support
data.Database.UseJSON = true

// Iteration 5: Enable prepared queries
data.Validation.EmitOptions.EmitPreparedQueries = true

// Iteration 6: Enable interface generation
data.Validation.EmitOptions.EmitInterface = true

// Test at each iteration
config, err := templates.NewHobbyTemplate().Generate(data)
if err != nil {
    log.Fatal(err)
}

// Verify config works
if config.SQL[0].Gen.Go.SQLPackage != "pgx/v5" {
    log.Fatal("SQL package incorrect")
}
```

**Pros:**
- ✅ Test each change incrementally
- ✅ Easier to debug if issues
- ✅ Can stop at any point
- ✅ Clear progression from simple to complex
- ✅ Learn template customization step-by-step

**Cons:**
- ❌ Takes longer to reach final config
- ❌ More iterations to manage
- ❌ May leave unfinished customization

### Pattern 5: Template Mixins

Create reusable customization functions:

```go
// mixin.go - Reusable customization functions
package mixin

// AddEnterpriseFeatures adds Enterprise template features to any template
func AddEnterpriseFeatures(data *generated.TemplateData) {
    data.Database.UseUUIDs = true
    data.Database.UseJSON = true
    data.Database.UseArrays = true
    data.Database.UseFullText = true
}

// AddAPISupport adds API template features to any template
func AddAPISupport(data *generated.TemplateData) {
    data.Validation.EmitOptions.EmitJSONTags = true
    data.Validation.EmitOptions.EmitPreparedQueries = true
    data.Validation.EmitOptions.EmitInterface = true
    data.Validation.EmitOptions.JSONTagsCaseStyle = "camel"
}

// AddStrictValidation adds strict validation to any template
func AddStrictValidation(data *generated.TemplateData) {
    data.Validation.StrictFunctions = true
    data.Validation.StrictOrderBy = true
    data.Validation.SafetyRules.NoSelectStar = true
    data.Validation.SafetyRules.RequireWhere = true
}

// main.go - Using mixins
func main() {
    // Start with Hobby template
    data := templates.NewHobbyTemplate().DefaultData()

    // Apply mixins
    mixin.AddEnterpriseFeatures(&data)
    mixin.AddAPISupport(&data)
    mixin.AddStrictValidation(&data)

    // Generate config
    config, err := templates.NewHobbyTemplate().Generate(data)
    if err != nil {
        log.Fatal(err)
    }
}
```

**Pros:**
- ✅ Reusable customization functions
- ✅ Clear separation of concerns
- ✅ Easy to compose features
- ✅ Maintainable mixin library
- ✅ Mix and match as needed

**Cons:**
- ❌ More complex code
- ❌ Requires understanding of all templates
- ❌ May create inconsistent combinations
- ❌ Harder to debug if issues arise

---

## Testing Your Customizations

### Unit Testing Generated Code

```go
// test/db_test.go
package db_test

import (
    "testing"
    "context"
    "github.com/stretchr/testify/assert"
    "myproject/db"
)

func TestCustomizedConfig(t *testing.T) {
    // Create custom config
    data := templates.NewHobbyTemplate().DefaultData()
    data.Database.UseUUIDs = true
    data.Database.UseJSON = true

    // Generate config
    config, err := templates.NewHobbyTemplate().Generate(data)
    assert.NoError(t, err)

    // Verify customization worked
    assert.Equal(t, "pgx/v5", config.SQL[0].Gen.Go.SQLPackage)
    assert.True(t, config.SQL[0].Gen.Go.JSONTags)

    // Verify type overrides
    hasUUIDOverride := false
    for _, override := range config.SQL[0].Gen.Go.Overrides {
        if override.DBType == "uuid" {
            hasUUIDOverride = true
            assert.Equal(t, "uuid.UUID", override.GoType)
            assert.Equal(t, "github.com/google/uuid", override.GoImportPath)
        }
    }
    assert.True(t, hasUUIDOverride, "UUID override should exist")
}
```

### Integration Testing Full Workflow

```go
// test/integration_test.go
package integration_test

import (
    "testing"
    "context"
    "os"
    "path/filepath"
    "github.com/stretchr/testify/require"
    "myproject/templates"
)

func TestCustomizationWorkflow(t *testing.T) {
    ctx := context.Background()

    // Step 1: Choose template
    template := templates.NewHobbyTemplate()

    // Step 2: Customize
    data := template.DefaultData()
    data.Database.UseUUIDs = true
    data.Database.UseJSON = true
    data.Database.URL = "file::memory:"  // Use in-memory database for testing

    // Step 3: Generate config
    config, err := template.Generate(data)
    require.NoError(t, err)

    // Step 4: Write config to file
    configPath := filepath.Join(t.TempDir(), "sqlc.yaml")
    err = config.Write(configPath)
    require.NoError(t, err)

    // Step 5: Generate Go code
    // (In real scenario, would run `sqlc generate`)

    // Step 6: Verify generated code
    require.FileExists(t, filepath.Join("db", "models.go"))
    require.FileExists(t, filepath.Join("db", "querier.go"))
}
```

---

## Troubleshooting

### Issue: Customizations Not Applied

**Symptoms:**
- Generated config doesn't show customizations
- Type overrides missing
- Rename rules not working

**Solutions:**
1. **Check Template Implementation** → Ensure template uses custom data fields
2. **Verify DefaultData()** → Check if template overrides your customizations
3. **Check Generate() Method** → Ensure template uses `data` correctly
4. **Test with Logging** → Add debug logging to see what template receives

```go
// Add logging to see what template receives
data := template.DefaultData()
log.Printf("Before customization: UseUUIDs=%v", data.Database.UseUUIDs)

data.Database.UseUUIDs = true
log.Printf("After customization: UseUUIDs=%v", data.Database.UseUUIDs)

config, err := template.Generate(data)
if err != nil {
    log.Printf("Generate failed: %v", err)
    log.Printf("Data was: %+v", data)
}
```

### Issue: Generated Code Won't Compile

**Symptoms:**
- Import errors for custom types
- Type mismatches
- Missing imports

**Solutions:**
1. **Check Type Overrides** → Ensure Go types exist and are imported
2. **Check Build Tags** → Ensure build tags are valid
3. **Check Import Paths** → Ensure Go module paths are correct
4. **Verify SQL Package** → Ensure sqlc package matches your database driver

```bash
# Check for compile errors
go build ./...

# Get verbose output
go build -v ./...
```

### Issue: Database Connection Fails

**Symptoms:**
- Connection refused
- Authentication failed
- Database not found

**Solutions:**
1. **Check Database URL** → Ensure connection string is correct
2. **Check Database Accessibility** → Ensure database is running and accessible
3. **Check SSL Settings** → Use `sslmode=disable` for local, `sslmode=require` for production
4. **Check Environment Variables** → Ensure they're set correctly
5. **Check Network** → Ensure network allows connections

```bash
# Test database connection
psql $DATABASE_URL  # PostgreSQL
mysql $DATABASE_URL   # MySQL
sqlite3 file:test.db  # SQLite

# Check if database is running
docker ps | grep postgres  # Check if PostgreSQL is running
```

### Issue: Type Overrides Not Working

**Symptoms:**
- Generated code uses wrong types
- Custom types not applied
- Default types used instead

**Solutions:**
1. **Check DBType Matching** → Ensure database type matches (uuid vs UUID vs text vs TEXT)
2. **Check Case Sensitivity** → Ensure DBType case matches exactly
3. **Check Template Implementation** → Ensure template applies type overrides
4. **Verify Generate() Method** → Ensure template uses `data.Database.TypeOverrides`

```go
// Debug type overrides
data := template.DefaultData()
data.Database.TypeOverrides = []config.Override{
    {
        DBType:       "uuid",  // Must match database column type exactly
        GoType:       "uuid.UUID",
        GoImportPath: "github.com/google/uuid",
    },
}

config, err := template.Generate(data)
if err != nil {
    log.Fatal(err)
}

// Check generated config
for _, override := range config.SQL[0].Gen.Go.Overrides {
    log.Printf("Override: %s → %s", override.DBType, override.GoType)
}
```

---

## Best Practices Summary

### 1. Start Simple, Add Gradually

Begin with the simplest template (Hobby/Testing), then add features incrementally.

### 2. Test Each Change

After each customization, generate config and verify it works before moving to the next change.

### 3. Use Environment Variables

Never hardcode database URLs or sensitive data in your configuration.

### 4. Document Your Customizations

Keep a file or comment block documenting all customizations you've made.

### 5. Leverage Existing Patterns

Mix and match features from multiple templates rather than building from scratch.

### 6. Follow Go Conventions

Use standard Go naming, directory structures, and best practices.

### 7. Test Generated Code

Always test the generated code to ensure it compiles and works as expected.

### 8. Version Control Your Config

Commit your sqlc.yaml files to git to track changes over time.

### 9. Use Production-Ready Features in Production

Enable strict validation, safety rules, and all production features.

### 10. Keep Development Flexible

Keep validation relaxed in development for faster iteration.

---

## Next Steps

1. Choose a template based on your needs (see [Template Usage Guide](./usage.md))
2. Review template defaults in [Template Comparison Matrix](./comparison.md)
3. Start with the closest template (see examples in [examples/](./examples/))
4. Customize incrementally (one change at a time)
5. Test after each customization
6. Document your customizations
7. Generate your Go code
8. Test your application thoroughly

---

## Additional Resources

- [Template Usage Guide](./usage.md) - Detailed guide for each template
- [Template Comparison Matrix](./comparison.md) - Compare features across all templates
- [Template Examples](./examples/) - Real-world examples for each template
- [sqlc Documentation](https://docs.sqlc.dev/) - Official sqlc documentation

---

**Last Updated:** 2026-02-05
