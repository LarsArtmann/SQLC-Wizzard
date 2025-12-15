package wizard_test

import (
	"fmt"
)

// MockUI is a mock implementation of wizard.UIInterface for testing
// TODO: Extract to internal/testing/mocks directory
// TODO: Consider using testify/mock or moq for generated mocks
// TODO: Add validation for UI state transitions
// TODO: Add thread safety if needed
type MockUI struct {
	WelcomeCalls     int
	StepHeaders      []string
	StepCompletions  []StepCompletion
	Sections         []string
	InfoMessages     []string
	ShouldFailRun    bool
	FailErrorMessage string
}

// StepCompletion represents a step completion notification
// TODO: Add timestamp for completion tracking
// TODO: Add duration measurement
type StepCompletion struct {
	Title   string
	Message string
}

// NewMockUI creates a fresh mock UI instance
// TODO: Add configurable initial state
// TODO: Add default behaviors
func NewMockUI() *MockUI {
	return &MockUI{
		StepHeaders:     make([]string, 0),
		StepCompletions: make([]StepCompletion, 0),
		Sections:        make([]string, 0),
		InfoMessages:    make([]string, 0),
	}
}

// ShowWelcome tracks welcome calls and handles failures
// TODO: Add validation for call sequence
// TODO: Add state validation
func (m *MockUI) ShowWelcome() {
	m.WelcomeCalls++
	if m.ShouldFailRun {
		panic(fmt.Errorf("UI operation failed: %s", m.FailErrorMessage))
	}
}

// ShowStepHeader tracks step headers
// TODO: Validate header format
// TODO: Prevent duplicate headers
func (m *MockUI) ShowStepHeader(title string) {
	m.StepHeaders = append(m.StepHeaders, title)
}

// ShowStepComplete tracks step completions
// TODO: Validate completion sequence
// TODO: Add completion status tracking
func (m *MockUI) ShowStepComplete(title, message string) {
	m.StepCompletions = append(m.StepCompletions, StepCompletion{
		Title:   title,
		Message: message,
	})
}

// ShowSection tracks section displays
// TODO: Validate section hierarchy
func (m *MockUI) ShowSection(title string) {
	m.Sections = append(m.Sections, title)
}

// ShowInfo tracks info messages
// TODO: Add message level validation
// TODO: Add message categorization
func (m *MockUI) ShowInfo(message string) {
	m.InfoMessages = append(m.InfoMessages, message)
}
