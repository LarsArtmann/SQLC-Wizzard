package wizard_test

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// MockStep is a mock implementation of wizard.StepInterface for testing
// TODO: Extract to internal/testing/mocks directory
// TODO: Add validation for step execution order
// TODO: Add state transition validation
// TODO: Add execution duration tracking.
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

// NewMockStep creates a fresh mock step instance
// TODO: Add configurable default behaviors
// TODO: Add validation constraints.
func NewMockStep() *MockStep {
	return &MockStep{}
}

// Execute tracks step execution and applies configured behavior
// TODO: Add input validation
// TODO: Add execution state tracking
// TODO: Add timing information.
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

// ValidateConfiguration tracks validation calls and applies configured behavior
// TODO: Add validation rule testing
// TODO: Add validation coverage tracking.
func (m *MockStep) ValidateConfiguration(data *generated.TemplateData) error {
	m.ValidateCalls++
	if m.ShouldValidateFail {
		return m.ValidateError
	}
	return nil
}

// MockTemplate is a mock implementation of templates.Template for testing
// TODO: Extract to internal/testing/mocks directory
// TODO: Add template feature validation
// TODO: Add generation result verification.
type MockTemplate struct {
	GenerateCalls int
	ShouldFail    bool
	FailError     error
	LastCallData  generated.TemplateData
}

// NewMockTemplate creates a fresh mock template instance
// TODO: Add configurable template types
// TODO: Add feature set configuration.
func NewMockTemplate() *MockTemplate {
	return &MockTemplate{}
}

// Generate tracks template generation and returns mock configuration
// TODO: Add result validation
// TODO: Add generation optimization tracking.
func (m *MockTemplate) Generate(data generated.TemplateData) (*config.SqlcConfig, error) {
	m.GenerateCalls++
	m.LastCallData = data

	if m.ShouldFail {
		return nil, m.FailError
	}

	// Return a basic valid config
	// TODO: Make this configurable based on input data
	// TODO: Add more realistic configuration generation
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

// DefaultData returns default template data for testing
// TODO: Make this configurable
// TODO: Add multiple default data sets for different scenarios.
func (m *MockTemplate) DefaultData() generated.TemplateData {
	return generated.TemplateData{
		ProjectType: generated.ProjectTypeMicroservice,
	}
}

// RequiredFeatures returns required features for the template
// TODO: Make this configurable
// TODO: Add feature validation.
func (m *MockTemplate) RequiredFeatures() []string {
	return []string{"basic"}
}

// Name returns the template name
// TODO: Make this configurable.
func (m *MockTemplate) Name() string {
	return "Mock Template"
}

// Description returns the template description
// TODO: Make this configurable.
func (m *MockTemplate) Description() string {
	return "A mock template for testing"
}
