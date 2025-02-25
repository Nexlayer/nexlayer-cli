// Package schema provides centralized schema validation for Nexlayer YAML configurations.
package schema

import (
	"fmt"
)

// ValidationErrorSeverity indicates the severity level of a validation error
type ValidationErrorSeverity string

const (
	ValidationErrorSeverityError   ValidationErrorSeverity = "error"
	ValidationErrorSeverityWarning ValidationErrorSeverity = "warning"
)

// ValidationError represents a validation error with context and suggestions
type ValidationError struct {
	Field       string   `json:"field"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions,omitempty"`
	Severity    string   `json:"severity"` // error, warning
}

// Error implements the error interface for ValidationError
func (e ValidationError) Error() string {
	base := fmt.Sprintf("%s: %s", e.Field, e.Message)
	if len(e.Suggestions) > 0 {
		base += "\nSuggestions:"
		for _, s := range e.Suggestions {
			base += fmt.Sprintf("\n- %s", s)
		}
	}
	return base
}

// ValidationContext provides context for validation
type ValidationContext struct {
	Config interface{}
}

// ValidatorFunc is a function that validates a field value
type ValidatorFunc func(field, value string, ctx *ValidationContext) []ValidationError

// ValidationRule represents a validation rule
type ValidationRule interface {
	Validate(field string, value interface{}, ctx *ValidationContext) []ValidationError
}

// ValidationFuncAdapter adapts a ValidatorFunc to implement ValidationRule
type ValidationFuncAdapter struct {
	ValidatorFunc ValidatorFunc
}

// Validate implements the ValidationRule interface
func (a ValidationFuncAdapter) Validate(field string, value interface{}, ctx *ValidationContext) []ValidationError {
	if strValue, ok := value.(string); ok {
		return a.ValidatorFunc(field, strValue, ctx)
	}
	return []ValidationError{
		{
			Field:    field,
			Message:  "value must be a string",
			Severity: string(ValidationErrorSeverityError),
		},
	}
}
