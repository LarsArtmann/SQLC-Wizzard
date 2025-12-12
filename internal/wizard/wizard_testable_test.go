package wizard_test

import (
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// MockUI is a mock implementation of UIInterface for testing
type MockUI struct {
	WelcomeCalls     int
	StepHeaders      []string
	StepCompletions  []StepCompletion
	Sections         []string
	InfoMessages     []string
	ShouldFailRun    bool
	FailErrorMessage string
}

type StepCompletion struct {
	Title   string
	Message string
}

func NewMockUI() *MockUI {
	return &MockUI{
		StepHeaders:     make([]string, 0),
		StepCompletions: make([]StepCompletion, 0),
		Sections:        make([]string, 0),
		InfoMessages:    make([]string, 0),
	}
}

func (m *MockUI) ShowWelcome() {
	m.WelcomeCalls++
	if m.ShouldFailRun {
		panic(fmt.Errorf("UI operation failed: %s", m.FailErrorMessage))
	}
}

func (m *MockUI) ShowStepHeader(title string) {
	m.StepHeaders = append(m.StepHeaders, title)
}

func (m *MockUI) ShowStepComplete(title, message string) {
	m.StepCompletions = append(m.StepCompletions, StepCompletion{
		Title:   title,
		Message: message,
	})
}

func (m *MockUI) ShowSection(title string) {
	m.Sections = append(m.Sections, title)
}

func (m *MockUI) ShowInfo(message string) {
	m.InfoMessages = append(m.InfoMessages, message)
}

// MockStep is a mock implementation of StepInterface for testing
type MockStep struct {
	ExecuteCalls       int
	ShouldFail         bool
	FailError          error
	LastCallData       *generated.TemplateData
	ExpectedProject    generated.ProjectType
	ValidateCalls      int
	ShouldValidateFail bool
	ValidateError      error
	ExecuteFunc        func(data *generated.TemplateData) error
}

func NewMockStep() *MockStep {
	return &MockStep{}
}

func (m *MockStep) Execute(data *generated.TemplateData) error {
	m.ExecuteCalls++
	m.LastCallData = data

	// Use custom execute function if provided
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(data)
	}

	if m.ShouldFail {
		return m.FailError
	}

	// Set some default data based on step type
	if data.ProjectType == "" {
		data.ProjectType = generated.ProjectTypeMicroservice
	}

	if data.Database.Engine == "" {
		data.Database.Engine = generated.DatabaseTypePostgreSQL
	}

	if data.Package.Name == "" {
		data.Package.Name = "testdb"
		data.Package.Path = "github.com/example/testdb"
	}

	return nil
}

// ValidateConfiguration implements ValidatableStepInterface
func (m *MockStep) ValidateConfiguration(data *generated.TemplateData) error {
	m.ValidateCalls++
	if m.ShouldValidateFail {
		return m.ValidateError
	}
	return nil
}

// MockTemplate is a mock implementation of TemplateInterface for testing
type MockTemplate struct {
	GenerateCalls int
	ShouldFail    bool
	FailError     error
	LastCallData  generated.TemplateData
}

func NewMockTemplate() *MockTemplate {
	return &MockTemplate{}
}

func (m *MockTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	m.GenerateCalls++
	m.LastCallData = data

	if m.ShouldFail {
		return nil, m.FailError
	}

	// Return a basic valid config
	return &config.SqlcConfig{
		Version: "2",
		SQL: []config.SQLConfig{
			{
				Engine:  "postgresql",
				Schema:  config.NewPathOrPaths([]string{"schema.sql"}),
				Queries: config.NewPathOrPaths([]string{"query.sql"}),
				Gen: config.GenConfig{
					Go: &config.GoGenConfig{
						Out:     "internal/db",
						Package: "db",
					},
				},
			},
		},
	}, nil
}

func (m *MockTemplate) DefaultData() generated.TemplateData {
	return generated.TemplateData{
		ProjectType: generated.ProjectTypeMicroservice,
	}
}

func (m *MockTemplate) RequiredFeatures() []string {
	return []string{"basic"}
}

