// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
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

// ValidateYAML performs validation of a YAML structure with auto-correction
func (v *Validator) ValidateYAML(config interface{}) []ValidationError {
	errors := []ValidationError{}

	// First, perform JSON Schema validation
	if v.schemaSource != nil {
		schemaErrors := v.validateWithJSONSchema(config)
		errors = append(errors, schemaErrors...)
	}

	// Then perform semantic validation with the registry
	ctx := &ValidationContext{Config: config}
	errors = append(errors, v.validateWithRegistry(config, ctx)...)

	// Attempt auto-correction for non-critical errors
	if nexConfig, ok := config.(*NexlayerYAML); ok {
		errors = append(errors, v.autoCorrectConfig(nexConfig)...)
	}

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

// isValidName checks if a string is a valid name (lowercase alphanumeric with hyphens)
func isValidName(name string) bool {
	return nameRegex.MatchString(name)
}

// isValidPodName checks if a string is a valid pod name
func isValidPodName(name string) bool {
	return podNameRegex.MatchString(name)
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
	return volumeSizeRegex.MatchString(size)
}

// isValidFileName checks if a string is a valid file name
func isValidFileName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9][a-zA-Z0-9\.\-_]*$`, name)
	return matched
}

// isValidEnvVarName checks if a string is a valid environment variable name
func isValidEnvVarName(name string) bool {
	return envVarNameRegex.MatchString(name)
}

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

// autoCorrectConfig attempts to fix common issues in the configuration
func (v *Validator) autoCorrectConfig(config *NexlayerYAML) []ValidationError {
	var errors []ValidationError

	// Auto-correct application name
	if config.Application.Name == "" {
		config.Application.Name = "nexlayer-app"
		errors = append(errors, ValidationError{
			Field:     "application.name",
			Message:   "Application name was empty, set to default value",
			Severity:  ValidationErrorSeverityWarning,
			AutoFixed: true,
		})
	}

	// Auto-correct pod configurations
	for i := range config.Application.Pods {
		pod := &config.Application.Pods[i]
		errors = append(errors, v.autoCorrectPod(pod)...)
	}

	return errors
}

// autoCorrectPod attempts to fix common issues in pod configuration
func (v *Validator) autoCorrectPod(pod *Pod) []ValidationError {
	var errors []ValidationError

	// Auto-correct pod name
	if pod.Name == "" {
		pod.Name = fmt.Sprintf("pod-%d", time.Now().Unix())
		errors = append(errors, ValidationError{
			Field:     "pod.name",
			Message:   "Pod name was empty, generated unique name",
			Severity:  ValidationErrorSeverityWarning,
			AutoFixed: true,
		})
	}

	// Ensure pod name is valid
	if !isValidPodName(pod.Name) {
		oldName := pod.Name
		pod.Name = sanitizePodName(pod.Name)
		errors = append(errors, ValidationError{
			Field:     "pod.name",
			Message:   fmt.Sprintf("Invalid pod name '%s' was sanitized to '%s'", oldName, pod.Name),
			Severity:  ValidationErrorSeverityWarning,
			AutoFixed: true,
		})
	}

	// Auto-correct service ports
	if len(pod.ServicePorts) == 0 {
		// Add default port based on common patterns
		defaultPort := getDefaultPortForImage(pod.Image)
		pod.ServicePorts = []ServicePort{{
			Name:       fmt.Sprintf("%s-port-1", pod.Name),
			Port:       defaultPort,
			TargetPort: defaultPort,
			Protocol:   "TCP",
		}}
		errors = append(errors, ValidationError{
			Field:     "pod.servicePorts",
			Message:   fmt.Sprintf("No service ports defined, added default port %d", defaultPort),
			Severity:  ValidationErrorSeverityWarning,
			AutoFixed: true,
		})
	}

	// Validate and auto-correct volume configurations
	for i := range pod.Volumes {
		vol := &pod.Volumes[i]
		if vol.Size == "" {
			vol.Size = getDefaultVolumeSize(pod.Image)
			errors = append(errors, ValidationError{
				Field:     fmt.Sprintf("pod.volumes[%d].size", i),
				Message:   fmt.Sprintf("Volume size was empty, set to default size %s", vol.Size),
				Severity:  ValidationErrorSeverityWarning,
				AutoFixed: true,
			})
		}
	}

	return errors
}

// getDefaultPortForImage returns a default port based on the image name
func getDefaultPortForImage(image string) int {
	imageLower := strings.ToLower(image)
	switch {
	case strings.Contains(imageLower, "nginx"):
		return 80
	case strings.Contains(imageLower, "postgres"):
		return 5432
	case strings.Contains(imageLower, "mysql"):
		return 3306
	case strings.Contains(imageLower, "redis"):
		return 6379
	case strings.Contains(imageLower, "mongo"):
		return 27017
	case strings.Contains(imageLower, "node"):
		return 3000
	default:
		return 8080
	}
}

// getDefaultVolumeSize returns a default volume size based on the image name
func getDefaultVolumeSize(image string) string {
	imageLower := strings.ToLower(image)
	switch {
	case strings.Contains(imageLower, "postgres"),
		strings.Contains(imageLower, "mysql"),
		strings.Contains(imageLower, "mongo"):
		return "10Gi"
	case strings.Contains(imageLower, "redis"):
		return "1Gi"
	default:
		return "1Gi"
	}
}

// sanitizePodName ensures a pod name follows Kubernetes naming conventions
func sanitizePodName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9-]`)
	name = re.ReplaceAllString(name, "-")

	// Remove consecutive hyphens
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}

	// Trim hyphens from start and end
	name = strings.Trim(name, "-")

	// Ensure it starts with a letter
	if name == "" || !('a' <= name[0] && name[0] <= 'z') {
		name = "pod-" + name
	}

	// Truncate if too long (63 characters is Kubernetes limit)
	if len(name) > 63 {
		name = name[:63]
		// Ensure it doesn't end with a hyphen
		name = strings.TrimRight(name, "-")
	}

	return name
}
