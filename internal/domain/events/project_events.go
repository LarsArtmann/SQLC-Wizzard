package events

import (
	"time"
)

// ProjectCreatedEventData represents typed event data for project creation
// Replaces 'any' type with specific typed structure
type ProjectCreatedEventData struct {
	ProjectID   string    `json:"project_id"`
	Name        string    `json:"name"`
	ProjectType string    `json:"project_type"`
	Database    string    `json:"database"`
	PackageName string    `json:"package_name"`
	PackagePath string    `json:"package_path"`
	OutputDir   string    `json:"output_dir"`
	CreatedBy   string    `json:"created_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// EventType returns the specific event type
func (pced *ProjectCreatedEventData) EventType() string {
	return "project.created"
}

// OccurredAt returns when the project was created
func (pced *ProjectCreatedEventData) OccurredAt() time.Time {
	return pced.CreatedAt
}

// AggregateID returns the project ID as the aggregate identifier
func (pced *ProjectCreatedEventData) AggregateID() string {
	return pced.ProjectID
}

// Data returns the event data as a safe map for JSON serialization
func (pced *ProjectCreatedEventData) Data() map[string]any {
	return map[string]any{
		"project_id":   pced.ProjectID,
		"name":         pced.Name,
		"project_type": pced.ProjectType,
		"database":     pced.Database,
		"package_name": pced.PackageName,
		"package_path": pced.PackagePath,
		"output_dir":   pced.OutputDir,
		"created_by":   pced.CreatedBy,
		"created_at":   pced.CreatedAt,
	}
}

// Validate validates the project created event data
func (pced *ProjectCreatedEventData) Validate() error {
	if pced.ProjectID == "" {
		return &EventValidationError{
			Code:    "MISSING_PROJECT_ID",
			Message: "Project ID is required for project created event",
		}
	}

	if pced.Name == "" {
		return &EventValidationError{
			Code:    "MISSING_NAME",
			Message: "Project name is required for project created event",
		}
	}

	if pced.ProjectType == "" {
		return &EventValidationError{
			Code:    "MISSING_PROJECT_TYPE",
			Message: "Project type is required for project created event",
		}
	}

	if pced.Database == "" {
		return &EventValidationError{
			Code:    "MISSING_DATABASE",
			Message: "Database type is required for project created event",
		}
	}

	if pced.CreatedAt.IsZero() {
		return &EventValidationError{
			Code:    "INVALID_CREATED_AT",
			Message: "Created at time is required for project created event",
		}
	}

	return nil
}

// NewProjectCreatedEventData creates a new project created event data with validation
func NewProjectCreatedEventData(projectID, name, projectType, database, packageName, packagePath, outputDir, createdBy string) (*ProjectCreatedEventData, error) {
	event := &ProjectCreatedEventData{
		ProjectID:   projectID,
		Name:        name,
		ProjectType: projectType,
		Database:    database,
		PackageName: packageName,
		PackagePath: packagePath,
		OutputDir:   outputDir,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	if err := event.Validate(); err != nil {
		return nil, err
	}

	return event, nil
}
