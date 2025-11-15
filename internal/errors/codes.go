package errors

import (
	"fmt"
	"strings"
)

// ErrorCode represents strongly-typed error codes
// Replaces string-based error codes with type safety
type ErrorCode string

const (
	// Migration Errors
	ErrorCodeMigrationFailed    ErrorCode = "MIGRATION_FAILED"
	ErrorCodeMigrationNotFound  ErrorCode = "MIGRATION_NOT_FOUND"
	ErrorCodeTooManyMigrations ErrorCode = "TOO_MANY_MIGRATIONS"
	
	// Schema Errors
	ErrorCodeSchemaNotFound      ErrorCode = "SCHEMA_NOT_FOUND"
	ErrorCodeSchemaValidation   ErrorCode = "SCHEMA_VALIDATION"
	ErrorCodeTableNotFound      ErrorCode = "TABLE_NOT_FOUND"
	ErrorCodeColumnNotFound     ErrorCode = "COLUMN_NOT_FOUND"
	
	// Event Errors
	ErrorCodeEventValidation     ErrorCode = "EVENT_VALIDATION"
	ErrorCodeEventNotFound       ErrorCode = "EVENT_NOT_FOUND"
	ErrorCodeInvalidEventType    ErrorCode = "INVALID_EVENT_TYPE"
	ErrorCodeEmptyAggregateID   ErrorCode = "EMPTY_AGGREGATE_ID"
	
	// Configuration Errors
	ErrorCodeConfigValidation    ErrorCode = "CONFIG_VALIDATION"
	ErrorCodeConfigNotFound     ErrorCode = "CONFIG_NOT_FOUND"
	ErrorCodeInvalidProjectType ErrorCode = "INVALID_PROJECT_TYPE"
	
	// System Errors
	ErrorCodeInternalServer     ErrorCode = "INTERNAL_SERVER"
	ErrorCodeTimeout           ErrorCode = "TIMEOUT"
	ErrorCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrorCodeNotFound          ErrorCode = "NOT_FOUND"
)

// IsValid returns true if ErrorCode is valid
func (ec ErrorCode) IsValid() bool {
	switch ec {
	case ErrorCodeMigrationFailed, ErrorCodeMigrationNotFound, ErrorCodeTooManyMigrations,
		 ErrorCodeSchemaNotFound, ErrorCodeSchemaValidation, ErrorCodeTableNotFound, ErrorCodeColumnNotFound,
		 ErrorCodeEventValidation, ErrorCodeEventNotFound, ErrorCodeInvalidEventType, ErrorCodeEmptyAggregateID,
		 ErrorCodeConfigValidation, ErrorCodeConfigNotFound, ErrorCodeInvalidProjectType,
		 ErrorCodeInternalServer, ErrorCodeTimeout, ErrorCodePermissionDenied, ErrorCodeNotFound:
		return true
	default:
		return false
	}
}

// IsMigrationError returns true if error is migration-related
func (ec ErrorCode) IsMigrationError() bool {
	migrationErrors := []ErrorCode{
		ErrorCodeMigrationFailed,
		ErrorCodeMigrationNotFound,
		ErrorCodeTooManyMigrations,
	}
	
	for _, err := range migrationErrors {
		if ec == err {
			return true
		}
	}
	return false
}

// IsSchemaError returns true if error is schema-related
func (ec ErrorCode) IsSchemaError() bool {
	schemaErrors := []ErrorCode{
		ErrorCodeSchemaNotFound,
		ErrorCodeSchemaValidation,
		ErrorCodeTableNotFound,
		ErrorCodeColumnNotFound,
	}
	
	for _, err := range schemaErrors {
		if ec == err {
			return true
		}
	}
	return false
}

// IsEventError returns true if error is event-related
func (ec ErrorCode) IsEventError() bool {
	eventErrors := []ErrorCode{
		ErrorCodeEventValidation,
		ErrorCodeEventNotFound,
		ErrorCodeInvalidEventType,
		ErrorCodeEmptyAggregateID,
	}
	
	for _, err := range eventErrors {
		if ec == err {
			return true
		}
	}
	return false
}

// IsConfigError returns true if error is configuration-related
func (ec ErrorCode) IsConfigError() bool {
	configErrors := []ErrorCode{
		ErrorCodeConfigValidation,
		ErrorCodeConfigNotFound,
		ErrorCodeInvalidProjectType,
	}
	
	for _, err := range configErrors {
		if ec == err {
			return true
		}
	}
	return false
}

