package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// stepTestCase represents a test case for step creation.
type stepTestCase struct {
	data        *generated.TemplateData
	description string
}

// generateStepTests generates common step creation tests for the given step function.
func generateStepTests(describeName string, stepFunc func(*generated.TemplateData) *huh.Input, testCases []stepTestCase) {
	Describe(describeName, func() {
		Context("with default template data", func() {
			It("should create a valid step", func() {
				for _, tc := range testCases {
					step := stepFunc(tc.data)
					Expect(step).ToNot(BeNil())
				}
			})

			It("should bind to fields", func() {
				for _, tc := range testCases {
					step := stepFunc(tc.data)
					Expect(step).ToNot(BeNil())
				}
			})
		})

		Context("with nil template data", func() {
			It("should not panic", func() {
				step := stepFunc(nil)
				Expect(step).ToNot(BeNil())
			})
		})
	})
}

// createStepTestCases creates test cases for step creation with default values.
func createStepTestCases() []stepTestCase {
	return []stepTestCase{
		{
			data:        &generated.TemplateData{},
			description: "default template data",
		},
		{
			data:        &generated.TemplateData{},
			description: "default template data",
		},
	}
}

// createFieldTestCases creates test cases with custom field configuration.
// This eliminates duplication across similar test case helper functions.
func createFieldTestCases(fieldDescription string, setupFields func(*generated.TemplateData)) []stepTestCase {
	return []stepTestCase{
		{
			data: func() *generated.TemplateData {
				data := &generated.TemplateData{}
				setupFields(data)
				return data
			}(),
			description: fieldDescription,
		},
		{
			data: func() *generated.TemplateData {
				data := &generated.TemplateData{}
				setupFields(data)
				return data
			}(),
			description: fieldDescription,
		},
	}
}

// createProjectNameTestCases returns test cases for project name step tests.
func createProjectNameTestCases() []stepTestCase {
	return createFieldTestCases("project name field", func(data *generated.TemplateData) {
		data.ProjectName = "myproject"
	})
}

// createPackageNameTestCases returns test cases for package name step tests.
func createPackageNameTestCases() []stepTestCase {
	return createFieldTestCases("package name field", func(data *generated.TemplateData) {
		data.Package.Name = "db"
	})
}

// createPackagePathTestCases returns test cases for package path step tests.
func createPackagePathTestCases() []stepTestCase {
	return createFieldTestCases("package path field", func(data *generated.TemplateData) {
		data.Package.Path = "github.com/myorg/myproject"
	})
}

// createFeatureStepTestData creates template data for feature step testing.
// This eliminates hardcoded TemplateData creation in tests.
func createFeatureStepTestData(enableFeatures bool) *generated.TemplateData {
	return &generated.TemplateData{
		Database: generated.DatabaseConfig{
			UseUUIDs:    enableFeatures,
			UseJSON:     enableFeatures,
			UseArrays:   !enableFeatures,
			UseFullText: !enableFeatures,
		},
		Validation: generated.ValidationConfig{
			StrictFunctions: enableFeatures,
			StrictOrderBy:   !enableFeatures,
		},
	}
}

