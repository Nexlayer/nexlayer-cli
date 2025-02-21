// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import "fmt"

// ErrorType represents the type of detection error
type ErrorType string

const (
	// ErrorTypeNotFound indicates the requested resource was not found
	ErrorTypeNotFound ErrorType = "NotFound"
	// ErrorTypeInvalid indicates invalid input or configuration
	ErrorTypeInvalid ErrorType = "Invalid"
	// ErrorTypeUnsupported indicates an unsupported operation or type
	ErrorTypeUnsupported ErrorType = "Unsupported"
	// ErrorTypeInternal indicates an internal error
	ErrorTypeInternal ErrorType = "Internal"
)

// DetectionError represents an error that occurred during detection
type DetectionError struct {
	Type    ErrorType
	Message string
	Cause   error
}

// Error implements the error interface
func (e *DetectionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// NewDetectionError creates a new DetectionError
func NewDetectionError(errType ErrorType, message string, cause error) error {
	return &DetectionError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}

// IsNotFound returns true if the error is a NotFound error
func IsNotFound(err error) bool {
	if detErr, ok := err.(*DetectionError); ok {
		return detErr.Type == ErrorTypeNotFound
	}
	return false
}

// IsInvalid returns true if the error is an Invalid error
func IsInvalid(err error) bool {
	if detErr, ok := err.(*DetectionError); ok {
		return detErr.Type == ErrorTypeInvalid
	}
	return false
}

// IsUnsupported returns true if the error is an Unsupported error
func IsUnsupported(err error) bool {
	if detErr, ok := err.(*DetectionError); ok {
		return detErr.Type == ErrorTypeUnsupported
	}
	return false
}

// IsInternal returns true if the error is an Internal error
func IsInternal(err error) bool {
	if detErr, ok := err.(*DetectionError); ok {
		return detErr.Type == ErrorTypeInternal
	}
	return false
}
