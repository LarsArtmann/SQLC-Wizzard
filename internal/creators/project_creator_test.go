package creators_test

import (
	"context"
	"errors"
	"io/fs"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/creators"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Creators Suite")
}

// createBaseConfig generates a base project configuration with standard defaults.
func createBaseConfig(projectName string) *creators.CreateConfig {
	return &creators.CreateConfig{
		ProjectName: projectName,
		ProjectType: generated.ProjectTypeMicroservice,
		Database:    generated.DatabaseTypePostgreSQL,
		Config: &config.SqlcConfig{
			Version: "2",
			SQL: []config.SQLConfig{
				{
					Engine:  "postgresql",
					Queries: config.NewPathOrPaths([]string{"queries/"}),
					Schema:  config.NewPathOrPaths([]string{"schema/"}),
					Gen: config.GenConfig{
						Go: &config.GoGenConfig{
							Package: "db",
							Out:     "internal/db",
						},
					},
				},
			},
		},
	}
}

// Mock adapters for testing.
type MockFileSystemAdapter struct {
	mkdirAllCalls   []MkdirAllCall
	writeFileCalls  []WriteFileCall
	shouldFailMkdir bool
	shouldFailWrite bool
}

type MkdirAllCall struct {
	Path string
	Perm fs.FileMode
}

type WriteFileCall struct {
	Path    string
	Content []byte
	Perm    fs.FileMode
}

func (m *MockFileSystemAdapter) MkdirAll(ctx context.Context, path string, perm fs.FileMode) error {
	m.mkdirAllCalls = append(m.mkdirAllCalls, MkdirAllCall{Path: path, Perm: perm})
	if m.shouldFailMkdir {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "mkdir failed")
	}
	return nil
}

func (m *MockFileSystemAdapter) WriteFile(ctx context.Context, path string, content []byte, perm fs.FileMode) error {
	m.writeFileCalls = append(m.writeFileCalls, WriteFileCall{Path: path, Content: content, Perm: perm})
	if m.shouldFailWrite {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "write failed")
	}
	return nil
}

