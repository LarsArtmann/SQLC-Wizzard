# SQLC-Wizzard Migration Guide

## Overview

This guide helps you migrate between different versions of SQLC-Wizzard and upgrade your existing SQLC projects.

## Version History

| Version | Release Date | Major Changes   | Migration Guide |
| ------- | ------------ | --------------- | --------------- |
| 1.0.0   | TBD          | Initial release | N/A             |
| 0.9.0   | 2024-01      | Beta features   | N/A             |

## Upgrading from Previous Versions

### Upgrading from 0.9.0 to 1.0.0

#### Breaking Changes

1. **Configuration Format**
   - Old: YAML with inline schemas
   - New: Structured yamlc format with separate files

2. **Template System**
   - Old: Single template with options
   - New: Multiple templates (hobby, microservice, enterprise, etc.)

3. **CLI Commands**
   - Old: `wizard init` only
   - New: `wizard init`, `wizard create`, `wizard validate`, `wizard doctor`

#### Migration Steps

**Step 1: Backup Your Project**

```bash
# Create a backup of your current project
cp -r my-sqlc-project my-sqlc-project.backup

# Or use git
git commit -m "backup before upgrade"
git branch backup-before-upgrade
```

**Step 2: Update sqlc.yaml**

Old format (0.9.0):

```yaml
version: "1"
sql:
  - name: "db"
    engine: "postgresql"
    queries: "queries/"
    schema: "schema/"
```

New format (1.0.0):

```yaml
version: "2"
sql:
  - name: "db"
    engine: "postgresql"
    queries: "./sql/queries"
    schema: "./sql/schema"
    gen:
      go:
        package: "db"
        out: "./internal/db"
        sql_package: "pgx/v5"
```

**Changes to make:**

- Change version from "1" to "2"
- Add `./` prefix to paths (or use absolute paths)
- Add `gen.go` section with configuration

**Step 3: Reorganize File Structure**

Old structure (0.9.0):

```
my-sqlc-project/
├── sqlc.yaml
├── queries/
├── schema/
└── db/
```

New structure (1.0.0):

```
my-sqlc-project/
├── sqlc.yaml
├── sql/
│   ├── queries/
│   └── schema/
└── internal/db/
```

**Migration:**

```bash
# Create new structure
mkdir -p sql/queries sql/schema

# Move files
mv queries/* sql/queries/
mv schema/* sql/schema/

# Verify structure
tree sql/
```

**Step 4: Update Go Imports**

If your Go code imports generated package:

Old (0.9.0):

```go
import "github.com/username/project/db"
```

New (1.0.0):

```go
import "github.com/username/project/internal/db"
```

**Migration:**

```bash
# Update imports throughout your codebase
find . -name "*.go" -exec sed -i '' 's|project/db|project/internal/db|g' {} +

# Verify with grep
grep -r "import.*db" --include="*.go"
```

**Step 5: Regenerate Code**

```bash
# Generate code with new format
sqlc generate

# Verify output
tree internal/db/
```

**Step 6: Run Tests**

```bash
# Run your test suite
go test ./...

# Fix any compilation errors
go build ./...
```

**Step 7: Verify Everything Works**

```bash
# Run your application
go run main.go

# Test database operations
# ... perform your tests ...
```

**Step 8: Clean Up**

```bash
# Once everything is working, remove backup
rm -rf my-sqlc-project.backup

# Or keep git branch
git branch -D backup-before-upgrade
```

## Common Migration Issues

### Issue: "sqlc: version 2 not supported"

**Cause:** Using old SQLC version

**Solution:**

```bash
# Update SQLC to latest version
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify version
sqlc version
```

### Issue: "path must be absolute or start with ./"

**Cause:** Using relative paths without `./` prefix

**Solution:**

```yaml
# ❌ OLD - doesn't work
queries: "sql/queries"

# ✅ NEW - works
queries: "./sql/queries"

# ✅ OR - use absolute path
queries: "/path/to/project/sql/queries"
```

### Issue: "package not found"

**Cause:** Package path mismatch between go.mod and sqlc.yaml

**Solution:**

```yaml
# In sqlc.yaml
gen:
  go:
    package: "db"  # Should match go.mod module + directory

# In go.mod
module github.com/username/project  # Should include "internal/db" path
```

