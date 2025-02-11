// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

var (
	validate *validator.Validate
	// Regex patterns for custom validations
	imagePattern     = regexp.MustCompile(`^(?:([^/]+)/)?(?:([^/]+)/)?([^/]+)(?:[:@][^/]+)?$`)
	volumeSizePattern = regexp.MustCompile(`^\d+[KMGT]i?$`)
	filenamePattern  = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
	envVarPattern    = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

func init() {
	// Initialize validator with struct validation enabled
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators
	_ = validate.RegisterValidation("image", validateImage)
	_ = validate.RegisterValidation("volumesize", validateVolumeSize)
	_ = validate.RegisterValidation("filename", validateFilename)
	_ = validate.RegisterValidation("envvar", validateEnvVar)
}

// ValidateNexlayerYAML validates the provided YAML configuration
func ValidateNexlayerYAML(yaml *types.NexlayerYAML) error {
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
	return nil
}

// Custom validators
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
	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only alphanumeric characters", field)
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
