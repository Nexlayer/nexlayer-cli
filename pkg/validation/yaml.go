// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation/schema"
)

var (
	validate *validator.Validate
	// Regex patterns for custom validations
	imagePattern     = regexp.MustCompile(`^(?:([^/]+)/)?(?:([^/]+)/)?([^/]+)(?:[:@][^/]+)?$`)
	volumeSizePattern = regexp.MustCompile(`^\d+[KMGT]i?$`)
	filenamePattern  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
	envVarPattern    = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	podNamePattern   = regexp.MustCompile(`^[a-z][a-z0-9\.\-]*$`)
)

func init() {
	// Initialize validator with struct validation enabled
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators
	_ = validate.RegisterValidation("image", validateImage)
	_ = validate.RegisterValidation("volumesize", validateVolumeSize)
	_ = validate.RegisterValidation("filename", validateFilename)
	_ = validate.RegisterValidation("envvar", validateEnvVar)
	_ = validate.RegisterValidation("podname", validatePodName)
}

// ValidateNexlayerYAML validates the provided YAML configuration
func ValidateNexlayerYAML(yaml *schema.NexlayerYAML) error {
	// First validate the struct
	if err := validate.Struct(yaml); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("invalid yaml structure: %w", err)
		}

		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		return fmt.Errorf("validation failed:\n%s", strings.Join(validationErrors, "\n"))
	}

	// Additional validation for Pod slice length
	if len(yaml.Application.Pods) == 0 {
		return fmt.Errorf("template must contain at least one pod")
	}

	// Validate pod names and service ports
	for _, pod := range yaml.Application.Pods {
		if pod.Name == "" {
			return fmt.Errorf("pod name cannot be empty")
		}
		if len(pod.ServicePorts) == 0 {
			return fmt.Errorf("pod '%s' must have at least one service port", pod.Name)
		}
	}

	// Additional validation for volumes
	for _, pod := range yaml.Application.Pods {
		for _, volume := range pod.Volumes {
			if !volumeSizePattern.MatchString(volume.Size) {
				return fmt.Errorf("invalid volume size '%s' for volume '%s' in pod '%s'", volume.Size, volume.Name, pod.Name)
			}
		}
	}

	return nil
}

// Custom validators
func validatePodName(fl validator.FieldLevel) bool {
	return podNamePattern.MatchString(fl.Field().String())
}

func validateImage(fl validator.FieldLevel) bool {
	return imagePattern.MatchString(fl.Field().String())
}

func validateVolumeSize(fl validator.FieldLevel) bool {
	return volumeSizePattern.MatchString(fl.Field().String())
}

func validateFilename(fl validator.FieldLevel) bool {
	return filenamePattern.MatchString(fl.Field().String())
}

func validateEnvVar(fl validator.FieldLevel) bool {
	return envVarPattern.MatchString(fl.Field().String())
}

// formatValidationError formats a validation error into a user-friendly message
func formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "podname":
		return fmt.Sprintf("Field '%s' must start with a lowercase letter and contain only lowercase alphanumeric characters, dots, or hyphens", field)
	case "image":
		return fmt.Sprintf("Field '%s' must be a valid Docker image reference", field)
	case "volumesize":
		return fmt.Sprintf("Field '%s' must be a valid volume size (e.g., '1Gi', '500Mi')", field)
	case "filename":
		return fmt.Sprintf("Field '%s' must be a valid filename", field)
	case "envvar":
		return fmt.Sprintf("Field '%s' must be a valid environment variable name", field)
	case "url":
		return fmt.Sprintf("Field '%s' must be a valid URL", field)
	case "hostname":
		return fmt.Sprintf("Field '%s' must be a valid hostname", field)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s", field, err.Param())
	case "max":
		return fmt.Sprintf("Field '%s' must not exceed %s", field, err.Param())
	case "startswith":
		return fmt.Sprintf("Field '%s' must start with '%s'", field, err.Param())
	default:
		return fmt.Sprintf("Field '%s' failed validation: %s", field, err.Tag())
	}
}
