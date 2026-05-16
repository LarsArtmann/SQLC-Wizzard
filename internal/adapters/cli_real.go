// Package adapters provides implementations for external system interfaces.
package adapters

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
)

// RealCLIAdapter provides actual CLI operations.
type RealCLIAdapter struct{}

// NewRealCLIAdapter creates a new real CLI adapter.
func NewRealCLIAdapter() *RealCLIAdapter {
	return &RealCLIAdapter{}
}

// RunCommand executes a CLI command.
func (a *RealCLIAdapter) RunCommand(
	ctx context.Context,
	cmd string,
	args ...string,
) (string, error) {
	command := exec.CommandContext(ctx, cmd, args...)

	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to run command %s with args %v: %w", cmd, args, err)
	}

	return string(output), nil
}

// CheckCommand checks if a command is available.
func (a *RealCLIAdapter) CheckCommand(_ context.Context, cmd string) error {
	_, err := exec.LookPath(cmd)

	return err //nolint:wrapcheck // external package error
}

// GetVersion returns version of a CLI tool.
func (a *RealCLIAdapter) GetVersion(ctx context.Context, cmd string) (string, error) {
	versionFlags := []string{"--version", "-v", "version"}

	for _, flag := range versionFlags {
		output, err := exec.CommandContext(ctx, cmd, flag).Output()
		if err == nil {
			return string(output), nil
		}
	}

	return "", fmt.Errorf(
		"could not determine version for %s: %w",
		cmd,
		apperrors.ErrCLIVersionUnknown,
	)
}

// Install installs a CLI tool.
func (a *RealCLIAdapter) Install(_ context.Context, cmd string) error {
	return fmt.Errorf(
		"auto-install not implemented for %s: %w",
		cmd,
		apperrors.ErrCLINotInstallable,
	)
}

// Println prints a message to output.
func (a *RealCLIAdapter) Println(_ context.Context, message string) error {
	fmt.Println(message) //nolint:forbidigo // User-facing output requires direct print

	return nil
}
