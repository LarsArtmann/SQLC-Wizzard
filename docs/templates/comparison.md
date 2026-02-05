# Template Comparison Matrix

This matrix provides side-by-side comparison of all 8 templates to help you choose the right one for your project.

---

## Quick Comparison Table

| Feature | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|----------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **Complexity** | Simple | Medium | High | Medium | Medium | Medium | Medium | Medium |
| **Database** | SQLite | PostgreSQL | PostgreSQL | PostgreSQL | PostgreSQL | PostgreSQL | PostgreSQL |
| **SQL Driver** | database/sql | pgx/v5 | pgx/v5 | pgx/v5 | pgx/v5 | pgx/v5 | pgx/v5 |
| **Output Path** | db | internal/db | internal/db | internal/db | internal/analytics | testdata/db | internal/db | internal/db |
| **Query Path** | db/queries | internal/db/queries | internal/db/queries | internal/db/queries | internal/analytics/queries | testdata/db/queries | internal/db/queries |
| **Schema Path** | db/schema | internal/db/schema | internal/db/schema | internal/db/schema | internal/analytics/schema | testdata/db/schema | internal/db/schema | internal/db/schema |
| **Default DB URL** | file:dev.db | ${DATABASE_URL} | ${DATABASE_URL} | ${DATABASE_URL} | ${ANALYTICS_DATABASE_URL} | file:testdata/test.db | ${DATABASE_URL} | ${DATABASE_URL} |

---

## Feature Comparison Matrix

### Core Features

| Feature | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|----------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **UUID Support** | ❌ | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ |
| **JSON Support** | ❌ | ✅ | ✅ | ✅ | ❌ | ✅ | ✅ |
| **Array Support** | ❌ | ❌ | ✅ | ✅ | ❌ | ✅ | ❌ |
| **Full-Text Search** | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ |

### Code Generation Options

| Option | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|---------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **Emit JSON Tags** | ❌ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Prepared Queries** | ❌ | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ |
| **Interface Generation** | ❌ | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ |
| **Empty Slices** | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Result Struct Pointers** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Params Struct Pointers** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Enum Validation** | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ |
| **All Enum Values** | ❌ | ❌ | ❌ | ✅ | ❌ | ❌ | ✅ |
| **JSON Tags Case Style** | snake | camel | camel | camel | snake | snake | camel | camel |

### Validation Options

| Validation | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|-----------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **Strict Functions** | ❌ | ❌ | ✅ | ❌ | ✅ | ❌ | ❌ |
| **Strict Order By** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **No Select Star** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Require Where** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **No Drop Table** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **No Truncate** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Require Limit** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ |

### Type Overrides

| Override Type | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|--------------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **UUID → uuid.UUID** | ❌ | ✅ | ✅ | ✅ | ❌ | ✅ | ❌ |
| **JSON → json.RawMessage** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Arrays → []string** | ❌ | ❌ | ✅ | ✅ | ❌ | ❌ | ❌ | ❌ |
| **Full-Text → string** | ❌ | ❌ | ✅ | ❌ | ❌ | ❌ | ❌ | ❌ |
| **Array Nullable** | - | - | ✅ | ✅ | - | - | ✅ | - |

### Rename Rules

| Rename Rule | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|-------------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **id → ID** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ✅ |
| **uuid → UUID** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **url → URL** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **uri → URI** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **json → JSON** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **api → API** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **http → HTTP** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **db → DB** | ❌ | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | ❌ |
| **otp → OTP** | ❌ | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |

### Database Engine Support

| Database | Hobby | Microservice | Enterprise | API First | Analytics | Testing | Multi Tenant | Library |
|----------|-------|-------------|------------|-----------|----------|--------|-------------|--------|
| **PostgreSQL** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **MySQL** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **SQLite** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## Detailed Template Analysis

### Hobby Template

**Summary:** Simple, lightweight template for personal projects and prototypes

**Strengths:**
- ✅ Minimal dependencies (database/sql)
- ✅ Easy to set up and use
- ✅ SQLite requires no database server
- ✅ Perfect for learning and prototyping
- ✅ Simple generated code

**Weaknesses:**
- ❌ No UUID support
- ❌ No JSON support
- ❌ No prepared queries
- ❌ No interface generation
- ❌ SQLite not production-ready for high scale

**Best For:** Personal projects, prototypes, learning sqlc, quick MVPs

**Not Best For:** Production systems, APIs, multi-user applications

**Estimated Lines of Generated Code:** ~50-100 lines (simple)

---

### Microservice Template

**Summary:** Production-ready template for API services and microservices

**Strengths:**
- ✅ pgx/v5 driver (high performance PostgreSQL)
- ✅ Prepared queries (security & performance)
- ✅ JSON tags (API standard)
- ✅ Interface generation (mocking & testing)
- ✅ Optimized for microservices

