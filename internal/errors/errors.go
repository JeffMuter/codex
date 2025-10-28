package errors

import (
	"errors"
	"fmt"
)

// Error types for different categories
var (
	// Configuration errors
	ErrConfigNotFound    = errors.New("configuration not found")
	ErrConfigInvalid     = errors.New("configuration is invalid")
	ErrConfigParseFailed = errors.New("failed to parse configuration")

	// Context gathering errors
	ErrRepositoryNotFound = errors.New("git repository not found")
	ErrNixConfigNotFound  = errors.New("nix configuration not found")
	ErrDotfilesNotFound   = errors.New("dotfiles not found")

	// Provider errors
	ErrProviderInvalid      = errors.New("provider configuration is invalid")
	ErrProviderUnavailable  = errors.New("provider is unavailable")
	ErrProviderAPIError     = errors.New("provider API error")
	ErrProviderRateLimit    = errors.New("provider rate limit exceeded")

	// Database errors
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrDatabaseQuery      = errors.New("database query failed")

	// Screenshot errors
	ErrScreenshotToolNotFound = errors.New("no screenshot tool found")
	ErrScreenshotCaptureFailed = errors.New("screenshot capture failed")
)

// CodexError wraps errors with additional context
type CodexError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *CodexError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CodexError) Unwrap() error {
	return e.Err
}

// New creates a new CodexError
func New(code, message string, err error) *CodexError {
	return &CodexError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Wrap wraps an error with additional context
func Wrap(err error, code, message string) *CodexError {
	return &CodexError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Error codes
const (
	// Configuration error codes
	CodeConfigNotFound    = "CONFIG_NOT_FOUND"
	CodeConfigInvalid     = "CONFIG_INVALID"
	CodeConfigParseFailed = "CONFIG_PARSE_FAILED"

	// Context gathering error codes
	CodeRepositoryNotFound = "REPOSITORY_NOT_FOUND"
	CodeNixConfigNotFound  = "NIX_CONFIG_NOT_FOUND"
	CodeDotfilesNotFound   = "DOTFILES_NOT_FOUND"

	// Provider error codes
	CodeProviderInvalid     = "PROVIDER_INVALID"
	CodeProviderUnavailable = "PROVIDER_UNAVAILABLE"
	CodeProviderAPIError    = "PROVIDER_API_ERROR"
	CodeProviderRateLimit   = "PROVIDER_RATE_LIMIT"

	// Database error codes
	CodeDatabaseConnection = "DATABASE_CONNECTION"
	CodeDatabaseQuery      = "DATABASE_QUERY"

	// Screenshot error codes
	CodeScreenshotToolNotFound  = "SCREENSHOT_TOOL_NOT_FOUND"
	CodeScreenshotCaptureFailed = "SCREENSHOT_CAPTURE_FAILED"
)

// IsConfigError checks if error is configuration-related
func IsConfigError(err error) bool {
	var codexErr *CodexError
	if errors.As(err, &codexErr) {
		switch codexErr.Code {
		case CodeConfigNotFound, CodeConfigInvalid, CodeConfigParseFailed:
			return true
		}
	}
	return false
}

// IsProviderError checks if error is provider-related
func IsProviderError(err error) bool {
	var codexErr *CodexError
	if errors.As(err, &codexErr) {
		switch codexErr.Code {
		case CodeProviderInvalid, CodeProviderUnavailable, CodeProviderAPIError, CodeProviderRateLimit:
			return true
		}
	}
	return false
}

// IsDatabaseError checks if error is database-related
func IsDatabaseError(err error) bool {
	var codexErr *CodexError
	if errors.As(err, &codexErr) {
		switch codexErr.Code {
		case CodeDatabaseConnection, CodeDatabaseQuery:
			return true
		}
	}
	return false
}
