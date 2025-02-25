package schema

import (
	"fmt"
	"net/url"
	"regexp"
)

// Common validation patterns
const (
	PodNamePattern    = `^[a-z][a-z0-9\.\-]*$`
	FileNamePattern   = `^[a-zA-Z0-9][a-zA-Z0-9\.\-_]*$`
	VolumeSizePattern = `^\d+[KMGT]i$`
	EnvVarPattern     = `^[a-zA-Z_][a-zA-Z0-9_]*$`
	PodRefPattern     = `<%\s*POD\.([\w\-\.]+)\.[\w\-\.]+\s*%>`
)

// validateName validates generic names
func validateName(field, value string, ctx *ValidationContext) []NewValidationError {
	matched, _ := regexp.MatchString(`^[a-z][a-z0-9\.\-]*$`, value)

	if !matched {
		return []NewValidationError{
			{
				Field:    field,
				Message:  fmt.Sprintf("invalid name: %s", value),
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"Names must start with a lowercase letter",
					"Use only lowercase letters, numbers, dots, and hyphens",
					"Example: my-service, api-v1, web.app",
				},
			},
		}
	}
	return nil
}

// validatePodName validates pod names
func validatePodName(field, value string, ctx *ValidationContext) []NewValidationError {
	if !isValidPodName(value) {
		return []NewValidationError{
			{
				Field:    field,
				Message:  fmt.Sprintf("invalid pod name: %s", value),
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"Pod names must start with a lowercase letter",
					"Use only lowercase letters, numbers, and hyphens",
					"Example: web-server, api-v1, db-postgres",
				},
			},
		}
	}
	return nil
}

// validateURL validates URLs
func validateURL(field, value string, ctx *ValidationContext) []NewValidationError {
	if value == "" {
		return nil
	}
	_, err := url.Parse(value)
	if err != nil {
		return []NewValidationError{
			{
				Field:    field,
				Message:  "invalid URL format",
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"URL must be in the format: protocol://domain.tld[:port][/path]",
					"Example: https://example.com or http://localhost:8080",
				},
			},
		}
	}
	return nil
}

// validateImageName validates Docker image names
func validateImageName(field, value string, ctx *ValidationContext) []NewValidationError {
	if !isValidImageName(value) {
		return []NewValidationError{
			{
				Field:    field,
				Message:  "invalid image format",
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"For private images: <% REGISTRY %>/path/image:tag",
					"For public images: [registry/]repository:tag",
					"Example private: <% REGISTRY %>/myapp/api:v1.0.0",
					"Example public: nginx:latest",
				},
			},
		}
	}
	return nil
}

// validateVolumeSize validates volume sizes
func validateVolumeSize(field, value string, ctx *ValidationContext) []NewValidationError {
	if !isValidVolumeSize(value) {
		return []NewValidationError{
			{
				Field:    field,
				Message:  "invalid volume size",
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"Volume size must be a number followed by a unit",
					"Valid units: Ki, Mi, Gi, Ti",
					"Example: 10Gi, 500Mi",
				},
			},
		}
	}
	return nil
}

// validateEnvVar validates environment variable names
func validateEnvVar(field, value string, ctx *ValidationContext) []NewValidationError {
	pattern := EnvVarPattern
	matched, _ := regexp.MatchString(pattern, value)

	if !matched {
		return []NewValidationError{
			{
				Field:    field,
				Message:  "invalid environment variable name",
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"Environment variable names must contain only letters, numbers, and underscores",
					"Names must not start with a number",
					"Example: APP_NAME, DB_PORT, DEBUG_MODE",
				},
			},
		}
	}
	return nil
}

// validateFileName validates filenames
func validateFileName(field, value string, ctx *ValidationContext) []NewValidationError {
	if !isValidFileName(value) {
		return []NewValidationError{
			{
				Field:    field,
				Message:  "invalid filename",
				Severity: string(ValidationErrorSeverityError),
				Suggestions: []string{
					"Filenames must contain only letters, numbers, dots, hyphens, and underscores",
					"Example: config.json, app-settings.yaml, data_file.txt",
				},
			},
		}
	}
	return nil
}
