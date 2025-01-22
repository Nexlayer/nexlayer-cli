// Package errors defines an error type for handling deployment errors.
package errors

import (
	"fmt"
	"strings"
)

// DeploymentError represents an error that occurred during deployment
type DeploymentError struct {
	Message     string
	Cause       error
	Suggestions []string
}

// NewDeploymentError creates a new deployment error
func NewDeploymentError(message string, cause error, suggestions ...string) *DeploymentError {
	return &DeploymentError{
		Message:     message,
		Cause:       cause,
		Suggestions: suggestions,
	}
}

// Error returns the error message
func (e *DeploymentError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	if e.Cause != nil {
		b.WriteString(fmt.Sprintf(": %v", e.Cause))
	}
	if len(e.Suggestions) > 0 {
		b.WriteString("\n\nSuggestions:")
		for _, s := range e.Suggestions {
			b.WriteString(fmt.Sprintf("\n- %s", s))
		}
	}
	return b.String()
}

// Unwrap returns the underlying error
func (e *DeploymentError) Unwrap() error {
	return e.Cause
}

// Is reports whether the target matches this error
func (e *DeploymentError) Is(target error) bool {
	t, ok := target.(*DeploymentError)
	if !ok {
		return false
	}
	return e.Message == t.Message
}
