// Package templates_test provides basic testing for template types
package templates_test

import (
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/templates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProjectType_Valid(t *testing.T) {
	result, err := templates.NewProjectType("microservice")

	require.NoError(t, err)
	assert.Equal(t, generated.ProjectTypeMicroservice, result)
}

func TestNewProjectType_Invalid(t *testing.T) {
	result, err := templates.NewProjectType("invalid")

	require.Error(t, err)
	assert.Empty(t, result)
	assert.True(t, apperrors.Is(err, apperrors.ErrInvalidType))
}

func TestNewDatabaseType_Valid(t *testing.T) {
	result, err := templates.NewDatabaseType("postgresql")

	require.NoError(t, err)
	assert.Equal(t, generated.DatabaseTypePostgreSQL, result)
}

func TestNewDatabaseType_Invalid(t *testing.T) {
	result, err := templates.NewDatabaseType("invalid")

	require.Error(t, err)
	assert.Empty(t, result)
	assert.True(t, apperrors.Is(err, apperrors.ErrInvalidType))
}

func TestMustNewProjectType_Valid(t *testing.T) {
	result := templates.MustNewProjectType("microservice")
	assert.Equal(t, generated.ProjectTypeMicroservice, result)
}

func TestMustNewProjectType_Invalid(t *testing.T) {
	assert.Panics(t, func() {
		templates.MustNewProjectType("invalid")
	})
}

func TestMustNewDatabaseType_Valid(t *testing.T) {
	result := templates.MustNewDatabaseType("postgresql")
	assert.Equal(t, generated.DatabaseTypePostgreSQL, result)
}

func TestMustNewDatabaseType_Invalid(t *testing.T) {
	assert.Panics(t, func() {
		templates.MustNewDatabaseType("invalid")
	})
}

var (
	validProjectTypes    = []string{"microservice", "hobby"}
	invalidProjectTypes  = []string{"invalid", ""}
	validDatabaseTypes   = []string{"postgresql", "mysql"}
	invalidDatabaseTypes = []string{"invalid", ""}
)

func TestIsValidProjectType(t *testing.T) {
	for _, v := range validProjectTypes {
		assert.True(t, templates.IsValidProjectType(v), "expected %q to be valid", v)
	}

	for _, v := range invalidProjectTypes {
		assert.False(t, templates.IsValidProjectType(v), "expected %q to be invalid", v)
	}
}

func TestIsValidDatabaseType(t *testing.T) {
	for _, v := range validDatabaseTypes {
		assert.True(t, templates.IsValidDatabaseType(v), "expected %q to be valid", v)
	}

	for _, v := range invalidDatabaseTypes {
		assert.False(t, templates.IsValidDatabaseType(v), "expected %q to be invalid", v)
	}
}

func TestProjectTypeConstants(t *testing.T) {
	assert.Equal(t, "hobby", string(templates.ProjectTypeHobby))
	assert.Equal(t, "microservice", string(templates.ProjectTypeMicroservice))
	assert.Equal(t, "enterprise", string(templates.ProjectTypeEnterprise))
}

func TestDatabaseTypeConstants(t *testing.T) {
	assert.Equal(t, "postgresql", string(templates.DatabaseTypePostgreSQL))
	assert.Equal(t, "mysql", string(templates.DatabaseTypeMySQL))
	assert.Equal(t, "sqlite", string(templates.DatabaseTypeSQLite))
}
