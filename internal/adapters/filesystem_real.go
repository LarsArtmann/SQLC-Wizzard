package adapters

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// RealFileSystemAdapter provides actual file system operations
type RealFileSystemAdapter struct{}

// NewRealFileSystemAdapter creates a new real file system adapter
func NewRealFileSystemAdapter() *RealFileSystemAdapter {
	return &RealFileSystemAdapter{}
}

// ReadFile reads a file
func (a *RealFileSystemAdapter) ReadFile(ctx context.Context, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}

// WriteFile writes a file
func (a *RealFileSystemAdapter) WriteFile(ctx context.Context, path string, data []byte, perm fs.FileMode) error {
	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	err := os.WriteFile(path, data, perm)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

// CreateDirectory creates a directory
func (a *RealFileSystemAdapter) CreateDirectory(ctx context.Context, path string, perm fs.FileMode) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}
	return nil
}

// MkdirAll creates a directory and all parent directories (alias for CreateDirectory)
func (a *RealFileSystemAdapter) MkdirAll(ctx context.Context, path string, perm fs.FileMode) error {
	return a.CreateDirectory(ctx, path, perm)
}

// Exists checks if a path exists
func (a *RealFileSystemAdapter) Exists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ListFiles lists files in a directory
func (a *RealFileSystemAdapter) ListFiles(ctx context.Context, dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return files, nil
}

// Remove removes a file or directory
func (a *RealFileSystemAdapter) Remove(ctx context.Context, path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to remove %s: %w", path, err)
	}
	return nil
}

// Copy copies a file or directory
func (a *RealFileSystemAdapter) Copy(ctx context.Context, src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat %s: %w", src, err)
	}

	if info.IsDir() {
		return a.copyDir(ctx, src, dst, info)
	}
	return a.copyFile(ctx, src, dst, info)
}

// TempDir creates a temporary directory
func (a *RealFileSystemAdapter) TempDir(ctx context.Context, prefix string) (string, error) {
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	return dir, nil
}

// copyFile copies a single file
func (a *RealFileSystemAdapter) copyFile(ctx context.Context, src, dst string, info fs.FileInfo) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return a.WriteFile(ctx, dst, data, info.Mode())
}

// copyDir copies a directory recursively
func (a *RealFileSystemAdapter) copyDir(ctx context.Context, src, dst string, info fs.FileInfo) error {
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if err := a.Copy(ctx, srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}
