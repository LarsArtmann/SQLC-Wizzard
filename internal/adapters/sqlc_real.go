package adapters

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/LarsArtmann/SQLC-Wizzard/pkg/config"
)

// RealSQLCAdapter provides actual sqlc operations.
type RealSQLCAdapter struct{}

// NewRealSQLCAdapter creates a new real SQLC adapter.
func NewRealSQLCAdapter() *RealSQLCAdapter {
	return &RealSQLCAdapter{}
}

// Generate generates Go code from SQL files.
func (a *RealSQLCAdapter) Generate(ctx context.Context, cfg *config.SqlcConfig) error {
	cmd := exec.CommandContext(ctx, "sqlc", "generate")
	cmd.Dir = filepath.Dir(".")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlc generate failed: %w", err)
	}

	return nil
}

// Validate validates sqlc configuration.
func (a *RealSQLCAdapter) Validate(ctx context.Context, cfg *config.SqlcConfig) error {
	cmd := exec.CommandContext(ctx, "sqlc", "validate")
	cmd.Dir = filepath.Dir(".")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlc validate failed: %w", err)
	}

	return nil
}

// Version returns sqlc version.
func (a *RealSQLCAdapter) Version(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "sqlc", "version")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("sqlc version failed: %w", err)
	}

	return string(output), nil
}

// CheckInstallation checks if sqlc is installed.
func (a *RealSQLCAdapter) CheckInstallation(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "sqlc", "--help")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlc check failed: %w", err)
	}

	return nil
}
