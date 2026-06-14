package creators_test

import (
	"context"
	"io/fs"

	"github.com/LarsArtmann/SQLC-Wizzard/internal/apperrors"
)

// MockFileSystemAdapter captures file system operations for testing.
type MockFileSystemAdapter struct {
	mkdirAllCalls   []MkdirAllCall
	writeFileCalls  []WriteFileCall
	callLog         []string
	shouldFailMkdir bool
	shouldFailWrite bool
}

// MkdirAllCall records a single MkdirAll call.
type MkdirAllCall struct {
	Path string
	Perm fs.FileMode
}

// WriteFileCall records a single WriteFile call.
type WriteFileCall struct {
	Path    string
	Content []byte
	Perm    fs.FileMode
}

func (m *MockFileSystemAdapter) MkdirAll(ctx context.Context, path string, perm fs.FileMode) error {
	m.mkdirAllCalls = append(m.mkdirAllCalls, MkdirAllCall{Path: path, Perm: perm})

	m.callLog = append(m.callLog, "mkdir:"+path)
	if m.shouldFailMkdir {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "mkdir failed")
	}

	return nil
}

func (m *MockFileSystemAdapter) WriteFile(
	ctx context.Context,
	path string,
	content []byte,
	perm fs.FileMode,
) error {
	m.writeFileCalls = append(
		m.writeFileCalls,
		WriteFileCall{Path: path, Content: content, Perm: perm},
	)

	m.callLog = append(m.callLog, "write:"+path)
	if m.shouldFailWrite {
		return apperrors.NewError(apperrors.ErrorCodeInternalServer, "write failed")
	}

	return nil
}

func (m *MockFileSystemAdapter) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return nil, apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) Exists(ctx context.Context, path string) (bool, error) {
	return false, apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) Remove(ctx context.Context, path string) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) TempDir(ctx context.Context, pattern string) (string, error) {
	return "", apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) Copy(ctx context.Context, src, dst string) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

func (m *MockFileSystemAdapter) CreateDirectory(
	ctx context.Context,
	path string,
	perm fs.FileMode,
) error {
	return m.MkdirAll(ctx, path, perm)
}

func (m *MockFileSystemAdapter) ListFiles(ctx context.Context, dir string) ([]string, error) {
	return nil, apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}

// MockCLIAdapter captures CLI output for testing.
type MockCLIAdapter struct {
	printedLines []string
}

func (m *MockCLIAdapter) Println(ctx context.Context, msg string) error {
	m.printedLines = append(m.printedLines, msg)

	return nil
}

func (m *MockCLIAdapter) Printf(format string, args ...any) {
	// Not needed for current tests
}

func (m *MockCLIAdapter) CheckCommand(ctx context.Context, name string) error {
	return nil
}

func (m *MockCLIAdapter) RunCommand(
	ctx context.Context,
	name string,
	args ...string,
) (string, error) {
	return "", nil
}

func (m *MockCLIAdapter) GetVersion(ctx context.Context, command string) (string, error) {
	return "", nil
}

func (m *MockCLIAdapter) Install(ctx context.Context, cmd string) error {
	return apperrors.NewError(apperrors.ErrorCodeInternalServer, "not implemented")
}
