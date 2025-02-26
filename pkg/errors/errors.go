// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
// Package errors provides standardized error types and utilities for the Nexlayer CLI.
package errors

import (
	"fmt"
	"strings"
)

// ErrorType represents the type of error
type ErrorType string

const (
	// Common error types
	ErrorTypeInput      ErrorType = "input_error"
	ErrorTypeValidation ErrorType = "validation_error"
	ErrorTypeNetwork    ErrorType = "network_error"
	ErrorTypePermission ErrorType = "permission_error"
	ErrorTypeConfig     ErrorType = "config_error"
	ErrorTypeInternal   ErrorType = "internal_error"
	ErrorTypeUnknown    ErrorType = "unknown_error"

	// Specific error types
	ErrorTypeInvalidPort   ErrorType = "invalid_port"
	ErrorTypeMissingImage  ErrorType = "missing_image"
	ErrorTypeInvalidVolume ErrorType = "invalid_volume"
	ErrorTypeInvalidName   ErrorType = "invalid_name"
	ErrorTypeUnsupportedOS ErrorType = "unsupported_os"
)

// Error represents a structured error with context
type Error struct {
	Type    ErrorType              `json:"type"`
	Message string                 `json:"message"`
	Code    string                 `json:"code,omitempty"`
	Field   string                 `json:"field,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
	Cause   error                  `json:"-"`
}

// Error implements the error interface
func (e *Error) Error() string {
	msg := e.Message
	if e.Field != "" {
		msg = fmt.Sprintf("%s: %s", e.Field, msg)
	}
	if e.Code != "" {
		msg = fmt.Sprintf("[%s] %s", e.Code, msg)
	}
	if e.Cause != nil {
		msg = fmt.Sprintf("%s: %s", msg, e.Cause.Error())
	}
	return msg
}

// Unwrap returns the underlying cause
func (e *Error) Unwrap() error {
	return e.Cause
}

// Is checks if this error is of the given type
func (e *Error) Is(target error) bool {
	if t, ok := target.(*Error); ok {
		return e.Type == t.Type
	}
	return false
}

// New creates a new error with the given type and message
func New(errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}

// Newf creates a new error with the given type and formatted message
func Newf(errType ErrorType, format string, args ...interface{}) *Error {
	return &Error{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
	}
}

// Wrap wraps an existing error with a new error type and message
func Wrap(err error, errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Cause:   err,
	}
}

// Wrapf wraps an existing error with a new error type and formatted message
func Wrapf(err error, errType ErrorType, format string, args ...interface{}) *Error {
	return &Error{
		Type:    errType,
		Message: fmt.Sprintf(format, args...),
		Cause:   err,
	}
}

// WithField adds a field name to the error
func (e *Error) WithField(field string) *Error {
	e.Field = field
	return e
}

// WithCode adds an error code to the error
func (e *Error) WithCode(code string) *Error {
	e.Code = code
	return e
}

// WithDetails adds additional context to the error
func (e *Error) WithDetails(details map[string]interface{}) *Error {
	e.Details = details
	return e
}

// IsErrorType checks if an error is of a specific type
func IsErrorType(err error, errType ErrorType) bool {
	var e *Error
	if ok := As(err, &e); ok {
		return e.Type == errType
	}
	return false
}

// As is a wrapper around errors.As to handle unwrapping properly
func As(err error, target interface{}) bool {
	if err == nil {
		return false
	}

	if target == nil {
		panic("errors: target cannot be nil")
	}

	val, ok := target.(*Error)
	if !ok {
		panic("errors: target must be *Error")
	}

	if e, ok := err.(*Error); ok {
		*val = *e
		return true
	}

	// Try to unwrap standard errors
	if errWithCause, ok := err.(interface{ Unwrap() error }); ok {
		return As(errWithCause.Unwrap(), target)
	}

	return false
}

// ErrorTypeFromString converts a string to an ErrorType
func ErrorTypeFromString(s string) ErrorType {
	switch strings.ToLower(s) {
	case "input_error", "input":
		return ErrorTypeInput
	case "validation_error", "validation":
		return ErrorTypeValidation
	case "network_error", "network":
		return ErrorTypeNetwork
	case "permission_error", "permission":
		return ErrorTypePermission
	case "config_error", "config":
		return ErrorTypeConfig
	case "internal_error", "internal":
		return ErrorTypeInternal
	case "invalid_port":
		return ErrorTypeInvalidPort
	case "missing_image":
		return ErrorTypeMissingImage
	case "invalid_volume":
		return ErrorTypeInvalidVolume
	case "invalid_name":
		return ErrorTypeInvalidName
	case "unsupported_os":
		return ErrorTypeUnsupportedOS
	default:
		return ErrorTypeUnknown
	}
}

// Convert common error types to standardized errors

// NewInvalidPortError creates a new invalid port error
func NewInvalidPortError(port interface{}, message string) *Error {
	if message == "" {
		message = fmt.Sprintf("Invalid port: %v", port)
	}

	details := map[string]interface{}{
		"port": port,
	}

	return New(ErrorTypeInvalidPort, message).WithDetails(details)
}

// NewMissingImageError creates a new missing image error
func NewMissingImageError(podName string) *Error {
	return Newf(ErrorTypeMissingImage, "Missing image for pod '%s'", podName).WithField("pod.image")
}

// NewInvalidVolumeError creates a new invalid volume error
func NewInvalidVolumeError(volume, reason string) *Error {
	return Newf(ErrorTypeInvalidVolume, "Invalid volume '%s': %s", volume, reason)
}

// NewInvalidNameError creates a new invalid name error
func NewInvalidNameError(name, field string) *Error {
	return Newf(ErrorTypeInvalidName, "Invalid name '%s'", name).WithField(field)
}

// NewUnsupportedOSError creates a new unsupported OS error
func NewUnsupportedOSError(os string) *Error {
	return Newf(ErrorTypeUnsupportedOS, "Unsupported operating system: %s", os)
}

// UserError creates a new user error
func UserError(message string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeInput,
		Message: message,
		Cause:   cause,
	}
}

// NetworkError creates a new network error
func NetworkError(message string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeNetwork,
		Message: message,
		Cause:   cause,
	}
}

// SystemError creates a new system error
func SystemError(message string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeUnknown,
		Message: message,
		Cause:   cause,
	}
}

// InternalError creates a new internal error
func InternalError(message string, cause error) *Error {
	return &Error{
		Type:    ErrorTypeInternal,
		Message: message,
		Cause:   cause,
	}
}