**Weaknesses:**
- ❌ No UUID support (need for distributed IDs)
- ❌ No array support (many-to-many relationships)
- ❌ No full-text search
- ❌ No strict validation (data integrity)

**Best For:** API services, microservices, REST/GraphQL backends

**Not Best For:** Analytics platforms, SaaS multi-tenant, simple CRUD

**Estimated Lines of Generated Code:** ~100-200 lines (medium)

---

### Enterprise Template

**Summary:** Strict, production-ready template for enterprise applications

**Strengths:**
- ✅ ALL features enabled (UUID, JSON, Arrays, Full-Text)
- ✅ pgx/v5 driver (high performance PostgreSQL)
- ✅ Strict validation (data integrity)
- ✅ All safety rules (no bad queries)
- ✅ Prepared queries
- ✅ Interface generation
- ✅ Comprehensive type overrides

**Weaknesses:**
- ❌ Most complex generated code
- ❌ Longer build times
- ❌ More dependencies (uuid, pgx)
- ❌ Overkill for simple projects

**Best For:** Production systems, enterprise applications, critical infrastructure

**Not Best For:** Personal projects, simple prototypes, learning sqlc

**Estimated Lines of Generated Code:** ~200-400 lines (complex)

---

### API First Template

**Summary:** API-optimized template with JSON serialization and camelCase naming

**Strengths:**
- ✅ pgx/v5 driver (high performance PostgreSQL)
- ✅ JSON tags with camelCase (API standard)
- ✅ Enum validation (type safety)
- ✅ All enum values generated (flexibility)
- ✅ Prepared queries
- ✅ Interface generation

**Weaknesses:**
- ❌ No UUID support
- ❌ No array support (many-to-many relationships)
- ❌ No strict validation
- ❌ No full-text search

**Best For:** REST/GraphQL APIs, API-first development, API services

**Not Best For:** Analytics platforms, SaaS multi-tenant, hobby projects

**Estimated Lines of Generated Code:** ~100-150 lines (medium)

---

### Analytics Template

**Summary:** Analytics-optimized template with array support, full-text search, and JSON storage

**Strengths:**
- ✅ Array support (time-series data)
- ✅ Full-text search (analytics queries)
- ✅ JSONB support (document storage)
- ✅ Strict validation (data quality)
- ✅ pgx/v5 driver (high performance PostgreSQL)
- ✅ Interface generation

**Weaknesses:**
- ❌ No UUID support (time-series data may use other IDs)
- ❌ No prepared queries
- ❌ No enum validation
- ❌ No all enum values

**Best For:** Data analytics, warehousing, reporting platforms

**Not Best For:** Simple CRUD apps, API services, hobby projects

**Estimated Lines of Generated Code:** ~100-200 lines (medium)

---

### Testing Template

**Summary:** Minimal template for test fixtures and mock data

**Strengths:**
- ✅ Minimal dependencies (database/sql)
- ✅ Test database paths (testdata/)
- ✅ Easy to set up and tear down
- ✅ No complex features (simple tests)
- ✅ Simple generated code

**Weaknesses:**
- ❌ No UUID support (tests use simple IDs)
- ❌ No JSON support (tests use simple types)
- ❌ No prepared queries
- ❌ No interface generation
- ❌ No advanced features

**Best For:** Test fixtures, mock data, integration tests

**Not Best For:** Production systems, API services, complex applications

**Estimated Lines of Generated Code:** ~50-100 lines (simple)

---

### Multi Tenant Template

**Summary:** SaaS-optimized template with UUID support and tenant isolation

**Strengths:**
- ✅ UUID support (for tenant IDs)
- ✅ Array support (tenant relationships)
- ✅ JSONB support (tenant metadata)
- ✅ Strict validation (tenant data integrity)
- ✅ pgx/v5 driver (high performance PostgreSQL)
- ✅ Prepared queries
- ✅ Interface generation

**Weaknesses:**
- ❌ No full-text search
- ❌ No enum validation
- ❌ No all enum values
- ❌ More complex than needed for simple multi-tenant

**Best For:** SaaS platforms, multi-tenant applications, B2B systems

**Not Best For:** Single-tenant apps, hobby projects, analytics platforms

**Estimated Lines of Generated Code:** ~150-250 lines (medium-complex)

---

### Library Template

**Summary:** Library-optimized template with minimal dependencies and interface generation

**Strengths:**
- ✅ Minimal dependencies (broad compatibility)
- ✅ pgx/v5 driver (high performance)
- ✅ Interface generation (library users can mock)
- ✅ Enum validation (type safety)
- ✅ All enum values generated (flexibility)
- ✅ JSON tags (library consumers)

**Weaknesses:**
- ❌ No UUID support (libraries may use different ID types)
- ❌ No array support
- ❌ No full-text search
- ❌ No strict validation (flexible for library consumers)
- ❌ No prepared queries (broad compatibility)

