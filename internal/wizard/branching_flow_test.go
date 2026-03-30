package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/wizard"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Branching Flow Context", func() {
	var ctx *wizard.FlowContext

	BeforeEach(func() {
		ctx = wizard.NewFlowContext()
	})

	Describe("Step Filtering", func() {
		Context("for hobby project type", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeHobby
				ctx.DatabaseType = generated.DatabaseTypeSQLite
			})

			It("should skip features step when SkipOptionalSteps is true", func() {
				ctx.SkipOptionalSteps = true
				steps := ctx.GetRequiredSteps()
				Expect(steps).NotTo(ContainElement(wizard.StepFeatures))
			})

			It("should include features step by default (simple wizard)", func() {
				// By default, hobby projects still show the features step
				// but with simplified options
				steps := ctx.GetRequiredSteps()
				Expect(steps).To(ContainElement(wizard.StepFeatures))
			})
		})

		Context("for testing project type", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeTesting
				ctx.DatabaseType = generated.DatabaseTypeSQLite
			})

			It("should skip features step when SkipOptionalSteps is true", func() {
				ctx.SkipOptionalSteps = true
				steps := ctx.GetRequiredSteps()
				Expect(steps).NotTo(ContainElement(wizard.StepFeatures))
			})

			It("should include features step by default", func() {
				steps := ctx.GetRequiredSteps()
				Expect(steps).To(ContainElement(wizard.StepFeatures))
			})
		})

		Context("for enterprise project type", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeEnterprise
				ctx.DatabaseType = generated.DatabaseTypePostgreSQL
			})

			It("should include features step", func() {
				steps := ctx.GetRequiredSteps()
				Expect(steps).To(ContainElement(wizard.StepFeatures))
			})

			It("should include advanced step", func() {
				steps := ctx.GetRequiredSteps()
				Expect(steps).To(ContainElement(wizard.StepAdvanced))
			})

			It("should include review step", func() {
				steps := ctx.GetRequiredSteps()
				Expect(steps).To(ContainElement(wizard.StepReview))
			})
		})

		Context("for microservice project type", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeMicroservice
				ctx.DatabaseType = generated.DatabaseTypePostgreSQL
			})

			It("should include features step", func() {
				steps := ctx.GetRequiredSteps()
				Expect(steps).To(ContainElement(wizard.StepFeatures))
			})

			It("should not include advanced step", func() {
				steps := ctx.GetRequiredSteps()
				Expect(steps).NotTo(ContainElement(wizard.StepAdvanced))
			})
		})
	})

	Describe("Database-Specific Features", func() {
		Context("for PostgreSQL", func() {
			BeforeEach(func() {
				ctx.DatabaseType = generated.DatabaseTypePostgreSQL
			})

			It("should enable UUIDs by default", func() {
				Expect(ctx.ShouldEnableUUIDs()).To(BeTrue())
			})

			It("should enable JSON by default", func() {
				Expect(ctx.ShouldEnableJSON()).To(BeTrue())
			})

			It("should enable arrays", func() {
				Expect(ctx.ShouldEnableArrays()).To(BeTrue())
			})

			It("should enable full-text search", func() {
				Expect(ctx.ShouldEnableFullText()).To(BeTrue())
			})

			It("should return 4 database-specific features", func() {
				features := ctx.GetDatabaseSpecificFeatures()
				Expect(features).To(HaveLen(4))
			})
		})

		Context("for MySQL", func() {
			BeforeEach(func() {
				ctx.DatabaseType = generated.DatabaseTypeMySQL
			})

			It("should enable UUIDs by default", func() {
				Expect(ctx.ShouldEnableUUIDs()).To(BeTrue())
			})

			It("should enable JSON by default", func() {
				Expect(ctx.ShouldEnableJSON()).To(BeTrue())
			})

			It("should not enable arrays", func() {
				Expect(ctx.ShouldEnableArrays()).To(BeFalse())
			})

			// MySQL has full-text but it's limited, so the implementation
			// only enables it for PostgreSQL by default
			It("should not enable full-text by default", func() {
				Expect(ctx.ShouldEnableFullText()).To(BeFalse())
			})

			It("should return 3 database-specific features", func() {
				features := ctx.GetDatabaseSpecificFeatures()
				Expect(features).To(HaveLen(3))
			})
		})

		Context("for SQLite", func() {
			BeforeEach(func() {
				ctx.DatabaseType = generated.DatabaseTypeSQLite
			})

			It("should not enable UUIDs by default", func() {
				Expect(ctx.ShouldEnableUUIDs()).To(BeFalse())
			})

			It("should not enable JSON by default", func() {
				Expect(ctx.ShouldEnableJSON()).To(BeFalse())
			})

			It("should not enable arrays", func() {
				Expect(ctx.ShouldEnableArrays()).To(BeFalse())
			})

			It("should not enable full-text search", func() {
				Expect(ctx.ShouldEnableFullText()).To(BeFalse())
			})

			It("should return 1 database-specific feature", func() {
				features := ctx.GetDatabaseSpecificFeatures()
				Expect(features).To(HaveLen(1))
			})
		})
	})

	Describe("Project Type Features", func() {
		Context("for enterprise", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeEnterprise
			})

			It("should show advanced features", func() {
				Expect(ctx.ShouldShowAdvancedFeatures()).To(BeTrue())
			})
		})

		Context("for multi-tenant", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeMultiTenant
			})

			It("should show advanced features", func() {
				Expect(ctx.ShouldShowAdvancedFeatures()).To(BeTrue())
			})
		})

		Context("for hobby", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeHobby
			})

			It("should not show advanced features", func() {
				Expect(ctx.ShouldShowAdvancedFeatures()).To(BeFalse())
			})

			It("should not show database features", func() {
				Expect(ctx.ShouldShowDatabaseFeatures()).To(BeFalse())
			})
		})

		Context("for microservice", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeMicroservice
			})

			It("should show database features", func() {
				Expect(ctx.ShouldShowDatabaseFeatures()).To(BeTrue())
			})
		})
	})

	Describe("Step Tracking", func() {
		It("should mark steps as completed", func() {
			Expect(ctx.IsStepCompleted(wizard.StepProjectType)).To(BeFalse())

			ctx.MarkStepCompleted(wizard.StepProjectType)

			Expect(ctx.IsStepCompleted(wizard.StepProjectType)).To(BeTrue())
		})

		It("should mark steps as skipped", func() {
			Expect(ctx.IsStepSkipped(wizard.StepFeatures)).To(BeFalse())

			ctx.MarkStepSkipped(wizard.StepFeatures)

			Expect(ctx.IsStepSkipped(wizard.StepFeatures)).To(BeTrue())
		})
	})

	Describe("Template Data Update", func() {
		It("should update context from template data", func() {
			data := &generated.TemplateData{
				ProjectType: generated.ProjectTypeMicroservice,
				Database: generated.DatabaseConfig{
					Engine:    generated.DatabaseTypePostgreSQL,
					UseUUIDs:  true,
					UseJSON:   true,
					UseArrays: true,
				},
				ProjectName: "test-project",
				Package: generated.PackageConfig{
					Name: "db",
				},
			}

			ctx.UpdateFromTemplateData(data)

			Expect(ctx.ProjectType).To(Equal(generated.ProjectTypeMicroservice))
			Expect(ctx.DatabaseType).To(Equal(generated.DatabaseTypePostgreSQL))
			Expect(ctx.ProjectName).To(Equal("test-project"))
			Expect(ctx.PackageName).To(Equal("db"))
			Expect(ctx.EnableUUIDs).To(BeTrue())
			Expect(ctx.EnableJSON).To(BeTrue())
			Expect(ctx.EnableArrays).To(BeTrue())
		})
	})
})

