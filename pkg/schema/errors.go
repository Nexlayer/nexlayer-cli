// Package schema provides centralized schema management for Nexlayer YAML configurations.
package schema

import (
	"encoding/json"
	"fmt"
	"strings"
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

// ValidationErrorSeverity defines the severity level of validation errors
type ValidationErrorSeverity string

const (
	// ValidationErrorSeverityError indicates a critical error
	ValidationErrorSeverityError ValidationErrorSeverity = "error"

	// ValidationErrorSeverityWarning indicates a non-critical warning
	ValidationErrorSeverityWarning ValidationErrorSeverity = "warning"
)

// ValidationErrorInfo contains additional metadata for a validation error
type ValidationErrorInfo struct {
	Category ValidationErrorCategory `json:"category,omitempty"`
}

// ValidationContextInfo provides contextual information for validation
type ValidationContextInfo struct {
	Path   string      `json:"path,omitempty"`
	Parent interface{} `json:"-"`
	Root   interface{} `json:"-"`
}

// ValidationReport represents a collection of validation errors and warnings
type ValidationReport struct {
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// HasErrors returns true if the validation report contains errors
func (r *ValidationReport) HasErrors() bool {
	return len(r.Errors) > 0
}

// HasWarnings returns true if the validation report contains warnings
func (r *ValidationReport) HasWarnings() bool {
	return len(r.Warnings) > 0
}

// IsValid returns true if the validation report contains no errors
func (r *ValidationReport) IsValid() bool {
	return !r.HasErrors()
}

// AddError adds an error to the validation report
func (r *ValidationReport) AddError(e ValidationError) {
	if e.Severity == string(ValidationErrorSeverityWarning) {
		r.Warnings = append(r.Warnings, e)
	} else {
		r.Errors = append(r.Errors, e)
	}
}

// AddErrors adds multiple errors to the validation report
func (r *ValidationReport) AddErrors(errors []ValidationError) {
	for _, e := range errors {
		r.AddError(e)
	}
}

// String returns a string representation of the validation report
func (r *ValidationReport) String() string {
	var sb strings.Builder

	if r.HasErrors() {
		sb.WriteString(fmt.Sprintf("Found %d validation errors:\n", len(r.Errors)))
		for i, e := range r.Errors {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, e.Error()))
		}
	}

	if r.HasWarnings() {
		if r.HasErrors() {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("Found %d validation warnings:\n", len(r.Warnings)))
		for i, e := range r.Warnings {
			sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, e.Error()))
		}
	}

	if !r.HasErrors() && !r.HasWarnings() {
		sb.WriteString("Validation passed with no errors or warnings.")
	}

	return sb.String()
}

// NewValidationReport creates a new validation report
func NewValidationReport() *ValidationReport {
	return &ValidationReport{
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationError, 0),
		Metadata: make(map[string]string),
	}
}

// CreateError creates a validation error with metadata
func CreateError(field, message string, severity string, suggestions ...string) ValidationError {
	return ValidationError{
		Field:       field,
		Message:     message,
		Severity:    severity,
		Suggestions: suggestions,
	}
}

// CreateRequiredError creates an error for a missing required field
func CreateRequiredError(field string) ValidationError {
	return CreateError(
		field,
		"Field is required",
		string(ValidationErrorSeverityError),
		"Add the required field to your configuration",
	)
}

// CreateFormatError creates an error for incorrectly formatted values
func CreateFormatError(field, format string, examples ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("Field must follow the format: %s", format),
	}
	for _, example := range examples {
		suggestions = append(suggestions, fmt.Sprintf("Example: %s", example))
	}
	return CreateError(
		field,
		"Invalid format",
		string(ValidationErrorSeverityError),
		suggestions...,
	)
}

// CreateReferenceError creates an error for invalid references
func CreateReferenceError(field, reference string, availableReferences ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("'%s' is not a valid reference", reference),
	}
	if len(availableReferences) > 0 {
		suggestions = append(suggestions, "Available references:")
		for _, ref := range availableReferences {
			suggestions = append(suggestions, fmt.Sprintf("- %s", ref))
		}
	}
	return CreateError(
		field,
		fmt.Sprintf("Invalid reference: %s", reference),
		string(ValidationErrorSeverityError),
		suggestions...,
	)
}

// CreateUnsupportedError creates an error for unsupported values
func CreateUnsupportedError(field, value string, supportedValues ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("'%s' is not supported", value),
	}
	if len(supportedValues) > 0 {
		suggestions = append(suggestions, "Supported values:")
		for _, val := range supportedValues {
			suggestions = append(suggestions, fmt.Sprintf("- %s", val))
		}
	}
	return CreateError(
		field,
		fmt.Sprintf("Unsupported value: %s", value),
		string(ValidationErrorSeverityError),
		suggestions...,
	)
}

// MarshalJSON marshals the ValidationReport to JSON
func (r *ValidationReport) MarshalJSON() ([]byte, error) {
	type Alias ValidationReport
	return json.Marshal(&struct {
		Valid bool `json:"valid"`
		*Alias
	}{
		Valid: r.IsValid(),
		Alias: (*Alias)(r),
	})
}
