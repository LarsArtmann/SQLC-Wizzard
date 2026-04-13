# Migration to Nix Flakes Proposal

**Project:** SQLC-Wizard  
**Date:** 2026-04-13  
**Status:** Draft  
**Author:** Auto-generated from tooling audit  

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [Current Tooling Inventory](#2-current-tooling-inventory)
3. [Pain Points](#3-pain-points)
4. [What Nix Flakes Solves](#4-what-nix-flakes-solves)
5. [Proposed Architecture](#5-proposed-architecture)
6. [Migration Plan](#6-migration-plan)
7. [Risk Analysis](#7-risk-analysis)
8. [Success Criteria](#8-success-criteria)
9. [Reference Implementation](#9-reference-implementation)
10. [Appendix: Full Tool Dependency Map](#appendix-full-tool-dependency-map)

---

## 1. Executive Summary

SQLC-Wizard currently relies on **5 separate tooling systems** (`just` + Makefile + shell scripts + Docker + CI YAML) to manage builds, tests, linting, and releases. This creates friction: onboarding requires manual installation of Go 1.26+, `just`, `golangci-lint`, `dupl`, `goreleaser`, Node.js/Bun (for TypeSpec), and `tsp` — all at compatible versions. Nix Flakes can collapse all of this into a single `flake.nix` that provides **reproducible, pinned, one-command developer environments** and **hermetic builds**.

This proposal outlines a phased migration that:
- Introduces Nix alongside existing tooling (no breakage)
- Replaces the Justfile's implicit dependency on globally-installed tools
- Eliminates version drift between CI and local development
- Enables `nix develop` → instant, correct dev shell on any machine

---

## 2. Current Tooling Inventory

### 2.1 Build Systems

| System | File | Purpose |
|--------|------|---------|
| **Just** | `justfile` | Primary task runner (build, test, lint, fmt, bench, etc.) |
| **Make** | `Makefile` | TypeSpec code generation only (`make types`) |
| **Shell scripts** | `scripts/check-cmd-single.sh` | Minimal check (placeholder) |
| **Docker** | `Dockerfile` | Multi-stage Alpine build for container images |
| **GoReleaser** | `.goreleaser.yml` | Cross-platform release (brew, Docker, GitHub Release) |
| **GitHub Actions** | `.github/workflows/ci-cd.yml` | CI/CD pipeline (test matrix, security scan, release) |

### 2.2 Language Runtimes

| Runtime | Version | Purpose |
|---------|---------|---------|
| **Go** | `1.26.1` (go.mod) | Primary language |
| **Go** | `1.24`, `1.25` (CI matrix) | Compatibility testing |
| **Node.js / Bun** | `bun.lock` + `package.json` | TypeSpec compiler (`@typespec/compiler ^1.10.0`) |

### 2.3 Development Tools (Currently Expected as Global Installs)

| Tool | Purpose | Install Method | Version Pinned? |
|------|---------|----------------|-----------------|
| `go` | Build & test | Manual / Homebrew | Yes (go.mod) |
| `golangci-lint` | 90+ linters | `go install ...@latest` | **No** — always latest |
| `dupl` | Duplicate code detection | `go install ...@latest` | **No** |
| `goreleaser` | Release automation | CI action / manual | **No** (CI uses `@v5`) |
| `just` | Task runner | Manual / Homebrew | **No** |
| `tsp` (TypeSpec) | Type generation | `npx @typespec/compiler` | Partially (package.json) |
| `git` | Version info via ldflags | System | **No** |
| `go-arch-lint` | Architecture enforcement | Manual | **No** |

### 2.4 Build Configuration

**Linter config** (`.golangci.yml`): 90+ linters enabled, Go 1.26.1 target, 5m timeout, experimental build tags (`goexperiment.goroutineleakprofile`, `goexperiment.jsonv2`, `goexperiment.simd`).

**Go module** (`go.mod`): Uses local `replace` directives:
```
replace github.com/LarsArtmann/SQLC-Wizzard/generated => ./generated
replace github.com/larsartmann/go-composable-business-types => /Users/larsartmann/.../go-composable-business-types
```
The second `replace` is a **local absolute path** — this is a migration concern (see §6.3).

**Version injection** (`cmd/sqlc-wizard/main.go`): Three variables set via `-ldflags` at build time:
- `main.Version` — git tag/describe
- `main.Commit` — git SHA
- `main.BuildDate` — UTC timestamp

---

## 3. Pain Points

### 3.1 "Works on My Machine" Syndrome

| Problem | Impact |
|---------|--------|
| `golangci-lint` installed via `@latest` — version drift between devs and CI | Linter passes locally, fails in CI (or vice versa) |
| No `toolchain` directive in go.mod — Go version enforcement relies on manual compliance | Subtle compiler behavior differences |
| `just` must be installed globally | New contributors must discover and install it |
| TypeSpec requires Bun/Node but this is undocumented in setup steps | `just generate-typespec` fails silently |

### 3.2 Duplicated Build Logic

The same build steps are defined in **4 separate places**:

| Step | justfile | Dockerfile | CI YAML | GoReleaser |
|------|----------|------------|---------|------------|
| `go mod download` | `deps` recipe | Build stage | CI step | `before.hooks` |
| `go build` with ldflags | `build` recipe | Build stage | Build job | `builds` section |
| `go test` | `test` recipe | — | Test job | — |
| `golangci-lint` | `lint` recipe | — | CI step | — |

Changes to ldflags or build flags must be replicated across all four. This has already led to inconsistencies (e.g., CI uses `github.ref_name` for version while justfile uses `git describe`).

### 3.3 CI Fragility

```yaml
# CI installs tools at @latest — unreproducible
- go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
- go install github.com/golangci/dupl@latest
```

A new `golangci-lint` release could break CI at any time with no advance warning.

### 3.4 Local Path Dependency

```
replace github.com/larsartmann/go-composable-business-types => /Users/larsartmann/projects/go-composable-business-types
```

This hard-coded absolute path makes the module uncompilable on any other machine. Nix cannot fix this directly, but the migration plan should address it.

---

## 4. What Nix Flakes Solves

### 4.1 Declarative, Pinned Development Environments

```bash
nix develop   # Enter a shell with: go 1.26, golangci-lint, just, goreleaser, bun, tsp — all pinned
```

Every tool at a known version. `flake.lock` ensures bit-for-bit reproducibility.

### 4.2 Hermetic Builds

```bash
nix build     # Produces a binary in ./result — no system Go required
```

The build cannot accidentally pick up system libraries or mismatched Go versions.

### 4.3 Single Source of Truth

`flake.nix` replaces the scattered dependency declarations:

| Current | After Nix |
|---------|-----------|
| "Install Go 1.26+ manually" | `go` pinned in flake.nix |
| "Install golangci-lint@latest" | `golangci-lint` pinned by nixpkgs hash |
| "Install just" | `just` pinned in flake.nix |
| CI `setup-go@v4` with matrix | `nix develop` in CI |
| Dockerfile `FROM golang:1.25-alpine` | `nix build` produces the binary |

### 4.4 Cross-Platform Consistency

Nix supports `x86_64-linux`, `aarch64-linux`, `x86_64-darwin`, `aarch64-darwin` natively. The flake can produce builds for all platforms from a single definition — replacing the GoReleaser matrix.

---

## 5. Proposed Architecture

### 5.1 File Structure

```
sqlc-wizard/
├── flake.nix              # Main flake: inputs, outputs (packages, devShells, apps, checks)
├── flake.lock             # Pinned dependency versions (auto-generated)
├── nix/
│   ├── packages/
│   │   └── default.nix    # buildGoModule derivation for sqlc-wizard
│   ├── devshells/
│   │   └── default.nix    # Development shell with all tools
│   └── overlays/
│       └── default.nix    # Custom overlays (e.g., specific golangci-lint version)
├── justfile               # Retained — but recipes delegate to nix where appropriate
├── Dockerfile             # Retained — but can use `nix build` output instead
└── .github/workflows/
    └── ci-cd.yml          # Updated to use `nix develop` for CI
```

### 5.2 Flake Inputs

```nix
inputs = {
  nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  flake-utils.url = "github:numtide/flake-utils";

  # Go tooling overlay (provides newer Go versions, golangci-lint, etc.)
  gopkgs.url = "github:sagikazarmark/go-flake";
  gopkgs.inputs.nixpkgs.follows = "nixpkgs";
};
```

### 5.3 Flake Outputs

| Output | Command | Purpose |
|--------|---------|---------|
| `packages.<system>.default` | `nix build` | Build the `sqlc-wizard` binary |
| `packages.<system>.docker` | `nix build .#docker` | Build Docker image via `dockerTools` |
| `devShells.<system>.default` | `nix develop` | Full dev environment |
| `apps.<system>.default` | `nix run` | Run the built binary |
| `checks.<system>` | `nix flake check` | Run tests, lints, vet |

---

## 6. Migration Plan

### Phase 0: Prerequisites (1 hour)

**Objective:** Resolve blocking issues before Nix introduction.

- [ ] **P0.1** Remove or externalize the local `replace` directive for `go-composable-business-types`. Options:
  - Publish the dependency to a Git repository and use a proper `require`
  - Use a Go workspace (`go.work`) with both projects checked out
  - Make it optional behind a build tag
- [ ] **P0.2** Ensure `generated/` module is self-contained and doesn't require local paths
- [ ] **P0.3** Install Nix on the development machine (if not already present):
  ```bash
  curl --proto '=https' --tlsv1.2 -sSf -L https://install.determinate.systems/nix | sh -s -- install
  ```
  This installs Nix with Flakes enabled by default (Determinate Systems installer).

### Phase 1: Add `flake.nix` Alongside Existing Tooling (2-4 hours)

**Objective:** Introduce Nix without breaking anything. Both systems coexist.

#### Step 1.1: Create `flake.nix`

```nix
{
  description = "SQLC-Wizard — Interactive CLI wizard for generating sqlc configurations";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        version = if (self ? shortRev) then self.shortRev else "dev";
        commit = if (self ? rev) then self.rev else "unknown";
        buildDate = self.lastModifiedDate or "1970-01-01T00:00:00Z";

        ldflags = [
          "-s"
          "-w"
          "-X main.Version=${version}"
          "-X main.Commit=${commit}"
          "-X main.BuildDate=${buildDate}"
        ];
      in
      {
        packages = {
          default = pkgs.buildGoModule {
            pname = "sqlc-wizard";
            inherit version;
            src = ./.;

            # Run: nix build . --no-link 2>&1 | grep 'got:' | cut -d: -f2-
            # Update this hash when go.sum changes
            vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

            subPackages = [ "cmd/sqlc-wizard" ];

            ldflags = ldflags;

            # Exclude generated typespec and local dev files from the build
            exclude = [ "generated" ];

            meta = with pkgs.lib; {
              description = "Interactive CLI wizard for generating sqlc configurations";
              homepage = "https://github.com/LarsArtmann/SQLC-Wizzard";
              license = licenses.mit;
              maintainers = [ ];
              platforms = platforms.unix;
            };
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go_1_26
            golangci-lint
            goreleaser
            just
            dupl
            go-arch-lint
            bun
            nodePackages.typescript
          ];

          shellHook = ''
            echo "SQLC-Wizard development environment"
            echo "Go: $(go version)"
            echo "golangci-lint: $(golangci-lint version)"
            echo "Just: $(just --version)"
          '';
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/sqlc-wizard";
        };
      }
    );
}
```

#### Step 1.2: Generate `flake.lock`

```bash
nix flake lock
```

#### Step 1.3: Fix `vendorHash`

```bash
# First build will fail with a hash mismatch — copy the "got:" hash
nix build . --no-link 2>&1 | grep 'got:'
# Update vendorHash in flake.nix with the actual hash
# Rebuild
nix build
```

#### Step 1.4: Verify

```bash
nix develop --command just build     # Dev shell + existing justfile works
nix build                              # Nix build produces binary
./result/bin/sqlc-wizard version       # Version info is correct
```

### Phase 2: Integrate Nix into CI (2-3 hours)

**Objective:** Replace ad-hoc `go install @latest` with deterministic Nix-based CI.

#### Step 2.1: Update `.github/workflows/ci-cd.yml`

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: DeterminateSystems/nix-installer-action@v16
      - uses: DeterminateSystems/magic-nix-cache-action@v9

      - name: Run tests
        run: nix develop --command go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Run linter
        run: nix develop --command golangci-lint run ./...

      - name: Build
        run: nix build

  # Release job can still use goreleaser for cross-platform builds
  # OR switch to nix-based cross-compilation
```

#### Step 2.2: Add `nix flake check` to CI

```yaml
      - name: Flake check
        run: nix flake check
```

This validates the flake schema, runs all `checks` outputs, and catches issues early.

### Phase 3: Add Nix Checks (1-2 hours)

**Objective:** Replace scattered quality checks with `nix flake check`.

```nix
checks = {
  test = pkgs.runCommandLocal "sqlc-wizard-tests" { } ''
    ${pkgs.go_1_26}/bin/go test -v -race ./...
    touch $out
  '';

  lint = pkgs.runCommandLocal "sqlc-wizard-lint" { } ''
    ${pkgs.golangci-lint}/bin/golangci-lint run ./...
    touch $out
  '';

  vet = pkgs.runCommandLocal "sqlc-wizard-vet" { } ''
    ${pkgs.go_1_26}/bin/go vet ./...
    touch $out
  '';
};
```

### Phase 4: Docker via Nix (1-2 hours)

**Objective:** Replace the Dockerfile with `pkgs.dockerTools.buildImage` for smaller, more reproducible images.

```nix
packages.docker = pkgs.dockerTools.buildImage {
  name = "sqlc-wizard";
  tag = "latest";

  copyToRoot = [
    self.packages.${system}.default
    pkgs.cacert
    pkgs.tzdata
  ];

  config = {
    Entrypoint = [ "/bin/sqlc-wizard" ];
    Cmd = [ "--help" ];
    Env = [ "TZ=UTC" ];
    Labels = {
      "org.opencontainers.image.title" = "SQLC-Wizard";
      "org.opencontainers.image.description" = "Interactive CLI wizard for sqlc configurations";
      "org.opencontainers.image.licenses" = "MIT";
    };
  };
};
```

### Phase 5: Update Justfile (1 hour)

**Objective:** Make the justfile delegate to Nix where beneficial while keeping ergonomic shortcuts.

```make
# Build using Nix (replaces manual go build)
build:
	nix build
	ln -sf $(pwd)/result/bin/sqlc-wizard bin/sqlc-wizard

# Enter development shell
shell:
	nix develop

# Run tests in Nix environment
test:
	nix develop --command go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linter in Nix environment
lint:
	nix develop --command golangci-lint run ./...

# Full verification
verify: build lint test
	nix flake check

# Update flake inputs
update:
	nix flake update
```

### Phase 6: Optional Enhancements (Future)

- [ ] **Use `flake-parts`** for modular flake composition (separate packages, devshells, checks into files)
- [ ] **Pre-commit hooks** via `pre-commit-hooks.nix` input (runs `gofmt`, `golangci-lint`, `go vet` on commit)
- [ ] **Nix-based cross-compilation** to replace GoReleaser for binary releases
- [ ] **Nix-based TypeSpec generation** by adding `nodejs` and `@typespec/compiler` to the dev shell
- [ ] **`direnv` integration** with `.envrc` containing `use flake` — automatic shell activation on `cd`
- [ ] **Nixpkgs-fmt** or **alejandra** as the `formatter` output for `.nix` files

---

## 7. Risk Analysis

### 7.1 Technical Risks

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| `vendorHash` breaks on every `go.sum` change | High | Low | Document the hash update workflow; consider `vendorHash = lib.fakeHash` pattern |
| `go_1_26` not yet in nixpkgs-unstable | Medium | High | Use `gopkgs` overlay or define a custom Go derivation; fallback: use `go_1_24` or overlay |
| Local `replace` directive in go.mod breaks Nix sandbox | High | Critical | Must resolve in Phase 0 before proceeding |
| TypeSpec/bun tooling not well-packaged in Nix | Medium | Medium | Keep TypeSpec generation outside Nix initially (Phase 1-5 skip this) |
| Learning curve for contributors unfamiliar with Nix | Medium | Medium | Keep justfile as primary interface; `nix develop` is transparent |
| `.golangci.yml` experimental build tags (`goexperiment.*`) | Low | Medium | Verify `golangci-lint` from nixpkgs supports these flags; may need overlay |

### 7.2 Compatibility Risks

| Risk | Mitigation |
|------|------------|
| Nix not available on all contributor machines | Justfile remains the primary interface; Nix is additive |
| CI runners need Nix installed | Use `DeterminateSystems/nix-installer-action` — well-maintained, fast |
| GoReleaser workflow changes | Phase 4 is optional; GoReleaser can coexist with Nix Docker builds |

### 7.3 What We Keep

The following tools and files are **not replaced** by Nix — they remain as-is:

| Tool | Reason |
|------|--------|
| `justfile` | Primary developer interface; Nix provides the env, just runs the tasks |
| `.golangci.yml` | Configuration stays the same; Nix only pins the binary version |
| `.goreleaser.yml` | Release automation continues; optionally replaced later |
| `Makefile` | TypeSpec generation; can be moved to Nix dev shell later |
| `Dockerfile` | Retained until Phase 4; Nix Docker image is a replacement, not a requirement |

---

## 8. Success Criteria

| Criterion | Measurement |
|-----------|-------------|
| `nix develop` provides a working dev shell with all tools | `go version`, `golangci-lint version`, `just --version` all succeed |
| `nix build` produces a working binary | `./result/bin/sqlc-wizard version` outputs correct version info |
| CI uses Nix for test/lint/build | No `go install @latest` in CI YAML |
| `nix flake check` passes | Tests, lints, and vet all green |
| No regression in existing workflow | `just dev` still works (now inside `nix develop`) |
| Onboarding reduced to 2 steps | `1. Install Nix`, `2. nix develop` |
| `flake.lock` committed | Reproducible builds across machines and CI |

---

## 9. Reference Implementation

### 9.1 Complete `flake.nix` (Production-Ready Template)

```nix
{
  description = "SQLC-Wizard — Interactive CLI wizard for generating sqlc configurations";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

    flake-utils.url = "github:numtide/flake-utils";

    gopkgs.url = "github:sagikazarmark/go-flake";
    gopkgs.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, gopkgs }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Version from git, fallback to "dev"
        version =
          if (self ? shortRev) then self.shortRev
          else "dev";

        commit =
          if (self ? rev) then self.rev
          else "unknown";

        buildDate =
          builtins.substring 0 10 (self.lastModifiedDate or "19700101T000000Z");

        ldflags = [
          "-s"
          "-w"
          "-X main.Version=${version}"
          "-X main.Commit=${commit}"
          "-X main.BuildDate=${buildDate}"
        ];

        # The main package
        sqlc-wizard = pkgs.buildGoModule {
          pname = "sqlc-wizard";
          inherit version;
          src = ./.;

          # UPDATE: After changing go.sum, run:
          #   nix build . --no-link 2>&1 | grep 'got:' | cut -d: -f2-
          # Paste the result here.
          vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";

          subPackages = [ "cmd/sqlc-wizard" ];

          ldflags = ldflags;

          meta = with pkgs.lib; {
            description = "Interactive CLI wizard for generating sqlc configurations";
            homepage = "https://github.com/LarsArtmann/SQLC-Wizzard";
            license = licenses.mit;
            platforms = platforms.unix;
            mainProgram = "sqlc-wizard";
          };
        };

        # Development tools (pinned versions from nixpkgs)
        devTools = with pkgs; [
          # Go toolchain
          go
          # Linting & quality
          golangci-lint
          dupl
          # Task runner
          just
          # Release
          goreleaser
          # Architecture linting
          go-arch-lint
          # TypeSpec / Node
          bun
        ];

      in
      {
        # --- Packages ---
        packages = {
          default = sqlc-wizard;

          docker = pkgs.dockerTools.buildImage {
            name = "ghcr.io/larsartmann/sqlc-wizard";
            tag = version;

            created = buildDate;

            copyToRoot = [
              sqlc-wizard
              pkgs.cacert
              pkgs.tzdata
              pkgs.busybox
            ];

            config = {
              User = "1000:1000";
              Entrypoint = [ "/bin/sqlc-wizard" ];
              Cmd = [ "--help" ];
              Env = [
                "SSL_CERT_FILE=/etc/ssl/certs/ca-bundle.crt"
                "TZ=UTC"
              ];
              Labels = {
                "org.opencontainers.image.title" = "SQLC-Wizard";
                "org.opencontainers.image.description" = "Interactive CLI wizard for sqlc configurations";
                "org.opencontainers.image.version" = version;
                "org.opencontainers.image.source" = "https://github.com/LarsArtmann/SQLC-Wizzard";
                "org.opencontainers.image.licenses" = "MIT";
                "org.opencontainers.image.vendor" = "Lars Artmann";
              };
            };
          };
        };

        # --- Development Shell ---
        devShells.default = pkgs.mkShell {
          buildInputs = devTools;

          shellHook = ''
            echo "🧙 SQLC-Wizard Development Environment"
            echo ""
            echo "  Go:             $(go version)"
            echo "  golangci-lint:  $(golangci-lint version --format short 2>/dev/null || echo 'available')"
            echo "  Just:           $(just --version)"
            echo ""
            echo "  Quick start:"
            echo "    just build    — Build the binary"
            echo "    just test     — Run all tests"
            echo "    just lint     — Run linters"
            echo "    just dev      — Full dev workflow"
            echo "    nix build     — Nix build"
          '';
        };

        # --- Apps ---
        apps.default = {
          type = "app";
          program = "${sqlc-wizard}/bin/sqlc-wizard";
        };

        # --- Checks ---
        checks = {
          build = sqlc-wizard;

          go-vet = pkgs.runCommandLocal "go-vet" { } ''
            cd ${./.}
            ${pkgs.go}/bin/go vet ./...
            touch $out
          '';
        };

        # --- Formatter ---
        formatter = pkgs.nixpkgs-fmt;
      }
    );
}
```

### 9.2 Updated `justfile` (Post-Migration)

```make
# SQLC-Wizard Justfile (Nix-aware)

default:
	@just --list

# Build the binary using Nix
build:
	@echo "Building sqlc-wizard (Nix)..."
	@mkdir -p bin
	nix build
	@cp -L result/bin/sqlc-wizard bin/sqlc-wizard
	@echo "Build complete: bin/sqlc-wizard"

# Run all tests (in Nix dev shell)
test:
	@echo "Running tests..."
	nix develop --command go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linters (in Nix dev shell)
lint:
	@echo "Running linters..."
	nix develop --command golangci-lint run ./...

# Find code duplicates
find-duplicates:
	@echo "Finding code duplicates..."
	nix develop --command dupl -t 100 -plumbing .

fd: find-duplicates
	@echo "Duplicate detection complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin coverage.txt coverage.html result
	nix develop --command go clean

# Format code
fmt:
	nix develop --command go fmt ./...
	nix develop --command gofmt -s -w .

# Run go vet
vet:
	nix develop --command go vet ./...

# Tidy go modules
tidy:
	nix develop --command go mod tidy

# Install dependencies
deps:
	nix develop --command go mod download

# Install sqlc-wizard locally
install-local: build
	@cp -L bin/sqlc-wizard $(shell go env GOPATH)/bin/sqlc-wizard
	@echo "Installed to GOPATH/bin"

# Full verification
verify: build lint test
	nix flake check

# Development workflow
dev: clean build test find-duplicates
	@echo "Development workflow complete"

# Enter Nix development shell
shell:
	nix develop

# Update Nix flake inputs
update:
	nix flake update

# Generate Go types from TypeSpec
generate-typespec:
	@echo "Generating types from TypeSpec..."
	@mkdir -p generated
	nix develop --command bash -c "tsp compile api/typespec.tsp --emit @typespec/openapi3 --output-dir tsp-output"
	@echo "TypeSpec compilation complete"

# Run benchmarks
bench:
	nix develop --command bash -c 'echo "=== Domain Benchmarks ===" && go test ./internal/domain -bench=. -benchmem && echo "" && echo "=== Adapter Benchmarks ===" && go test ./internal/adapters -bench=. -benchmem'
```

### 9.3 Updated CI Workflow

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
    tags: ["v*"]
  pull_request:
    branches: [main]

jobs:
  check:
    name: Flake Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: DeterminateSystems/nix-installer-action@v16
      - uses: DeterminateSystems/magic-nix-cache-action@v9

      - name: Check flake
        run: nix flake check --no-build

  test:
    name: Test
    runs-on: ubuntu-latest
    needs: check
    steps:
      - uses: actions/checkout@v4

      - uses: DeterminateSystems/nix-installer-action@v16
      - uses: DeterminateSystems/magic-nix-cache-action@v9

      - name: Run tests
        run: nix develop --command go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.txt

  lint:
    name: Lint
    runs-on: ubuntu-latest
    needs: check
    steps:
      - uses: actions/checkout@v4

      - uses: DeterminateSystems/nix-installer-action@v16
      - uses: DeterminateSystems/magic-nix-cache-action@v9

      - name: Run linter
        run: nix develop --command golangci-lint run ./...

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [test, lint]
    steps:
      - uses: actions/checkout@v4

      - uses: DeterminateSystems/nix-installer-action@v16
      - uses: DeterminateSystems/magic-nix-cache-action@v9

      - name: Build
        run: nix build

      - name: Test binary
        run: ./result/bin/sqlc-wizard version

  release:
    name: Release
    runs-on: ubuntu-latest
    needs: [build]
    if: startsWith(github.ref, 'refs/tags/v')
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: DeterminateSystems/nix-installer-action@v16

      - name: Run GoReleaser
        run: nix develop --command goreleaser release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

---

## Appendix: Full Tool Dependency Map

### Current External Dependencies

```
sqlc-wizard (Go binary)
├── Go 1.26.1 (language runtime)
├── golangci-lint (90+ linters)
├── dupl (duplicate detection)
├── go-arch-lint (architecture enforcement)
├── goreleaser (release automation)
├── just (task runner)
├── git (version injection)
├── bun + @typespec/compiler (type generation)
└── Docker (container builds)
```

### Post-Migration Nix Dependency Tree

```
flake.nix
├── nixpkgs (pinned via flake.lock)
│   ├── go (version from nixpkgs or overlay)
│   ├── golangci-lint
│   ├── dupl
│   ├── just
│   ├── goreleaser
│   ├── go-arch-lint
│   ├── bun
│   └── nixpkgs-fmt
├── flake-utils (cross-platform support)
├── gopkgs (Go-specific overlay, optional)
└── buildGoModule (derivation builder)
    ├── vendorHash (pinned go.sum)
    ├── ldflags (version injection)
    └── subPackages = ["cmd/sqlc-wizard"]
```

### Migration Effort Summary

| Phase | Duration | Risk | Reversible? |
|-------|----------|------|-------------|
| Phase 0: Prerequisites | 1h | High (local replace) | N/A — must fix regardless |
| Phase 1: Add flake.nix | 2-4h | Low | Yes — just delete flake.nix |
| Phase 2: Nix in CI | 2-3h | Medium | Yes — revert CI YAML |
| Phase 3: Nix checks | 1-2h | Low | Yes — remove checks |
| Phase 4: Docker via Nix | 1-2h | Low | Yes — keep Dockerfile |
| Phase 5: Update justfile | 1h | Low | Yes — git revert |
| Phase 6: Enhancements | Ongoing | Low | Each is independent |
| **Total (Phases 1-5)** | **7-11h** | | |

---

## Recommended Next Steps

1. **Resolve the local `replace` directive** (Phase 0) — this is a blocker regardless of Nix adoption
2. **Create `flake.nix`** (Phase 1) — start with dev shell only; packages can come later
3. **Try `nix develop`** — verify all tools are available and versions are compatible
4. **Iterate on `vendorHash`** — get the build working before touching CI
5. **Gradually migrate CI** — start with a single job, expand as confidence grows

The key insight: **Nix Flakes are additive, not disruptive.** The existing justfile, Dockerfile, and CI continue to work. Nix simply pins what was previously unpinned.
