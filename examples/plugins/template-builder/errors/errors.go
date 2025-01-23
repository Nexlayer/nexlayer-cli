package errors

import (
	"fmt"
	"strings"
)

// ErrorCode represents a unique error code for each type of error
type ErrorCode string

const (
	// Template related errors
	ErrTemplateNotFound    ErrorCode = "TEMPLATE_NOT_FOUND"
	ErrTemplateInvalid     ErrorCode = "TEMPLATE_INVALID"
	ErrTemplateGeneration  ErrorCode = "TEMPLATE_GENERATION_FAILED"
	
	// Project related errors
	ErrProjectNotFound     ErrorCode = "PROJECT_NOT_FOUND"
	ErrProjectInvalid      ErrorCode = "PROJECT_INVALID"
	
	// Configuration related errors
	ErrConfigInvalid       ErrorCode = "CONFIG_INVALID"
	ErrConfigNotFound      ErrorCode = "CONFIG_NOT_FOUND"
	
	// Registry related errors
	ErrRegistryUnavailable ErrorCode = "REGISTRY_UNAVAILABLE"
	ErrRegistryAuth        ErrorCode = "REGISTRY_AUTH_FAILED"
)

// CLIError represents a structured error with context
type CLIError struct {
	Code    ErrorCode
	Message string
	Err     error
	Context map[string]interface{}
}

func (e *CLIError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[%s] %s", e.Code, e.Message))
	
	if e.Err != nil {
		sb.WriteString(fmt.Sprintf(": %v", e.Err))
	}
	
	if len(e.Context) > 0 {
		sb.WriteString("\nContext:")
		for k, v := range e.Context {
			sb.WriteString(fmt.Sprintf("\n  %s: %v", k, v))
		}
	}
	
	return sb.String()
}

// NewError creates a new CLIError
func NewError(code ErrorCode, message string, err error) *CLIError {
	return &CLIError{
		Code:    code,
		Message: message,
		Err:     err,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context to the error
func (e *CLIError) WithContext(key string, value interface{}) *CLIError {
	e.Context[key] = value
	return e
}
