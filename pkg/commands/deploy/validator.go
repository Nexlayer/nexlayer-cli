package deploy

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
)

// validateDeployConfig validates the deployment configuration
func validateDeployConfig(yaml *schema.NexlayerYAML) error {
	if yaml == nil {
		return fmt.Errorf("deployment configuration is required")
	}

	// Use the centralized validation system
	validator := validation.NewValidator(true)
	errors := validator.ValidateYAML(yaml)

	if len(errors) > 0 {
		// Return the first error for backward compatibility
		return fmt.Errorf("validation failed: %s: %s", errors[0].Field, errors[0].Message)
	}

	return nil
}

// validatePod validates a pod configuration using the centralized validation system
func validatePod(pod schema.Pod) error {
	// Create a minimal YAML with just the pod to validate
	yaml := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: "temp",
			Pods: []schema.Pod{pod},
		},
	}

	validator := validation.NewValidator(true)
	errors := validator.ValidateYAML(yaml)

	if len(errors) > 0 {
		// Return the first error for backward compatibility
		return fmt.Errorf("pod validation failed: %s: %s", errors[0].Field, errors[0].Message)
	}

	return nil
}
