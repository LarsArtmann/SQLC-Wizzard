package testing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/utils"
)

// TDDTemplate represents a TDD test template
type TDDTemplate struct {
	Name        string
	Description string
	TestName    string
	RedPhase    string
	GreenPhase  string
	Refactor    string
}

// TTDScaffold provides TDD workflow scaffolding
type TTDScaffold struct {
	templates map[string]*TDDTemplate
}

// NewTTDScaffold creates a new TDD scaffolder
func NewTTDScaffold() *TTDScaffold {
	return &TTDScaffold{
		templates: map[string]*TDDTemplate{
			"command": {
				Name:        "Command Pattern",
				Description: "TDD template for command implementation",
				TestName:    "TestCommand_Execute",
				RedPhase:    "// RED: Write failing test\nfunc (c *Command) Execute() error {\n\treturn fmt.Errorf(\"not implemented\")\n}",
				GreenPhase:  "// GREEN: Make test pass\nfunc (c *Command) Execute() error {\n\t// Basic implementation to pass test\n\treturn c.validate()\n}",
				Refactor:    "// REFACTOR: Clean up\nfunc (c *Command) Execute() error {\n\tif err := c.validate(); err != nil {\n\t\treturn err\n\t}\n\treturn c.execute()\n}",
			},
			"service": {
				Name:        "Service Pattern",
				Description: "TDD template for service implementation",
				TestName:    "TestService_Process",
				RedPhase:    "// RED: Write failing test\nfunc (s *Service) Process(data interface{}) error {\n\treturn fmt.Errorf(\"not implemented\")\n}",
				GreenPhase:  "// GREEN: Make test pass\nfunc (s *Service) Process(data interface{}) error {\n\t// Basic processing\n\treturn s.basicProcess(data)\n}",
				Refactor:    "// REFACTOR: Clean up\nfunc (s *Service) Process(data interface{}) error {\n\tif err := s.validate(data); err != nil {\n\t\treturn err\n\t}\n\treturn s.process(data)\n}",
			},
			"repository": {
				Name:        "Repository Pattern",
				Description: "TDD template for repository implementation",
				TestName:    "TestRepository_FindByID",
				RedPhase:    "// RED: Write failing test\nfunc (r *Repository) FindByID(id string) (interface{}, error) {\n\treturn nil, fmt.Errorf(\"not implemented\")\n}",
				GreenPhase:  "// GREEN: Make test pass\nfunc (r *Repository) FindByID(id string) (interface{}, error) {\n\t// Basic lookup\n\treturn r.basicFind(id)\n}",
				Refactor:    "// REFACTOR: Clean up\nfunc (r *Repository) FindByID(id string) (interface{}, error) {\n\tif id == \"\" {\n\t\treturn nil, fmt.Errorf(\"id required\")\n\t}\n\treturn r.find(id)\n}",
			},
		},
	}
}

// GenerateTestFile creates a TDD test file from template
func (s *TTDScaffold) GenerateTestFile(templateType, filePath, testName string) error {
	tmpl, exists := s.templates[templateType]
	if !exists {
		return fmt.Errorf("template type '%s' not found", templateType)
	}

	testContent := s.generateTestContent(tmpl, testName)

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write test file
	return os.WriteFile(filePath, []byte(testContent), 0644)
}

// generateTestContent creates the complete test file content
func (s *TTDScaffold) generateTestContent(tmpl *TDDTemplate, testName string) string {
	content := fmt.Sprintf(`// %s - %s
// Generated with TDD scaffolding - follow RED -> GREEN -> REFACTOR workflow

package %s

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RED PHASE: Write failing test first
func %s(t *testing.T) {
	// Write this test to FAIL first - this is the RED phase
	// Test the behavior you want to implement
	
	// Example assertion that will fail initially
	result, err := YourImplementation()
	require.Error(t, err) // This should fail initially
	assert.Nil(t, result)    // This should fail initially
}

// GREEN PHASE: Make the test pass with minimal implementation
// After this test passes, move to the next phase

// REFACTOR PHASE: Clean up the implementation
// Maintain green tests while improving code quality
`, tmpl.Name, tmpl.Description, filepath.Base(filepath.Dir(filePath)), testName)

	return content
}

// CreateTDDWorkflow creates a complete TDD workflow directory
func (s *TTDScaffold) CreateTDDWorkflow(featureName, packagePath string) error {
	workflowDir := filepath.Join(packagePath, featureName+"_tdd")
	
	// Create workflow directory structure
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return err
	}

	// Create README with TDD instructions
	readme := fmt.Sprintf(`# %s - TDD Workflow

## 🚦 RED -> GREEN -> REFACTOR

### Phase 1: RED
1. Write a failing test
2. Run the test to ensure it fails
3. Verify the failure makes sense

### Phase 2: GREEN  
1. Write minimal code to make the test pass
2. Run the test to ensure it passes
3. No refactoring yet - just make it work

### Phase 3: REFACTOR
1. Clean up the implementation
2. Ensure test still passes
3. Improve code quality while maintaining functionality

## 📋 Checklist Before Commit
- [ ] Test passes
- [ ] Code follows project standards
- [ ] No duplication introduced
- [ ] Error handling is comprehensive

## 🎯 Success Criteria
- Test coverage > 80%%
- Code passes all linting
- Implementation follows SOLID principles
- Error cases are handled properly
`, featureName)

	if err := os.WriteFile(filepath.Join(workflowDir, "README.md"), []byte(readme), 0644); err != nil {
		return err
	}

	// Create test template
	testTemplate := s.generateTestContent(s.templates["command"], "TestFeature_Execute")
	testFile := filepath.Join(workflowDir, "feature_test.go")
	
	return os.WriteFile(testFile, []byte(testTemplate), 0644)
}

// ListTemplates returns available TDD templates
func (s *TTDScaffold) ListTemplates() []string {
	var names []string
	for name := range s.templates {
		names = append(names, name)
	}
	return names
}

// ValidateTDDImplementation checks if TDD principles are followed
func ValidateTDDImplementation(packagePath string) error {
	// Check for test files
	testFiles, err := filepath.Glob(filepath.Join(packagePath, "*_test.go"))
	if err != nil {
		return err
	}

	if len(testFiles) == 0 {
		return fmt.Errorf("no test files found in %s", packagePath)
	}

	// Check for TDD workflow markers
	for _, testFile := range testFiles {
		content, err := os.ReadFile(testFile)
		if err != nil {
			return err
		}

		testContent := string(content)
		if !containsTDDPhases(testContent) {
			return fmt.Errorf("test file %s missing TDD phase comments", testFile)
		}
	}

	return nil
}

// containsTDDPhases checks if TDD workflow phases are present
func containsTDDPhases(content string) bool {
	return strings.Contains(content, "RED") &&
		strings.Contains(content, "GREEN") &&
		strings.Contains(content, "REFACTOR")
}