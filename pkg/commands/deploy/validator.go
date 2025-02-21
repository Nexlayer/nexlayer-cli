// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/template"
)

// ValidationError represents a single validation error with field path and suggestions
type ValidationError struct {
	Field       string
	Message     string
	Suggestions []string
}

// Validator holds the configuration and collects validation errors
type Validator struct {
	config *template.NexlayerYAML
	errors []ValidationError
}

// NewValidator creates a new Validator instance
func NewValidator(config *template.NexlayerYAML) *Validator {
	return &Validator{config: config}
}

// Validate performs the full validation of the NexlayerYAML configuration
func (v *Validator) Validate() error {
	if v.config == nil {
		v.errors = append(v.errors, ValidationError{
			Field:   "",
			Message: "deployment configuration is required",
		})
		return v.formatErrors()
	}

	v.validateApplication()
	v.validateRegistryLogin()
	v.validatePods()

	if len(v.errors) > 0 {
		return v.formatErrors()
	}
	return nil
}

// validateApplication checks the application-level fields
func (v *Validator) validateApplication() {
	if v.config.Application.Name == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   "application.name",
			Message: "application name is required",
		})
	}
	// Optional: application.url is not validated as it's optional per schema
}

// validateRegistryLogin ensures registry login is correctly configured if present
func (v *Validator) validateRegistryLogin() {
	rl := v.config.Application.RegistryLogin
	if rl != nil {
		if rl.Registry == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.registry",
				Message: "registry hostname is required when registryLogin is present",
			})
		}
		if rl.Username == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.username",
				Message: "registry username is required when registryLogin is present",
			})
		}
		if rl.PersonalAccessToken == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.personalAccessToken",
				Message: "registry personal access token is required when registryLogin is present",
			})
		}
	}
}

// validatePods checks all pod configurations
func (v *Validator) validatePods() {
	if len(v.config.Application.Pods) == 0 {
		v.errors = append(v.errors, ValidationError{
			Field:   "application.pods",
			Message: "at least one pod is required",
		})
		return
	}

	podNames := make(map[string]bool)
	for i, pod := range v.config.Application.Pods {
		podNames[pod.Name] = true
		v.validatePod(i, pod)
	}

	// Validate pod references in environment variables
	for i, pod := range v.config.Application.Pods {
		for _, varEnv := range pod.Vars {
			if strings.Contains(varEnv.Value, ".pod") {
				refPod := extractPodName(varEnv.Value)
				if refPod != "" && !podNames[refPod] {
					v.errors = append(v.errors, ValidationError{
						Field:   fmt.Sprintf("pods[%d].vars[%s]", i, varEnv.Key),
						Message: fmt.Sprintf("referenced pod '%s' not found", refPod),
						Suggestions: []string{
							"Check the pod name in the reference",
							"Ensure the referenced pod is defined in the configuration",
						},
					})
				}
			}
		}
	}
}

// validatePod checks an individual pod's configuration
func (v *Validator) validatePod(index int, pod template.Pod) {
	if pod.Name == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].name", index),
			Message: "pod name is required",
		})
	} else if !isValidName(pod.Name) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].name", index),
			Message: "pod name must follow Kubernetes naming conventions",
			Suggestions: []string{
				"Use a name like 'my-pod' or 'web-service'",
			},
		})
	}

	// Validate pod type if specified
	if pod.Type != "" && !isValidPodType(pod.Type) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].type", index),
			Message: fmt.Sprintf("invalid pod type: %s", pod.Type),
			Suggestions: []string{
				"Valid types: nextjs, react, node, python, go, postgres, redis, mongodb",
			},
		})
	}

	if pod.Image == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].image", index),
			Message: "pod image is required",
		})
	} else if strings.Contains(pod.Image, "<% REGISTRY %>") {
		if !strings.HasPrefix(pod.Image, "<% REGISTRY %>/") {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].image", index),
				Message: "private images must start with '<% REGISTRY %>/'",
				Suggestions: []string{
					"Example: <% REGISTRY %>/myapp/backend:v1.0.0",
				},
			})
		}
	}

	// Validate volumes
	for j, volume := range pod.Volumes {
		if volume.Name == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].volumes[%d].name", index, j),
				Message: "volume name is required",
			})
		} else if !isValidName(volume.Name) {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].volumes[%d].name", index, j),
				Message: "volume name must be lowercase, alphanumeric, or contain '-'",
			})
		}
		if volume.Size == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].volumes[%d].size", index, j),
				Message: "volume size is required",
			})
		}
		if volume.Path == "" || !strings.HasPrefix(volume.Path, "/") {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].volumes[%d].path", index, j),
				Message: "path must start with '/'",
				Suggestions: []string{
					fmt.Sprintf("Change to '/%s'", volume.Path),
				},
			})
		}
	}

	// Validate secrets
	for j, secret := range pod.Secrets {
		if secret.Name == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].secrets[%d].name", index, j),
				Message: "secret name is required",
			})
		}
		if secret.Data == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].secrets[%d].data", index, j),
				Message: "secret data is required",
			})
		}
		if secret.Path == "" || !strings.HasPrefix(secret.Path, "/") {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].secrets[%d].path", index, j),
				Message: "path must start with '/'",
			})
		}
		if secret.FileName == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].secrets[%d].fileName", index, j),
				Message: "fileName is required for secrets",
			})
		}
	}

	// Validate service ports
	if len(pod.ServicePorts) == 0 {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts", index),
			Message: "at least one service port is required",
		})
	} else {
		for j, port := range pod.ServicePorts {
			v.validateServicePort(index, j, port)
		}
	}
}

