# SQLC-Wizard Justfile

# Default recipe
default:
	@just --list

# Build the binary
build:
	@echo "Building sqlc-wizard..."
	@mkdir -p bin
	go build -ldflags "-X main.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')" -o bin/sqlc-wizard cmd/sqlc-wizard/main.go
	@echo "Build complete: bin/sqlc-wizard"

# Run all tests
test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

# Run linters
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Find code duplicates (alias for fd)
find-duplicates:
	@echo "Finding code duplicates..."
	@if command -v dupl >/dev/null 2>&1; then \
		echo "Using dupl tool to find duplicates..."; \
		dupl -t 100 -plumbing . || echo "No duplicates found or dupl encountered an issue"; \
	else \
		echo "dupl not installed. Installing..."; \
		go install github.com/golangci/dupl@latest; \
		echo "Running duplicate detection..."; \
		dupl -t 100 -plumbing . || echo "No duplicates found or dupl encountered an issue"; \
	fi

# Native alias for find-duplicates - improved with better error handling
fd: find-duplicates
	@echo "Duplicate detection complete!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin coverage.txt coverage.html
	go clean

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w .

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Tidy go modules
tidy:
	@echo "Tidying go modules..."
	go mod tidy

# Install dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Install sqlc-wizard locally to GOPATH/bin
install-local: build
	@echo "Installing sqlc-wizard to GOPATH/bin..."
	@go install -ldflags "-X main.Version=$(shell git describe --tags --always --dirty 2>/dev/null || echo 'dev')" ./cmd/sqlc-wizard
	@echo "Installation complete! Run 'sqlc-wizard --help' to verify."

# Run all verification steps (build, lint, test)
verify: build lint test

# Run performance benchmarks
bench:
	@echo "Running performance benchmarks..."
	@echo "=== Domain Events Benchmarks ==="
	@go test ./internal/domain -bench=. -benchmem
	@echo ""
	@echo "=== Adapter Benchmarks ==="
	@go test ./internal/adapters -bench=. -benchmem

# Run performance benchmarks with profiling
bench-profile:
	@echo "Running performance benchmarks with profiling..."
	@mkdir -p profile
	@go test ./internal/domain -bench=. -benchmem -cpuprofile=profile/domain-cpu.prof -memprofile=profile/domain-mem.prof
	@go test ./internal/adapters -bench=. -benchmem -cpuprofile=profile/adapter-cpu.prof -memprofile=profile/adapter-mem.prof
	@echo "Profiles saved in profile/ directory"

# Generate Go types from TypeSpec
generate-typespec:
	@echo "Generating types from TypeSpec..."
	@mkdir -p generated
	@tsp compile api/typespec.tsp --emit @typespec/openapi3 --output-dir tsp-output
	@echo "TypeSpec compilation complete: tsp-output/"

# Development workflow (clean, build, test, find duplicates)
dev: clean build test find-duplicates
	@echo "Development workflow complete"