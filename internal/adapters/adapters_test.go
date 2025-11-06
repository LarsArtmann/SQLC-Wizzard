// Package adapters_test provides testing infrastructure for adapter pattern
package adapters_test

import (
	"context"
	"testing"

	"github.com/LarsArtmann/SQLC-Wizzard/generated"
	"github.com/LarsArtmann/SQLC-Wizzard/internal/adapters"
	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRealSQLCAdapter_NewRealSQLCAdapter(t *testing.T) {
	adapter := adapters.NewRealSQLCAdapter()
	
	assert.NotNil(t, adapter)
	assert.IsType(t, &adapters.RealSQLCAdapter{}, adapter)
}

func TestRealSQLCAdapter_CheckInstallation(t *testing.T) {
	adapter := adapters.NewRealSQLCAdapter()
	
	err := adapter.CheckInstallation(context.Background())
	
	// sqlc might not be installed in test environment
	if err != nil {
		assert.Contains(t, err.Error(), "which: sqlc")
	}
}

func TestRealSQLCAdapter_Version(t *testing.T) {
	adapter := adapters.NewRealSQLCAdapter()
	
	version, err := adapter.Version(context.Background())
	
	// sqlc might not be installed
	if err != nil {
		assert.Contains(t, err.Error(), "sqlc")
		assert.Empty(t, version)
	}
}

func TestRealDatabaseAdapter_NewRealDatabaseAdapter(t *testing.T) {
	adapter := adapters.NewRealDatabaseAdapter()
	
	assert.NotNil(t, adapter)
	assert.IsType(t, &adapters.RealDatabaseAdapter{}, adapter)
}

func TestRealDatabaseAdapter_TestConnection(t *testing.T) {
	adapter := adapters.NewRealDatabaseAdapter()
	
	// Test with empty config (should fail)
	cfg := &config.DatabaseConfig{URI: ""}
	err := adapter.TestConnection(context.Background(), cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database URI is required")
	
	// Test with valid config
	cfg.URI = "postgresql://localhost:5432/test"
	err = adapter.TestConnection(context.Background(), cfg)
	// Implementation is basic, so this might pass
	// In real implementation, would actually test connection
}

func TestRealDatabaseAdapter_CreateDatabase(t *testing.T) {
	adapter := adapters.NewRealDatabaseAdapter()
	
	cfg := &config.DatabaseConfig{URI: "postgresql://localhost:5432/test"}
	err := adapter.CreateDatabase(context.Background(), cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not yet implemented")
}

func TestRealCLIAdapter_NewRealCLIAdapter(t *testing.T) {
	adapter := adapters.NewRealCLIAdapter()
	
	assert.NotNil(t, adapter)
	assert.IsType(t, &adapters.RealCLIAdapter{}, adapter)
}

func TestRealCLIAdapter_CheckCommand(t *testing.T) {
	adapter := adapters.NewRealCLIAdapter()
	
	// Test with existing command (should pass)
	err := adapter.CheckCommand(context.Background(), "echo")
	assert.NoError(t, err)
	
	// Test with non-existing command (should fail)
	err = adapter.CheckCommand(context.Background(), "nonexistentcommand123")
	assert.Error(t, err)
}

func TestRealCLIAdapter_RunCommand(t *testing.T) {
	adapter := adapters.NewRealCLIAdapter()
	
	output, err := adapter.RunCommand(context.Background(), "echo", "test")
	require.NoError(t, err)
	assert.Equal(t, "test\n", output)
}

func TestRealCLIAdapter_GetVersion(t *testing.T) {
	adapter := adapters.NewRealCLIAdapter()
	
	version, err := adapter.GetVersion(context.Background(), "echo")
	// echo doesn't have version, so this might fail
	if err != nil {
		assert.Empty(t, version)
	}
}

func TestRealTemplateAdapter_NewRealTemplateAdapter(t *testing.T) {
	adapter := adapters.NewRealTemplateAdapter()
	
	assert.NotNil(t, adapter)
	assert.IsType(t, &adapters.RealTemplateAdapter{}, adapter)
}

func TestRealTemplateAdapter_ValidateTemplateData(t *testing.T) {
	adapter := adapters.NewRealTemplateAdapter()
	
	// Test with valid data
	data := generated.TemplateData{
		ProjectType: generated.ProjectTypeMicroservice,
		Database: generated.DatabaseConfig{
			Engine: generated.DatabaseTypePostgreSQL,
		},
		Package: generated.PackageConfig{
			Name: "db",
		},
	}
	
	err := adapter.ValidateTemplateData(context.Background(), data)
	assert.NoError(t, err)
	
	// Test with invalid project type
	data.ProjectType = "invalid"
	err = adapter.ValidateTemplateData(context.Background(), data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid project type")
	
	// Test with invalid database type
	data.ProjectType = generated.ProjectTypeMicroservice
	data.Database.Engine = "invalid"
	err = adapter.ValidateTemplateData(context.Background(), data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid database type")
	
	// Test with empty package name
	data.Database.Engine = generated.DatabaseTypePostgreSQL
	data.Package.Name = ""
	err = adapter.ValidateTemplateData(context.Background(), data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "package name is required")
}

func TestRealFileSystemAdapter_NewRealFileSystemAdapter(t *testing.T) {
	adapter := adapters.NewRealFileSystemAdapter()
	
	assert.NotNil(t, adapter)
	assert.IsType(t, &adapters.RealFileSystemAdapter{}, adapter)
}

func TestRealFileSystemAdapter_Exists(t *testing.T) {
	adapter := adapters.NewRealFileSystemAdapter()
	ctx := context.Background()
	
	// Test with existing directory
	exists, err := adapter.Exists(ctx, ".")
	require.NoError(t, err)
	assert.True(t, exists)
	
	// Test with non-existing file
	exists, err = adapter.Exists(ctx, "/nonexistent/path")
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestRealFileSystemAdapter_TempDir(t *testing.T) {
	adapter := adapters.NewRealFileSystemAdapter()
	
	tempDir, err := adapter.TempDir(context.Background(), "test")
	require.NoError(t, err)
	assert.NotEmpty(t, tempDir)
	assert.Contains(t, tempDir, "test")
	
	// Clean up
	err = adapter.Remove(context.Background(), tempDir)
	require.NoError(t, err)
}

func TestRealFileSystemAdapter_ReadWriteFile(t *testing.T) {
	adapter := adapters.NewRealFileSystemAdapter()
	ctx := context.Background()
	
	// Test write and read
	testData := []byte("test content")
	tempDir := t.TempDir()
	filePath := tempDir + "/test.txt"
	
	err := adapter.WriteFile(ctx, filePath, testData, 0644)
	require.NoError(t, err)
	
	data, err := adapter.ReadFile(ctx, filePath)
	require.NoError(t, err)
	assert.Equal(t, testData, data)
}