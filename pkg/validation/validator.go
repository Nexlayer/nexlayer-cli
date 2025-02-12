// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

// ValidationError represents a validation error with context
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
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
			Field:   "application.name",
			Message: "application name is required",
		})
	} else if !isValidName(yaml.Application.Name) {
		errors = append(errors, ValidationError{
			Field:   "application.name",
			Message: "application name must be lowercase alphanumeric with dashes only",
		})
	}

	// Validate pods
	if len(yaml.Application.Pods) == 0 {
		errors = append(errors, ValidationError{
			Field:   "application.pods",
			Message: "at least one pod configuration is required",
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
			Field:   prefix + ".name",
			Message: "pod name is required",
		})
	} else if !isValidName(pod.Name) {
		errors = append(errors, ValidationError{
			Field:   prefix + ".name",
			Message: "pod name must be lowercase alphanumeric with dashes only",
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
	if len(pod.Ports) == 0 {
		errors = append(errors, ValidationError{
			Field:   prefix + ".servicePorts",
			Message: "at least one service port is required",
		})
	}

	for _, port := range pod.Ports {
		if port.ServicePort < 1 || port.ServicePort > 65535 {
			errors = append(errors, ValidationError{
				Field:   prefix + ".ports.servicePort",
				Message: fmt.Sprintf("invalid service port number: %d (must be between 1 and 65535)", port.ServicePort),
			})
		}
		if port.ContainerPort < 1 || port.ContainerPort > 65535 {
			errors = append(errors, ValidationError{
				Field:   prefix + ".ports.containerPort",
				Message: fmt.Sprintf("invalid container port number: %d (must be between 1 and 65535)", port.ContainerPort),
			})
		}
		if port.Name == "" {
			errors = append(errors, ValidationError{
				Field:   prefix + ".ports.name",
				Message: "port name is required",
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
	match, _ := regexp.MatchString("^[a-z0-9][a-z0-9-]*[a-z0-9]$", name)
	return match
}

func isValidImageName(image string) bool {
	// Basic image name validation
	// Format: [registry/]repository[:tag]
	parts := strings.Split(image, "/")
	if len(parts) > 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
	}

	// Check tag format if present
	if strings.Contains(parts[len(parts)-1], ":") {
		tagParts := strings.Split(parts[len(parts)-1], ":")
		if len(tagParts) != 2 || tagParts[0] == "" || tagParts[1] == "" {
			return false
		}
	}

	return true
}

func isValidVolumeSize(size string) bool {
	match, _ := regexp.MatchString("^[0-9]+(Ki|Mi|Gi|Ti)$", size)
	return match
}
