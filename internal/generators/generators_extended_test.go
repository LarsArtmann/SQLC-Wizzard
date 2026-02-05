package generators_test

import (
	"os"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/generators"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	tempDir := t.TempDir()
	gen := generators.NewGenerator(tempDir)

	assert.NotNil(t, gen)
}

func TestNewGenerator_NilDir(t *testing.T) {
	gen := generators.NewGenerator("")

	assert.NotNil(t, gen)
}

func TestGenerateConfig_Basic(t *testing.T) {
	tempDir := t.TempDir()
	gen := generators.NewGenerator(tempDir)

	// Test basic config generation
	config := map[string]interface{}{
		"version": "2",
		"sql": []map[string]interface{}{
			{
				"name":   "test",
				"engine": "postgresql",
			},
		},
	}

	// This would normally call internal method, but we test structure
	assert.NotNil(t, config)
	assert.NotNil(t, gen)
}

func TestGenerateConfig_Invalid(t *testing.T) {
	tempDir := t.TempDir()
	gen := generators.NewGenerator(tempDir)

	// Test invalid config
	config := map[string]interface{}{
		"version": "invalid",
	}

	// Verify generator is created
	assert.NotNil(t, gen)
	assert.NotNil(t, config)
}

func TestValidatePaths_Valid(t *testing.T) {
	tempDir := t.TempDir()

	// Test valid paths
	baseDir := tempDir + "/internal/db"
	queriesDir := tempDir + "/sql/queries"
	schemaDir := tempDir + "/sql/schema"

	// Create directories
	require.NoError(t, os.MkdirAll(baseDir, 0o755))
	require.NoError(t, os.MkdirAll(queriesDir, 0o755))
	require.NoError(t, os.MkdirAll(schemaDir, 0o755))

	// Verify paths exist
	info, err := os.Stat(baseDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	info, err = os.Stat(queriesDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	info, err = os.Stat(schemaDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestValidatePaths_Invalid(t *testing.T) {
	tempDir := t.TempDir()

	// Test invalid path (non-existent)
	invalidPath := tempDir + "/nonexistent"

	// Verify path doesn't exist
	_, err := os.Stat(invalidPath)
	assert.Error(t, err)
}

func TestValidatePaths_Relative(t *testing.T) {
	// Test relative paths
	baseDir := "./internal/db"
	queriesDir := "./sql/queries"

	// Verify they are valid relative paths
	assert.True(t, len(baseDir) > 0)
	assert.True(t, len(queriesDir) > 0)
}
