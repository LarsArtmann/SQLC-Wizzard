# Domain Language

A shared vocabulary for `generating production-ready sqlc v2 configurations via an interactive CLI wizard` — grounded in the repo's own code and docs.

## Glossary

| Term | Definition | Where it lives |
| ---- | ---------- | -------------- |
| Wizard | The interactive `sqlc-wizard init` flow that asks about project type, database, and package path. | `internal/wizard/` |
| sqlc.yaml | The generated sqlc v2 configuration with smart defaults, plus scaffolded `schema/` and `queries/` directories. | `internal/creators/` |
| ProjectType | Template enum driving defaults: Hobby, Microservice, Enterprise, API-First, Analytics, Testing, Multi-tenant, Library. | `generated/types.go` |
| DatabaseType | Supported target databases: PostgreSQL, MySQL, SQLite. | `generated/types.go` |
| Safety Policy | Best-practice safety rules applied to generated configurations. | `internal/domain/safety_policy.go` |
| Emit Options | Fine-grained codegen toggles: emit modes, struct pointers, null handling, JSON tag styles, enum generation. | `internal/domain/emit_modes.go` |
| Config Migration | The `migrate` command that upgrades existing sqlc.yaml configurations. | `internal/migration/` |
| Commands | The CLI surface: `init`, `validate`, `generate`, `doctor`, `migrate`. | `internal/commands/` |

> Verified against the repo's README and source on 2026-11-21.
