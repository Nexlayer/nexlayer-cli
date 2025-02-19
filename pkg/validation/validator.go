// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
	"gopkg.in/yaml.v2"
)

var (
	// DefaultValidator is the package-level validator instance
	DefaultValidator = NewValidator(false)
)

// ValidateYAMLString validates a YAML string against the Nexlayer schema
func ValidateYAMLString(yamlContent string) ([]ValidationError, error) {
	var config schema.NexlayerYAML
	if err := yaml.Unmarshal([]byte(yamlContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return DefaultValidator.ValidateYAML(&config), nil
}

// ValidateYAMLBytes validates a YAML byte slice against the Nexlayer schema
func ValidateYAMLBytes(yamlBytes []byte) ([]ValidationError, error) {
	var config schema.NexlayerYAML
	if err := yaml.Unmarshal(yamlBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return DefaultValidator.ValidateYAML(&config), nil
}

// ValidationError represents a validation error with context and suggestions
type ValidationError struct {
	Field       string
	Message     string
	Suggestions []string
	Severity    string // error, warning
}

func (e ValidationError) Error() string {
	base := fmt.Sprintf("%s: %s", e.Field, e.Message)
	if len(e.Suggestions) > 0 {
		base += "\nSuggestions:"
		for _, s := range e.Suggestions {
			base += fmt.Sprintf("\n  - %s", s)
		}
	}
	return base
}

// Validator provides YAML configuration validation
type Validator struct {
	strict bool
}

// NewValidator creates a new validator instance
func NewValidator(strict bool) *Validator {
	return &Validator{
		strict: strict,
	}
}

// ValidateYAML performs comprehensive validation of a Nexlayer YAML configuration
func (v *Validator) ValidateYAML(yaml *schema.NexlayerYAML) []ValidationError {
	var errors []ValidationError

	// Validate application
	if yaml.Application.Name == "" {
		errors = append(errors, ValidationError{
			Field:    "application.name",
			Message:  "Application name is required",
			Severity: "error",
			Suggestions: []string{
				"Add 'name' field under 'application' section",
				"Use a descriptive name that reflects your application's purpose",
			},
		})
	} else if !isValidName(yaml.Application.Name) {
		errors = append(errors, ValidationError{
			Field:    "application.name",
			Message:  "Invalid application name format",
			Severity: "error",
			Suggestions: []string{
				"Use only lowercase letters, numbers, and dashes",
				"Start with a letter",
				"Example: my-app-123",
			},
		})
	}

	// Validate pods
	if len(yaml.Application.Pods) == 0 {
		errors = append(errors, ValidationError{
			Field:    "application.pods",
			Message:  "No pod configurations found",
			Severity: "error",
			Suggestions: []string{
				"Add at least one pod under 'pods' section",
				"Each pod should define an 'image' and 'servicePorts'",
				"Example:\n  pods:\n    - name: web\n      image: nginx:latest\n      servicePorts:\n        - 80",
			},
		})
	}

	for i, pod := range yaml.Application.Pods {
		podErrors := v.validatePod(pod, i)
		errors = append(errors, podErrors...)
	}

	// Validate registry credentials if needed
	if yaml.Application.RegistryLogin != nil {
		credErrors := v.validateRegistryCredentials(*yaml.Application.RegistryLogin)
		errors = append(errors, credErrors...)
	}

	return errors
}

// validatePod validates a single pod configuration
func (v *Validator) validatePod(pod schema.Pod, index int) []ValidationError {
	var errors []ValidationError
	prefix := fmt.Sprintf("application.pods[%d]", index)

	// Validate required fields
	if pod.Name == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Pod name is required",
			Severity: "error",
			Suggestions: []string{
				"Add 'name' field to pod configuration",
				"Use a descriptive name that reflects the pod's purpose",
				"Example: web-server, api, database",
			},
		})
	} else if !isValidName(pod.Name) {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Invalid pod name format",
			Severity: "error",
			Suggestions: []string{
				"Use only lowercase letters, numbers, and dashes",
				"Start with a letter",
				"Example: web-app-1, api-server, redis-cache",
			},
		})
	}

	if pod.Image == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".image",
			Message: "image is required",
		})
	} else if !isValidImageName(pod.Image) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".image",
			Message: "invalid image format. Expected format: [registry/]repository[:tag]",
		})
	}

	// Validate volumes
	for j, vol := range pod.Volumes {
		volErrors := v.validateVolume(vol, fmt.Sprintf("%s.volumes[%d]", prefix, j))
		errors = append(errors, volErrors...)
	}

	// Validate service ports
	if len(pod.ServicePorts) == 0 {
		errors = append(errors, ValidationError{
			Field:   prefix + ".servicePorts",
			Message: "at least one service port is required",
		})
	}

	for _, port := range pod.ServicePorts {
		if port < 1 || port > 65535 {
			errors = append(errors, ValidationError{
				Field:   prefix + ".servicePorts",
				Message: fmt.Sprintf("invalid port number: %d (must be between 1 and 65535)", port),
			})
		}
	}

	return errors
}

// validateVolume validates a volume configuration
func (v *Validator) validateVolume(vol schema.Volume, prefix string) []ValidationError {
	var errors []ValidationError

	if vol.Name == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".name",
			Message: "volume name is required",
		})
	} else if !isValidName(vol.Name) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".name",
			Message: "volume name must be lowercase alphanumeric with dashes only",
		})
	}

	if vol.Size == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".size",
			Message: "volume size is required",
		})
	} else if !isValidVolumeSize(vol.Size) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".size",
			Message: "invalid volume size format. Expected format: number + unit (e.g., 1Gi, 500Mi)",
		})
	}

	if vol.MountPath == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".mountPath",
			Message: "volume mount path is required",
		})
	}

	return errors
}

// validateRegistryCredentials validates registry login credentials
func (v *Validator) validateRegistryCredentials(creds schema.RegistryLogin) []ValidationError {
	var errors []ValidationError
	prefix := "application.registryLogin"

	if creds.Registry == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".registry",
			Message: "registry URL is required",
		})
	}

	if creds.Username == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".username",
			Message: "registry username is required",
		})
	}

	if creds.PersonalAccessToken == "" {
		errors = append(errors, ValidationError{
			Field:   prefix + ".personalAccessToken",
			Message: "registry personal access token is required",
		})
	}

	return errors
}

// Helper functions for validation
func isValidName(name string) bool {
	return regexp.MustCompile(`^[a-z][a-z0-9-]*[a-z0-9]$`).MatchString(name)
}

func isValidImageName(image string) bool {
	// Allow template variables
	if strings.Contains(image, "<%") && strings.Contains(image, "%>") {
		return true
	}

	// Split image name into parts
	parts := strings.Split(image, ":")
	if len(parts) > 2 {
		return false // More than one colon
	}

	// Validate repository name
	repo := parts[0]
	if repo == "" || strings.HasPrefix(repo, "/") || strings.HasSuffix(repo, "/") {
		return false
	}

	// Check path components (max 3: registry/namespace/repository)
	pathParts := strings.Split(repo, "/")
	if len(pathParts) > 3 {
		return false
	}

	// Validate tag if present
	if len(parts) == 2 && parts[1] == "" {
		return false // Empty tag after colon
	}

	return true
}

func isValidVolumeSize(size string) bool {
	return regexp.MustCompile(`^\d+[KMGT]i$`).MatchString(size)
}
