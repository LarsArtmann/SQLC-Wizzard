package config

import (
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"gopkg.in/yaml.v3"
)

// PathOrPaths represents a field that can be either a single path string
// or an array of path strings in YAML. This provides type-safe handling
// of the flexible queries/schema fields in sqlc.yaml.
//
// This type eliminates the need for interface{} and provides compile-time
// type safety with proper YAML marshaling/unmarshaling.
type PathOrPaths struct {
	paths []string
}

// NewPathOrPaths creates a PathOrPaths from a slice of strings.
func NewPathOrPaths(paths []string) PathOrPaths {
	if paths == nil {
		paths = []string{}
	}
	return PathOrPaths{paths: paths}
}

// NewSinglePath creates a PathOrPaths from a single string.
func NewSinglePath(path string) PathOrPaths {
	return PathOrPaths{paths: []string{path}}
}

// Strings returns the paths as a slice of strings.
// This is the primary accessor for the underlying data.
func (p PathOrPaths) Strings() []string {
	if p.paths == nil {
		return []string{}
	}
	return p.paths
}

// UnmarshalYAML implements yaml.Unmarshaler to handle both string and []string.
// This allows sqlc.yaml to use either:
//
//	queries: "path/to/queries"
//
// or:
//
//	queries:
//	  - "path/to/queries1"
//	  - "path/to/queries2"
func (p *PathOrPaths) UnmarshalYAML(value *yaml.Node) error {
	// Try to unmarshal as a slice of strings first
	var paths []string
	if err := value.Decode(&paths); err == nil {
		p.paths = paths
		return nil
	}

	// If that fails, try as a single string
	var singlePath string
	if err := value.Decode(&singlePath); err == nil {
		p.paths = []string{singlePath}
		return nil
	}

	return apperrors.Newf(
		apperrors.ErrorCodeInvalidValue,
		"path_or_paths must be either a string or array of strings, got: %v", value.Kind,
	)
}

// MarshalYAML implements yaml.Marshaler to output as []string.
// We always output as an array for consistency, even if there's only one path.
func (p PathOrPaths) MarshalYAML() (any, error) {
	return p.paths, nil
}

// IsEmpty returns true if there are no paths.
func (p PathOrPaths) IsEmpty() bool {
	return len(p.paths) == 0
}

// First returns the first path, or empty string if no paths exist.
// Useful when you know there's typically only one path.
func (p PathOrPaths) First() string {
	if len(p.paths) == 0 {
		return ""
	}
	return p.paths[0]
}
