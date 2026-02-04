package creators

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// CreateConfig contains configuration for project creation.
type CreateConfig struct {
	ProjectName     string
	ProjectType     generated.ProjectType
	Database        generated.DatabaseType
	TemplateData    generated.TemplateData
	Config          *config.SqlcConfig
	IncludeAuth     bool
	IncludeFrontend bool
	Force           bool
}

// ProjectCreator handles creating complete project structures.
type ProjectCreator struct {
	fs  adapters.FileSystemAdapter
	cli adapters.CLIAdapter
}

// NewProjectCreator creates a new project creator.
func NewProjectCreator(fs adapters.FileSystemAdapter, cli adapters.CLIAdapter) *ProjectCreator {
	return &ProjectCreator{
		fs:  fs,
		cli: cli,
	}
}

// CreateProject creates a complete project structure.
func (pc *ProjectCreator) CreateProject(ctx context.Context, config *CreateConfig) error {
	_ = pc.cli.Println(ctx, "🏗️  Creating project structure...")

	// Create directory structure
	if err := pc.createDirectoryStructure(ctx, config); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate sqlc.yaml
	if err := pc.generateSQLCConfig(ctx, config); err != nil {
		return fmt.Errorf("failed to generate sqlc.yaml: %w", err)
	}

	// Generate database schema
	if err := pc.generateDatabaseSchema(ctx, config); err != nil {
		return fmt.Errorf("failed to generate database schema: %w", err)
	}

	// TODO: Full project scaffolding is not yet implemented
	// See GitHub issues for roadmap:
	// - Database schema generation
	// - Query file generation
	// - Migration file generation
	// - Go module structure
	// - Docker configuration
	// - Makefile generation
	// - Development scripts
	// - README generation
	// - Project-specific files
	//
	// For now, ProjectCreator only generates:
	// 1. Directory structure
	// 2. sqlc.yaml configuration file
	//
	// Additional scaffolding will be added based on user feedback and demand.

	return nil
}

