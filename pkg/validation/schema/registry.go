package schema

import (
	"regexp"
	"strings"
)

// ValidationRegistry maintains a registry of validation functions and rules
type ValidationRegistry struct {
	validators map[string]ValidationRule
}

// NewValidationRegistry creates a new validation registry with default validators
func NewValidationRegistry() *ValidationRegistry {
	r := &ValidationRegistry{
		validators: make(map[string]ValidationRule),
	}

	// Register default validators
	r.registerDefaultValidators()
	return r
}

// Register adds a validator to the registry
func (r *ValidationRegistry) Register(name string, validator interface{}) {
	switch v := validator.(type) {
	case ValidationRule:
		r.validators[name] = v
	case ValidatorFunc:
		r.validators[name] = ValidationFuncAdapter{ValidatorFunc: v}
	}
}

// Get retrieves a validator from the registry
func (r *ValidationRegistry) Get(name string) (ValidationRule, bool) {
	v, ok := r.validators[name]
	return v, ok
}

// Validate runs a named validator against a value
func (r *ValidationRegistry) Validate(name, field string, value interface{}, ctx *ValidationContext) []ValidationError {
	validator, ok := r.Get(name)
	if !ok {
		return []ValidationError{
			{
				Field:    field,
				Message:  "no validator found for " + name,
				Severity: string(ValidationErrorSeverityError),
			},
		}
	}
	return validator.Validate(field, value, ctx)
}

// registerDefaultValidators registers all built-in validators
func (r *ValidationRegistry) registerDefaultValidators() {
	// Pod name validator
	r.Register("podname", validatePodName)

	// URL validator
	r.Register("url", validateURL)

	// Image name validator
	r.Register("image", validateImageName)

	// Volume size validator
	r.Register("volumesize", validateVolumeSize)

	// Environment variable name validator
	r.Register("envvar", validateEnvVar)

	// Filename validator
	r.Register("filename", validateFileName)

	// Name validator (generic)
	r.Register("name", validateName)
}

// Helper validation functions

// isValidName checks if a string is a valid name (lowercase alphanumeric with hyphens)
func isValidName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9\.\-]*$`, name)
	return matched
}

// isValidPodName checks if a string is a valid pod name
func isValidPodName(name string) bool {
	return isValidName(name)
}

// isValidVolumeName checks if a string is a valid volume name
func isValidVolumeName(name string) bool {
	return isValidName(name)
}

// isValidImageName checks if a string is a valid image name
func isValidImageName(image string) bool {
	// Allow template variables
	if strings.Contains(image, "<%") && strings.Contains(image, "%>") {
		return true
	}

	// Basic docker image validation (registry/repo:tag)
	matched, _ := regexp.MatchString(`^([a-zA-Z0-9\-\.]+(\.[a-zA-Z0-9\-\.]+)*(:[0-9]+)?/)?[a-zA-Z0-9\-\.]+(/[a-zA-Z0-9\-\.]+)*:[a-zA-Z0-9\-\.]+$`, image)
	return matched
}

// isValidVolumeSize checks if a string is a valid volume size
func isValidVolumeSize(size string) bool {
	matched, _ := regexp.MatchString(`^\d+([KMGT]i)?$`, size)
	return matched
}

// isValidFileName checks if a string is a valid file name
func isValidFileName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9\.\-_]*$`, name)
	return matched
}
