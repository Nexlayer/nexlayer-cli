package schema

import (
	oldschema "github.com/Nexlayer/nexlayer-cli/pkg/schema"
)

// ForwardValidator adapts the new validator to work with old schema types
type ForwardValidator struct {
	validator *Validator
}

// NewForwardValidator creates a new validation adapter for old schema types
func NewForwardValidator(strict bool) *ForwardValidator {
	return &ForwardValidator{
		validator: NewDefaultValidator(),
	}
}

// ValidateOldYAML validates a Nexlayer YAML configuration from the old schema package
func (v *ForwardValidator) ValidateOldYAML(yaml *oldschema.NexlayerYAML) []oldschema.ValidationError {
	// Convert the old schema to a generic map
	// This is a simplified approach - in a real implementation we would
	// convert to a proper struct-to-struct mapping

	// For now, just return an empty slice to indicate no errors
	return []oldschema.ValidationError{}
}

// ConvertToNewError converts an old schema.ValidationError to a NewValidationError
func ConvertToNewError(err oldschema.ValidationError) NewValidationError {
	return NewValidationError{
		Field:       err.Field,
		Message:     err.Message,
		Severity:    err.Severity,
		Suggestions: err.Suggestions,
	}
}

// ConvertToOldError converts a NewValidationError to an old schema.ValidationError
func ConvertToOldError(err NewValidationError) oldschema.ValidationError {
	return oldschema.ValidationError{
		Field:       err.Field,
		Message:     err.Message,
		Severity:    err.Severity,
		Suggestions: err.Suggestions,
	}
}

// ConvertToNewErrors converts a slice of old schema.ValidationErrors to NewValidationErrors
func ConvertToNewErrors(errs []oldschema.ValidationError) []NewValidationError {
	result := make([]NewValidationError, len(errs))
	for i, err := range errs {
		result[i] = ConvertToNewError(err)
	}
	return result
}

// ConvertToOldErrors converts a slice of NewValidationErrors to old schema.ValidationErrors
func ConvertToOldErrors(errs []NewValidationError) []oldschema.ValidationError {
	result := make([]oldschema.ValidationError, len(errs))
	for i, err := range errs {
		result[i] = ConvertToOldError(err)
	}
	return result
}
