// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	// Regex patterns for custom validations
	appNamePattern    = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	podNamePattern    = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*[a-z0-9]$`)
	volumeNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)
	secretNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)
	imagePattern      = regexp.MustCompile(`^(?:([^/]+)/)?(?:([^/]+)/)?([^/]+)(?:[:@][^/]+)?$`)
	volumeSizePattern = regexp.MustCompile(`^\d+[KMGT]i?$`)
	filenamePattern   = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._-]*$`)
	envVarPattern     = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom validators
	_ = validate.RegisterValidation("appname", validateAppName)
	_ = validate.RegisterValidation("podname", validatePodName)
	_ = validate.RegisterValidation("volumename", validateVolumeName)
	_ = validate.RegisterValidation("secretname", validateSecretName)
	_ = validate.RegisterValidation("image", validateImage)
	_ = validate.RegisterValidation("volumesize", validateVolumeSize)
	_ = validate.RegisterValidation("filename", validateFilename)
	_ = validate.RegisterValidation("envvar", validateEnvVar)
	_ = validate.RegisterValidation("podtype", validatePodType)
}

// Validator handles all template validation
type Validator struct{}

// NewValidator creates a new template validator
func NewValidator() *Validator {
	return &Validator{}
}

// Validate validates a complete template
func (v *Validator) Validate(yaml *NexlayerYAML) error {
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
		return fmt.Errorf("at least one pod must be specified")
	}

	// Additional validation for volumes
	for _, pod := range yaml.Application.Pods {
		for _, volume := range pod.Volumes {
			if !volumeSizePattern.MatchString(volume.Size) {
				return fmt.Errorf("invalid volume size '%s' for volume '%s' in pod '%s'", volume.Size, volume.Name, pod.Name)
			}
		}
	}

	// Validate port configurations
	usedPorts := make(map[int]string)
	for _, pod := range yaml.Application.Pods {
		// Validate service ports are unique across all pods
		for _, port := range pod.ServicePorts {
			if existingPod, exists := usedPorts[port]; exists {
				return fmt.Errorf("service port %d is already used by pod '%s'", port, existingPod)
			}
			usedPorts[port] = pod.Name
		}
	}

	return nil
}

// Custom validators
func validateAppName(fl validator.FieldLevel) bool {
	return appNamePattern.MatchString(fl.Field().String())
}

func validatePodName(fl validator.FieldLevel) bool {
	return podNamePattern.MatchString(fl.Field().String())
}

func validateVolumeName(fl validator.FieldLevel) bool {
	return volumeNamePattern.MatchString(fl.Field().String())
}

func validateSecretName(fl validator.FieldLevel) bool {
	return secretNamePattern.MatchString(fl.Field().String())
}

func validateImage(fl validator.FieldLevel) bool {
	img := fl.Field().String()
	// Allow <% REGISTRY %> template variable
	if strings.Contains(img, "<% REGISTRY %>") {
		img = strings.ReplaceAll(img, "<% REGISTRY %>", "example.com")
	}
	return imagePattern.MatchString(img)
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

func validatePodType(fl validator.FieldLevel) bool {
	podType := PodType(fl.Field().String())
	switch podType {
	// Frontend types
	case Frontend, React, NextJS, Vue:
		return true

	// Backend types
	case Backend, Express, Django, FastAPI,
		Node, Python, Golang, Java:
		return true

	// Database types
	case Database, MongoDB, Postgres, Redis,
		MySQL, Clickhouse:
		return true

	// Message Queue types
	case RabbitMQ, Kafka:
		return true

	// Storage types
	case Minio, Elastic:
		return true

	// Web Server types
	case Nginx, Traefik:
		return true

	// AI/ML types
	case LLM, Ollama, HFModel, VertexAI, Jupyter:
		return true

	default:
		return false
	}
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
	case "podtype":
		return fmt.Sprintf("Field '%s' must be a valid pod type", field)
	case "startswith":
		return fmt.Sprintf("Field '%s' must start with '%s'", field, err.Param())
	default:
		return fmt.Sprintf("Field '%s' failed validation: %s", field, err.Tag())
	}
}
