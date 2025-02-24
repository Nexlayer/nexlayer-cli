// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
// Package errors provides centralized error handling for the Nexlayer CLI.
package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Kind represents the type of error
type Kind string

const (
	// Error types
	KindValidation Kind = "validation"
	KindConfig     Kind = "config"
	KindAPI        Kind = "api"
	KindRuntime    Kind = "runtime"
	KindSystem     Kind = "system"
	KindCommand    Kind = "command"
)

// NexError represents a structured error with context
type NexError struct {
	Type        Kind           `json:"type"`
	Message     string         `json:"message"`
	Field       string         `json:"field,omitempty"`
	Suggestions []string       `json:"suggestions,omitempty"`
	Context     map[string]any `json:"context,omitempty"`
	Cause       error          `json:"cause,omitempty"`
}

// Error implements the error interface
func (e *NexError) Error() string {
	var msg strings.Builder
	msg.WriteString(fmt.Sprintf("[%s] %s", e.Type, e.Message))

	if e.Field != "" {
		msg.WriteString(fmt.Sprintf(" (field: %s)", e.Field))
	}

	if len(e.Suggestions) > 0 {
		msg.WriteString("\nSuggestions:")
		for _, s := range e.Suggestions {
			msg.WriteString(fmt.Sprintf("\n  • %s", s))
		}
	}

	if e.Cause != nil {
		msg.WriteString(fmt.Sprintf("\nCaused by: %v", e.Cause))
	}

	return msg.String()
}

// MarshalJSON implements json.Marshaler
func (e *NexError) MarshalJSON() ([]byte, error) {
	type Alias NexError
	return json.Marshal(&struct {
		*Alias
		Cause string `json:"cause,omitempty"`
	}{
		Alias: (*Alias)(e),
		Cause: e.Cause.Error(),
	})
}

// ValidationError creates a new validation error
func ValidationError(message string, field string, suggestions ...string) *NexError {
	return &NexError{
		Type:        KindValidation,
		Message:     message,
		Field:       field,
		Suggestions: suggestions,
	}
}

// ConfigError creates a new configuration error
func ConfigError(message string, cause error, suggestions ...string) *NexError {
	return &NexError{
		Type:        KindConfig,
		Message:     message,
		Cause:       cause,
		Suggestions: suggestions,
	}
}

// APIError creates a new API error
func APIError(message string, cause error, context map[string]any) *NexError {
	return &NexError{
		Type:    KindAPI,
		Message: message,
		Cause:   cause,
		Context: context,
	}
}

// RuntimeError creates a new runtime error
func RuntimeError(message string, cause error) *NexError {
	return &NexError{
		Type:    KindRuntime,
		Message: message,
		Cause:   cause,
	}
}

// NewSystemError creates a new system error
func NewSystemError(message string, cause error, suggestions ...string) *NexError {
	return &NexError{
		Type:        KindSystem,
		Message:     message,
		Cause:       cause,
		Suggestions: suggestions,
	}
}

// CommandError creates a new command error
func CommandError(message string, suggestions ...string) *NexError {
	return &NexError{
		Type:        KindCommand,
		Message:     message,
		Suggestions: suggestions,
	}
}

// WithContext adds context to an error
func (e *NexError) WithContext(context map[string]any) *NexError {
	e.Context = context
	return e
}

// WithSuggestions adds suggestions to an error
func (e *NexError) WithSuggestions(suggestions ...string) *NexError {
	e.Suggestions = append(e.Suggestions, suggestions...)
	return e
}

// WithCause adds a cause to an error
func (e *NexError) WithCause(cause error) *NexError {
	e.Cause = cause
	return e
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if e, ok := err.(*NexError); ok {
		return e.Type == KindValidation
	}
	return false
}

// IsConfigError checks if an error is a configuration error
func IsConfigError(err error) bool {
	if e, ok := err.(*NexError); ok {
		return e.Type == KindConfig
	}
	return false
}

// IsAPIError checks if an error is an API error
func IsAPIError(err error) bool {
	if e, ok := err.(*NexError); ok {
		return e.Type == KindAPI
	}
	return false
}

// FormatError formats an error for display
func FormatError(err error) string {
	if e, ok := err.(*NexError); ok {
		return e.Error()
	}
	return err.Error()
}
