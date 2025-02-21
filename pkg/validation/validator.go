// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"gopkg.in/yaml.v3"
)

var (
	// DefaultValidator is the package-level validator instance
	DefaultValidator = NewValidator(false)
)

// ValidateYAMLString validates a YAML string against the Nexlayer schema
func ValidateYAMLString(yamlContent string) ([]ValidationError, error) {
	var config types.NexlayerYAML
	if err := yaml.Unmarshal([]byte(yamlContent), &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return DefaultValidator.ValidateYAML(&config), nil
}

// ValidateYAMLBytes validates a YAML byte slice against the Nexlayer schema
func ValidateYAMLBytes(yamlBytes []byte) ([]ValidationError, error) {
	var config types.NexlayerYAML
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

// ValidateYAML performs basic validation of a Nexlayer YAML configuration
func (v *Validator) ValidateYAML(yaml *types.NexlayerYAML) []ValidationError {
	var errors []ValidationError

	// Validate application name
	if yaml.Application.Name == "" {
		errors = append(errors, ValidationError{
			Field:    "application.name",
			Message:  "Application name is required",
			Severity: "error",
		})
	} else if !isValidName(yaml.Application.Name) {
		errors = append(errors, ValidationError{
			Field:    "application.name",
			Message:  "Invalid application name format",
			Severity: "error",
			Suggestions: []string{
				"Must start with a lowercase letter",
				"Can include only alphanumeric characters, '-', '.'",
				"Example: my-app.v1",
			},
		})
	}

	// Validate pods
	if len(yaml.Application.Pods) == 0 {
		errors = append(errors, ValidationError{
			Field:    "application.pods",
			Message:  "At least one pod configuration is required",
			Severity: "error",
		})
	}

	for i, pod := range yaml.Application.Pods {
		podErrors := v.validatePod(pod, i)
		errors = append(errors, podErrors...)
	}

	return errors
}

// validatePod performs basic validation of a pod configuration
func (v *Validator) validatePod(pod types.Pod, index int) []ValidationError {
	var errors []ValidationError
	prefix := fmt.Sprintf("application.pods[%d]", index)

	// Validate required fields
	if pod.Name == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Pod name is required",
			Severity: "error",
		})
	} else if !isValidName(pod.Name) {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Invalid pod name format",
			Severity: "error",
			Suggestions: []string{
				"Must start with a lowercase letter",
				"Can include only alphanumeric characters, '-', '.'",
				"Example: web-server.v1",
			},
		})
	}

	if pod.Image == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".image",
			Message:  "Image is required",
			Severity: "error",
		})
	} else if !isValidImageName(pod.Image) {
		errors = append(errors, ValidationError{
			Field:    prefix + ".image",
			Message:  "Invalid image format",
			Severity: "error",
			Suggestions: []string{
				"For private images: <% REGISTRY %>/path/image:tag",
				"For public images: [registry/]repository:tag",
				"Example private: <% REGISTRY %>/myapp/api:v1.0.0",
				"Example public: nginx:latest",
			},
		})
	}

	// Validate service ports
	if len(pod.ServicePorts) == 0 {
		errors = append(errors, ValidationError{
			Field:    prefix + ".servicePorts",
			Message:  "At least one service port is required",
			Severity: "error",
		})
	}

	for j, port := range pod.ServicePorts {
		portErrors := v.validateServicePort(port, fmt.Sprintf("%s.servicePorts[%d]", prefix, j))
		errors = append(errors, portErrors...)
	}

	// Validate volumes if present
	for j, volume := range pod.Volumes {
		volumeErrors := v.validateVolume(volume, fmt.Sprintf("%s.volumes[%d]", prefix, j))
		errors = append(errors, volumeErrors...)
	}

	return errors
}

// validateServicePort performs basic validation of a service port configuration
func (v *Validator) validateServicePort(port types.ServicePort, prefix string) []ValidationError {
	var errors []ValidationError

	if port.Name == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Port name is required",
			Severity: "error",
		})
	}

	if port.Port < 1 || port.Port > 65535 {
		errors = append(errors, ValidationError{
			Field:    prefix + ".port",
			Message:  fmt.Sprintf("Invalid port number: %d (must be between 1 and 65535)", port.Port),
			Severity: "error",
		})
	}

	if port.TargetPort < 1 || port.TargetPort > 65535 {
		errors = append(errors, ValidationError{
			Field:    prefix + ".targetPort",
			Message:  fmt.Sprintf("Invalid target port number: %d (must be between 1 and 65535)", port.TargetPort),
			Severity: "error",
		})
	}

	return errors
}

