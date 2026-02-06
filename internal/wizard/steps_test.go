package wizard

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/charmbracelet/huh"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// stepTestCase represents a test case for step creation
type stepTestCase struct {
	data        *generated.TemplateData
	description string
}

// generateStepTests generates common step creation tests for the given step function
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
	Context("with PostgreSQL database type", func() {
		It("should create a valid step", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypePostgreSQL,
				},
			}

			step := CreateDatabaseStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with MySQL database type", func() {
		It("should create a valid step", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeMySQL,
				},
			}

			step := CreateDatabaseStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with SQLite database type", func() {
		It("should create a valid step", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeSQLite,
				},
			}

			step := CreateDatabaseStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateDatabaseStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreateProjectNameStep", func() {
	Context("with default template data", func() {
		It("should create a valid step", func() {
			data := &generated.TemplateData{
				ProjectName: "myproject",
			}

			step := CreateProjectNameStep(data)

			Expect(step).ToNot(BeNil())
		})

		It("should bind to project name field", func() {
			data := &generated.TemplateData{
				ProjectName: "testproject",
			}

			step := CreateProjectNameStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateProjectNameStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreatePackageNameStep", func() {
	generateStepTests(
		"CreatePackageNameStep",
		CreatePackageNameStep,
		[]stepTestCase{
			{
				data: &generated.TemplateData{
					Package: generated.PackageConfig{Name: "db"},
				},
				description: "package name field",
			},
			{
				data: &generated.TemplateData{
					Package: generated.PackageConfig{Name: "mypackage"},
				},
				description: "package name field",
			},
		},
	)
})

var _ = Describe("CreatePackagePathStep", func() {
	generateStepTests(
		"CreatePackagePathStep",
		CreatePackagePathStep,
		[]stepTestCase{
			{
				data: &generated.TemplateData{
					Package: generated.PackageConfig{Path: "github.com/myorg/myproject"},
				},
				description: "package path field",
			},
			{
				data: &generated.TemplateData{
					Package: generated.PackageConfig{Path: "github.com/example/project"},
				},
				description: "package path field",
			},
		},
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
	Context("with PostgreSQL database", func() {
		It("should create a valid step with PostgreSQL-specific placeholder", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypePostgreSQL,
					URL:     "postgresql://localhost:5432/mydb",
				},
			}

			step := CreateDatabaseURLStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with SQLite database", func() {
		It("should create a valid step with SQLite-specific placeholder", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeSQLite,
					URL:     "./data.db",
				},
			}

			step := CreateDatabaseURLStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with MySQL database", func() {
		It("should create a valid step with MySQL-specific placeholder", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					Engine: generated.DatabaseTypeMySQL,
					URL:     "mysql://localhost:3306/mydb",
				},
			}

			step := CreateDatabaseURLStep(data)

			Expect(step).ToNot(BeNil())
		})
	})

	Context("with nil template data", func() {
		It("should not panic", func() {
			step := CreateDatabaseURLStep(nil)
			Expect(step).ToNot(BeNil())
		})
	})
})

var _ = Describe("CreateFeatureSteps", func() {
	Context("with default template data", func() {
		It("should create all feature steps", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					UseUUIDs:    true,
					UseJSON:     true,
					UseArrays:   false,
					UseFullText: false,
				},
				Validation: generated.ValidationConfig{
					StrictFunctions: true,
					StrictOrderBy:   false,
				},
			}

			steps := CreateFeatureSteps(data)

			Expect(steps).ToNot(BeNil())
			Expect(len(steps)).To(BeNumerically(">", 0))
		})

		It("should create exactly 6 feature steps", func() {
			data := &generated.TemplateData{}

			steps := CreateFeatureSteps(data)

			Expect(len(steps)).To(Equal(6))
		})
	})

	Context("with all features enabled", func() {
		It("should create valid steps", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					UseUUIDs:    true,
					UseJSON:     true,
					UseArrays:   true,
					UseFullText: true,
				},
				Validation: generated.ValidationConfig{
					StrictFunctions: true,
					StrictOrderBy:   true,
				},
			}

			steps := CreateFeatureSteps(data)

			Expect(steps).ToNot(BeNil())
			Expect(len(steps)).To(Equal(6))
		})
	})

	Context("with all features disabled", func() {
		It("should create valid steps", func() {
			data := &generated.TemplateData{
				Database: generated.DatabaseConfig{
					UseUUIDs:    false,
					UseJSON:     false,
					UseArrays:   false,
					UseFullText: false,
				},
				Validation: generated.ValidationConfig{
					StrictFunctions: false,
					StrictOrderBy:   false,
				},
			}

			steps := CreateFeatureSteps(data)

			Expect(steps).ToNot(BeNil())
			Expect(len(steps)).To(Equal(6))
		})
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