// createDirectoryStructure creates the basic directory structure.
func (pc *ProjectCreator) createDirectoryStructure(ctx context.Context, config *CreateConfig) error {
	_ = pc.cli.Println(ctx, "📁 Creating directory structure...")

	dirs := []string{
		"db/schema",
		"db/migrations",
		"internal/db",
		"internal/db/queries",
		"cmd/server",
		"pkg/config",
		"scripts",
		"test",
		"docs",
	}

	// Add project-specific directories based on project type
	// Note: Some project types may not be in generated types yet
	switch config.ProjectType {
	case generated.ProjectTypeMicroservice:
		dirs = append(dirs, "api", "internal/api", "internal/handlers")
	case generated.ProjectTypeHobby:
		dirs = append(dirs, "internal/handlers")
	case generated.ProjectTypeEnterprise:
		dirs = append(dirs, "api", "internal/api", "internal/handlers")
	case generated.ProjectTypeAPIFirst:
		dirs = append(dirs, "api", "internal/api", "internal/handlers")
	case generated.ProjectTypeAnalytics:
		dirs = append(dirs, "internal/analytics", "internal/processors")
	case generated.ProjectTypeTesting:
		dirs = append(dirs, "test/integration", "test/e2e")
	case generated.ProjectTypeMultiTenant:
		dirs = append(dirs, "api", "internal/api", "internal/tenants")
	case generated.ProjectTypeLibrary:
		dirs = append(dirs, "examples", "internal/testutil")
	default:
		return fmt.Errorf("unsupported project type: %s", config.ProjectType)
	}

	for _, dir := range dirs {
		if err := pc.fs.MkdirAll(ctx, dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateSQLCConfig generates the sqlc.yaml file.
func (pc *ProjectCreator) generateSQLCConfig(ctx context.Context, cfg *CreateConfig) error {
	_ = pc.cli.Println(ctx, "⚙️  Generating sqlc.yaml...")

	// Defensive check: ensure config is not nil before marshalling
	if cfg.Config == nil {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "sqlc config is nil: cannot marshal empty configuration to yaml")
	}

	// Convert config to YAML using the marshaller
	yamlContent, err := config.Marshal(cfg.Config)
	if err != nil {
		return fmt.Errorf("failed to convert config to YAML: %w", err)
	}

	return pc.fs.WriteFile(ctx, "sqlc.yaml", yamlContent, 0o644)
}

// generateDatabaseSchema creates database schema file.
func (pc *ProjectCreator) generateDatabaseSchema(ctx context.Context, cfg *CreateConfig) error {
	_ = pc.cli.Println(ctx, "🗄️  Generating database schema...")

	// Build schema content based on project configuration
	// Note: TemplateData must be properly populated from CreateConfig
	templateData := generated.TemplateData{
		ProjectName: cfg.ProjectName,
		ProjectType: cfg.ProjectType,
	}
	schemaContent := pc.buildSchemaSQL(templateData)

	return pc.fs.WriteFile(ctx, "schema.sql", []byte(schemaContent), 0o644)
}

// generateQueryFiles creates SQL query files for sqlc.
func (pc *ProjectCreator) generateQueryFiles(ctx context.Context, cfg *CreateConfig) error {
	_ = pc.cli.Println(ctx, "🔍 Generating SQL query files...")

	// Create queries directory
	if err := pc.fs.MkdirAll(ctx, "queries", 0o755); err != nil {
		return fmt.Errorf("failed to create queries directory: %w", err)
	}

	// Build template data from config
	templateData := generated.TemplateData{
		ProjectName: cfg.ProjectName,
		ProjectType: cfg.ProjectType,
	}

	// Generate basic queries based on database schema
	usersQueries := pc.buildUsersQueries(templateData)
	if err := pc.fs.WriteFile(ctx, "queries/users.sql", []byte(usersQueries), 0o644); err != nil {
		return fmt.Errorf("failed to write users queries: %w", err)
	}

	// Add project-specific queries
	switch templateData.ProjectType {
	case generated.ProjectTypeMicroservice:
		microserviceQueries := pc.buildMicroserviceQueries(templateData)
		if err := pc.fs.WriteFile(ctx, "queries/microservice.sql", []byte(microserviceQueries), 0o644); err != nil {
			return fmt.Errorf("failed to write microservice queries: %w", err)
		}
	case generated.ProjectTypeEnterprise:
		enterpriseQueries := pc.buildEnterpriseQueries(templateData)
		if err := pc.fs.WriteFile(ctx, "queries/enterprise.sql", []byte(enterpriseQueries), 0o644); err != nil {
			return fmt.Errorf("failed to write enterprise queries: %w", err)
		}
	case generated.ProjectTypeAPIFirst:
		apiQueries := pc.buildAPIQueries(templateData)
		if err := pc.fs.WriteFile(ctx, "queries/api.sql", []byte(apiQueries), 0o644); err != nil {
			return fmt.Errorf("failed to write api queries: %w", err)
		}
	case generated.ProjectTypeHobby:
		// Use basic users queries for hobby projects
	case generated.ProjectTypeAnalytics:
		analyticsQueries := pc.buildAnalyticsQueries(templateData)
		if err := pc.fs.WriteFile(ctx, "queries/analytics.sql", []byte(analyticsQueries), 0o644); err != nil {
			return fmt.Errorf("failed to write analytics queries: %w", err)
		}
	case generated.ProjectTypeTesting:
		testingQueries := pc.buildTestingQueries(templateData)
		if err := pc.fs.WriteFile(ctx, "queries/testing.sql", []byte(testingQueries), 0o644); err != nil {
			return fmt.Errorf("failed to write testing queries: %w", err)
		}
	case generated.ProjectTypeMultiTenant:
		tenantQueries := pc.buildMultiTenantQueries(templateData)
		if err := pc.fs.WriteFile(ctx, "queries/tenant.sql", []byte(tenantQueries), 0o644); err != nil {
			return fmt.Errorf("failed to write tenant queries: %w", err)
		}
	case generated.ProjectTypeLibrary:
		// Use basic users queries for library projects
	default:
		return fmt.Errorf("unsupported project type: %s", templateData.ProjectType)
	}

	return nil
}

// generateGoModuleStructure creates basic Go module structure.
func (pc *ProjectCreator) generateGoModuleStructure(ctx context.Context, cfg *CreateConfig) error {
	_ = pc.cli.Println(ctx, "📦 Generating Go module structure...")

	// Create go.mod
	goModContent := pc.buildGoMod(cfg.TemplateData)
	if err := pc.fs.WriteFile(ctx, "go.mod", []byte(goModContent), 0o644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Create main.go for executable projects
	if cfg.TemplateData.ProjectType == generated.ProjectTypeAPIFirst ||
		cfg.TemplateData.ProjectType == generated.ProjectTypeMicroservice {
		mainGoContent := pc.buildMainGo(cfg.TemplateData)
		if err := pc.fs.WriteFile(ctx, "main.go", []byte(mainGoContent), 0o644); err != nil {
			return fmt.Errorf("failed to write main.go: %w", err)
		}
	}

	// Create basic package structure
	if err := pc.fs.MkdirAll(ctx, "internal/db", 0o755); err != nil {
		return fmt.Errorf("failed to create db package directory: %w", err)
	}

	// Create db package file
	dbGoContent := pc.buildDBPackage(cfg.TemplateData)
	if err := pc.fs.WriteFile(ctx, "internal/db/db.go", []byte(dbGoContent), 0o644); err != nil {
		return fmt.Errorf("failed to write db package: %w", err)
	}

	return nil
}

// buildSchemaSQL creates SQL schema content.
func (pc *ProjectCreator) buildSchemaSQL(data generated.TemplateData) string {
	schema := "-- Database schema for " + data.ProjectName + "\n"
	schema += "-- Generated by SQLC-Wizard\n\n"

	// Basic user table (common in most projects)
	schema += pc.createUserTable(data)

	// Add project-specific tables based on project type
	switch data.ProjectType {
	case generated.ProjectTypeMicroservice:
		schema += pc.createMicroserviceTables(data)
	case generated.ProjectTypeEnterprise:
		schema += pc.createEnterpriseTables(data)
	case generated.ProjectTypeAPIFirst:
		schema += pc.createAPIFirstTables(data)
	case generated.ProjectTypeHobby:
		schema += pc.createHobbyTables(data)
	case generated.ProjectTypeAnalytics:
		schema += pc.createAnalyticsTables(data)
	case generated.ProjectTypeTesting:
		schema += pc.createTestingTables(data)
	case generated.ProjectTypeMultiTenant:
		schema += pc.createMultiTenantTables(data)
	case generated.ProjectTypeLibrary:
		schema += pc.createLibraryTables(data)
	default:
		panic(fmt.Sprintf("unsupported project type: %s", data.ProjectType))
	}

	schema += "\n-- Indexes for performance\n"
	schema += pc.createBasicIndexes(data)

	return schema
}

// createUserTable creates users table.
func (pc *ProjectCreator) createUserTable(data generated.TemplateData) string {
	return `
-- Users table for authentication and user management
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add basic indexes for users
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);
`
}

// createMicroserviceTables creates tables for microservice projects.
func (pc *ProjectCreator) createMicroserviceTables(data generated.TemplateData) string {
	return `
-- API tokens table for microservice authentication
CREATE TABLE api_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_api_tokens_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_api_tokens_user_id ON api_tokens(user_id);
CREATE INDEX idx_api_tokens_expires_at ON api_tokens(expires_at);
`
}

// createEnterpriseTables creates tables for enterprise projects.
func (pc *ProjectCreator) createEnterpriseTables(data generated.TemplateData) string {
	return `
-- Audit log table for enterprise compliance
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    resource_id VARCHAR(255),
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fk_audit_logs_user FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
`
}

// createAPIFirstTables creates tables for API-first projects.
func (pc *ProjectCreator) createAPIFirstTables(data generated.TemplateData) string {
	return `
-- Rate limiting table for API-first projects
CREATE TABLE rate_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_identifier VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    max_requests INTEGER NOT NULL DEFAULT 100,
    window_seconds INTEGER NOT NULL DEFAULT 60,
    current_requests INTEGER DEFAULT 0,
    window_start TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(client_identifier, endpoint)
);

CREATE INDEX idx_rate_limits_client ON rate_limits(client_identifier);
CREATE INDEX idx_rate_limits_window ON rate_limits(window_start);
`
}

// createBasicIndexes creates basic performance indexes.
func (pc *ProjectCreator) createBasicIndexes(data generated.TemplateData) string {
	return `
-- Basic performance indexes
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at ON users(created_at);
`
}

// NOTE: Additional scaffolding methods will be implemented based on demand
// See the TODO in CreateProject for planned features

// buildUsersQueries creates basic user management queries.
func (pc *ProjectCreator) buildUsersQueries(data generated.TemplateData) string {
	return `-- name: CreateUser :one
INSERT INTO users (email, password_hash, full_name) 
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users 
SET full_name = COALESCE(sqlc.narg('full_name'), full_name),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');
`
}

// buildAnalyticsQueries creates analytics query templates.
func (pc *ProjectCreator) buildAnalyticsQueries(data generated.TemplateData) string {
	return `-- name: CreateEvent :one
INSERT INTO events (event_type, payload, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetEventByID :one
SELECT * FROM events WHERE id = $1;

-- name: ListEvents :many
SELECT * FROM events
WHERE created_at >= NOW() - INTERVAL '24 hours'
ORDER BY created_at DESC;
`
}

// buildTestingQueries creates testing query templates.
func (pc *ProjectCreator) buildTestingQueries(data generated.TemplateData) string {
	return `-- name: CreateFixture :one
INSERT INTO fixtures (name, data)
VALUES ($1, $2)
RETURNING *;

-- name: GetFixture :one
SELECT * FROM fixtures WHERE name = $1;

-- name: ListFixtures :many
SELECT * FROM fixtures ORDER BY name;
`
}

// buildMultiTenantQueries creates multi-tenant query templates.
func (pc *ProjectCreator) buildMultiTenantQueries(data generated.TemplateData) string {
	return `-- name: CreateTenant :one
INSERT INTO tenants (name, domain)
VALUES ($1, $2)
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants WHERE id = $1;

-- name: GetTenantByDomain :one
SELECT * FROM tenants WHERE domain = $1;

-- name: ListTenants :many
SELECT * FROM tenants ORDER BY created_at DESC;
`
}

// createHobbyTables creates tables for hobby projects.
func (pc *ProjectCreator) createHobbyTables(data generated.TemplateData) string {
	return `
-- Settings table for hobby project configuration
CREATE TABLE settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key VARCHAR(255) NOT NULL,
    value TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, key)
);

CREATE INDEX idx_settings_user_id ON settings(user_id);
`
}

// createAnalyticsTables creates tables for analytics projects.
func (pc *ProjectCreator) createAnalyticsTables(data generated.TemplateData) string {
	return `
-- Events table for analytics data
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(255) NOT NULL,
    payload JSONB,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_created_at ON events(created_at);

-- Aggregations table for pre-computed analytics
CREATE TABLE aggregations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    metric_name VARCHAR(255) NOT NULL,
    metric_value NUMERIC NOT NULL,
    aggregation_period VARCHAR(50) NOT NULL,
    computed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_aggregations_metric ON aggregations(metric_name);
`
}

// createTestingTables creates tables for testing projects.
func (pc *ProjectCreator) createTestingTables(data generated.TemplateData) string {
	return `
-- Fixtures table for test data
CREATE TABLE fixtures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Test runs table for tracking test execution
CREATE TABLE test_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    duration_ms INTEGER,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_test_runs_status ON test_runs(status);
CREATE INDEX idx_test_runs_started_at ON test_runs(started_at);
`
}

// createMultiTenantTables creates tables for multi-tenant projects.
func (pc *ProjectCreator) createMultiTenantTables(data generated.TemplateData) string {
	return `
-- Tenants table for multi-tenancy support
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    domain VARCHAR(255) NOT NULL UNIQUE,
    settings JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_tenants_domain ON tenants(domain);

-- Add tenant_id column to users
ALTER TABLE users ADD COLUMN tenant_id UUID REFERENCES tenants(id) ON DELETE CASCADE;

CREATE INDEX idx_users_tenant_id ON users(tenant_id);
`
}

// createLibraryTables creates tables for library projects.
func (pc *ProjectCreator) createLibraryTables(data generated.TemplateData) string {
	return `
-- Examples table for library examples
CREATE TABLE examples (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    code_snippet TEXT NOT NULL,
    language VARCHAR(50),
    difficulty VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_examples_language ON examples(language);
CREATE INDEX idx_examples_difficulty ON examples(difficulty);
`
}

// buildGoMod creates go.mod content.
func (pc *ProjectCreator) buildGoMod(data generated.TemplateData) string {
	return `module ` + data.ProjectName + `

go 1.21

require (
	github.com/jackc/pgx/v5 v5.4.3
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/text v0.12.0 // indirect
)
`
}

// buildDBPackage creates basic database package.
func (pc *ProjectCreator) buildDBPackage(data generated.TemplateData) string {
	return `package db

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

// Connect creates a new database connection
func Connect() (*DB, error) {
	dsn := "postgres://localhost/` + data.ProjectName + `?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
`
}

// NOTE: Additional scaffolding methods will be implemented based on demand
// See TODO in CreateProject for planned features

// buildMicroserviceQueries creates microservice-specific SQL queries.
func (pc *ProjectCreator) buildMicroserviceQueries(data generated.TemplateData) string {
	return `-- name: CreateAPIToken :one
INSERT INTO api_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAPITokenByHash :one
SELECT * FROM api_tokens WHERE token_hash = $1;

-- name: ListAPITokensForUser :many
SELECT * FROM api_tokens
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: RevokeAPIToken :exec
DELETE FROM api_tokens WHERE id = $1;

-- name: CleanupExpiredTokens :exec
DELETE FROM api_tokens WHERE expires_at < NOW();`
}

// buildEnterpriseQueries creates enterprise-specific SQL queries.
func (pc *ProjectCreator) buildEnterpriseQueries(data generated.TemplateData) string {
	return `-- name: CreateAuditLog :one
INSERT INTO audit_logs (user_id, action, resource_type, resource_id, old_values, new_values, ip_address, user_agent)
VALUES (sqlc.narg('user_id'), $1, $2, $3, sqlc.narg('old_values'), sqlc.narg('new_values'), sqlc.narg('ip_address'), sqlc.narg('user_agent'))
RETURNING *;

-- name: ListAuditLogs :many
SELECT * FROM audit_logs
WHERE (sqlc.narg('user_id')::UUID IS NULL OR user_id = sqlc.narg('user_id')::UUID)
  AND (sqlc.narg('resource_type')::TEXT IS NULL OR resource_type = sqlc.narg('resource_type')::TEXT)
  AND (sqlc.narg('action')::TEXT IS NULL OR action = sqlc.narg('action')::TEXT)
ORDER BY created_at DESC
LIMIT sqlc.narg('limit')::INTEGER OFFSET sqlc.narg('offset')::INTEGER;

-- name: GetAuditLogByID :one
SELECT * FROM audit_logs WHERE id = $1;`
}

// buildAPIQueries creates API-first specific SQL queries.
func (pc *ProjectCreator) buildAPIQueries(data generated.TemplateData) string {
	return `-- name: CheckRateLimit :one
SELECT * FROM rate_limits
WHERE client_identifier = $1 AND endpoint = $2
ORDER BY window_start DESC
LIMIT 1;

-- name: CreateRateLimit :one
INSERT INTO rate_limits (client_identifier, endpoint, max_requests, window_seconds, current_requests, window_start)
VALUES ($1, $2, $3, $4, 1, NOW())
ON CONFLICT (client_identifier, endpoint) DO UPDATE SET
    current_requests = 1,
    window_start = NOW()
RETURNING *;

-- name: IncrementRateLimit :one
UPDATE rate_limits
SET current_requests = current_requests + 1
WHERE id = $1 AND current_requests < max_requests
RETURNING *;

-- name: CleanupExpiredRateLimits :exec
DELETE FROM rate_limits
WHERE window_start + (window_seconds || ' seconds')::INTERVAL < NOW();`
}

// buildMainGo creates main.go for executable projects.
func (pc *ProjectCreator) buildMainGo(data generated.TemplateData) string {
	return `package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"` + data.Package.Path + `/internal/db"
)

func main() {
	// Initialize database connection
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create HTTP server
	mux := http.NewServeMux()

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Add API routes
	// TODO: Add your API routes here

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped")
}`
}