func (m *MockTemplate) Name() string {
	return "Mock Template"
}

func (m *MockTemplate) Description() string {
	return "A mock template for testing"
}

var _ = Describe("Wizard with Dependency Injection", func() {
	var (
		mockUI       *MockUI
		mockSteps    map[string]*MockStep
		mockTemplate *MockTemplate
		wizardDeps   wizard.WizardDependencies
		wiz          *wizard.Wizard
	)

	BeforeEach(func() {
		mockUI = NewMockUI()
		mockSteps = map[string]*MockStep{
			"projectType": NewMockStep(),
			"database":    NewMockStep(),
			"details":     NewMockStep(),
			"features":    NewMockStep(),
			"output":      NewMockStep(),
		}
		mockTemplate = NewMockTemplate()

		wizardDeps = wizard.WizardDependencies{
			UI:          mockUI,
			ProjectType: mockSteps["projectType"],
			Database:    mockSteps["database"],
			Details:     mockSteps["details"],
			Features:    mockSteps["features"],
			Output:      mockSteps["output"],
			TemplateFunc: func(projectType templates.ProjectType) (wizard.TemplateInterface, error) {
				return mockTemplate, nil
			},
		}

		wiz = wizard.NewTestableWizard(wizardDeps)
	})

	AfterEach(func() {
		// Verify no unexpected calls
		// This will help catch issues in test setup
	})

	Describe("Wizard Initialization", func() {
		It("should create wizard with injected dependencies", func() {
			Expect(wiz).NotTo(BeNil())
			result := wiz.GetResult()
			Expect(result).NotTo(BeNil())
			Expect(result.GenerateQueries).To(BeTrue())
			Expect(result.GenerateSchema).To(BeTrue())
		})
	})

	Describe("Wizard Run Method", func() {
		Context("when all steps execute successfully", func() {
			It("should execute all wizard steps in correct order", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())

				// Verify UI calls
				Expect(mockUI.WelcomeCalls).To(Equal(1))
				Expect(len(mockUI.StepHeaders)).To(Equal(5))
				Expect(mockUI.StepHeaders).To(Equal([]string{
					"Project Type",
					"Database",
					"Project Details",
					"Features",
					"Output Configuration",
				}))

				// Verify step execution
				for stepName, step := range mockSteps {
					Expect(step.ExecuteCalls).To(Equal(1), "Step %s should be called once", stepName)
					Expect(step.LastCallData).NotTo(BeNil(), "Step %s should receive data", stepName)
				}

				// Verify template generation
				Expect(mockTemplate.GenerateCalls).To(Equal(1))
				Expect(result.Config).NotTo(BeNil())
				Expect(result.Config.Version).To(Equal("2"))
			})

			It("should show welcome message and summary", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())

				// Verify welcome and summary
				Expect(mockUI.WelcomeCalls).To(Equal(1))
				Expect(len(mockUI.Sections)).To(BeNumerically(">", 0))
				Expect(len(mockUI.InfoMessages)).To(BeNumerically(">", 0))
			})

			It("should pass consistent template data through all steps", func() {
				_, err := wiz.Run()

				Expect(err).To(BeNil())

				// Verify all steps received the same data object
				originalData := mockSteps["projectType"].LastCallData
				for stepName, step := range mockSteps {
					Expect(step.LastCallData).To(Equal(originalData),
						"Step %s should receive same data object", stepName)
				}
			})

			It("should populate wizard result with generated config and template data", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Config).NotTo(BeNil())
				Expect(result.Config.Version).To(Equal("2"))
				Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
				Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
			})
		})

		Context("when a step fails", func() {
			BeforeEach(func() {
				mockSteps["database"].ShouldFail = true
				mockSteps["database"].FailError = fmt.Errorf("database step failed")
			})

			It("should return error and stop execution", func() {
				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("step 'Database' failed"))
				Expect(result).To(BeNil())

				// Verify only steps before failure executed
				Expect(mockSteps["projectType"].ExecuteCalls).To(Equal(1))
				Expect(mockSteps["database"].ExecuteCalls).To(Equal(1))
				Expect(mockSteps["details"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["features"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["output"].ExecuteCalls).To(Equal(0))
			})
		})

		Context("when template generation fails", func() {
			BeforeEach(func() {
				mockTemplate.ShouldFail = true
				mockTemplate.FailError = fmt.Errorf("template generation failed")
			})

			It("should return error and not show summary", func() {
				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("config generation failed"))
				Expect(result).To(BeNil())

				// All steps should have executed
				for stepName, step := range mockSteps {
					Expect(step.ExecuteCalls).To(Equal(1), "Step %s should have been called", stepName)
				}

				// Template generation attempted
				Expect(mockTemplate.GenerateCalls).To(Equal(1))
			})
		})

		Context("when UI operations fail", func() {
			BeforeEach(func() {
				mockUI.ShouldFailRun = true
				mockUI.FailErrorMessage = "UI operation failed"
			})

			It("should panic and stop execution", func() {
				Expect(func() {
					wiz.Run()
				}).To(Panic())
			})
		})
	})

	Describe("Template Data Flow", func() {
		It("should start with default template data", func() {
			_, err := wiz.Run()

			Expect(err).To(BeNil())

			// Verify default data
			data := mockSteps["projectType"].LastCallData
			Expect(data.Package.Name).To(Equal("myproject"))
			Expect(data.Package.Path).To(Equal("github.com/myorg/myproject"))
			Expect(data.Database.UseUUIDs).To(BeTrue())
			Expect(data.Database.UseJSON).To(BeTrue())
			Expect(data.Database.UseArrays).To(BeFalse())
			Expect(data.Database.UseFullText).To(BeFalse())
			Expect(data.Output.BaseDir).To(Equal("./internal/db"))
			Expect(data.Output.QueriesDir).To(Equal("./sql/queries"))
			Expect(data.Output.SchemaDir).To(Equal("./sql/schema"))
		})

		It("should allow steps to modify template data", func() {
			// Create a step that modifies data
			customStep := NewMockStep()
			customStep.ExpectedProject = generated.ProjectTypeEnterprise

			wizardDeps.ProjectType = customStep
			wiz = wizard.NewTestableWizard(wizardDeps)

			result, err := wiz.Run()

			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())

			// Verify data was modified
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
		})
	})

	Describe("Configuration Generation", func() {
		It("should use correct template for project type", func() {
			// Create a custom step that changes project type to enterprise
			customStep := &MockStep{}
			customStep.ExecuteFunc = func(data *generated.TemplateData) error {
				customStep.ExecuteCalls++
				customStep.LastCallData = data
				data.ProjectType = generated.ProjectTypeEnterprise
				return nil
			}

			wizardDeps.ProjectType = customStep
			wiz = wizard.NewTestableWizard(wizardDeps)

			_, err := wiz.Run()

			Expect(err).To(BeNil())
			Expect(mockTemplate.GenerateCalls).To(Equal(1))
			Expect(mockTemplate.LastCallData.ProjectType).To(Equal(generated.ProjectTypeEnterprise))
		})

		It("should pass complete template data to template generation", func() {
			_, err := wiz.Run()

			Expect(err).To(BeNil())
			Expect(mockTemplate.LastCallData.ProjectType).ToNot(BeEmpty())
			Expect(mockTemplate.LastCallData.Package.Name).ToNot(BeEmpty())
			Expect(mockTemplate.LastCallData.Database.Engine).ToNot(BeEmpty())
		})
	})

	Describe("Wizard Execution with Mocked Steps", func() {
		Context("when all steps succeed", func() {
			It("should execute wizard and generate complete result", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())
				Expect(result.Config).NotTo(BeNil())
				Expect(result.Config.Version).To(Equal("2"))
				Expect(len(result.Config.SQL)).To(Equal(1))
				Expect(result.Config.SQL[0].Engine).To(Equal("postgresql"))
				Expect(result.Config.SQL[0].Gen.Go.Out).To(Equal("internal/db"))
				Expect(result.Config.SQL[0].Gen.Go.Package).To(Equal("db"))

				// Verify template data was properly populated
				Expect(result.TemplateData.ProjectType).ToNot(BeEmpty())
				Expect(result.TemplateData.Database.Engine).ToNot(BeEmpty())
				Expect(result.TemplateData.Package.Name).ToNot(BeEmpty())
				Expect(result.TemplateData.Package.Path).ToNot(BeEmpty())
			})

			It("should execute all wizard steps with correct data flow", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())

				// Verify that all steps were called exactly once
				for stepName, step := range mockSteps {
					Expect(step.ExecuteCalls).To(Equal(1),
						"Step %s should be called exactly once", stepName)
					Expect(step.LastCallData).NotTo(BeNil(),
						"Step %s should receive template data", stepName)
				}

				// Verify template generation was called
				Expect(mockTemplate.GenerateCalls).To(Equal(1))
				Expect(mockTemplate.LastCallData.ProjectType).ToNot(BeEmpty())
			})

			It("should show proper UI sequence", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())

				// Verify UI call sequence
				Expect(mockUI.WelcomeCalls).To(Equal(1))
				Expect(len(mockUI.StepHeaders)).To(Equal(5))
				Expect(len(mockUI.StepCompletions)).To(Equal(5))
				Expect(len(mockUI.Sections)).To(BeNumerically(">", 0))
				Expect(len(mockUI.InfoMessages)).To(BeNumerically(">", 0))

				// Verify correct step order
				Expect(mockUI.StepHeaders).To(Equal([]string{
					"Project Type",
					"Database",
					"Project Details",
					"Features",
					"Output Configuration",
				}))
			})

			It("should populate wizard result correctly", func() {
				result, err := wiz.Run()

				Expect(err).To(BeNil())
				Expect(result).NotTo(BeNil())

				// Verify wizard result properties
				Expect(result.GenerateQueries).To(BeTrue())
				Expect(result.GenerateSchema).To(BeTrue())
				Expect(result.Config).NotTo(BeNil())
				Expect(result.TemplateData.ProjectType).ToNot(BeEmpty())
			})

			It("should handle complete wizard execution cycle", func() {
				// Test multiple runs to ensure state is properly managed
				for i := 0; i < 3; i++ {
					result, err := wiz.Run()

					Expect(err).To(BeNil(),
						fmt.Sprintf("Wizard run %d should succeed", i+1))
					Expect(result).NotTo(BeNil(),
						fmt.Sprintf("Wizard run %d should return result", i+1))
					Expect(result.Config).NotTo(BeNil(),
						fmt.Sprintf("Wizard run %d should generate config", i+1))

					// Reset mock calls for next iteration
					mockUI.WelcomeCalls = 0
					mockUI.StepHeaders = make([]string, 0)
					mockUI.StepCompletions = make([]StepCompletion, 0)
					for _, step := range mockSteps {
						step.ExecuteCalls = 0
					}
					mockTemplate.GenerateCalls = 0
				}
			})
		})

		Context("when steps fail", func() {
			It("should handle project type step failure", func() {
				mockSteps["projectType"].ShouldFail = true
				mockSteps["projectType"].FailError = fmt.Errorf("project type selection failed")

				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("step 'Project Type' failed"))
				Expect(result).To(BeNil())

				// Verify only project type step was called
				Expect(mockSteps["projectType"].ExecuteCalls).To(Equal(1))
				Expect(mockSteps["database"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["details"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["features"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["output"].ExecuteCalls).To(Equal(0))
			})

			It("should handle database step failure", func() {
				mockSteps["database"].ShouldFail = true
				mockSteps["database"].FailError = fmt.Errorf("database configuration failed")

				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("step 'Database' failed"))
				Expect(result).To(BeNil())

				// Verify steps up to database were called
				Expect(mockSteps["projectType"].ExecuteCalls).To(Equal(1))
				Expect(mockSteps["database"].ExecuteCalls).To(Equal(1))
				Expect(mockSteps["details"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["features"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["output"].ExecuteCalls).To(Equal(0))
			})

			It("should handle output validation failure", func() {
				mockSteps["output"].ShouldValidateFail = true
				mockSteps["output"].ValidateError = fmt.Errorf("output validation failed")

				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("invalid output configuration"))
				Expect(result).To(BeNil())

				// All steps should execute but template generation should fail
				for stepName, step := range mockSteps {
					Expect(step.ExecuteCalls).To(Equal(1),
						"Step %s should have been called", stepName)
				}
				Expect(mockTemplate.GenerateCalls).To(Equal(0))
			})

			It("should handle template generation failure", func() {
				mockTemplate.ShouldFail = true
				mockTemplate.FailError = fmt.Errorf("template generation failed")

				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("config generation failed"))
				Expect(result).To(BeNil())

				// All steps should execute but template generation should fail
				for stepName, step := range mockSteps {
					Expect(step.ExecuteCalls).To(Equal(1),
						"Step %s should have been called", stepName)
				}
				Expect(mockTemplate.GenerateCalls).To(Equal(1))
			})

			It("should handle multiple step failures gracefully", func() {
				// Set multiple steps to fail - first one should stop execution
				mockSteps["projectType"].ShouldFail = true
				mockSteps["projectType"].FailError = fmt.Errorf("first failure")
				mockSteps["database"].ShouldFail = true // This should never be reached
				mockSteps["database"].FailError = fmt.Errorf("second failure")

				// First failure should stop execution
				result, err := wiz.Run()

				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(ContainSubstring("first failure"))
				Expect(result).To(BeNil())

				// Only first step should have been called
				Expect(mockSteps["projectType"].ExecuteCalls).To(Equal(1))
				Expect(mockSteps["database"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["details"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["features"].ExecuteCalls).To(Equal(0))
				Expect(mockSteps["output"].ExecuteCalls).To(Equal(0))
			})
		})

		Context("when UI operations fail", func() {
			It("should handle welcome failure", func() {
				mockUI.ShouldFailRun = true
				mockUI.FailErrorMessage = "welcome display failed"

				Expect(func() {
					wiz.Run()
				}).To(Panic())
			})
		})
	})

	Describe("Advanced Wizard Configuration", func() {
		It("should handle complex template data scenarios", func() {
			// Create custom steps with specific data modifications
			projectTypeStep := &MockStep{}
			projectTypeStep.ExecuteFunc = func(data *generated.TemplateData) error {
				data.ProjectType = generated.ProjectTypeEnterprise
				data.ProjectName = "enterprise-app"
				return nil
			}

			databaseStep := &MockStep{}
			databaseStep.ExecuteFunc = func(data *generated.TemplateData) error {
				data.Database.Engine = generated.DatabaseTypePostgreSQL
				data.Database.UseUUIDs = true
				data.Database.UseJSON = true
				data.Database.UseFullText = true
				return nil
			}

			detailsStep := &MockStep{}
			detailsStep.ExecuteFunc = func(data *generated.TemplateData) error {
				data.Package.Name = "enterprise_db"
				data.Package.Path = "github.com/company/enterprise"
				data.Package.BuildTags = "enterprise,postgres"
				return nil
			}

			featuresStep := &MockStep{}
			featuresStep.ExecuteFunc = func(data *generated.TemplateData) error {
				data.Validation.EmitOptions.EmitInterface = true
				data.Validation.EmitOptions.EmitJSONTags = true
				data.Validation.SafetyRules.NoSelectStar = true
				data.Validation.SafetyRules.RequireWhere = true
				return nil
			}

			outputStep := &MockStep{}
			outputStep.ExecuteFunc = func(data *generated.TemplateData) error {
				data.Output.BaseDir = "./internal/enterprise/db"
				data.Output.QueriesDir = "./sql/enterprise/queries"
				data.Output.SchemaDir = "./sql/enterprise/schema"
				return nil
			}

			wizardDeps.ProjectType = projectTypeStep
			wizardDeps.Database = databaseStep
			wizardDeps.Details = detailsStep
			wizardDeps.Features = featuresStep
			wizardDeps.Output = outputStep
			wiz = wizard.NewTestableWizard(wizardDeps)

			result, err := wiz.Run()

			Expect(err).To(BeNil())
			Expect(result).NotTo(BeNil())

			// Verify complex template data was correctly populated
			Expect(result.TemplateData.ProjectType).To(Equal(generated.ProjectTypeEnterprise))
			Expect(result.TemplateData.ProjectName).To(Equal("enterprise-app"))
			Expect(result.TemplateData.Database.Engine).To(Equal(generated.DatabaseTypePostgreSQL))
			Expect(result.TemplateData.Database.UseUUIDs).To(BeTrue())
			Expect(result.TemplateData.Database.UseJSON).To(BeTrue())
			Expect(result.TemplateData.Database.UseFullText).To(BeTrue())
			Expect(result.TemplateData.Package.Name).To(Equal("enterprise_db"))
			Expect(result.TemplateData.Package.Path).To(Equal("github.com/company/enterprise"))
			Expect(result.TemplateData.Package.BuildTags).To(Equal("enterprise,postgres"))
			Expect(result.TemplateData.Validation.EmitOptions.EmitInterface).To(BeTrue())
			Expect(result.TemplateData.Validation.EmitOptions.EmitJSONTags).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.NoSelectStar).To(BeTrue())
			Expect(result.TemplateData.Validation.SafetyRules.RequireWhere).To(BeTrue())
			Expect(result.TemplateData.Output.BaseDir).To(Equal("./internal/enterprise/db"))
			Expect(result.TemplateData.Output.QueriesDir).To(Equal("./sql/enterprise/queries"))
			Expect(result.TemplateData.Output.SchemaDir).To(Equal("./sql/enterprise/schema"))
		})

		It("should handle all project type scenarios", func() {
			projectTypes := []generated.ProjectType{
				generated.ProjectTypeHobby,
				generated.ProjectTypeMicroservice,
				generated.ProjectTypeEnterprise,
				generated.ProjectTypeAPIFirst,
				generated.ProjectTypeAnalytics,
				generated.ProjectTypeTesting,
				generated.ProjectTypeMultiTenant,
				generated.ProjectTypeLibrary,
			}

			for _, projectType := range projectTypes {
				// Create step that sets specific project type
				projectTypeStep := &MockStep{}
				projectTypeStep.ExecuteFunc = func(data *generated.TemplateData) error {
					data.ProjectType = projectType
					data.ProjectName = fmt.Sprintf("%s-project", projectType)
					return nil
				}

				wizardDeps.ProjectType = projectTypeStep
				wiz = wizard.NewTestableWizard(wizardDeps)

				result, err := wiz.Run()

				Expect(err).To(BeNil(),
					fmt.Sprintf("Wizard should succeed for project type %s", projectType))
				Expect(result).NotTo(BeNil(),
					fmt.Sprintf("Wizard should return result for project type %s", projectType))
				Expect(result.TemplateData.ProjectType).To(Equal(projectType),
					fmt.Sprintf("Project type should be %s", projectType))
				Expect(result.TemplateData.ProjectName).To(ContainSubstring(string(projectType)),
					fmt.Sprintf("Project name should contain %s", projectType))
			}
		})

		It("should handle all database type scenarios", func() {
			databaseTypes := []generated.DatabaseType{
				generated.DatabaseTypePostgreSQL,
				generated.DatabaseTypeMySQL,
				generated.DatabaseTypeSQLite,
			}

			for _, dbType := range databaseTypes {
				// Create step that sets specific database type
				databaseStep := &MockStep{}
				databaseStep.ExecuteFunc = func(data *generated.TemplateData) error {
					data.Database.Engine = dbType
					return nil
				}

				wizardDeps.Database = databaseStep
				wiz = wizard.NewTestableWizard(wizardDeps)

				result, err := wiz.Run()

				Expect(err).To(BeNil(),
					fmt.Sprintf("Wizard should succeed for database type %s", dbType))
				Expect(result).NotTo(BeNil(),
					fmt.Sprintf("Wizard should return result for database type %s", dbType))
				Expect(result.TemplateData.Database.Engine).To(Equal(dbType),
					fmt.Sprintf("Database type should be %s", dbType))
			}
		})
	})
})