// validateServicePort validates an individual service port configuration
func (v *Validator) validateServicePort(podIndex, portIndex int, port template.ServicePort) {
	if port.Name == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts[%d].name", podIndex, portIndex),
			Message: "port name is required",
			Suggestions: []string{
				"Use descriptive names like 'http', 'api', or 'metrics'",
			},
		})
	} else if !isValidName(port.Name) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts[%d].name", podIndex, portIndex),
			Message: "port name must be lowercase alphanumeric with hyphens",
		})
	}

	if port.Port < 1 || port.Port > 65535 {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts[%d].port", podIndex, portIndex),
			Message: fmt.Sprintf("invalid port number: %d (must be between 1 and 65535)", port.Port),
		})
	}

	if port.TargetPort < 1 || port.TargetPort > 65535 {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts[%d].targetPort", podIndex, portIndex),
			Message: fmt.Sprintf("invalid target port number: %d (must be between 1 and 65535)", port.TargetPort),
		})
	}

	if port.Protocol != "" && !isValidProtocol(port.Protocol) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts[%d].protocol", podIndex, portIndex),
			Message: fmt.Sprintf("invalid protocol: %s", port.Protocol),
			Suggestions: []string{
				"Valid protocols: TCP, UDP, SCTP",
			},
		})
	}
}

// isValidPodType checks if the pod type is supported
func isValidPodType(podType string) bool {
	validTypes := map[string]bool{
		"nextjs":   true,
		"react":    true,
		"node":     true,
		"python":   true,
		"go":       true,
		"postgres": true,
		"redis":    true,
		"mongodb":  true,
	}
	return validTypes[podType]
}

// isValidProtocol checks if the protocol is supported
func isValidProtocol(protocol string) bool {
	validProtocols := map[string]bool{
		"TCP":  true,
		"UDP":  true,
		"SCTP": true,
	}
	return validProtocols[protocol]
}

// Helper functions

// isValidName checks if a name follows Kubernetes naming conventions
func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	// Must start with a lowercase letter
	if !strings.Contains("abcdefghijklmnopqrstuvwxyz", string(name[0])) {
		return false
	}
	// Can only contain lowercase letters, numbers, '-', and '.'
	for _, c := range name {
		if !strings.Contains("abcdefghijklmnopqrstuvwxyz0123456789-.", string(c)) {
			return false
		}
	}
	return true
}

// extractPodName extracts the pod name from a pod reference in an environment variable
func extractPodName(value string) string {
	parts := strings.Split(value, ".pod.")
	if len(parts) != 2 {
		return ""
	}
	return parts[0]
}

// formatErrors formats all validation errors into a single error message
func (v *Validator) formatErrors() error {
	var messages []string
	for _, err := range v.errors {
		msg := fmt.Sprintf("Error in %s: %s", err.Field, err.Message)
		if len(err.Suggestions) > 0 {
			msg += "\nSuggestions:"
			for _, suggestion := range err.Suggestions {
				msg += fmt.Sprintf("\n  - %s", suggestion)
			}
		}
		messages = append(messages, msg)
	}
	return fmt.Errorf("validation failed:\n%s", strings.Join(messages, "\n\n"))
}

// ValidatePod validates a single pod configuration
func ValidatePod(pod template.Pod) error {
	validator := NewValidator(&template.NexlayerYAML{
		Application: template.Application{
			Name: "temp",
			Pods: []template.Pod{pod},
		},
	})
	return validator.Validate()
}
