package creators

import (
	"context"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
)

// CreatorConfig represents common configuration for all creators
type CreatorConfig struct {
	ProjectName  string
	ProjectType  generated.ProjectType
	Database     generated.DatabaseType
	TemplateData generated.TemplateData
	OutputPath   string
	Force        bool
}

// Creator defines common interface for all project creators
// Uses generics for type-safe configuration handling
type Creator[T any] interface {
	// Create executes the creation operation with type-safe config
	Create(ctx context.Context, config T) error

	// CanHandle returns true if this creator can handle the given config
	CanHandle(config T) bool

	// Dependencies returns any other creators that must run first
	Dependencies() []string
}

// Result represents the result of a creation operation
type Result struct {
	Success bool
	Message string
	Files   []string
	Errors  []error
}

// IsSuccess returns true if the operation was successful
func (r *Result) IsSuccess() bool {
	return r.Success && len(r.Errors) == 0
}

// AddError adds an error to the result
func (r *Result) AddError(err error) {
	if err != nil {
		r.Errors = append(r.Errors, err)
		r.Success = false
	}
}

// AddFile adds a created file to the result
func (r *Result) AddFile(path string) {
	r.Files = append(r.Files, path)
}

// FileSystemCreator defines interface for creators that work with file system
type FileSystemCreator interface {
	Creator[CreatorConfig]

	// FileSystem returns the file system adapter
	FileSystem() adapters.FileSystemAdapter

	// SetFileSystem sets the file system adapter (for testing)
	SetFileSystem(fs adapters.FileSystemAdapter)
}

// TemplateCreator defines interface for creators that generate templates
type TemplateCreator interface {
	FileSystemCreator

	// Template returns the template name
	Template() string

	// SupportedProjectTypes returns project types this creator supports
	SupportedProjectTypes() []generated.ProjectType
}

// ScaffoldCreator defines interface for high-level scaffolding operations
type ScaffoldCreator interface {
	Creator[CreatorConfig]

	// ScaffoldOrder returns the execution order for this creator
	ScaffoldOrder() int

	// Prerequisites returns any prerequisites for this creator
	Prerequisites() []string
}
