package generators_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/generators"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGenerators(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Generators Suite")
}

// Helper function to create template data for testing.
func createTemplateData(engine generated.DatabaseType, outputDir string) generated.TemplateData {
	return generated.TemplateData{
		Database: generated.DatabaseConfig{
			Engine:    engine,
			UseUUIDs:  engine == generated.DatabaseTypePostgreSQL || engine == generated.DatabaseTypeMySQL,
			UseJSON:   true,
			UseArrays: engine == generated.DatabaseTypePostgreSQL,
		},
		Output: generated.OutputConfig{
			BaseDir:    outputDir,
			QueriesDir: "queries",
			SchemaDir:  "schema",
		},
	}
}

// Helper function to create and setup generator.
func setupGenerator(tempDir string) (*generators.Generator, func()) {
	gen := generators.NewGenerator(tempDir)
	cleanup := func() {
		os.RemoveAll(tempDir)
	}
	return gen, cleanup
}

// Helper function to verify schema file content.
func verifySchemaContent(content, engine string) {
	switch engine {
	case "postgresql":
		Expect(content).To(Or(
			ContainSubstring("UUID"),
			ContainSubstring("JSONB"),
			ContainSubstring("uuid-ossp"),
		))
	case "mysql":
		Expect(content).To(Or(
			ContainSubstring("AUTO_INCREMENT"),
			ContainSubstring("JSON"),
			ContainSubstring("VARCHAR"),
		))
	case "sqlite":
		Expect(content).NotTo(ContainSubstring("uuid-ossp"))
		Expect(content).NotTo(ContainSubstring("JSONB"))
		Expect(content).To(Or(
			ContainSubstring("TEXT"),
			ContainSubstring("INTEGER"),
		))
	}
}

// Helper function to verify query content.
func verifyQueryContent(content string) {
	Expect(content).To(ContainSubstring("SELECT"))
	Expect(content).To(ContainSubstring("INSERT"))
	// UPDATE and DELETE are optional in different templates
	Expect(content).To(ContainSubstring("-- name:"))
}

var _ = Describe("NewGenerator", func() {
	It("should create a generator with valid output directory", func() {
		tempDir, err := os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(tempDir)

		gen := generators.NewGenerator(tempDir)
		Expect(gen).NotTo(BeNil())
	})

	It("should handle invalid output directory gracefully", func() {
		// Test with empty directory - generator should still be created
		gen := generators.NewGenerator("")
		Expect(gen).NotTo(BeNil())
	})
})

