package schema

// ForwardValidator adapts the new validator to work with old schema types
type ForwardValidator struct {
	validator *Validator
	strict    bool
}

// NexlayerOldYAML represents the old schema structure to avoid direct imports
type NexlayerOldYAML struct {
	// Simplified structure just for compatibility
	Application struct {
		Name string
		Pods []interface{}
	}
}

// OldValidationError represents the old validation error structure
type OldValidationError struct {
	Field       string
	Message     string
	Severity    string
	Suggestions []string
}

// NewForwardValidator creates a new validation adapter for old schema types
func NewForwardValidator(strict bool) *ForwardValidator {
	return &ForwardValidator{
		validator: NewDefaultValidator(),
		strict:    strict,
	}
}

// ValidateOldYAML validates a Nexlayer YAML configuration from the old schema package
func (v *ForwardValidator) ValidateOldYAML(yaml *NexlayerOldYAML) []OldValidationError {
	// This implementation would normally delegate to old schema validator
	// Since we can't import it directly, return empty result for now
	// Actual implementation will be fixed when addressing import structure
	return []OldValidationError{}
}

// ConvertToNewError converts an old schema.ValidationError to a NewValidationError
func ConvertToNewError(err OldValidationError) NewValidationError {
	return NewValidationError{
		Field:       err.Field,
		Message:     err.Message,
		Severity:    err.Severity,
		Suggestions: err.Suggestions,
	}
}

// ConvertToOldError converts a NewValidationError to an old schema.ValidationError
func ConvertToOldError(err NewValidationError) OldValidationError {
	return OldValidationError{
		Field:       err.Field,
		Message:     err.Message,
		Severity:    err.Severity,
		Suggestions: err.Suggestions,
	}
}

// ConvertToNewErrors converts a slice of old schema.ValidationErrors to NewValidationErrors
func ConvertToNewErrors(errs []OldValidationError) []NewValidationError {
	result := make([]NewValidationError, len(errs))
	for i, err := range errs {
		result[i] = ConvertToNewError(err)
	}
	return result
}

// ConvertToOldErrors converts a slice of NewValidationErrors to old schema.ValidationErrors
func ConvertToOldErrors(errs []NewValidationError) []OldValidationError {
	result := make([]OldValidationError, len(errs))
	for i, err := range errs {
		result[i] = ConvertToOldError(err)
	}
	return result
}
