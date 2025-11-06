package domain

import (
	"github.com/LarsArtmann/SQLC-Wizzard/generated"
)

// ProjectCreated represents a domain event when a project is created
type ProjectCreated struct {
	ProjectID   string                      `json:"project_id"`
	Name        string                      `json:"name"`
	ProjectType generated.ProjectType        `json:"project_type"`
	Database    generated.DatabaseType       `json:"database"`
	CreatedAt   string                      `json:"created_at"`
}

// ConfigValidated represents a domain event when configuration is validated
type ConfigValidated struct {
	ProjectID       string   `json:"project_id"`
	IsValid         bool     `json:"is_valid"`
	ValidationErrors []string `json:"validation_errors,omitempty"`
	ValidatedAt     string   `json:"validated_at"`
}

// FilesGenerated represents a domain event when files are generated
type FilesGenerated struct {
	ProjectID    string   `json:"project_id"`
	OutputFiles  []string `json:"output_files"`
	GeneratedAt  string   `json:"generated_at"`
	TemplateData *generated.TemplateData `json:"template_data,omitempty"`
}

// ProjectValidationError represents a domain event when validation fails
type ProjectValidationError struct {
	ProjectID string `json:"project_id"`
	Field     string `json:"field"`
	Message   string `json:"message"`
	ErrorAt   string `json:"error_at"`
}

// EventType constants for all domain events
const (
	EventTypeProjectCreated       = "ProjectCreated"
	EventTypeConfigValidated      = "ConfigValidated"
	EventTypeFilesGenerated       = "FilesGenerated"
	EventTypeProjectValidationError = "ProjectValidationError"
)