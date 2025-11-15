package events

import (
	"time"
)

// FilesGeneratedEventData represents typed event data for file generation
// Replaces 'any' type with specific typed structure
type FilesGeneratedEventData struct {
	ProjectID      string    `json:"project_id"`
	TemplateType   string    `json:"template_type"`
	GeneratedFiles []string   `json:"generated_files"`
	OutputDir      string    `json:"output_dir"`
	FileCount      int       `json:"file_count"`
	TotalSize      int64     `json:"total_size"`
	GeneratedBy    string    `json:"generated_by,omitempty"`
	GeneratedAt    time.Time `json:"generated_at"`
}

// EventType returns specific event type
func (fged *FilesGeneratedEventData) EventType() string {
	return "files.generated"
}

// OccurredAt returns when files were generated
func (fged *FilesGeneratedEventData) OccurredAt() time.Time {
	return fged.GeneratedAt
}

// AggregateID returns the project ID as the aggregate identifier
func (fged *FilesGeneratedEventData) AggregateID() string {
	return fged.ProjectID
}

// Data returns the event data as a safe map for JSON serialization
func (fged *FilesGeneratedEventData) Data() map[string]interface{} {
	return map[string]interface{}{
		"project_id":      fged.ProjectID,
		"template_type":   fged.TemplateType,
		"generated_files": fged.GeneratedFiles,
		"output_dir":      fged.OutputDir,
		"file_count":      fged.FileCount,
		"total_size":      fged.TotalSize,
		"generated_by":    fged.GeneratedBy,
		"generated_at":    fged.GeneratedAt,
	}
}

// Validate validates the files generated event data
func (fged *FilesGeneratedEventData) Validate() error {
	if fged.ProjectID == "" {
		return &EventValidationError{
			Code:    "MISSING_PROJECT_ID",
			Message: "Project ID is required for files generated event",
		}
	}
	
	if fged.TemplateType == "" {
		return &EventValidationError{
			Code:    "MISSING_TEMPLATE_TYPE",
			Message: "Template type is required for files generated event",
		}
	}
	
	if len(fged.GeneratedFiles) == 0 {
		return &EventValidationError{
			Code:    "NO_GENERATED_FILES",
			Message: "At least one generated file is required",
		}
	}
	
	if fged.OutputDir == "" {
		return &EventValidationError{
			Code:    "MISSING_OUTPUT_DIR",
			Message: "Output directory is required for files generated event",
		}
	}
	
	if fged.FileCount <= 0 {
		return &EventValidationError{
			Code:    "INVALID_FILE_COUNT",
			Message: "File count must be greater than zero",
		}
	}
	
	if fged.GeneratedAt.IsZero() {
		return &EventValidationError{
			Code:    "INVALID_GENERATED_AT",
			Message: "Generated at time is required for files generated event",
		}
	}
	
	return nil
}

// NewFilesGeneratedEventData creates a new files generated event data with validation
func NewFilesGeneratedEventData(projectID, templateType, outputDir, generatedBy string, generatedFiles []string, totalSize int64) (*FilesGeneratedEventData, error) {
	if len(generatedFiles) == 0 {
		return nil, &EventValidationError{
			Code:    "EMPTY_FILES_LIST",
			Message: "Generated files list cannot be empty",
		}
	}
	
	if totalSize < 0 {
		return nil, &EventValidationError{
			Code:    "NEGATIVE_TOTAL_SIZE",
			Message: "Total size cannot be negative",
		}
	}
	
	event := &FilesGeneratedEventData{
		ProjectID:      projectID,
		TemplateType:   templateType,
		GeneratedFiles: generatedFiles,
		OutputDir:      outputDir,
		FileCount:      len(generatedFiles),
		TotalSize:      totalSize,
		GeneratedBy:    generatedBy,
		GeneratedAt:    time.Now(),
	}
	
	if err := event.Validate(); err != nil {
		return nil, err
	}
	
	return event, nil
}