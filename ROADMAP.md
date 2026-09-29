# Roadmap

Long-term direction and raw ideas. Bounded, actionable work lives in TODO_LIST.md; shipped features live in FEATURES.md.

## Themes

### 1. Trustworthy core

The wizard is the product; its correctness must be provable.

- Raise wizard, commands, and adapters test coverage toward the 70-80% band
- Golden/snapshot testing culture for every generated artifact (sqlc.yaml, SQL starters)
- Integration tests that run real `sqlc` against generated configs
- Performance benchmarks and regression gates for template generation

### 2. Frictionless adoption

A user should go from install to generated project in minutes.

- More real-world examples (microservice with Docker Compose, enterprise with audit tables, analytics, api-first)
- Video tutorials / animated walkthroughs
- Shell completions (bash/zsh/fish)
- `--dry-run` preview mode for all mutating commands
- Interactive tutorial mode inside the wizard

### 3. Extensibility

Templates and rules should be extensible without forking.

- Plugin architecture for custom templates
- Template versioning with compatibility checks and migration tools
- Template linting tool (consistency, naming, feature coverage)
- Configuration presets (save/load wizard answers as reusable profiles)
- Config export/import across YAML/JSON/JSONC

### 4. Reach

- Documentation website with interactive template comparison
- IDE integrations (VS Code / GoLand plugin surfaces)
- Web-based configuration generator
- Homebrew formula and additional distribution channels
- Kubernetes/Helm deployment guide for enterprise users
- Security audit pass (OWASP-style review of generated configs and file handling)
- Error-message i18n
- Load testing with 100+ table schemas

## Explicit non-goals

- Becoming a general-purpose sqlc wrapper: SQLC-Wizard generates and validates configuration; it does not replace the `sqlc` CLI
- Runtime query execution, connection pooling, or database driver management
- Supporting non-sqlc code generators