var _ = Describe("Branching Policy", func() {
	var policy *wizard.DefaultBranchingPolicy
	var ctx *wizard.FlowContext

	BeforeEach(func() {
		policy = wizard.NewDefaultBranchingPolicy()
		ctx = wizard.NewFlowContext()
	})

	Describe("Step Visibility", func() {
		Context("with hobby project", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeHobby
			})

			It("should show project type step", func() {
				Expect(policy.ShouldShowStep(wizard.StepProjectType, ctx)).To(BeTrue())
			})

			It("should not show features step", func() {
				Expect(policy.ShouldShowStep(wizard.StepFeatures, ctx)).To(BeFalse())
			})
		})

		Context("with enterprise project", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeEnterprise
			})

			It("should show features step", func() {
				Expect(policy.ShouldShowStep(wizard.StepFeatures, ctx)).To(BeTrue())
			})

			It("should show advanced step", func() {
				Expect(policy.ShouldShowStep(wizard.StepAdvanced, ctx)).To(BeTrue())
			})

			It("should show review step", func() {
				Expect(policy.ShouldShowStep(wizard.StepReview, ctx)).To(BeTrue())
			})
		})
	})

	Describe("Feature Visibility", func() {
		Context("with PostgreSQL database", func() {
			BeforeEach(func() {
				ctx.DatabaseType = generated.DatabaseTypePostgreSQL
			})

			It("should show uuid feature", func() {
				Expect(policy.ShouldShowFeature("uuid", ctx)).To(BeTrue())
			})

			It("should show array feature", func() {
				Expect(policy.ShouldShowFeature("array", ctx)).To(BeTrue())
			})
		})

		Context("with SQLite database", func() {
			BeforeEach(func() {
				ctx.DatabaseType = generated.DatabaseTypeSQLite
			})

			It("should not show uuid feature", func() {
				Expect(policy.ShouldShowFeature("uuid", ctx)).To(BeFalse())
			})

			It("should not show array feature", func() {
				Expect(policy.ShouldShowFeature("array", ctx)).To(BeFalse())
			})

			It("should show json feature", func() {
				Expect(policy.ShouldShowFeature("json", ctx)).To(BeTrue())
			})
		})
	})

	Describe("Feature Defaults", func() {
		Context("for enterprise project", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeEnterprise
				ctx.DatabaseType = generated.DatabaseTypePostgreSQL
			})

			It("should default interface to true", func() {
				Expect(policy.GetFeatureDefault("interface", ctx)).To(BeTrue())
			})

			It("should default json_tags to true", func() {
				Expect(policy.GetFeatureDefault("json_tags", ctx)).To(BeTrue())
			})
		})

		Context("for hobby project", func() {
			BeforeEach(func() {
				ctx.ProjectType = generated.ProjectTypeHobby
				ctx.DatabaseType = generated.DatabaseTypeSQLite
			})

			It("should default interface to false", func() {
				Expect(policy.GetFeatureDefault("interface", ctx)).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Simple Branching Policy", func() {
	var policy *wizard.SimpleBranchingPolicy
	var ctx *wizard.FlowContext

	BeforeEach(func() {
		policy = wizard.NewSimpleBranchingPolicy()
		ctx = wizard.NewFlowContext()
	})

	It("should always show all steps", func() {
		Expect(policy.ShouldShowStep(wizard.StepProjectType, ctx)).To(BeTrue())
		Expect(policy.ShouldShowStep(wizard.StepFeatures, ctx)).To(BeTrue())
		Expect(policy.ShouldShowStep(wizard.StepAdvanced, ctx)).To(BeTrue())
	})

	It("should always show all features", func() {
		Expect(policy.ShouldShowFeature("uuid", ctx)).To(BeTrue())
		Expect(policy.ShouldShowFeature("array", ctx)).To(BeTrue())
	})

	It("should always return false for feature defaults", func() {
		Expect(policy.GetFeatureDefault("uuid", ctx)).To(BeFalse())
		Expect(policy.GetFeatureDefault("json", ctx)).To(BeFalse())
	})
})