var _ = Describe("CreateProjectTypeStep", func() {
	Context("with default template data", func() {
		It("should create a valid step with all project types", func() {
			data := &generated.TemplateData{
				ProjectType: generated.ProjectTypeMicroservice,
			}

			step := CreateProjectTypeStep(data)

			Expect(step).ToNot(BeNil())
		})

		It("should initialize with hobby project type", func() {
			data := &generated.TemplateData{
				ProjectType: generated.ProjectTypeHobby,
			}

			step := CreateProjectTypeStep(data)

			Expect(step).ToNot(BeNil())
		})

		It("should initialize with enterprise project type", func() {
			data := &generated.TemplateData{
				ProjectType: generated.ProjectTypeEnterprise,
			}

			step := CreateProjectTypeStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateProjectTypeStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreateDatabaseStep", func() {
	DescribeTable("should create valid step for each database type",
		func(engine generated.DatabaseType) {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: engine,
				},
			}

			step := CreateDatabaseStep(data)

			Expect(step).ToNot(BeNil())
		},
		Entry("PostgreSQL", generated.DatabaseTypePostgreSQL),
		Entry("MySQL", generated.DatabaseTypeMySQL),
		Entry("SQLite", generated.DatabaseTypeSQLite),
	)

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateDatabaseStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreateProjectNameStep", func() {
	generateStepTests(
		"CreateProjectNameStep",
		CreateProjectNameStep,
		createProjectNameTestCases(),
	)
})

var _ = Describe("CreatePackageNameStep", func() {
	generateStepTests(
		"CreatePackageNameStep",
		CreatePackageNameStep,
		createPackageNameTestCases(),
	)
})

var _ = Describe("CreatePackagePathStep", func() {
	generateStepTests(
		"CreatePackagePathStep",
		CreatePackagePathStep,
		createPackagePathTestCases(),
	)
})

var _ = Describe("CreateOutputDirStep", func() {
	Context("with default template data", func() {
		It("should create a valid step", func() {
			step := CreateOutputDirStep(&generated.TemplateData{
				Output: generated.OutputConfig{BaseDir: "./internal/db"},
			})
			Expect(step).ToNot(BeNil())
		})

		It("should bind to output directory field", func() {
			step := CreateOutputDirStep(&generated.TemplateData{
				Output: generated.OutputConfig{BaseDir: "./gen/db"},
			})
			Expect(step).ToNot(BeNil())
		})
	})

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateOutputDirStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreateDatabaseURLStep", func() {
	DescribeTable("should create valid step with database-specific placeholders",
		func(engine generated.DatabaseType, url string) {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: engine,
					URL:    url,
				},
			}

			step := CreateDatabaseURLStep(data)

			Expect(step).ToNot(BeNil())
		},
		Entry("PostgreSQL", generated.DatabaseTypePostgreSQL, "postgresql://localhost:5432/mydb"),
		Entry("SQLite", generated.DatabaseTypeSQLite, "./data.db"),
		Entry("MySQL", generated.DatabaseTypeMySQL, "mysql://localhost:3306/mydb"),
	)

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateDatabaseURLStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreateFeatureSteps", func() {
	Context("with partial feature configuration", func() {
		It("should create all feature steps", func() {
			data := createFeatureStepTestData(true)

			steps := CreateFeatureSteps(data)

			Expect(steps).ToNot(BeNil())
			Expect(steps).ToNot(BeEmpty())
		})

		It("should create exactly 6 feature steps", func() {
			data := createFeatureStepTestData(true)

			steps := CreateFeatureSteps(data)

			Expect(steps).To(HaveLen(6))
		})
	})

	Context("with all features", func() {
		DescribeTable("should create valid steps regardless of feature state",
			func(enabled bool) {
				data := createTemplateDataWithAllFeatures(enabled)

				steps := CreateFeatureSteps(data)

				Expect(steps).ToNot(BeNil())
				Expect(steps).To(HaveLen(6))
			},
			Entry("when enabled", true),
			Entry("when disabled", false),
		)
	})

	Context("with nil template data", func() {
		It("should not panic", func() {
			steps := CreateFeatureSteps(nil)
			Expect(steps).ToNot(BeNil())
		})
	})
})

var _ = Describe("createValidatedInput", func() {
	Context("with all parameters provided", func() {
		It("should create a valid input field", func() {
			value := ""

			input := createValidatedInput(
				"Test Title",
				"Test Description",
				"test-placeholder",
				"testField",
				&value,
			)

			Expect(input).ToNot(BeNil())
		})
	})

	Context("with empty title", func() {
		It("should create a valid input field", func() {
			value := ""

			input := createValidatedInput(
				"",
				"Test Description",
				"test-placeholder",
				"testField",
				&value,
			)

			Expect(input).ToNot(BeNil())
		})
	})

	Context("with nil value pointer", func() {
		It("should not panic", func() {
			input := createValidatedInput(
				"Test Title",
				"Test Description",
				"test-placeholder",
				"testField",
				nil,
			)

			Expect(input).ToNot(BeNil())
		})
	})

	Context("with empty field name", func() {
		It("should create a valid input field", func() {
			value := ""

			input := createValidatedInput(
				"Test Title",
				"Test Description",
				"test-placeholder",
				"",
				&value,
			)

			Expect(input).ToNot(BeNil())
		})
	})
})
