package schema

import (
	"encoding/json"
	"fmt"
)

// SchemaSource represents a source for JSON Schema
type SchemaSource interface {
	GetSchemaJSON() string
}

// StringSchemaSource provides a schema from a string
type StringSchemaSource struct {
	schema string
}

// NewStringSchemaSource creates a new string schema source
func NewStringSchemaSource(schema string) *StringSchemaSource {
	return &StringSchemaSource{schema: schema}
}

// GetSchemaJSON returns the schema JSON
func (s *StringSchemaSource) GetSchemaJSON() string {
	return s.schema
}

// Validator provides schema validation for Nexlayer YAML configurations
type Validator struct {
	strict       bool
	registry     *ValidationRegistry
	schemaSource SchemaSource
}

// NewValidator creates a new validator with the specified settings
func NewValidator(strict bool, schemaSource SchemaSource) *Validator {
	return &Validator{
		strict:       strict,
		registry:     NewValidationRegistry(),
		schemaSource: schemaSource,
	}
}

// NewDefaultValidator creates a new validator with default settings
func NewDefaultValidator() *Validator {
	// Create validator using the built-in JSON schema (imported from schema.go)
	schemaSource := NewStringSchemaSource(SchemaV2)
	return NewValidator(true, schemaSource)
}

// ValidateYAML performs validation of a YAML structure
func (v *Validator) ValidateYAML(config interface{}) []ValidationError {
	errors := []ValidationError{} // Initialize with empty slice, not nil

	// First, perform JSON Schema validation
	if v.schemaSource != nil {
		schemaErrors := v.validateWithJSONSchema(config)
		errors = append(errors, schemaErrors...)
	}

	// Then perform semantic validation with the registry
	ctx := &ValidationContext{Config: config}
	errors = append(errors, v.validateWithRegistry(config, ctx)...)

	return errors
}

// validateWithJSONSchema validates a configuration using JSON Schema
func (v *Validator) validateWithJSONSchema(config interface{}) []ValidationError {
	errors := []ValidationError{} // Initialize with empty slice, not nil

	// In a real implementation, we would use a JSON schema validation library
	// For now, this is just a placeholder to show how it would work
	schemaJSON := v.schemaSource.GetSchemaJSON()
	_ = schemaJSON // Use the variable to avoid linter errors

	// For demonstration, we validate that the config can be marshaled to JSON
	_, err := json.Marshal(config)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:    "schema",
			Message:  "Configuration cannot be converted to JSON: " + err.Error(),
			Severity: string(ValidationErrorSeverityError),
		})
	}

	// In a real implementation, this would validate against the schema
	// For now, we'll just return an empty list of errors
	return errors
}

// validateWithRegistry performs semantic validation using the validator registry
func (v *Validator) validateWithRegistry(config interface{}, ctx *ValidationContext) []ValidationError {
	errors := []ValidationError{} // Initialize with empty slice, not nil

	// In a real implementation, we would walk through the config structure
	// and validate each field using the appropriate validator
	// For now, we'll just return an empty list of errors

	return errors
}

// CreateError creates a validation error
func CreateError(field, message, severity string, suggestions ...string) ValidationError {
	return ValidationError{
		Field:       field,
		Message:     message,
		Severity:    severity,
		Suggestions: suggestions,
	}
}

// CreateRequiredError creates a validation error for a required field
func CreateRequiredError(field string) ValidationError {
	return CreateError(
		field,
		"field is required",
		string(ValidationErrorSeverityError),
		"Add the required field to your configuration",
	)
}

// CreateFormatError creates a validation error for incorrectly formatted values
func CreateFormatError(field, format string, examples ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("Format should be: %s", format),
	}
	for _, example := range examples {
		suggestions = append(suggestions, fmt.Sprintf("Example: %s", example))
	}
	return CreateError(
		field,
		"invalid format",
		string(ValidationErrorSeverityError),
		suggestions...,
	)
}

// CreateReferenceError creates a validation error for invalid references
func CreateReferenceError(field, ref string, available ...string) ValidationError {
	suggestions := []string{
		fmt.Sprintf("'%s' is not a valid reference", ref),
	}
	if len(available) > 0 {
		suggestions = append(suggestions, "Available options:")
		for _, option := range available {
			suggestions = append(suggestions, fmt.Sprintf("- %s", option))
		}
	}
	return CreateError(
		field,
		fmt.Sprintf("invalid reference: %s", ref),
		string(ValidationErrorSeverityError),
		suggestions...,
	)
}