// IsSystemError returns true if error is system-related
func (ec ErrorCode) IsSystemError() bool {
	systemErrors := []ErrorCode{
		ErrorCodeInternalServer,
		ErrorCodeTimeout,
		ErrorCodePermissionDenied,
		ErrorCodeNotFound,
	}
	
	for _, err := range systemErrors {
		if ec == err {
			return true
		}
	}
	return false
}

// Category returns error category for grouping
func (ec ErrorCode) Category() string {
	switch {
	case ec.IsMigrationError():
		return "migration"
	case ec.IsSchemaError():
		return "schema"
	case ec.IsEventError():
		return "event"
	case ec.IsConfigError():
		return "configuration"
	case ec.IsSystemError():
		return "system"
	default:
		return "unknown"
	}
}

// Severity returns default severity for error code
func (ec ErrorCode) Severity() ErrorSeverity {
	switch ec {
	case ErrorCodeMigrationFailed, ErrorCodeSchemaValidation, ErrorCodeEventValidation, ErrorCodeConfigValidation:
		return ErrorSeverityError
	case ErrorCodeMigrationNotFound, ErrorCodeTableNotFound, ErrorCodeColumnNotFound, ErrorCodeEventNotFound, ErrorCodeConfigNotFound, ErrorCodeNotFound:
		return ErrorSeverityWarning
	case ErrorCodeTimeout:
		return ErrorSeverityCritical
	case ErrorCodePermissionDenied:
		return ErrorSeverityError
	case ErrorCodeInternalServer, ErrorCodeTooManyMigrations:
		return ErrorSeverityCritical
	default:
		return ErrorSeverityError
	}
}

// HTTPStatus returns appropriate HTTP status code
func (ec ErrorCode) HTTPStatus() int {
	switch ec {
	case ErrorCodeMigrationNotFound, ErrorCodeTableNotFound, ErrorCodeColumnNotFound, ErrorCodeEventNotFound, ErrorCodeConfigNotFound, ErrorCodeNotFound:
		return 404
	case ErrorCodeConfigValidation, ErrorCodeSchemaValidation, ErrorCodeEventValidation, ErrorCodeInvalidProjectType:
		return 400
	case ErrorCodePermissionDenied:
		return 403
	case ErrorCodeTimeout:
		return 408
	case ErrorCodeTooManyMigrations:
		return 413 // Request Entity Too Large
	case ErrorCodeInternalServer:
		return 500
	default:
		return 500
	}
}

// LogMessage returns appropriate log message
func (ec ErrorCode) LogMessage() string {
	switch ec {
	case ErrorCodeMigrationFailed:
		return "Migration operation failed"
	case ErrorCodeMigrationNotFound:
		return "Migration not found"
	case ErrorCodeSchemaValidation:
		return "Schema validation failed"
	case ErrorCodeEventValidation:
		return "Event validation failed"
	case ErrorCodeConfigValidation:
		return "Configuration validation failed"
	case ErrorCodeTimeout:
		return "Operation timed out"
	case ErrorCodePermissionDenied:
		return "Permission denied"
	case ErrorCodeNotFound:
		return "Resource not found"
	case ErrorCodeInternalServer:
		return "Internal server error"
	default:
		return "Unknown error occurred"
	}
}

// UserMessage returns user-friendly message
func (ec ErrorCode) UserMessage() string {
	switch ec {
	case ErrorCodeMigrationFailed:
		return "Migration failed. Please check your database connection and try again."
	case ErrorCodeMigrationNotFound:
		return "Migration not found. Please verify the migration ID."
	case ErrorCodeSchemaValidation:
		return "Schema validation failed. Please check your SQL syntax."
	case ErrorCodeEventValidation:
		return "Invalid event data. Please check your configuration."
	case ErrorCodeConfigValidation:
		return "Configuration validation failed. Please check your settings."
	case ErrorCodeTimeout:
		return "Operation timed out. Please try again."
	case ErrorCodePermissionDenied:
		return "Permission denied. You don't have access to this resource."
	case ErrorCodeNotFound:
		return "Resource not found. Please check your request."
	case ErrorCodeInternalServer:
		return "Internal server error. Please try again later."
	default:
		return "An error occurred. Please try again."
	}
}

