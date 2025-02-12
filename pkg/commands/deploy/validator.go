package deploy

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

// validateDeployConfig validates the deployment configuration
func validateDeployConfig(yaml *schema.NexlayerYAML) error {
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
func validatePod(pod schema.Pod) error {
	if pod.Name == "" {
		return fmt.Errorf("pod name is required")
	}
	if pod.Image == "" {
		return fmt.Errorf("pod image is required")
	}
	if len(pod.Ports) == 0 {
		return fmt.Errorf("pod service ports are required")
	}

	return nil
}