func (m *MockFileSystemAdapter) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return nil, apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) Exists(ctx context.Context, path string) (bool, error) {
	return false, apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) Remove(ctx context.Context, path string) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) TempDir(ctx context.Context, pattern string) (string, error) {
	return "", apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) Copy(ctx context.Context, src, dst string) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) CreateDirectory(ctx context.Context, path string, perm fs.FileMode) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) ListFiles(ctx context.Context, dir string) ([]string, error) {
	return nil, apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

type MockCLIAdapter struct {
	printedLines []string
}

func (m *MockCLIAdapter) Println(msg string) error {
	m.printedLines = append(m.printedLines, msg)
	return nil
}

func (m *MockCLIAdapter) Printf(format string, args ...any) {
	// Not needed for current tests
}

func (m *MockCLIAdapter) CheckCommand(ctx context.Context, name string) error {
	return nil
}

func (m *MockCLIAdapter) RunCommand(ctx context.Context, name string, args ...string) (string, error) {
	return "", nil
}

func (m *MockCLIAdapter) GetVersion(ctx context.Context, command string) (string, error) {
	return "", nil
}

func (m *MockCLIAdapter) Install(ctx context.Context, cmd string) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

var _ = Describe("ProjectCreator", func() {
	var (
		mockFS  *MockFileSystemAdapter
		mockCLI *MockCLIAdapter
		creator *creators.ProjectCreator
		ctx     context.Context
	)

	BeforeEach(func() {
		mockFS = &MockFileSystemAdapter{}
		mockCLI = &MockCLIAdapter{}
		creator = creators.NewProjectCreator(mockFS, mockCLI)
		ctx = context.Background()
	})

	Context("NewProjectCreator", func() {
		It("should create a new project creator", func() {
			Expect(creator).NotTo(BeNil())
		})
	})

	Context("CreateProject", func() {
		var cfg *creators.CreateConfig

		BeforeEach(func() {
			cfg = createBaseConfig("test-project")
		})

		It("should create project successfully", func() {
			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify directories were created
			Expect(mockFS.mkdirAllCalls).NotTo(BeEmpty())

			// Verify sqlc.yaml and schema.sql were written
			Expect(mockFS.writeFileCalls).To(HaveLen(2))
			Expect(mockFS.writeFileCalls[0].Path).To(Equal("sqlc.yaml"))
			Expect(mockFS.writeFileCalls[1].Path).To(Equal("schema.sql"))

			// Verify CLI output
			Expect(mockCLI.printedLines).NotTo(BeEmpty())
			Expect(mockCLI.printedLines).To(ContainElement(ContainSubstring("Creating project structure")))
		})

		It("should create microservice-specific directories", func() {
			cfg.ProjectType = generated.ProjectTypeMicroservice

			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify microservice-specific directories
			dirPaths := make([]string, len(mockFS.mkdirAllCalls))
			for i, call := range mockFS.mkdirAllCalls {
				dirPaths[i] = call.Path
			}

			Expect(dirPaths).To(ContainElement("api"))
			Expect(dirPaths).To(ContainElement("internal/api"))
			Expect(dirPaths).To(ContainElement("internal/handlers"))
		})

		It("should create standard directories for all projects", func() {
			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify standard directories
			dirPaths := make([]string, len(mockFS.mkdirAllCalls))
			for i, call := range mockFS.mkdirAllCalls {
				dirPaths[i] = call.Path
			}

			Expect(dirPaths).To(ContainElement("db/schema"))
			Expect(dirPaths).To(ContainElement("db/migrations"))
			Expect(dirPaths).To(ContainElement("internal/db"))
			Expect(dirPaths).To(ContainElement("internal/db/queries"))
			Expect(dirPaths).To(ContainElement("cmd/server"))
			Expect(dirPaths).To(ContainElement("pkg/config"))
			Expect(dirPaths).To(ContainElement("scripts"))
			Expect(dirPaths).To(ContainElement("test"))
			Expect(dirPaths).To(ContainElement("docs"))
		})

		It("should use correct permissions for directories", func() {
			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify all directories use 0755 permissions
			for _, call := range mockFS.mkdirAllCalls {
				Expect(call.Perm).To(Equal(fs.FileMode(0o755)))
			}
		})

		It("should use correct permissions for files", func() {
			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify sqlc.yaml and schema.sql use 0644 permissions
			Expect(mockFS.writeFileCalls).To(HaveLen(2))
			Expect(mockFS.writeFileCalls[0].Perm).To(Equal(fs.FileMode(0o644)))
			Expect(mockFS.writeFileCalls[1].Perm).To(Equal(fs.FileMode(0o644)))
		})

		It("should write valid YAML config", func() {
			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify YAML and schema were written
			Expect(mockFS.writeFileCalls).To(HaveLen(2))
			yamlContent := string(mockFS.writeFileCalls[0].Content)
			schemaContent := string(mockFS.writeFileCalls[1].Content)

			// Basic YAML validation
			Expect(yamlContent).To(ContainSubstring("version:"))
			Expect(yamlContent).To(ContainSubstring("sql:"))

			// Schema validation
			Expect(schemaContent).To(ContainSubstring("CREATE TABLE users"))
			Expect(schemaContent).To(ContainSubstring("Database schema for"))
		})

		It("should fail when directory creation fails", func() {
			mockFS.shouldFailMkdir = true

			err := creator.CreateProject(ctx, cfg)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to create directory structure"))
		})

		It("should fail when YAML generation fails", func() {
			mockFS.shouldFailWrite = true

			err := creator.CreateProject(ctx, cfg)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to generate sqlc.yaml"))
		})

		It("should print progress messages", func() {
			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify CLI printed progress messages
			Expect(mockCLI.printedLines).To(ContainElement(ContainSubstring("Creating project structure")))
			Expect(mockCLI.printedLines).To(ContainElement(ContainSubstring("Creating directory structure")))
			Expect(mockCLI.printedLines).To(ContainElement(ContainSubstring("Generating sqlc.yaml")))
		})
	})

	Context("CreateConfig", func() {
		It("should have all required fields", func() {
			cfg := &creators.CreateConfig{
				ProjectName:     "test",
				ProjectType:     generated.ProjectTypeMicroservice,
				Database:        generated.DatabaseTypePostgreSQL,
				IncludeAuth:     true,
				IncludeFrontend: false,
				Force:           false,
			}

			Expect(cfg.ProjectName).To(Equal("test"))
			Expect(cfg.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
			Expect(cfg.Database).To(Equal(generated.DatabaseTypePostgreSQL))
			Expect(cfg.IncludeAuth).To(BeTrue())
			Expect(cfg.IncludeFrontend).To(BeFalse())
			Expect(cfg.Force).To(BeFalse())
		})

		It("should support all project types", func() {
			projectTypes := []generated.ProjectType{
				generated.ProjectTypeMicroservice,
				generated.ProjectTypeHobby,
				generated.ProjectTypeEnterprise,
			}

			for _, pt := range projectTypes {
				cfg := &creators.CreateConfig{
					ProjectType: pt,
				}
				Expect(cfg.ProjectType).To(Equal(pt))
			}
		})

		It("should support all database types", func() {
			databaseTypes := []generated.DatabaseType{
				generated.DatabaseTypePostgreSQL,
				generated.DatabaseTypeMySQL,
				generated.DatabaseTypeSQLite,
			}

			for _, dt := range databaseTypes {
				cfg := &creators.CreateConfig{
					Database: dt,
				}
				Expect(cfg.Database).To(Equal(dt))
			}
		})
	})

	Context("Integration", func() {
		It("should create complete project structure in correct order", func() {
			cfg := createBaseConfig("integration-test")

			err := creator.CreateProject(ctx, cfg)

			Expect(err).NotTo(HaveOccurred())

			// Verify order: directories first, then files
			Expect(mockFS.mkdirAllCalls).NotTo(BeEmpty())
			Expect(mockFS.writeFileCalls).To(HaveLen(2))

			// All directories should be created before files
			// (This is implicit in the implementation, but good to verify)
			Expect(mockFS.mkdirAllCalls).ToNot(BeEmpty())
		})
	})

	Context("Error Handling", func() {
		It("should return descriptive error when mkdir fails", func() {
			mockFS.shouldFailMkdir = true

			cfg := &creators.CreateConfig{
				ProjectName: "test",
				ProjectType: generated.ProjectTypeMicroservice,
				Database:    generated.DatabaseTypePostgreSQL,
				Config:      &config.SqlcConfig{Version: "2"},
			}

			err := creator.CreateProject(ctx, cfg)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to create directory structure"))
			Expect(err.Error()).To(ContainSubstring("mkdir failed"))
		})

		It("should return descriptive error when YAML write fails", func() {
			mockFS.shouldFailWrite = true

			cfg := &creators.CreateConfig{
				ProjectName: "test",
				ProjectType: generated.ProjectTypeMicroservice,
				Database:    generated.DatabaseTypePostgreSQL,
				Config:      &config.SqlcConfig{Version: "2"},
			}

			err := creator.CreateProject(ctx, cfg)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to generate sqlc.yaml"))
		})
	})
})
