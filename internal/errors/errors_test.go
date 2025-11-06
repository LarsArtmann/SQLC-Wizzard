// Package errors_test provides basic testing for error handling
package errors_test

import (
	"testing"

	apperrors "github.com/LarsArtmann/SQLC-Wizzard/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_Error(t *testing.T) {
	err := apperrors.New(apperrors.ErrInvalidConfig, "test error")
	
	require.Error(t, err)
	assert.Contains(t, err.Error(), "INVALID_CONFIG")
	assert.Contains(t, err.Error(), "test error")
}

func TestValidationError_Helper(t *testing.T) {
	err := apperrors.ValidationError("field", "invalid")
	
	require.Error(t, err)
	assert.True(t, apperrors.Is(err, apperrors.ErrInvalidType))
	assert.Contains(t, err.Error(), "invalid value for field")
}

func TestFileNotFoundError_Helper(t *testing.T) {
	path := "/test/path"
	err := apperrors.FileNotFoundError(path)
	
	require.Error(t, err)
	assert.True(t, apperrors.Is(err, apperrors.ErrFileNotFound))
	assert.Contains(t, err.Error(), "file not found")
	assert.Contains(t, err.Error(), path)
}

func TestTemplateNotFoundError_Helper(t *testing.T) {
	template := "missing-template"
	err := apperrors.TemplateNotFoundError(template)
	
	require.Error(t, err)
	assert.True(t, apperrors.Is(err, apperrors.ErrTemplateNotFound))
	assert.Contains(t, err.Error(), "template not found")
	assert.Contains(t, err.Error(), template)
}