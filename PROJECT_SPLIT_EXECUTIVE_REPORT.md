# Project Split Executive Report: SQLC-Wizard Refactoring

## 1. Executive Summary

This report outlines a proposed refactoring of the existing `SQLC-Wizard` project into a set of highly focused, independent Go modules and applications. The current monolithic structure, while functional, presents challenges in terms of maintainability, reusability, and independent development. By splitting the project along logical domain boundaries, we aim to enhance modularity, improve testability, facilitate parallel development, and enable wider reuse of core components.

## 2. Current Project Overview (SQLC-Wizard)

The `SQLC-Wizard` is an interactive CLI tool designed to generate `sqlc.yaml` configurations. Its key features include:

- An interactive TUI wizard (`charmbracelet/huh`).
- Generation of type-safe `sqlc.yaml` configurations with validation.
- Support for multiple project templates and database types.
- Various CLI commands (`init`, `validate`, `doctor`, `generate`, `migrate`).
- Generation of Go types from TypeSpec for compile-time safety.

The existing architecture follows Domain-Driven Design (DDD) principles with a layered structure, but its implementation as a single Go module limits the benefits of this design.

## 3. Proposed Project Split

The `SQLC-Wizard` project can be logically decomposed into five distinct, highly focused Go projects:

### 3.1. `sqlc-config-schema` (Go Module)

- **Responsibility:** Defines the canonical `sqlc.yaml` configuration types, associated validation logic, and utilities for serialization/deserialization. This module serves as the single source of truth for the `sqlc` configuration structure.
- **Current Mapping:** Primarily `pkg/config/` and related error definitions within `internal/errors/` that pertain to configuration parsing and validation.
- **Dependencies:** Minimal, likely only standard Go libraries and potentially common utility libraries for validation. It explicitly does _not_ depend on `sqlc` itself, only on the schema _for_ `sqlc` configurations.

### 3.2. `sqlc-template-engine` (Go Module)

- **Responsibility:** Encapsulates the core logic for managing, selecting, and generating `sqlc.yaml` configurations based on predefined templates. It will provide the `Template` interface and various template implementations for different project types (e.g., microservice, hobby). This module consumes the `sqlc-config-schema`.
- **Current Mapping:** `internal/templates/`, `internal/domain/`, and the raw SQL template files located in `templates/`.
- **Dependencies:** `sqlc-config-schema` (as a Go module), `github.com/samber/lo`, `github.com/samber/mo`, and the `sqlc-typespec-generated` module.

### 3.3. `sqlc-wizard-ui` (Go Module)

- **Responsibility:** Implements the interactive Text User Interface (TUI) wizard functionality. This module focuses solely on user interaction, gathering input through `charmbracelet/huh` components, and presenting feedback. It outputs a structured set of parameters that can be used by the `sqlc-template-engine` or `sqlc-cli-core`.
- **Current Mapping:** `internal/wizard/`.
- **Dependencies:** `charmbracelet/huh`, `sqlc-config-schema` (for understanding config structure for UI elements), and potentially an interface to `sqlc-template-engine` if it needs to dynamically present template options or validate input against template requirements.

### 3.4. `sqlc-cli-core` (Go Application)

- **Responsibility:** This is the executable CLI application. It handles command-line argument parsing (`spf13/cobra`), orchestrates the flow between the `sqlc-wizard-ui`, `sqlc-template-engine`, and `sqlc-typespec-generated` modules, and manages the actual file generation (`internal/generators`). It acts as the application layer, connecting UI to domain logic and infrastructure.
- **Current Mapping:** `cmd/sqlc-wizard/`, `internal/commands/`, `internal/generators/`, `internal/adapters/`, and the main application entry point.
- **Dependencies:** `sqlc-wizard-ui`, `sqlc-template-engine`, `sqlc-config-schema`, `spf13/cobra`, and other CLI-specific libraries.

### 3.5. `sqlc-typespec-generated` (Go Module - Potentially Separate Repository)

- **Responsibility:** Contains the TypeSpec definitions and the generated Go types (e.g., enums like `ProjectType`, `DatabaseType`). This module is a foundational dependency that provides compile-time type safety across other modules. Given its generated nature and fundamental role, it could be a standalone repository if these types are intended to be consumed by other projects outside this direct ecosystem. For internal use within a monorepo, a dedicated Go module is sufficient.
- **Current Mapping:** `generated/`.
- **Dependencies:** None, as it is a source of generated types. The TypeSpec generation process itself is external to this module's runtime dependencies.

## 4. Benefits of Modularization

The proposed project split offers several significant advantages:

- **Clearer Boundaries and Separation of Concerns:** Each project has a distinct, well-defined responsibility, making it easier to understand, manage, and evolve individual components.
- **Enhanced Maintainability:** Changes within one module (e.g., UI updates in `sqlc-wizard-ui`) are isolated and less likely to introduce regressions in unrelated core logic or configuration schema.
- **Increased Reusability:** Core components like `sqlc-config-schema` and `sqlc-template-engine` become independently consumable Go modules. This allows other tools, services, or even different UI implementations to leverage the underlying logic without pulling in the entire CLI.
- **Facilitated Independent Development:** Different teams or developers can work on distinct sub-projects concurrently, reducing merge conflicts and accelerating development cycles.
- **Simplified Testing:** Each module can be unit-tested and integration-tested in isolation, leading to more focused test suites, faster test execution, and higher confidence in the correctness of individual components.
- **Optimized Dependency Management:** Each sub-project declares only its direct dependencies, resulting in smaller, more manageable `go.mod` files and reduced build times.
- **Improved Scalability:** The modular structure provides a clearer path for future expansion and the addition of new features without increasing the complexity of a single large codebase.

## 5. Implementation Considerations

- **Monorepo vs. Polyrepo:** While the projects are logically separate, an initial implementation might benefit from a monorepo approach (e.g., using Go Workspaces) to simplify development and testing across modules, especially for the `sqlc-typespec-generated` module. As the project matures, specific modules could be extracted into their own repositories if external consumption becomes a primary driver.
- **Dependency Updates:** Establishing clear dependency graphs between the new modules will be crucial to manage updates and releases effectively.
- **Migration Strategy:** A phased migration approach would minimize disruption. Start by extracting the most independent modules (e.g., `sqlc-config-schema`, `sqlc-typespec-generated`), then refactor the dependent modules.

## 6. Conclusion

Splitting the `SQLC-Wizard` project into these focused modules represents a strategic investment in the long-term health, scalability, and reusability of the codebase. This approach aligns with best practices for large-scale software development and will yield significant benefits in terms of maintainability, development efficiency, and architectural clarity.