### Issue: "no such file or directory"

**Cause:** File structure changed, code still references old paths

**Solution:**

```bash
# Find all files referencing old paths
grep -r "queries/" --include="*.go"

# Update imports
find . -name "*.go" -exec sed -i '' 's|queries/|sql/queries/|g' {} +
```

## Rollback Plan

If migration fails, rollback to previous version:

**Option 1: Use Git**

```bash
# Discard changes
git reset --hard HEAD

# Or revert to backup branch
git checkout backup-before-upgrade
```

**Option 2: Use Backup**

```bash
# Restore from backup
rm -rf my-sqlc-project
cp -r my-sqlc-project.backup my-sqlc-project
```

**Option 3: Selective Rollback**

```bash
# Keep new configuration but revert code changes
git reset --hard HEAD~1  # Reset commits but keep config changes
```

## Post-Migration Checklist

- [ ] All tests passing
- [ ] Application compiles without errors
- [ ] Database operations work correctly
- [ ] No warnings from SQLC
- [ ] Code generation produces expected output
- [ ] File structure matches new format
- [ ] Documentation updated (if applicable)
- [ ] Team members notified of changes
- [ ] CI/CD pipeline updated
- [ ] Backup removed (after verification)

## Advanced Migration Scenarios

### Migrating Multiple Projects

If you have multiple SQLC projects to migrate:

```bash
# List all projects
find . -name "sqlc.yaml" -type f

# Run migration script for each
for project in $(find . -name "sqlc.yaml" -type f); do
    cd $(dirname "$project")
    echo "Migrating: $project"
    # ... migration steps ...
done
```

### Custom Template Migration

If you have a custom SQLC configuration:

**Option 1: Use New Template System**

```bash
# Find closest matching template
wizard init

# Customize generated config
# Edit sqlc.yaml to match your needs
```

**Option 2: Keep Existing Config**

```bash
# Update version field in sqlc.yaml
sed -i '' 's/version: "1"/version: "2"/' sqlc.yaml

# Test if still works
sqlc generate
```

### CI/CD Pipeline Migration

Update your CI/CD configuration for new format:

**GitHub Actions Example:**

```yaml
name: Test SQLC Project

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.21"
      - name: Install SQLC
        run: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
      - name: Generate code
        run: sqlc generate
      - name: Run tests
        run: go test ./...
```

## Getting Help

### Documentation

- [User Guide](./USER_GUIDE.md)
- [Best Practices](./BEST_PRACTICES.md)
- [Troubleshooting](./TROUBLESHOOTING.md)

### Community

- [GitHub Issues](https://github.com/LarsArtmann/SQLC-Wizzard/issues)
- [GitHub Discussions](https://github.com/LarsArtmann/SQLC-Wizzard/discussions)
- [SQLC Documentation](https://docs.sqlc.dev/)

### Support

If you encounter issues not covered in this guide:

1. Check the [Troubleshooting](./TROUBLESHOOTING.md) guide
2. Search [GitHub Issues](https://github.com/LarsArtmann/SQLC-Wizzard/issues)
3. Ask a question in [GitHub Discussions](https://github.com/LarsArtmann/SQLC-Wizzard/discussions)
4. Create a new issue with detailed information

## Best Practices for Migration

1. **Always Backup First**

   ```bash
   cp -r project project.backup
   ```

2. **Test Migration on Copy**

   ```bash
   cp -r project project.test
   cd project.test
   # ... perform migration ...
   ```

3. **Migrate Incrementally**
   - Update one component at a time
   - Test after each change
   - Commit working state

4. **Document Changes**
   - Keep track of what you changed
   - Note any issues encountered
   - Document solutions for future reference

5. **Verify Thoroughly**
   - Run all tests
   - Test all major features
   - Check error logs
   - Performance test

## Summary

Migration should be straightforward if you follow this guide:

✅ **Backup** your project before starting
✅ **Update** sqlc.yaml to new format
✅ **Restructure** directories
✅ **Update** imports
✅ **Regenerate** code
✅ **Test** everything
✅ **Clean up** backup

Following these steps will ensure a smooth transition to the new version of SQLC-Wizzard.

## License

MIT License - See [LICENSE](./LICENSE) for details
