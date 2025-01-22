// Formatted with gofmt -s
package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType int

const (
	ValidationError ErrorType = iota
	ConfigError
	DeploymentError
	NetworkError
	AuthenticationError
)

// CLIError represents a structured error with context
type CLIError struct {
	Type    ErrorType
	Message string
	Cause   error
	Hints   []string
}

func (e *CLIError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("❌ %s\n", e.Message))

	if e.Cause != nil {
		sb.WriteString(fmt.Sprintf("   Cause: %v\n", e.Cause))
	}

	if len(e.Hints) > 0 {
		sb.WriteString("\n💡 Quick fixes:\n")
		for _, hint := range e.Hints {
			sb.WriteString(fmt.Sprintf("   • %s\n", hint))
		}
	}

	return sb.String()
}

// NewValidationError creates a new validation error with helpful hints
func NewValidationError(msg string, cause error, hints ...string) *CLIError {
	return &CLIError{
		Type:    ValidationError,
		Message: msg,
		Cause:   cause,
		Hints:   hints,
	}
}

// NewConfigError creates a new configuration error with helpful hints
func NewConfigError(msg string, cause error, hints ...string) *CLIError {
	return &CLIError{
		Type:    ConfigError,
		Message: msg,
		Cause:   cause,
		Hints:   hints,
	}
}

// NewDeploymentError creates a new deployment error with helpful hints
func NewDeploymentError(msg string, cause error, hints ...string) *CLIError {
	return &CLIError{
		Type:    DeploymentError,
		Message: msg,
		Cause:   cause,
		Hints:   hints,
	}
}

// NewAuthError creates a new authentication error with helpful hints
func NewAuthError(msg string, cause error, hints ...string) *CLIError {
	return &CLIError{
		Type:    AuthenticationError,
		Message: msg,
		Cause:   cause,
		Hints: append([]string{
			"Make sure NEXLAYER_AUTH_TOKEN is set in your environment",
			"Run 'nexlayer auth login' to authenticate",
		}, hints...),
	}
}
