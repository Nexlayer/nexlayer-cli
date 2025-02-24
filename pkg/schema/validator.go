package schema

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// ValidationError represents a validation error with context and suggestions
type ValidationError struct {
	Field       string   `json:"field"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions,omitempty"`
	Severity    string   `json:"severity"` // error, warning
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

// ValidateYAML performs validation of a Nexlayer YAML configuration
func (v *Validator) ValidateYAML(yaml *NexlayerYAML) []ValidationError {
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

	// Check for duplicate pod names
	podNames := make(map[string]bool)
	for i, pod := range yaml.Application.Pods {
		if podNames[pod.Name] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].name", i),
				Message: fmt.Sprintf("duplicate pod name: %s", pod.Name),
				Suggestions: []string{
					"Each pod must have a unique name",
					fmt.Sprintf("Rename one of the pods with name '%s'", pod.Name),
				},
			})
		}
		podNames[pod.Name] = true
		errors = append(errors, v.validatePod(pod, i)...)
	}

	// Validate pod references in environment variables
	for i, pod := range yaml.Application.Pods {
		for _, varEnv := range pod.Vars {
			refs := extractPodReferences(varEnv.Value)
			for _, ref := range refs {
				if !podNames[ref] {
					suggestion := findClosestPodName(ref, podNames)
					err := ValidationError{
						Field:   fmt.Sprintf("pods[%d].vars[%s]", i, varEnv.Key),
						Message: fmt.Sprintf("referenced pod '%s' not found", ref),
					}
					if suggestion != "" {
						err.Suggestions = []string{
							fmt.Sprintf("Did you mean '%s'?", suggestion),
							fmt.Sprintf("Available pods: %s", strings.Join(getAvailablePods(podNames), ", ")),
						}
					} else {
						err.Suggestions = []string{
							fmt.Sprintf("Available pods: %s", strings.Join(getAvailablePods(podNames), ", ")),
						}
					}
					errors = append(errors, err)
				}
			}
		}
	}

	return errors
}

// validatePod performs validation of a pod configuration
func (v *Validator) validatePod(pod Pod, index int) []ValidationError {
	var errors []ValidationError
	prefix := fmt.Sprintf("application.pods[%d]", index)

	// Validate required fields
	if pod.Name == "" {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  "Pod name is required",
			Severity: "error",
		})
	} else if !isValidPodName(pod.Name) {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  fmt.Sprintf("invalid pod name: %s", pod.Name),
			Severity: "error",
			Suggestions: []string{
				"Pod names must start with a lowercase letter",
				"Use only lowercase letters, numbers, and hyphens",
				"Example: web-server, api-v1, db-postgres",
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

	// Check for duplicate port names and numbers
	portNames := make(map[string]bool)
	portNumbers := make(map[int]bool)
	for j, port := range pod.ServicePorts {
		if port.Name == "" {
			errors = append(errors, ValidationError{
				Field:    fmt.Sprintf("%s.servicePorts[%d].name", prefix, j),
				Message:  "Port name is required",
				Severity: "error",
				Suggestions: []string{
					"Use descriptive names like 'http', 'api', or 'metrics'",
				},
			})
		} else if !isValidName(port.Name) {
			errors = append(errors, ValidationError{
				Field:    fmt.Sprintf("%s.servicePorts[%d].name", prefix, j),
				Message:  "Port name must be lowercase alphanumeric with hyphens",
				Severity: "error",
			})
		} else if portNames[port.Name] {
			errors = append(errors, ValidationError{
				Field:    fmt.Sprintf("%s.servicePorts[%d].name", prefix, j),
				Message:  fmt.Sprintf("duplicate port name: %s", port.Name),
				Severity: "error",
			})
		}
		portNames[port.Name] = true

		if port.Port < 1 || port.Port > 65535 {
			errors = append(errors, ValidationError{
				Field:    fmt.Sprintf("%s.servicePorts[%d].port", prefix, j),
				Message:  fmt.Sprintf("invalid port number: %d (must be between 1 and 65535)", port.Port),
				Severity: "error",
			})
		} else if portNumbers[port.Port] {
			errors = append(errors, ValidationError{
				Field:    fmt.Sprintf("%s.servicePorts[%d].port", prefix, j),
				Message:  fmt.Sprintf("duplicate port number: %d", port.Port),
				Severity: "error",
			})
		}
		portNumbers[port.Port] = true

		if port.TargetPort != 0 && (port.TargetPort < 1 || port.TargetPort > 65535) {
			errors = append(errors, ValidationError{
				Field:    fmt.Sprintf("%s.servicePorts[%d].targetPort", prefix, j),
				Message:  fmt.Sprintf("invalid target port number: %d (must be between 1 and 65535)", port.TargetPort),
				Severity: "error",
			})
		}
	}

	// Validate volumes if present
	if len(pod.Volumes) > 0 {
		volumeNames := make(map[string]bool)
		for j, volume := range pod.Volumes {
			errors = append(errors, v.validateVolume(volume, fmt.Sprintf("%s.volumes[%d]", prefix, j), volumeNames)...)
		}
	}

	return errors
}

// validateVolume performs validation of a volume configuration
func (v *Validator) validateVolume(volume Volume, prefix string, volumeNames map[string]bool) []ValidationError {
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
	} else if volumeNames[volume.Name] {
		errors = append(errors, ValidationError{
			Field:    prefix + ".name",
			Message:  fmt.Sprintf("duplicate volume name: %s", volume.Name),
			Severity: "error",
		})
	}
	volumeNames[volume.Name] = true

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

// Helper functions for validation

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

func isValidPodName(name string) bool {
	return isValidName(name)
}

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
	for _, r := range size {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func extractPodReferences(value string) []string {
	re := regexp.MustCompile(`([a-z][a-z0-9-]*).pod`)
	matches := re.FindAllStringSubmatch(value, -1)
	refs := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			refs = append(refs, match[1])
		}
	}
	return refs
}

func findClosestPodName(ref string, podNames map[string]bool) string {
	minDist := len(ref)
	var closest string
	for name := range podNames {
		dist := levenshteinDistance(ref, name)
		if dist < minDist {
			minDist = dist
			closest = name
		}
	}
	if minDist <= len(ref)/2 {
		return closest
	}
	return ""
}

func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func getAvailablePods(podNames map[string]bool) []string {
	pods := make([]string, 0, len(podNames))
	for name := range podNames {
		pods = append(pods, name)
	}
	sort.Strings(pods)
	return pods
}
