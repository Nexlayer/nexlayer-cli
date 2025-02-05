// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
// Package errors defines an error type for handling deployment errors.
package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ValidationContext contains structured information about a validation error
type ValidationContext struct {
	Field           string   `json:"field,omitempty"`            // The field that failed validation
	ExpectedType    string   `json:"expected_type,omitempty"`    // Expected type/format
	ActualValue     string   `json:"actual_value,omitempty"`     // Actual value received
	MissingVar      string   `json:"missing_var,omitempty"`      // Name of missing environment variable
	AllowedValues   []string `json:"allowed_values,omitempty"`   // List of allowed values
	ResolutionHints []string `json:"resolution_hints,omitempty"` // Hints for fixing the error
	Example         string   `json:"example,omitempty"`          // Example of correct usage
}

// DeploymentError represents an error that occurred during deployment
type DeploymentError struct {
	Message   string
	Cause     error
	ErrorType string
	Context   *ValidationContext
}

// NewDeploymentError creates a new deployment error
func NewDeploymentError(message string, cause error) *DeploymentError {
	return &DeploymentError{
		Message:   message,
		Cause:     cause,
		ErrorType: "DeploymentError",
	}
}

// NewValidationError creates a new validation error with structured context
func NewValidationError(message string, context *ValidationContext) *DeploymentError {
	return &DeploymentError{
		Message:   message,
		ErrorType: "ValidationError",
		Context:   context,
	}
}

// Error returns the error message
func (e *DeploymentError) Error() string {
	var b strings.Builder
	b.WriteString(e.Message)
	if e.Cause != nil {
		b.WriteString(fmt.Sprintf(": %v", e.Cause))
	}
	if e.Context != nil && len(e.Context.ResolutionHints) > 0 {
		b.WriteString("\n\nSuggestions:")
		for _, s := range e.Context.ResolutionHints {
			b.WriteString(fmt.Sprintf("\n- %s", s))
		}
		if e.Context.Example != "" {
			b.WriteString(fmt.Sprintf("\n\nExample:\n%s", e.Context.Example))
		}
	}
	return b.String()
}

// MarshalJSON implements json.Marshaler interface
func (e *DeploymentError) MarshalJSON() ([]byte, error) {
	errorContext := map[string]interface{}{
		"type":    e.ErrorType,
		"message": e.Message,
	}

	if e.Context != nil {
		errorContext["validation_context"] = e.Context
	}

	if e.Cause != nil {
		errorContext["cause"] = e.Cause.Error()
	}

	return json.Marshal(map[string]interface{}{
		"error_context": errorContext,
	})
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
	return t.Message == e.Message
}
