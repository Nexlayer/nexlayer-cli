package deploy

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

// validateDeployConfig validates the deployment configuration
func validateDeployConfig(yaml *types.NexlayerYAML) error {
	if yaml == nil {
		return fmt.Errorf("deployment configuration is required")
	}

	if yaml.Application.Name == "" {
		return fmt.Errorf("application name is required")
	}

	if len(yaml.Application.Pods) == 0 {
		return fmt.Errorf("at least one pod is required")
	}

	for _, pod := range yaml.Application.Pods {
		if err := validatePod(pod); err != nil {
			return fmt.Errorf("invalid pod %s: %v", pod.Name, err)
		}
	}

	return nil
}

// validatePod validates a pod configuration
func validatePod(pod types.Pod) error {
	if pod.Name == "" {
		return fmt.Errorf("pod name is required")
	}
	if pod.Type == "" {
		return fmt.Errorf("pod type is required")
	}
	if pod.Image == "" {
		return fmt.Errorf("pod image is required")
	}

	// Validate pod type
	validTypes := []string{"postgres", "mysql", "mongodb", "redis", "nginx", "react", "express", "fastapi", "django", "vue", "angular", "llm"}
	valid := false
	for _, t := range validTypes {
		if pod.Type == t {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid pod type: %s", pod.Type)
	}

	return nil
}
