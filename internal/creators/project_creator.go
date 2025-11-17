package creators

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// CreateConfig contains configuration for project creation
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

// ProjectCreator handles creating complete project structures
type ProjectCreator struct {
	fs  adapters.FileSystemAdapter
	cli adapters.CLIAdapter
}

// NewProjectCreator creates a new project creator
func NewProjectCreator(fs adapters.FileSystemAdapter, cli adapters.CLIAdapter) *ProjectCreator {
	return &ProjectCreator{
		fs:  fs,
		cli: cli,
	}
}

// CreateProject creates a complete project structure
func (pc *ProjectCreator) CreateProject(config *CreateConfig) error {
	pc.cli.Println("🏗️  Creating project structure...")

	// Create directory structure
	if err := pc.createDirectoryStructure(config); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Generate sqlc.yaml
	if err := pc.generateSQLCConfig(config); err != nil {
		return fmt.Errorf("failed to generate sqlc.yaml: %w", err)
	}

	// Generate database schemas
	if err := pc.generateSchemas(config); err != nil {
		return fmt.Errorf("failed to generate schemas: %w", err)
	}

	// Generate query files
	if err := pc.generateQueries(config); err != nil {
		return fmt.Errorf("failed to generate queries: %w", err)
	}

	// Generate migration files
	if err := pc.generateMigrations(config); err != nil {
		return fmt.Errorf("failed to generate migrations: %w", err)
	}

	// Generate Go module and basic structure
	if err := pc.generateGoStructure(config); err != nil {
		return fmt.Errorf("failed to generate Go structure: %w", err)
	}

	// Generate Docker configuration
	if err := pc.generateDockerConfig(config); err != nil {
		return fmt.Errorf("failed to generate Docker configuration: %w", err)
	}

	// Generate Makefile
	if err := pc.generateMakefile(config); err != nil {
		return fmt.Errorf("failed to generate Makefile: %w", err)
	}

	// Generate development scripts
	if err := pc.generateDevScripts(config); err != nil {
		return fmt.Errorf("failed to generate development scripts: %w", err)
	}

	// Generate README
	if err := pc.generateREADME(config); err != nil {
		return fmt.Errorf("failed to generate README: %w", err)
	}

	// Generate project-specific files
	if err := pc.generateProjectSpecificFiles(config); err != nil {
		return fmt.Errorf("failed to generate project-specific files: %w", err)
	}

	return nil
}

