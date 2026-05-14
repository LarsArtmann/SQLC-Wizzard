# Dependency Graph — SQLC-Wizard

**Date:** 2026-05-13

---

## Current State (Pre-Modularization)

### Module Dependency Diagram

```
┌─────────────────────────────────────────────────────────┐
│                    Root Module                           │
│              github.com/LarsArtmann/SQLC-Wizzard         │
│                                                          │
│  ┌──────────────┐                                        │
│  │ cmd/sqlc-    │                                        │
│  │ wizard       │──→ internal/commands                   │
│  └──────────────┘                                        │
│         │                                                │
│         ▼                                                │
│  ┌──────────────┐                                        │
│  │ commands     │──→ adapters, creators, generators,     │
│  │              │    templates, wizard, ui, apperrors,   │
│  │              │    generated, pkg/config               │
│  └──────┬───────┘                                        │
│         │                                                │
│    ┌────┴────────────────────────────────────┐           │
│    ▼              ▼            ▼              ▼           │
│  ┌──────┐  ┌──────────┐  ┌──────┐  ┌──────────┐        │
│  │wizard│  │ creators  │  │gen-  │  │ adapters │        │
│  │      │  │           │  │erators│  │          │        │
│  └──┬───┘  └─────┬─────┘  └──┬───┘  └────┬─────┘        │
│     │            │           │            │               │
│     ▼            ▼           ▼            ▼               │
│  ┌──────────────────────────────────────────┐            │
│  │         templates (strategy pattern)      │            │
│  └───────────────────┬──────────────────────┘            │
│                      │                                    │
│                      ▼                                    │
│  ┌──────────┐  ┌───────────┐  ┌───────┐                  │
│  │validation│  │ pkg/config │  │domain │                  │
│  └─────┬────┘  └─────┬─────┘  └───┬───┘                  │
│        │             │             │                       │
│        └──────┬──────┘             │                       │
│               ▼                    ▼                       │
│        ┌─────────────┐    ┌──────────────┐                │
│        │ apperrors   │    │  generated/  │ ◄── replace    │
│        └─────────────┘    │  (submodule) │    directive    │
│                           └──────────────┘                │
│                                                          │
│  ┌──────────┐  ┌──────┐  ┌──────┐  ┌──────┐             │
│  │migration │  │schema│  │  ui  │  │utils │             │
│  │(no deps) │  │(none)│  │(none)│  │(none)│             │
│  └──────────┘  └──────┘  └──────┘  └──────┘             │
└─────────────────────────────────────────────────────────┘

┌─────────────────────┐    ┌─────────────────────────┐
│ examples/hobby-     │    │ generated/              │
│ project (standalone)│    │ module: sqlc-wizard-    │
│ No deps on root     │    │ types (no deps)         │
└─────────────────────┘    └─────────────────────────┘
```

### Internal Package Coupling Matrix

| Package | generated | apperrors | domain | validation | schema | templates | config | adapters | commands | creators | generators | wizard | ui | utils | migration | testing |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| **apperrors** | | | | | | | | | | | | | | | | |
| **domain** | ✅ | | | | | | | | | | | | | | | |
| **validation** | ✅ | | ✅ | | | | | | | | | | | | | |
| **schema** | | | | | | | | | | | | | | | | |
| **templates** | ✅ | ✅ | | ✅ | | | ✅ | | | | | | | | | |
| **config** | ✅ | ✅ | | | | | | | | | | | | | | |
| **adapters** | ✅ | ✅ | | | ✅ | ✅ | ✅ | | | | | | | | ✅ | |
| **commands** | ✅ | ✅ | | | | ✅ | ✅ | ✅ | | ✅ | ✅ | ✅ | ✅ | | | |
| **creators** | ✅ | ✅ | | | | | ✅ | ✅ | | | | | | | | |
| **generators** | | | | | | ✅ | ✅ | | | | | | | | | |
| **wizard** | ✅ | ✅ | | | ✅ | ✅ | ✅ | | | | | | ✅ | | | |
| **ui** | | | | | | | | | | | | | | | | |
| **utils** | | | | | | | | | | | | | | | | |
| **migration** | | | | | | | | | | | | | | | | |
| **testing** | ✅ | | ✅ | | | | ✅ | | | | | | | | | |
| **cmd/sqlc-wizard** | | | | | | | | | ✅ | | | | | | | |

### Coupling Metrics

| Metric | Value | Assessment |
|---|---|---|
| Total internal packages | 16 | Medium |
| Packages with ≥5 internal deps | 3 (commands, adapters, wizard) | God-package territory |
| Leaf packages (0 deps) | 5 (apperrors, schema, ui, utils, migration) | Good |
| Max dependency depth | 6 (cmd → commands → ... → generated) | Acceptable |
| Layer violations | 1 (pkg/config → internal/apperrors) | Must fix |
| God-packages (>15 exports) | 7 | Needs attention |

---

## Proposed State (Post-Modularization)

### Module Dependency Diagram

```
┌──────────────────┐
│   generated/     │     Layer 0 — Pure types, no deps
│ sqlc-wizard-types│
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│     core/        │     Layer 1 — Domain logic
│ apperrors        │
│ domain           │     Depends on: generated
│ validation       │
│ schema           │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│    config/       │     Layer 2 — YAML config types
│ SqlcConfig       │
│ Parser/Validator │     Depends on: generated, core
└────────┬─────────┘
         │
         ▼
┌──────────────────────────────────────────────┐
│               Root (app)                      │  Layer 3 — Application
│ cmd/sqlc-wizard                               │
│ internal/templates, wizard, adapters,         │  Depends on: generated,
│   commands, creators, generators,             │    core, config
│   migration, ui, utils, testing               │
└──────────────────────────────────────────────┘

┌──────────────────┐
│ examples/hobby-  │     Standalone — no change
│ project          │
└──────────────────┘
```

### Post-Modularization Coupling (Within Root Module)

After extraction, the root module's internal coupling is reduced:

| Package | Removed deps (moved to core/config) | Remaining deps |
|---|---|---|
| **templates** | `apperrors`, `validation` | `generated`, `config` |
| **adapters** | `apperrors`, `schema`, `migration` | `generated`, `templates`, `config` |
| **commands** | `apperrors` | `adapters`, `creators`, `generators`, `templates`, `wizard`, `ui`, `generated`, `config` |
| **wizard** | `apperrors`, `schema` | `generated`, `templates`, `ui`, `config` |
| **creators** | `apperrors` | `adapters`, `generated`, `config` |

### Impact Summary

| Metric | Before | After | Change |
|---|---|---|---|
| Root go.mod direct deps | 15+ | ~12 (4 moved to core/config) | -20% |
| Layer violations | 1 (pkg → internal) | 0 | Fixed |
| Independent modules | 3 (root, generated, examples) | 4 (root, generated, core, config) | +1 |
| Max module depth | N/A | 3 levels | Enforced DAG |
| go.work | None | Yes | Added |