**Best For:** Reusable Go libraries, SDKs, packages

**Not Best For:** Full applications, production systems, multi-tenant SaaS

**Estimated Lines of Generated Code:** ~100-150 lines (medium)

---

## Complexity Ranking

1. **Hobby** (Simplest) - ~50-100 LOC, minimal features
2. **Testing** (Simple) - ~50-100 LOC, test-optimized
3. **Library** (Low) - ~100-150 LOC, interface-focused
4. **API First** (Low-Medium) - ~100-150 LOC, JSON-focused
5. **Microservice** (Medium) - ~100-200 LOC, API-optimized
6. **Analytics** (Medium) - ~100-200 LOC, data-heavy
7. **Multi Tenant** (Medium-High) - ~150-250 LOC, SaaS-optimized
8. **Enterprise** (Most Complex) - ~200-400 LOC, all features

---

## Use Case Recommendations

### For Personal Projects

**Recommended:** Hobby Template
**Alternative:** Testing Template (for testing features)
**Why:** Minimal setup, SQLite, easy to learn

### For API Services

**Recommended:** Microservice Template
**Alternative:** API First Template (for JSON-first APIs)
**Why:** Prepared queries, pgx/v5, interface generation

### For SaaS/Multi-Tenant Applications

**Recommended:** Multi Tenant Template
**Alternative:** Enterprise Template (if strict validation needed)
**Why:** UUID support, tenant isolation optimization

### For Production/Enterprise Systems

**Recommended:** Enterprise Template
**Alternative:** Multi Tenant Template (if SaaS)
**Why:** Strict validation, all safety rules, comprehensive features

### For Data Analytics

**Recommended:** Analytics Template
**Why:** Array support, full-text search, JSONB storage

### For Reusable Libraries

**Recommended:** Library Template
**Why:** Minimal dependencies, interface generation, broad compatibility

### For Test Fixtures

**Recommended:** Testing Template
**Why:** Test database paths, minimal features, easy setup

### For Learning/Prototyping

**Recommended:** Hobby Template
**Alternative:** Testing Template
**Why:** SQLite, minimal setup, simple generated code

---

## Migration Guide

### From Hobby to Microservice

```yaml
# Add UUID support
data.Database.UseUUIDs = true

# Add JSON support
data.Database.UseJSON = true

# Change database engine
data.Database.Engine = templates.DatabaseTypePostgreSQL

# Enable prepared queries
data.Validation.EmitOptions.EmitPreparedQueries = true

# Enable interface generation
data.Validation.EmitOptions.EmitInterface = true
```

### From Hobby to Enterprise

```yaml
# Add ALL features
data.Database.UseUUIDs = true
data.Database.UseJSON = true
data.Database.UseArrays = true
data.Database.UseFullText = true

# Enable ALL strict validation
data.Validation.StrictFunctions = true
data.Validation.StrictOrderBy = true
data.Validation.SafetyRules.NoSelectStar = true
data.Validation.SafetyRules.RequireWhere = true
data.Validation.SafetyRules.NoDropTable = true
data.Validation.SafetyRules.NoTruncate = true
data.Validation.SafetyRules.RequireLimit = true

# Enable ALL code gen options
data.Validation.EmitOptions.EmitResultStructPointers = true
data.Validation.EmitOptions.EmitParamsStructPointers = true
```

### From Microservice to Enterprise

```yaml
# Add missing features
data.Database.UseUUIDs = true
data.Database.UseArrays = true
data.Database.UseFullText = true

# Enable strict validation
data.Validation.StrictFunctions = true
data.Validation.StrictOrderBy = true

# Enable safety rules
data.Validation.SafetyRules.NoSelectStar = true
data.Validation.SafetyRules.RequireWhere = true
# etc.
```

---

## Summary

- **8 Templates** available for different use cases
- **Hobby & Testing**: Simple, SQLite, minimal features
- **Microservice & API First**: API-optimized, pgx/v5, prepared queries
- **Enterprise**: Production-ready, all features, strict validation
- **Analytics**: Data-heavy, arrays, full-text search
- **Multi Tenant**: SaaS-optimized, UUID support, tenant isolation
- **Library**: Minimal deps, interface generation, broad compatibility

**Choose template based on:**
1. Project type (personal, API, SaaS, enterprise, etc.)
2. Complexity requirements (simple, medium, complex)
3. Database needs (SQLite, PostgreSQL, MySQL)
4. Feature requirements (UUID, JSON, arrays, full-text)
5. Validation strictness (flexible vs strict)
6. Performance requirements (prepared queries, interfaces)
7. Integration needs (library consumers, internal use)

**Next:** Review [Template Usage Guide](./usage.md) for detailed usage instructions.

---

**Last Updated:** 2026-02-05
