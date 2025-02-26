package schema

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidationError represents a validation error with suggestions for fixing it
type ValidationError struct {
	Field       string   `json:"field"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions"`
}

// Validate validates a NexlayerYAML configuration
func Validate(yaml *NexlayerYAML) error {
	if yaml == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Validate application name
	if yaml.Application.Name == "" {
		return fmt.Errorf("application name is required")
	}

	// Validate application URL if provided
	if yaml.Application.URL != "" {
		if !strings.HasPrefix(yaml.Application.URL, "http://") && !strings.HasPrefix(yaml.Application.URL, "https://") {
			return fmt.Errorf("application URL must start with http:// or https://")
		}
	}

	// Validate registry login if provided
	if yaml.Application.RegistryLogin != nil {
		if yaml.Application.RegistryLogin.Registry == "" {
			return fmt.Errorf("registry is required when registry login is provided")
		}
		if yaml.Application.RegistryLogin.Username == "" {
			return fmt.Errorf("username is required when registry login is provided")
		}
		if yaml.Application.RegistryLogin.PersonalAccessToken == "" {
			return fmt.Errorf("personal access token is required when registry login is provided")
		}
	}

	// Validate pods
	if len(yaml.Application.Pods) == 0 {
		return fmt.Errorf("at least one pod is required")
	}

	for _, pod := range yaml.Application.Pods {
		if err := validatePod(pod); err != nil {
			return fmt.Errorf("invalid pod %q: %v", pod.Name, err)
		}
	}

	return nil
}

// validatePod validates a single pod configuration
func validatePod(pod Pod) error {
	// Validate pod name
	if pod.Name == "" {
		return fmt.Errorf("pod name is required")
	}
	if !isValidPodName(pod.Name) {
		return fmt.Errorf("pod name must start with a letter and contain only lowercase letters, numbers, and hyphens")
	}

	// Validate image
	if pod.Image == "" {
		return fmt.Errorf("image is required")
	}

	// Validate path if provided
	if pod.Path != "" && !strings.HasPrefix(pod.Path, "/") {
		return fmt.Errorf("path must start with /")
	}

	// Validate service ports
	if len(pod.ServicePorts) == 0 {
		return fmt.Errorf("at least one service port is required")
	}

	// Validate volumes if provided
	for _, volume := range pod.Volumes {
		if err := validateVolume(volume); err != nil {
			return fmt.Errorf("invalid volume %q: %v", volume.Name, err)
		}
	}

	// Validate secrets if provided
	for _, secret := range pod.Secrets {
		if err := validateSecret(secret); err != nil {
			return fmt.Errorf("invalid secret %q: %v", secret.Name, err)
		}
	}

	// Validate environment variables if provided
	for _, envVar := range pod.Vars {
		if err := validateEnvVar(envVar); err != nil {
			return fmt.Errorf("invalid environment variable %q: %v", envVar.Key, err)
		}
	}

	return nil
}

// validateVolume validates a volume configuration
func validateVolume(volume Volume) error {
	if volume.Name == "" {
		return fmt.Errorf("volume name is required")
	}
	if !isValidName(volume.Name) {
		return fmt.Errorf("volume name must contain only lowercase letters, numbers, and hyphens")
	}
	if volume.Path == "" {
		return fmt.Errorf("volume path is required")
	}
	if !strings.HasPrefix(volume.Path, "/") {
		return fmt.Errorf("volume path must start with /")
	}
	if volume.Size != "" && !isValidVolumeSize(volume.Size) {
		return fmt.Errorf("invalid volume size format (e.g., '1Gi', '500Mi')")
	}
	return nil
}

// validateSecret validates a secret configuration
func validateSecret(secret Secret) error {
	if secret.Name == "" {
		return fmt.Errorf("secret name is required")
	}
	if !isValidName(secret.Name) {
		return fmt.Errorf("secret name must contain only lowercase letters, numbers, and hyphens")
	}
	if secret.Data == "" {
		return fmt.Errorf("secret data is required")
	}
	if secret.Path == "" {
		return fmt.Errorf("secret path is required")
	}
	if !strings.HasPrefix(secret.Path, "/") {
		return fmt.Errorf("secret path must start with /")
	}
	if secret.FileName == "" {
		return fmt.Errorf("secret file name is required")
	}
	return nil
}

// validateEnvVar validates an environment variable
func validateEnvVar(envVar EnvVar) error {
	if envVar.Key == "" {
		return fmt.Errorf("environment variable key is required")
	}
	if !isValidEnvVarName(envVar.Key) {
		return fmt.Errorf("invalid environment variable name format")
	}
	if envVar.Value == "" {
		return fmt.Errorf("environment variable value is required")
	}
	return nil
}

// Helper functions for validation

var (
	podNameRegex    = regexp.MustCompile(`^[a-z][a-z0-9-]*$`)
	nameRegex       = regexp.MustCompile(`^[a-z0-9-]+$`)
	envVarNameRegex = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)
	volumeSizeRegex = regexp.MustCompile(`^\d+[KMGT]i?$`)
)

func isValidPodName(name string) bool {
	return podNameRegex.MatchString(name)
}

func isValidName(name string) bool {
	return nameRegex.MatchString(name)
}

func isValidEnvVarName(name string) bool {
	return envVarNameRegex.MatchString(name)
}

func isValidVolumeSize(size string) bool {
	return volumeSizeRegex.MatchString(size)
}
