// Package schema provides centralized schema management for Nexlayer YAML configurations.
package schema

import (
	"fmt"
	"regexp"
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

// Validator is a deprecated validator for backwards compatibility
type Validator struct {
	strict bool
}

// NewValidator creates a new validator instance
// Deprecated: Use the new validation package instead
func NewValidator(strict bool) *Validator {
	return &Validator{
		strict: strict,
	}
}

// ValidateYAML performs validation of a YAML configuration
// This is a simplified implementation for backwards compatibility
func (v *Validator) ValidateYAML(yaml *NexlayerYAML) []ValidationError {
	var errors []ValidationError

	// Basic validation for demonstration
	if yaml == nil {
		errors = append(errors, ValidationError{
			Field:       "yaml",
			Message:     "YAML configuration cannot be nil",
			Severity:    "error",
			Suggestions: []string{"Provide a valid YAML configuration"},
		})
		return errors
	}

	// Application validation
	if yaml.Application.Name == "" {
		errors = append(errors, ValidationError{
			Field:       "application.name",
			Message:     "Application name is required",
			Severity:    "error",
			Suggestions: []string{"Provide a valid application name"},
		})
	}

	// URL validation if present
	if yaml.Application.URL != "" {
		urlPattern := regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)
		if !urlPattern.MatchString(yaml.Application.URL) {
			errors = append(errors, ValidationError{
				Field:    "application.url",
				Message:  "Invalid domain format",
				Severity: "error",
				Suggestions: []string{
					"Domain must contain only lowercase letters, numbers, and hyphens",
					"Domain must start and end with a letter or number",
					"Domain parts must be separated by periods",
					"Example: my-app.example.com",
				},
			})
		}
	}

	return errors
}

// ValidateYAML performs validation of a YAML configuration
// This is a simplified implementation for backwards compatibility
func ValidateYAML(yaml *NexlayerYAML) []ValidationError {
	validator := NewValidator(true)
	return validator.ValidateYAML(yaml)
}