// createDirectoryStructure creates the basic directory structure
func (pc *ProjectCreator) createDirectoryStructure(config *CreateConfig) error {
	pc.cli.Println("📁 Creating directory structure...")

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

	// Add project-specific directories
	switch config.ProjectType {
	case generated.ProjectTypeFullstack:
		dirs = append(dirs, "web", "web/src", "web/public", "internal/api")
	case generated.ProjectTypeAPIFirst:
		dirs = append(dirirs, "api", "internal/api", "internal/handlers")
	case generated.ProjectTypeLibrary:
		dirs = append(dirs, "examples", "internal/testutil")
	}

	for _, dir := range dirs {
		if err := pc.fs.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// generateSQLCConfig generates the sqlc.yaml file
func (pc *ProjectCreator) generateSQLCConfig(config *CreateConfig) error {
	pc.cli.Println("⚙️  Generating sqlc.yaml...")

	// Convert config to YAML and write
	yamlContent, err := config.Config.ToYAML()
	if err != nil {
		return fmt.Errorf("failed to convert config to YAML: %w", err)
	}

	return pc.fs.WriteFile("sqlc.yaml", []byte(yamlContent), 0o644)
}

// generateSchemas generates database schema files
func (pc *ProjectCreator) generateSchemas(config *CreateConfig) error {
	pc.cli.Println("🗄️  Generating database schemas...")

	// Generate users table schema
	usersSchema := pc.generateUsersSchema(config.Database)
	if err := pc.fs.WriteFile("db/schema/001_users.sql", []byte(usersSchema), 0o644); err != nil {
		return err
	}

	// Add project-specific schemas
	if config.ProjectType == generated.ProjectTypeMicroservice || config.ProjectType == generated.ProjectTypeFullstack {
		postsSchema := pc.generatePostsSchema(config.Database)
		if err := pc.fs.WriteFile("db/schema/002_posts.sql", []byte(postsSchema), 0o644); err != nil {
			return err
		}
	}

	return nil
}

// generateQueries generates SQL query files
func (pc *ProjectCreator) generateQueries(config *CreateConfig) error {
	pc.cli.Println("🔍 Generating query files...")

	// Generate users queries
	usersQueries := pc.generateUsersQueries(config.Database)
	if err := pc.fs.WriteFile("internal/db/queries/users.sql", []byte(usersQueries), 0o644); err != nil {
		return err
	}

	return nil
}

// generateMigrations generates migration files
func (pc *ProjectCreator) generateMigrations(config *CreateConfig) error {
	pc.cli.Println("🔄 Generating migration files...")

	// Generate up migration
	upMigration := pc.generateUpMigration(config.Database)
	if err := pc.fs.WriteFile("db/migrations/001_create_tables.up.sql", []byte(upMigration), 0o644); err != nil {
		return err
	}

	// Generate down migration
	downMigration := pc.generateDownMigration(config.Database)
	if err := pc.fs.WriteFile("db/migrations/001_create_tables.down.sql", []byte(downMigration), 0o644); err != nil {
		return err
	}

	return nil
}

// generateGoStructure generates Go module and basic structure
func (pc *ProjectCreator) generateGoStructure(config *CreateConfig) error {
	pc.cli.Println("🐹 Generating Go structure...")

	// Generate go.mod
	goMod := fmt.Sprintf(`module %s

go 1.21

require (
	github.com/jackc/pgx/v5 v5.4.3
	github.com/lib/pq v1.10.9
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)`, config.ProjectName)

	if err := pc.fs.WriteFile("go.mod", []byte(goMod), 0o644); err != nil {
		return err
	}

	// Generate main.go for server
	mainGo := pc.generateMainGo(config)
	if err := pc.fs.WriteFile("cmd/server/main.go", []byte(mainGo), 0o644); err != nil {
		return err
	}

	return nil
}

// generateDockerConfig generates Docker configuration
func (pc *ProjectCreator) generateDockerConfig(config *CreateConfig) error {
	pc.cli.Println("🐳 Generating Docker configuration...")

	// Generate Dockerfile
	dockerfile := pc.generateDockerfile(config)
	if err := pc.fs.WriteFile("Dockerfile", []byte(dockerfile), 0o644); err != nil {
		return err
	}

	// Generate docker-compose.yml
	dockerCompose := pc.generateDockerCompose(config)
	if err := pc.fs.WriteFile("docker-compose.yml", []byte(dockerCompose), 0o644); err != nil {
		return err
	}

	return nil
}

// generateMakefile generates Makefile with common tasks
func (pc *ProjectCreator) generateMakefile(config *CreateConfig) error {
	pc.cli.Println("🔨 Generating Makefile...")

	makefile := pc.generateMakefileContent(config)
	return pc.fs.WriteFile("Makefile", []byte(makefile), 0o644)
}

// generateDevScripts generates development scripts
func (pc *ProjectCreator) generateDevScripts(config *CreateConfig) error {
	pc.cli.Println("📜 Generating development scripts...")

	// Generate dev.sh
	devScript := pc.generateDevScript(config)
	if err := pc.fs.WriteFile("scripts/dev.sh", []byte(devScript), 0o755); err != nil {
		return err
	}

	// Generate migrate.sh
	migrateScript := pc.generateMigrateScript(config)
	if err := pc.fs.WriteFile("scripts/migrate.sh", []byte(migrateScript), 0o755); err != nil {
		return err
	}

	return nil
}

// generateREADME generates README.md
func (pc *ProjectCreator) generateREADME(config *CreateConfig) error {
	pc.cli.Println("📖 Generating README...")

	readme := pc.generateReadmeContent(config)
	return pc.fs.WriteFile("README.md", []byte(readme), 0o644)
}

// generateProjectSpecificFiles generates project-specific additional files
func (pc *ProjectCreator) generateProjectSpecificFiles(config *CreateConfig) error {
	switch config.ProjectType {
	case generated.ProjectTypeFullstack:
		return pc.generateFullstackFiles(config)
	case generated.ProjectTypeAPIFirst:
		return pc.generateAPIFiles(config)
	case generated.ProjectTypeLibrary:
		return pc.generateLibraryFiles(config)
	}
	return nil
}
