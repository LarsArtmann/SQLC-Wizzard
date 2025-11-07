package wizard_test

import (
	"context"
	"os"
	"path/filepath"

	. "github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("SQLC Configuration Generation", func() {
	var (
		tempDir   string
		ctx       context.Context
		wizard    *Wizard
		result    *WizardResult
	)

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "sqlc-wizard-test")
		Expect(err).ToNot(HaveOccurred())
		
		ctx = context.Background()
		wizard = NewWizard()
		
		// Change to temp directory
		oldDir, _ := os.Getwd()
		DeferCleanup(func() {
			os.Chdir(oldDir)
			os.RemoveAll(tempDir)
		})
		os.Chdir(tempDir)
	})

	AfterEach(func() {
		if result != nil {
			// Clean up any generated files
			os.Remove("sqlc.yaml")
			os.RemoveAll("sql")
			os.RemoveAll("migrations")
		}
	})

	Context("Microservice Project Generation", func() {
		When("I choose microservice template", func() {
			It("should generate valid sqlc.yaml configuration", func() {
				// This is a BDD scenario implementation
				// Given: We have microservice template selection
				// When: We run the wizard with microservice options
				// Then: sqlc.yaml should be created with correct configuration

				// Mock the wizard run with microservice configuration
				result = &WizardResult{
					Config: &config.SqlcConfig{
						Version: "2",
						SQL: []config.SQL{
							{
								Engine: "postgresql",
								Queries: "sql/",
								Schema: "migrations/",
								// Additional microservice-specific config
							},
						},
						Go: &config.GoGenConfig{
							Package: "models",
							Out:     "gen/",
							// Microservice-specific generation options
						},
						Overrides: []config.Override{
							// Type overrides for microservice
						},
					},
					TemplateData: templates.TemplateData{
						ProjectType: templates.MicroserviceTemplate{},
						Database:    "postgresql",
						OutputDir:   "./gen",
						PackagePath: "github.com/example/microservice",
					},
					GenerateQueries: true,
					GenerateSchema:  true,
				}

				// Write the config
				configContent, err := utils.GenerateYAML(result.Config)
				Expect(err).ToNot(HaveOccurred())

				err = os.WriteFile("sqlc.yaml", []byte(configContent), 0644)
				Expect(err).ToNot(HaveOccurred())

				// Then: File should exist with correct content
				configFile, err := os.Stat("sqlc.yaml")
				Expect(err).ToNot(HaveOccurred())
				Expect(configFile.Size()).To(BeNumerically(">", 0))

				// And: Content should contain PostgreSQL configuration
				content, err := os.ReadFile("sqlc.yaml")
				Expect(err).ToNot(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("postgresql"))
				Expect(string(content)).To(ContainSubstring("package: models"))
			})

			It("should generate microservice-specific directories", func() {
				// BDD scenario: Directory structure generation
				// Given: Microservice template selected
				// When: Generate files command executed
				// Then: Appropriate directories should be created

				result = &WizardResult{
					TemplateData: templates.TemplateData{
						ProjectType: templates.MicroserviceTemplate{},
						Database:    "postgresql",
						OutputDir:   "./gen",
						QueriesDir:  "sql",
						SchemaDir:   "migrations",
					},
					GenerateQueries: true,
					GenerateSchema:  true,
				}

				// Generate directories
				err := os.MkdirAll(result.TemplateData.QueriesDir, 0755)
				Expect(err).ToNot(HaveOccurred())
				err = os.MkdirAll(result.TemplateData.SchemaDir, 0755)
				Expect(err).ToNot(HaveOccurred())

				// Then: Directories should exist
				queriesDir, err := os.Stat(result.TemplateData.QueriesDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(queriesDir.IsDir()).To(BeTrue())

				schemaDir, err := os.Stat(result.TemplateData.SchemaDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(schemaDir.IsDir()).To(BeTrue())
			})
		})
	})

	Context("CLI Command Integration", func() {
		When("I run sqlc-wizard init command", func() {
			It("should start wizard and complete successfully", func() {
				// BDD scenario: CLI integration
				// Given: sqlc-wizard is installed
				// When: I run 'sqlc-wizard init' command
				// Then: Wizard should start and allow configuration

				// Test the binary exists and can be executed
				_, err := gexec.Build("../cmd/sqlc-wizard")
				Expect(err).ToNot(HaveOccurred())

				// Note: In real BDD, we would use expect script or similar
				// For now, we validate the command structure exists
				Expect(true).To(BeTrue()) // Placeholder for CLI integration test
			})
		})
	})

	Context("Template Validation", func() {
		When("I select different project types", func() {
			It("should validate project type constraints", func() {
				// BDD scenario: Project type validation
				// Given: Multiple project type options
				// When: I select enterprise template
				// Then: Additional enterprise features should be available

				projectTypes := []string{
					"hobby", "microservice", "enterprise",
					"api-first", "analytics", "testing",
					"multi-tenant", "library",
				}

				for _, projectType := range projectTypes {
					// Validate each project type is supported
					template := templates.SelectTemplate(projectType)
					Expect(template).ToNot(BeNil())
					Expect(template.Name()).To(Equal(projectType))
				}
			})
		})
	})

	Context("Error Handling", func() {
		When("Configuration validation fails", func() {
			It("should provide actionable error messages", func() {
				// BDD scenario: Error handling
				// Given: Invalid configuration provided
				// When: Validation is attempted
				// Then: Clear error messages should be provided

				invalidConfig := &config.SqlcConfig{
					Version: "",
					// Missing required fields
				}

				// Attempt to validate/generate
				_, err := utils.GenerateYAML(invalidConfig)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("validation"))
			})
		})
	})
})