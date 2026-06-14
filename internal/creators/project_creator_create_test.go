package creators_test

import (
	"context"
	"io/fs"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/creators"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
			Expect(
				mockCLI.printedLines,
			).To(ContainElement(ContainSubstring("Creating project structure")))
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
			Expect(
				mockCLI.printedLines,
			).To(ContainElement(ContainSubstring("Creating project structure")))
			Expect(
				mockCLI.printedLines,
			).To(ContainElement(ContainSubstring("Creating directory structure")))
			Expect(
				mockCLI.printedLines,
			).To(ContainElement(ContainSubstring("Generating sqlc.yaml")))
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

		testEnumAssignment("ProjectType", allProjectTypes,
			func(cfg *creators.CreateConfig, pt generated.ProjectType) { cfg.ProjectType = pt },
			func(cfg *creators.CreateConfig) generated.ProjectType { return cfg.ProjectType })

		testEnumAssignment("Database", allDatabaseTypes,
			func(cfg *creators.CreateConfig, dt generated.DatabaseType) { cfg.Database = dt },
			func(cfg *creators.CreateConfig) generated.DatabaseType { return cfg.Database })
	})
})
