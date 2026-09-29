# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

## [Unreleased]

### Added

- CI quality gate: test with race detector and coverage profile, golangci-lint, and `dupl` duplication check on every push/PR (`.github/workflows/ci-cd.yml`)
- Release automation: GoReleaser workflow on version tags with GHCR image publish, Cosign signing, and SBOM (`.github/workflows/release.yml`, `.goreleaser.yml`); tags `v0.1.0` and `v0.2.0` cut
- Nix flake: reproducible build, dev shell, and runnable app (`flake.nix`); formatting migrated to dprint; TypeSpec toolchain to pnpm
- Template system: 8 project templates (hobby, microservice, enterprise, api-first, analytics, testing, multi-tenant, library) with `BaseTemplate`/`ConfiguredTemplate` shared infrastructure and zero-value-safe constructors
- Template documentation set: usage guide, comparison matrix, customization guide, and 8 per-template examples (`docs/templates/`)
- Documentation suite: user guide, tutorial, best practices, troubleshooting, migration guide, advanced features (`docs/`)
- Working hobby example projects (`examples/hobby-project/`, `examples/hobby-sqlite/`)
- CONTRIBUTING guide
- Error package `internal/apperrors` (renamed from `internal/errors` to avoid standard-library collision; 26 packages updated)

### Changed

- File-size refactor: 11 oversized files split into 24 focused files; every flagged file now within the 350-line policy (2026-06-14)
- Code deduplication pass: clone groups reduced from 10 to 3; table-driven test conversions and shared test helpers (2026-01-14)
- Templates migrated from `BaseTemplate` embedding to `ConfiguredTemplate` embedding where appropriate; MicroserviceTemplate intentionally kept on `BaseTemplate` (custom `Generate()` logic)

### Fixed

- 8 failing template `DefaultData` tests: the assertion helper compared plain style names (`"camel"`/`"snake"`) with `assert.JSONEq`, which requires JSON on both sides; corrected to `assert.Equal` in `internal/testing/assertions.go` (2026-09-29)

## [0.2.0] - 2026-05-17

### Changed

- Tagged from the same commit as v0.1.0 (`bda4efe`) while validating the release pipeline (brew/scoop upload skipped until tokens are configured)

## [0.1.0] - 2026-01-01

### Added

- Initial release