// ErrorGroup represents multiple error codes
type ErrorGroup []ErrorCode

// Contains checks if group contains error code
func (eg ErrorGroup) Contains(code ErrorCode) bool {
	for _, c := range eg {
		if c == code {
			return true
		}
	}
	return false
}

// String returns string representation of error group
func (eg ErrorGroup) String() string {
	if len(eg) == 0 {
		return "[]"
	}
	
	var codes []string
	for _, code := range eg {
		codes = append(codes, string(code))
	}
	
	return fmt.Sprintf("[%s]", strings.Join(codes, ", "))
}

// GetMigrationErrors returns all migration error codes
func GetMigrationErrors() ErrorGroup {
	return ErrorGroup{
		ErrorCodeMigrationFailed,
		ErrorCodeMigrationNotFound,
		ErrorCodeTooManyMigrations,
	}
}

// GetSchemaErrors returns all schema error codes
func GetSchemaErrors() ErrorGroup {
	return ErrorGroup{
		ErrorCodeSchemaNotFound,
		ErrorCodeSchemaValidation,
		ErrorCodeTableNotFound,
		ErrorCodeColumnNotFound,
	}
}

// GetEventErrors returns all event error codes
func GetEventErrors() ErrorGroup {
	return ErrorGroup{
		ErrorCodeEventValidation,
		ErrorCodeEventNotFound,
		ErrorCodeInvalidEventType,
		ErrorCodeEmptyAggregateID,
	}
}

// GetConfigErrors returns all configuration error codes
func GetConfigErrors() ErrorGroup {
	return ErrorGroup{
		ErrorCodeConfigValidation,
		ErrorCodeConfigNotFound,
		ErrorCodeInvalidProjectType,
	}
}

// GetSystemErrors returns all system error codes
func GetSystemErrors() ErrorGroup {
	return ErrorGroup{
		ErrorCodeInternalServer,
		ErrorCodeTimeout,
		ErrorCodePermissionDenied,
		ErrorCodeNotFound,
	}
}

// GetAllErrorCodes returns all valid error codes
func GetAllErrorCodes() ErrorGroup {
	return ErrorGroup{
		// Migration
		ErrorCodeMigrationFailed,
		ErrorCodeMigrationNotFound,
		ErrorCodeTooManyMigrations,
		
		// Schema
		ErrorCodeSchemaNotFound,
		ErrorCodeSchemaValidation,
		ErrorCodeTableNotFound,
		ErrorCodeColumnNotFound,
		
		// Event
		ErrorCodeEventValidation,
		ErrorCodeEventNotFound,
		ErrorCodeInvalidEventType,
		ErrorCodeEmptyAggregateID,
		
		// Configuration
		ErrorCodeConfigValidation,
		ErrorCodeConfigNotFound,
		ErrorCodeInvalidProjectType,
		
		// System
		ErrorCodeInternalServer,
		ErrorCodeTimeout,
		ErrorCodePermissionDenied,
		ErrorCodeNotFound,
	}
}

// ErrorSeverity represents error severity levels
type ErrorSeverity string

const (
	ErrorSeverityInfo    ErrorSeverity = "info"
	ErrorSeverityWarning ErrorSeverity = "warning"
	ErrorSeverityError   ErrorSeverity = "error"
	ErrorSeverityCritical ErrorSeverity = "critical"
)

// IsValid returns true if ErrorSeverity is valid
func (es ErrorSeverity) IsValid() bool {
	switch es {
	case ErrorSeverityInfo, ErrorSeverityWarning, ErrorSeverityError, ErrorSeverityCritical:
		return true
	default:
		return false
	}
}

// Level returns numeric level for severity comparison
func (es ErrorSeverity) Level() int {
	switch es {
	case ErrorSeverityInfo:
		return 1
	case ErrorSeverityWarning:
		return 2
	case ErrorSeverityError:
		return 3
	case ErrorSeverityCritical:
		return 4
	default:
		return 0
	}
}

// HigherThan returns true if this severity is higher than another
func (es ErrorSeverity) HigherThan(other ErrorSeverity) bool {
	return es.Level() > other.Level()
}

// LowerOrEqual returns true if this severity is lower or equal to another
func (es ErrorSeverity) LowerOrEqual(other ErrorSeverity) bool {
	return es.Level() <= other.Level()
}