// validateVolume performs validation of a volume configuration
func (v *Validator) validateVolume(volume types.Volume, prefix string) []ValidationError {
	var errors []ValidationError

	if volume.Name == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Volume name is required",
			Severity: "error",
		})
	} else if !isValidVolumeName(volume.Name) {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Invalid volume name format",
			Severity: "error",
			Suggestions: []string{
				"Must start with a lowercase letter",
				"Can include only lowercase letters, numbers, and hyphens",
				"Example: data-volume-1",
			},
		})
	}

	if volume.Path == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".path",
			Message:  "Volume path is required",
			Severity: "error",
			Suggestions: []string{
				"Must start with '/'",
				"Example: /var/lib/data",
			},
		})
	} else if !strings.HasPrefix(volume.Path, "/") {
		errors = append(errors, ValidationError{
			Field:    prefix + ".path",
			Message:  "Volume path must start with '/'",
			Severity: "error",
		})
	}

	if volume.Size != "" && !isValidVolumeSize(volume.Size) {
		errors = append(errors, ValidationError{
			Field:    prefix + ".size",
			Message:  "Invalid volume size format",
			Severity: "error",
			Suggestions: []string{
				"Use a positive integer with a valid unit (Ki, Mi, Gi, Ti)",
				"Example: 1Gi",
				"Example: 500Mi",
			},
		})
	}

	return errors
}

// isValidName checks if a name follows Nexlayer schema requirements
// - must start with a lowercase letter
// - can include only alphanumeric characters, '-', '.'
func isValidName(name string) bool {
	if name == "" {
		return false
	}

	// Must start with a lowercase letter
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}

	// Only allow lowercase letters, numbers, hyphens, and dots
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.') {
			return false
		}
	}

	return true
}

// isValidVolumeName checks if a volume name follows Nexlayer schema requirements
// - must start with a lowercase letter
// - can include only lowercase letters, numbers, and hyphens
func isValidVolumeName(name string) bool {
	if name == "" {
		return false
	}

	// Must start with a lowercase letter
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}

	// Only allow lowercase letters, numbers, and hyphens
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}

	return true
}

// isValidSecretName checks if a secret name follows Nexlayer schema requirements
// - must start with a lowercase letter
// - can include only lowercase letters, numbers, and hyphens
func isValidSecretName(name string) bool {
	if name == "" {
		return false
	}

	// Must start with a lowercase letter
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}

	// Only allow lowercase letters, numbers, and hyphens
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}

	return true
}

// isValidImageName checks if a Docker image name is valid for Nexlayer
// Format: <% REGISTRY %>/path/image:tag for private images
// Format: standard Docker image name for public images
func isValidImageName(image string) bool {
	if image == "" {
		return false
	}

	// Handle private registry images
	if strings.Contains(image, "<% REGISTRY %>") {
		parts := strings.Split(image, ":")
		if len(parts) != 2 {
			return false // Must have a tag
		}
		repo := strings.TrimPrefix(parts[0], "<% REGISTRY %>/")
		if repo == "" || strings.HasPrefix(repo, "/") || strings.HasSuffix(repo, "/") {
			return false // Invalid path after registry
		}
		return true
	}

	// Handle public images
	parts := strings.Split(image, ":")
	if len(parts) > 2 {
		return false // Too many colons
	}

	// Check repository part
	repo := parts[0]
	if strings.HasPrefix(repo, "/") || strings.HasSuffix(repo, "/") {
		return false // Cannot start or end with slash
	}

	// Count slashes (max 2 for registry/repository)
	if strings.Count(repo, "/") > 2 {
		return false
	}

	// Check each component
	components := strings.Split(repo, "/")
	for _, comp := range components {
		if comp == "" {
			return false // Empty component
		}
		// Allow letters, numbers, dots, and dashes in each component
		for _, r := range comp {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '.' || r == '-') {
				return false
			}
		}
	}

	// Check tag if present
	if len(parts) == 2 {
		tag := parts[1]
		if tag == "" {
			return false // Empty tag
		}
		// Allow letters, numbers, dots, and dashes in tag
		for _, r := range tag {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
				(r >= '0' && r <= '9') || r == '.' || r == '-') {
				return false
			}
		}
	}

	return true
}

// isValidVolumeSize checks if a volume size is valid
// Format: number + unit (Ki, Mi, Gi, Ti)
func isValidVolumeSize(size string) bool {
	if size == "" {
		return false
	}

	// Must end with valid unit
	validUnits := []string{"Ki", "Mi", "Gi", "Ti"}
	hasValidUnit := false
	for _, unit := range validUnits {
		if strings.HasSuffix(size, unit) {
			hasValidUnit = true
			size = strings.TrimSuffix(size, unit)
			break
		}
	}
	if !hasValidUnit {
		return false
	}

	// Remaining part must be a positive integer
	number, err := strconv.Atoi(size)
	if err != nil || number <= 0 {
		return false
	}

	return true
}
