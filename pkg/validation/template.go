// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package validation

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/compose/components"
	"github.com/xeipuuv/gojsonschema"
)

var (
	// Regex for validating pod names (lowercase alphanumeric, '-', '.')
	podNameRegex = regexp.MustCompile(`^[a-z][a-z0-9\.\-]*$`)
	
	// Valid volume size units
	validSizeUnits = []string{"Ki", "Mi", "Gi", "Ti"}
)

// ValidateTemplate performs comprehensive validation of a Nexlayer template
// according to the v2.0 schema specification.
func ValidateTemplate(template *components.Template) error {
	// First validate against JSON Schema
	if err := validateAgainstSchema(template); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// Additional semantic validation
	// Validate application name
	if template.Application.Name == "" {
		return fmt.Errorf("application name is required")
	}

	// Validate registry login if provided
	if template.Application.RegistryLogin != nil {
		if err := validateRegistryLogin(template.Application.RegistryLogin); err != nil {
			return fmt.Errorf("invalid registry login: %w", err)
		}
	}

	// Validate pods
	if len(template.Application.Pods) == 0 {
		return fmt.Errorf("at least one pod is required")
	}

	usedPorts := make(map[int]string) // Track used service ports
	for i, pod := range template.Application.Pods {
		if err := validatePod(pod, usedPorts); err != nil {
			return fmt.Errorf("invalid pod %q (#%d): %w", pod.Name, i+1, err)
		}
	}

	return nil
}

// validateRegistryLogin validates registry authentication details
func validateRegistryLogin(login *components.RegistryLogin) error {
	if login.Registry == "" {
		return fmt.Errorf("registry URL is required")
	}
	if login.Username == "" {
		return fmt.Errorf("username is required")
	}
	if login.PersonalAccessToken == "" {
		return fmt.Errorf("personal access token is required")
	}
	return nil
}

// validatePod validates a single pod configuration
func validatePod(pod components.Pod, usedPorts map[int]string) error {
	// Validate pod name
	if pod.Name == "" {
		return fmt.Errorf("pod name is required")
	}
	if !podNameRegex.MatchString(pod.Name) {
		return fmt.Errorf("invalid pod name %q: must start with lowercase letter and contain only alphanumeric, '-', or '.'", pod.Name)
	}

	// Validate image
	if pod.Image == "" {
		return fmt.Errorf("image is required")
	}

	// Validate service ports
	if len(pod.ServicePorts) == 0 {
		return fmt.Errorf("at least one service port is required")
	}
	for _, port := range pod.ServicePorts {
		if port <= 0 || port > 65535 {
			return fmt.Errorf("invalid port number %d: must be between 1 and 65535", port)
		}
		if existingPod, exists := usedPorts[port]; exists {
			return fmt.Errorf("port %d is already in use by pod %q", port, existingPod)
		}
		usedPorts[port] = pod.Name
	}

	// Validate volumes if present
	for _, vol := range pod.Volumes {
		if err := validateVolume(vol); err != nil {
			return fmt.Errorf("invalid volume %q: %w", vol.Name, err)
		}
	}

	// Validate secrets if present
	for _, secret := range pod.Secrets {
		if err := validateSecret(secret); err != nil {
			return fmt.Errorf("invalid secret %q: %w", secret.Name, err)
		}
	}

	return nil
}

// validateAgainstSchema validates the template against the JSON Schema
func validateAgainstSchema(template *components.Template) error {
	// Use embedded schema
	schemaLoader := gojsonschema.NewStringLoader(SchemaV2)

	// Convert template to JSON for validation
	templateJSON, err := json.Marshal(template)
	if err != nil {
		return fmt.Errorf("failed to marshal template to JSON: %w", err)
	}

	documentLoader := gojsonschema.NewBytesLoader(templateJSON)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var errors []string
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return fmt.Errorf("schema validation failed:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

// validateVolume validates a volume configuration
func validateVolume(vol components.Volume) error {
	if vol.Name == "" {
		return fmt.Errorf("volume name is required")
	}
	if vol.Size == "" {
		return fmt.Errorf("volume size is required")
	}
	if vol.MountPath == "" {
		return fmt.Errorf("volume mount path is required")
	}

	// Validate size format
	valid := false
	for _, unit := range validSizeUnits {
		if strings.HasSuffix(vol.Size, unit) {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid size %q: must end with one of %v", vol.Size, validSizeUnits)
	}

	return nil
}

// validateSecret validates a secret configuration
func validateSecret(secret components.Secret) error {
	if secret.Name == "" {
		return fmt.Errorf("secret name is required")
	}
	if secret.Data == "" {
		return fmt.Errorf("secret data is required")
	}
	if secret.MountPath == "" {
		return fmt.Errorf("secret mount path is required")
	}
	if secret.FileName == "" {
		return fmt.Errorf("secret file name is required")
	}
	return nil
}
