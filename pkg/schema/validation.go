// Package schema provides centralized schema management for Nexlayer YAML configurations.
package schema

import (
	"fmt"
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

// ValidateYAML performs validation of a YAML configuration
func ValidateYAML(yaml *NexlayerYAML) []ValidationError {
	// This is a placeholder that will be replaced by the new validation package
	// For now, just return an empty slice to avoid breaking existing code
	return []ValidationError{}
}
