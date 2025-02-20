// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
)

// validateDeployConfig validates the deployment configuration
func validateDeployConfig(yaml *template.NexlayerYAML) error {
	if yaml == nil {
		return fmt.Errorf("deployment configuration is required")
	}

	// Use the centralized validation system
	errors := validation.ValidateYAML(yaml)
	if len(errors) > 0 {
		// Format all validation errors
		var errMsg string
		for _, err := range errors {
			errMsg += fmt.Sprintf("\n- %s", err.Message)
			if err.Suggestion != "" {
				errMsg += fmt.Sprintf("\n  Suggestion: %s", err.Suggestion)
			}
		}
		return fmt.Errorf("validation failed:%s", errMsg)
	}

	return nil
}

// validatePod validates a pod configuration using the centralized validation system
func validatePod(pod template.Pod) error {
	// Create a minimal YAML with just the pod to validate
	yaml := &template.NexlayerYAML{
		Application: template.Application{
			Name: "temp",
			Pods: []template.Pod{pod},
		},
	}

	errors := validation.ValidateYAML(yaml)
	if len(errors) > 0 {
		// Return the first error for backward compatibility
		return fmt.Errorf("pod validation failed: %s", errors[0].Message)
	}

	return nil
}
