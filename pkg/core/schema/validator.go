// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ValidationContext provides context for validation
type ValidationContext struct {
	Config interface{}
}

// ValidatorFunc is a function that validates a field value
type ValidatorFunc func(field, value string, ctx *ValidationContext) []ValidationError

// ValidationRule represents a validation rule
type ValidationRule interface {
	Validate(field string, value interface{}, ctx *ValidationContext) []ValidationError
}

// ValidationFuncAdapter adapts a ValidatorFunc to implement ValidationRule
type ValidationFuncAdapter struct {
	ValidatorFunc ValidatorFunc
}

// Validate implements the ValidationRule interface
func (a ValidationFuncAdapter) Validate(field string, value interface{}, ctx *ValidationContext) []ValidationError {
	if strValue, ok := value.(string); ok {
		return a.ValidatorFunc(field, strValue, ctx)
	}
	return []ValidationError{makeValidationError(field, "value must be a string", ValidationErrorSeverityError)}
}

// ValidationRegistry maintains a registry of validation functions and rules
type ValidationRegistry struct {
	validators map[string]ValidationRule
}

// NewValidationRegistry creates a new validation registry with default validators
func NewValidationRegistry() *ValidationRegistry {
	r := &ValidationRegistry{
		validators: make(map[string]ValidationRule),
	}
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
		return []ValidationError{makeValidationError(field, "no validator found for "+name, ValidationErrorSeverityError)}
	}
	return validator.Validate(field, value, ctx)
}

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
		errors = append(errors, makeValidationError("schema", "Configuration cannot be converted to JSON: "+err.Error(), ValidationErrorSeverityError))
	}

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

// Helper functions for validation

var (
	podNameRegex    = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	nameRegex       = regexp.MustCompile(`^[a-z0-9-]+$`)
	envVarNameRegex = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	volumeSizeRegex = regexp.MustCompile(`^\d+[KMGT]i?$`)
)

func isValidName(name string) bool {
	return nameRegex.MatchString(name)
}

func isValidPodName(name string) bool {
	return podNameRegex.MatchString(name)
}

func isValidImageName(image string) bool {
	// Allow template variables
	if strings.Contains(image, "<%") && strings.Contains(image, "%>") {
		return true
	}

	// Basic docker image validation (registry/repo:tag)
	matched, _ := regexp.MatchString(`^([a-zA-Z0-9\-\.]+(\.[a-zA-Z0-9\-\.]+)*(:[0-9]+)?/)?[a-zA-Z0-9\-\.]+(/[a-zA-Z0-9\-\.]+)*:[a-zA-Z0-9\-\.]+$`, image)
	return matched
}

func isValidVolumeSize(size string) bool {
	return volumeSizeRegex.MatchString(size)
}

func isValidFileName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9\.\-_]*$`, name)
	return matched
}

func isValidEnvVarName(name string) bool {
	return envVarNameRegex.MatchString(name)
}

// Validator functions that return []ValidationError

// validateName checks if a string is a valid name (lowercase alphanumeric with hyphens)
func validateName(field, value string) []ValidationError {
	if !isValidName(value) {
		return []ValidationError{makeValidationError(field, "must contain only lowercase letters, numbers, and hyphens", ValidationErrorSeverityError)}
	}
	return nil
}

// validatePodName checks if a string is a valid pod name
func validatePodName(field, value string) []ValidationError {
	if !isValidPodName(value) {
		return []ValidationError{makeValidationError(field, "must start with a letter and contain only lowercase letters, numbers, and hyphens", ValidationErrorSeverityError)}
	}
	return nil
}

// validateURL checks if a string is a valid URL
func validateURL(field, value string) []ValidationError {
	if !strings.HasPrefix(value, "http://") && !strings.HasPrefix(value, "https://") {
		return []ValidationError{makeValidationError(field, "must start with http:// or https://", ValidationErrorSeverityError)}
	}
	return nil
}

// validateImageName checks if a string is a valid image name
func validateImageName(field, value string) []ValidationError {
	if !isValidImageName(value) {
		return []ValidationError{makeValidationError(field, "invalid image name format", ValidationErrorSeverityError,
			"Format: [registry/]repo[:tag]",
			"Example: nginx:latest",
			"Example: docker.io/library/postgres:14")}
	}
	return nil
}

// validateVolumeSize checks if a string is a valid volume size
func validateVolumeSize(field, value string) []ValidationError {
	if !isValidVolumeSize(value) {
		return []ValidationError{makeValidationError(field, "invalid volume size format", ValidationErrorSeverityError,
			"Format: <number>[KMGT]i",
			"Example: 1Gi",
			"Example: 500Mi")}
	}
	return nil
}

// validateEnvVar checks if a string is a valid environment variable name
func validateEnvVar(field, value string) []ValidationError {
	if !isValidEnvVarName(value) {
		return []ValidationError{makeValidationError(field, "must start with a letter and contain only uppercase letters, numbers, and underscores", ValidationErrorSeverityError)}
	}
	return nil
}

// validateFileName checks if a string is a valid file name
func validateFileName(field, value string) []ValidationError {
	if !isValidFileName(value) {
		return []ValidationError{makeValidationError(field, "must start with a letter or number and contain only letters, numbers, dots, hyphens, and underscores", ValidationErrorSeverityError)}
	}
	return nil
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
