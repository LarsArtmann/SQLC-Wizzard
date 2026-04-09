package adapters

import (
	"io/fs"
)

// File permission constants used across adapters.
// These use standard Unix permission values.

const (
	// DefaultDirPermissions is the standard permission for created directories (rwxr-xr-x).
	DefaultDirPermissions fs.FileMode = 0o755

	// DefaultFilePermissions is the standard permission for created files (rw-r--r--).
	DefaultFilePermissions fs.FileMode = 0o644

	// TestFilePermissions is the permission used for test files.
	TestFilePermissions fs.FileMode = 0o644
)