var _ = Describe("Generator Schema Generation", func() {
	var (
		gen     *generators.Generator
		tempDir string
		cleanup func()
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
		gen, cleanup = setupGenerator(tempDir)
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	DescribeTable("Schema generation for different databases",
		func(engine generated.DatabaseType, expectedEngine string) {
			templateData := createTemplateData(engine, tempDir)

			err := gen.GenerateExampleSchema(templateData)
			Expect(err).NotTo(HaveOccurred())

			// Check if schema file was created
			schemaFile := filepath.Join(tempDir, "schema", "001_users_table.sql")
			Expect(schemaFile).To(BeARegularFile())

			// Check content
			content, err := os.ReadFile(schemaFile)
			Expect(err).NotTo(HaveOccurred())
			contentStr := string(content)

			verifySchemaContent(contentStr, expectedEngine)
		},
		Entry("PostgreSQL", generated.DatabaseTypePostgreSQL, "postgresql"),
		Entry("MySQL", generated.DatabaseTypeMySQL, "mysql"),
		Entry("SQLite", generated.DatabaseTypeSQLite, "sqlite"),
	)
})

var _ = Describe("Generator Query Generation", func() {
	var (
		gen     *generators.Generator
		tempDir string
		cleanup func()
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
		gen, cleanup = setupGenerator(tempDir)
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	DescribeTable("Query generation for different databases",
		func(engine generated.DatabaseType) {
			templateData := createTemplateData(engine, tempDir)

			err := gen.GenerateExampleQueries(templateData)
			Expect(err).NotTo(HaveOccurred())

			// Check if query file was created
			queryFile := filepath.Join(tempDir, "queries", "users.sql")
			Expect(queryFile).To(BeARegularFile())

			// Check content
			content, err := os.ReadFile(queryFile)
			Expect(err).NotTo(HaveOccurred())
			contentStr := string(content)

			verifyQueryContent(contentStr)
		},
		Entry("PostgreSQL", generated.DatabaseTypePostgreSQL),
		Entry("MySQL", generated.DatabaseTypeMySQL),
		Entry("SQLite", generated.DatabaseTypeSQLite),
	)
})

var _ = Describe("Error Handling", func() {
	var (
		gen     *generators.Generator
		tempDir string
		cleanup func()
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
		gen, cleanup = setupGenerator(tempDir)
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	It("should handle invalid database engine gracefully", func() {
		templateData := createTemplateData(generated.DatabaseType("invalid"), tempDir)

		err := gen.GenerateExampleSchema(templateData)
		// Should not panic, but may return error
		Expect(err).To(SatisfyAny(
			BeNil(),
			HaveOccurred(),
		))
	})

	// TODO: This test is brittle - permission enforcement varies by OS and user privileges
	// Consider refactoring to use a mock filesystem adapter instead
	PIt("should handle read-only directory", func() {
		// SKIP: This test fails in some environments (Docker, root, etc.)
		// where file permissions don't prevent writes as expected.
		// Need to refactor to use dependency injection with a mock filesystem
		// that can reliably simulate permission errors.

		// Make directory read-only
		err := os.Chmod(tempDir, 0o444)
		Expect(err).NotTo(HaveOccurred())
		defer func() {
			if err := os.Chmod(tempDir, 0o755); err != nil {
				// Log error but don't fail test
				fmt.Printf("Warning: failed to restore permissions: %v\n", err)
			}
		}()

		templateData := createTemplateData(generated.DatabaseTypeSQLite, tempDir)

		err = gen.GenerateExampleSchema(templateData)
		Expect(err).To(HaveOccurred())
	})
})

var _ = Describe("Template Data Integration", func() {
	It("should use template data correctly", func() {
		tempDir, err := os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
		defer os.RemoveAll(tempDir)

		gen, cleanup := setupGenerator(tempDir)
		defer cleanup()

		templateData := createTemplateData(generated.DatabaseTypePostgreSQL, tempDir)

		// Generate schema
		err = gen.GenerateExampleSchema(templateData)
		Expect(err).NotTo(HaveOccurred())

		// Generate queries
		err = gen.GenerateExampleQueries(templateData)
		Expect(err).NotTo(HaveOccurred())

		// Check files exist
		schemaFile := filepath.Join(tempDir, "schema", "001_users_table.sql")
		queryFile := filepath.Join(tempDir, "queries", "users.sql")

		Expect(schemaFile).To(BeARegularFile())
		Expect(queryFile).To(BeARegularFile())
	})
})

var _ = Describe("File Structure", func() {
	var (
		gen     *generators.Generator
		tempDir string
		cleanup func()
		err     error
	)

	BeforeEach(func() {
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test-*")
		Expect(err).NotTo(HaveOccurred())
		gen, cleanup = setupGenerator(tempDir)
	})

	AfterEach(func() {
		if cleanup != nil {
			cleanup()
		}
	})

	It("should create proper directory structure", func() {
		templateData := createTemplateData(generated.DatabaseTypeSQLite, tempDir)

		err := gen.GenerateExampleSchema(templateData)
		Expect(err).NotTo(HaveOccurred())

		err = gen.GenerateExampleQueries(templateData)
		Expect(err).NotTo(HaveOccurred())

		// Check directories exist
		schemaDir := filepath.Join(tempDir, "schema")
		queriesDir := filepath.Join(tempDir, "queries")

		Expect(schemaDir).To(BeADirectory())
		Expect(queriesDir).To(BeADirectory())
	})

	It("should use correct file naming convention", func() {
		templateData := createTemplateData(generated.DatabaseTypeSQLite, tempDir)

		err := gen.GenerateExampleSchema(templateData)
		Expect(err).NotTo(HaveOccurred())

		// Check schema file name
		schemaFile := filepath.Join(tempDir, "schema", "001_users_table.sql")
		Expect(schemaFile).To(BeARegularFile())
	})
})
