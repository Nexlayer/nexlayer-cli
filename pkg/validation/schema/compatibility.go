package schema

import (
	oldschema "github.com/Nexlayer/nexlayer-cli/pkg/schema"
)

// ForwardValidator adapts the new validator to work with old schema types
type ForwardValidator struct {
	validator *Validator
	strict    bool
}

// NewForwardValidator creates a new validation adapter for old schema types
func NewForwardValidator(strict bool) *ForwardValidator {
	return &ForwardValidator{
		validator: NewDefaultValidator(),
		strict:    strict,
	}
}

// ValidateOldYAML validates a Nexlayer YAML configuration from the old schema package
func (v *ForwardValidator) ValidateOldYAML(yaml *oldschema.NexlayerYAML) []oldschema.ValidationError {
	// Since we now have a direct implementation in the schema package,
	// we just delegate to it for simplicity and to avoid import cycles

	// This approach allows us to fully implement validation in the schema package
	// while still maintaining the new validation structure in pkg/validation/schema
	validator := oldschema.NewValidator(v.strict)
	return validator.ValidateYAML(yaml)
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
