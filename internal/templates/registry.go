package templates

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
)

// Registry manages available templates.
type Registry struct {
	templates map[ProjectType]Template
}

// NewRegistry creates a new template registry.
func NewRegistry() *Registry {
	r := &Registry{
		templates: make(map[ProjectType]Template),
	}

	// Register built-in templates
	r.Register(NewMicroserviceTemplate())
	// TODO: Register other templates
	// r.Register(NewHobbyTemplate())
	// r.Register(NewEnterpriseTemplate())

	return r
}

// Register adds a template to the registry.
func (r *Registry) Register(tmpl Template) {
	r.templates[ProjectType(tmpl.Name())] = tmpl
}

// Get retrieves a template by project type.
func (r *Registry) Get(projectType ProjectType) (Template, error) {
	tmpl, ok := r.templates[projectType]
	if !ok {
		return nil, apperrors.TemplateNotFoundError(string(projectType))
	}
	return tmpl, nil
}

// List returns all available templates.
func (r *Registry) List() []Template {
	templates := make([]Template, 0, len(r.templates))
	for _, tmpl := range r.templates {
		templates = append(templates, tmpl)
	}
	return templates
}

// HasTemplate checks if a template exists.
func (r *Registry) HasTemplate(projectType ProjectType) bool {
	_, ok := r.templates[projectType]
	return ok
}

// DefaultRegistry returns the default template registry.
var defaultRegistry *Registry

func init() {
	defaultRegistry = NewRegistry()
}

// GetTemplate is a convenience function to get a template from the default registry.
func GetTemplate(projectType ProjectType) (Template, error) {
	return defaultRegistry.Get(projectType)
}

// ListTemplates returns all templates from the default registry.
func ListTemplates() []Template {
	return defaultRegistry.List()
}
