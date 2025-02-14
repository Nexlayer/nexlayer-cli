// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package errors

import (
	"fmt"
	"runtime"
)

// ErrorType represents the category of error
type ErrorType int

const (
	// ErrorTypeUser represents user-caused errors (invalid input, etc.)
	ErrorTypeUser ErrorType = iota
	// ErrorTypeSystem represents system-level errors (IO, permissions, etc.)
	ErrorTypeSystem
	// ErrorTypeNetwork represents network-related errors
	ErrorTypeNetwork
	// ErrorTypeInternal represents internal application errors
	ErrorTypeInternal
)

// String returns the string representation of the error type
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeUser:
		return "user"
	case ErrorTypeSystem:
		return "system"
	case ErrorTypeNetwork:
		return "network"
	case ErrorTypeInternal:
		return "internal"
	default:
		return "unknown"
	}
}

// Error represents a structured error with context
type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Cause   error     `json:"cause,omitempty"`
	File    string    `json:"file,omitempty"`
	Line    int       `json:"line,omitempty"`
	Stack   []string  `json:"stack,omitempty"`
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// NewError creates a new error with context
func NewError(errType ErrorType, message string, cause error) *Error {
	_, file, line, _ := runtime.Caller(1)
	stack := make([]string, 0)

	// Capture stack trace
	for i := 1; i < 5; i++ { // Limit to 4 levels
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}

	return &Error{
		Type:    errType,
		Message: message,
		Cause:   cause,
		File:    file,
		Line:    line,
		Stack:   stack,
	}
}

// UserError creates a new user error
func UserError(message string, cause error) *Error {
	return NewError(ErrorTypeUser, message, cause)
}

// SystemError creates a new system error
func SystemError(message string, cause error) *Error {
	return NewError(ErrorTypeSystem, message, cause)
}

// NetworkError creates a new network error
func NetworkError(message string, cause error) *Error {
	return NewError(ErrorTypeNetwork, message, cause)
}

// InternalError creates a new internal error
func InternalError(message string, cause error) *Error {
	return NewError(ErrorTypeInternal, message, cause)
}
