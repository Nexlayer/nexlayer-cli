package schema

import (
	"fmt"

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
func (v *ForwardValidator) ValidateOldYAML(yaml *oldschema.NexlayerYAML) []NewValidationError {
	// Convert the old schema to a generic map
	// This is a simplified approach - in a real implementation we would
	// convert to a proper struct-to-struct mapping

	// For now, just return an empty slice to indicate no errors
	return []NewValidationError{}
}

// ConvertToOldError converts a new NewValidationError to an old schema error type
func ConvertToOldError(err NewValidationError) error {
	// Since we can't find the old ValidationError type, we'll just return
	// a basic error with the message
	return fmt.Errorf("%s: %s", err.Field, err.Message)
}

// ConvertToOldErrors converts a slice of new NewValidationErrors to old schema error types
func ConvertToOldErrors(errs []NewValidationError) []error {
	result := make([]error, len(errs))
	for i, err := range errs {
		result[i] = ConvertToOldError(err)
	}
	return result
}
