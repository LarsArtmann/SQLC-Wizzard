package adapters

import (
	"context"
	"fmt"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// RealTemplateAdapter provides actual template operations.
type RealTemplateAdapter struct{}

// NewRealTemplateAdapter creates a new real template adapter.
func NewRealTemplateAdapter() *RealTemplateAdapter {
	return &RealTemplateAdapter{}
}

// GetTemplate retrieves a template by type.
func (a *RealTemplateAdapter) GetTemplate(projectType generated.ProjectType) (templates.Template, error) {
	switch projectType {
	case generated.ProjectTypeMicroservice:
		return templates.NewMicroserviceTemplate(), nil
	case generated.ProjectTypeHobby:
		return templates.NewMicroserviceTemplate(), nil // Use microservice as fallback
	case generated.ProjectTypeEnterprise:
		return templates.NewMicroserviceTemplate(), nil // Use microservice as fallback
	default:
		return nil, fmt.Errorf("Unknown project type: %s", projectType)
	}
}

// GenerateConfig generates configuration from template data.
func (a *RealTemplateAdapter) GenerateConfig(ctx context.Context, data generated.TemplateData) (*config.SqlcConfig, error) {
	template, err := a.GetTemplate(data.ProjectType)
	if err != nil {
		return nil, err
	}

	return template.Generate(data)
}

// GenerateFiles generates files from template.
func (a *RealTemplateAdapter) GenerateFiles(ctx context.Context, data generated.TemplateData, outputDir string) ([]string, error) {
	// For now, return empty slice as GenerateFiles is not implemented
	_ = data
	_ = outputDir
	return []string{}, nil
}

// ValidateTemplateData validates template data.
func (a *RealTemplateAdapter) ValidateTemplateData(ctx context.Context, data generated.TemplateData) error {
	// Basic validation
	if data.ProjectType == "" {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "Project type is required")
	}

	if !data.ProjectType.IsValid() {
		return fmt.Errorf("invalid project type: %s", data.ProjectType)
	}

	if data.Database.Engine == "" {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "Database engine is required")
	}

	if !data.Database.Engine.IsValid() {
		return fmt.Errorf("invalid database type: %s", data.Database.Engine)
	}

	// Validate package name
	if data.Package.Name == "" {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "package name is required")
	}

	// For now, just return nil as full Validate is not implemented
	return nil
}

// ListTemplates returns all available templates.
func (a *RealTemplateAdapter) ListTemplates(ctx context.Context) ([]templates.Template, error) {
	templates := []templates.Template{
		templates.NewMicroserviceTemplate(),
	}

	return templates, nil
}
