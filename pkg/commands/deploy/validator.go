// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
	"gopkg.in/yaml.v3"
)

// validateDeployConfig validates the deployment configuration
func validateDeployConfig(yamlConfig *template.NexlayerYAML) error {
	if yamlConfig == nil {
		return fmt.Errorf("deployment configuration is required")
	}

	// Convert YAML struct to string for validation
	yamlBytes, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	// Use the centralized validation system
	errors, err := validation.ValidateYAMLString(string(yamlBytes))
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if len(errors) > 0 {
		// Format all validation errors
		var errMsg string
		for _, err := range errors {
			errMsg += fmt.Sprintf("\n- %s", err.Message)
			if len(err.Suggestions) > 0 {
				for _, suggestion := range err.Suggestions {
					errMsg += fmt.Sprintf("\n  Suggestion: %s", suggestion)
				}
			}
		}
		return fmt.Errorf("validation failed:%s", errMsg)
	}

	return nil
}

// validatePod validates a pod configuration using the centralized validation system
func validatePod(pod template.Pod) error {
	// Create a minimal YAML with just the pod to validate
	yamlConfig := &template.NexlayerYAML{
		Application: template.Application{
			Name: "temp",
			Pods: []template.Pod{pod},
		},
	}

	// Convert YAML struct to string for validation
	yamlBytes, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	errors, err := validation.ValidateYAMLString(string(yamlBytes))
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	if len(errors) > 0 {
		// Return the first error for backward compatibility
		return fmt.Errorf("pod validation failed: %s", errors[0].Message)
	}

	return nil
}
