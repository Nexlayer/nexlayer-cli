// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package schema provides centralized schema management for Nexlayer YAML configurations.
package schema

import (
	"fmt"
)

// ValidationErrorCategory defines the category of validation errors
type ValidationErrorCategory string

const (
	// ValidationErrorCategoryRequired indicates a required field is missing
	ValidationErrorCategoryRequired ValidationErrorCategory = "required"

	// ValidationErrorCategoryFormat indicates an invalid format
	ValidationErrorCategoryFormat ValidationErrorCategory = "format"

	// ValidationErrorCategoryReference indicates a reference to a non-existent resource
	ValidationErrorCategoryReference ValidationErrorCategory = "reference"

	// ValidationErrorCategoryConflict indicates a conflict between fields
	ValidationErrorCategoryConflict ValidationErrorCategory = "conflict"

	// ValidationErrorCategoryUnsupported indicates an unsupported value
	ValidationErrorCategoryUnsupported ValidationErrorCategory = "unsupported"
)

// ValidationErrorSeverity represents the severity level of a validation error
type ValidationErrorSeverity string

const (
	// ValidationErrorSeverityError indicates a critical error
	ValidationErrorSeverityError ValidationErrorSeverity = "error"

	// ValidationErrorSeverityWarning indicates a non-critical warning
	ValidationErrorSeverityWarning ValidationErrorSeverity = "warning"

	// ValidationErrorSeverityInfo indicates additional information
	ValidationErrorSeverityInfo ValidationErrorSeverity = "info"
)

// ValidationError represents a validation error with severity and suggestions
type ValidationError struct {
	Field       string                  `json:"field"`
	Message     string                  `json:"message"`
	Suggestions []string                `json:"suggestions,omitempty"`
	Severity    ValidationErrorSeverity `json:"severity"`
	Info        *ValidationErrorInfo    `json:"info,omitempty"`
	AutoFixed   bool                    `json:"auto_fixed,omitempty"`
}

// ValidationErrorInfo provides additional context for validation errors
type ValidationErrorInfo struct {
	Category ValidationErrorCategory `json:"category"`
	Details  map[string]interface{}  `json:"details,omitempty"`
}

// Error implements the error interface
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

// Helper functions for creating validation errors

// makeValidationError creates a new validation error
func makeValidationError(field, message string, severity ValidationErrorSeverity, suggestions ...string) ValidationError {
	return ValidationError{
		Field:       field,
		Message:     message,
		Severity:    severity,
		Suggestions: suggestions,
	}
}

// MakeRequiredError creates a validation error for a required field
func MakeRequiredError(field string) ValidationError {
	return makeValidationError(
		field,
		"field is required",
		ValidationErrorSeverityError,
		"Add the required field to your configuration",
	)
}

// MakeFormatError creates a validation error for incorrectly formatted values
func MakeFormatError(field, format string, examples ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("Format should be: %s", format),
	}
	for _, example := range examples {
		suggestions = append(suggestions, fmt.Sprintf("Example: %s", example))
	}
	return makeValidationError(
		field,
		"invalid format",
		ValidationErrorSeverityError,
		suggestions...,
	)
}

// MakeReferenceError creates a validation error for invalid references
func MakeReferenceError(field, ref string, available ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("'%s' is not a valid reference", ref),
	}
	if len(available) > 0 {
		suggestions = append(suggestions, "Available options:")
		for _, option := range available {
			suggestions = append(suggestions, fmt.Sprintf("- %s", option))
		}
	}
	return makeValidationError(
		field,
		fmt.Sprintf("invalid reference: %s", ref),
		ValidationErrorSeverityError,
		suggestions...,
	)
}

// Validate is a helper function to validate a configuration using the default validator
func Validate(config interface{}) []ValidationError {
	validator := NewDefaultValidator()
	return validator.ValidateYAML(config)
}
