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

// runSQLCCommand executes a sqlc subcommand in the current working directory and
// wraps any failure with the supplied action label (e.g. "generate", "validate").
func (a *RealSQLCAdapter) runSQLCCommand(ctx context.Context, subcommand, action string) error {
	cmd := exec.CommandContext(ctx, "sqlc", subcommand)
	cmd.Dir = filepath.Dir(".")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("sqlc %s failed: %w", action, err)
	}

	return nil
}

// Generate generates Go code from SQL files.
func (a *RealSQLCAdapter) Generate(ctx context.Context, cfg *config.SqlcConfig) error {
	return a.runSQLCCommand(ctx, "generate", "generate")
}

// Validate validates sqlc configuration.
func (a *RealSQLCAdapter) Validate(ctx context.Context, cfg *config.SqlcConfig) error {
	return a.runSQLCCommand(ctx, "validate", "validate")
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

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("sqlc check failed: %w", err)
	}

	return nil
}
