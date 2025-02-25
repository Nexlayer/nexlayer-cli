// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
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
			Suggestions: []string{
				"Add 'name' field under 'application' in nexlayer.yaml",
				"Run 'nexlayer init' to generate a valid configuration",
			},
		})
	} else if !isValidName(v.config.Application.Name) {
		v.errors = append(v.errors, ValidationError{
			Field:   "application.name",
			Message: "application name must follow Nexlayer platform naming conventions",
			Suggestions: []string{
				"Use lowercase letters, numbers, and hyphens",
				"Must start with a letter",
				"Example: my-app-v1, web-service, api-backend",
			},
		})
	}

	// Validate URL if provided
	if v.config.Application.URL != "" && !isValidURL(v.config.Application.URL) {
		v.errors = append(v.errors, ValidationError{
			Field:   "application.url",
			Message: "invalid URL format",
			Suggestions: []string{
				"Use a valid domain name (e.g., example.com)",
				"Only alphanumeric characters, dots, and hyphens are allowed",
			},
		})
	}
}

// validateRegistryLogin ensures registry login is correctly configured if present
func (v *Validator) validateRegistryLogin() {
	rl := v.config.Application.RegistryLogin
	if rl != nil {
		if rl.Registry == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.registry",
				Message: "registry hostname is required when registryLogin is present",
				Suggestions: []string{
					"Add 'registry' field with the hostname",
					"Example: docker.io, ghcr.io",
				},
			})
		} else if !isValidRegistryHost(rl.Registry) {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.registry",
				Message: "invalid registry hostname format",
				Suggestions: []string{
					"Use a valid hostname (e.g., docker.io, ghcr.io)",
					"Only alphanumeric characters, dots, and hyphens are allowed",
				},
			})
		}

		if rl.Username == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.username",
				Message: "registry username is required when registryLogin is present",
			})
		}

		if rl.Password == "" {
			v.errors = append(v.errors, ValidationError{
				Field:   "application.registryLogin.password",
				Message: "registry password is required when registryLogin is present",
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
			Suggestions: []string{
				"Add at least one pod configuration",
				"Run 'nexlayer init' to generate a valid configuration",
			},
		})
		return
	}

	// Check for duplicate pod names
	podNames := make(map[string]bool)
	for i, pod := range v.config.Application.Pods {
		if podNames[pod.Name] {
			v.errors = append(v.errors, ValidationError{
				Field:   fmt.Sprintf("pods[%d].name", i),
				Message: fmt.Sprintf("duplicate pod name: %s", pod.Name),
				Suggestions: []string{
					"Each pod must have a unique name",
					fmt.Sprintf("Rename one of the pods with name '%s'", pod.Name),
				},
			})
		}
		podNames[pod.Name] = true

		// Convert PodYAML to Pod for validation
		validationPod := template.Pod{
			Name:         pod.Name,
			Image:        pod.Image,
			ServicePorts: make([]template.ServicePort, 0),
			Vars:         make([]template.EnvVar, 0),
			Volumes:      make([]template.Volume, 0),
		}
		v.validatePod(i, validationPod)
	}

	// Note: Environment variable validation is handled by the template package's own validation
}

// validatePod validates a single pod configuration
func (v *Validator) validatePod(index int, pod template.Pod) {
	// Validate required fields
	if pod.Name == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].name", index),
			Message: "pod name is required",
		})
	} else if !isValidPodName(pod.Name) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].name", index),
			Message: fmt.Sprintf("invalid pod name: %s", pod.Name),
			Suggestions: []string{
				"Pod names must start with a lowercase letter",
				"Use only lowercase letters, numbers, and hyphens",
				"Example: web-server, api-v1, db-postgres",
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

	// Validate service ports
	if len(pod.ServicePorts) == 0 {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].servicePorts", index),
			Message: "at least one service port is required",
		})
	} else {
		// Check for duplicate port names and numbers
		portNames := make(map[string]bool)
		portNumbers := make(map[int]bool)
		for j, port := range pod.ServicePorts {
			if port.Name == "" {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].name", index, j),
					Message: "port name is required",
					Suggestions: []string{
						"Use descriptive names like 'http', 'api', or 'metrics'",
					},
				})
			} else if !isValidName(port.Name) {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].name", index, j),
					Message: "port name must be lowercase alphanumeric with hyphens",
				})
			} else if portNames[port.Name] {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].name", index, j),
					Message: fmt.Sprintf("duplicate port name: %s", port.Name),
				})
			}
			portNames[port.Name] = true

			if port.Port < 1 || port.Port > 65535 {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].port", index, j),
					Message: fmt.Sprintf("invalid port number: %d (must be between 1 and 65535)", port.Port),
				})
			} else if portNumbers[port.Port] {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].port", index, j),
					Message: fmt.Sprintf("duplicate port number: %d", port.Port),
				})
			}
			portNumbers[port.Port] = true

			if port.TargetPort != 0 && (port.TargetPort < 1 || port.TargetPort > 65535) {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].targetPort", index, j),
					Message: fmt.Sprintf("invalid target port number: %d (must be between 1 and 65535)", port.TargetPort),
				})
			}

			if port.Protocol != "" && !isValidProtocol(port.Protocol) {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].servicePorts[%d].protocol", index, j),
					Message: fmt.Sprintf("invalid protocol: %s", port.Protocol),
					Suggestions: []string{
						"Valid protocols: TCP, UDP, SCTP",
					},
				})
			}
		}
	}

	// Validate volumes if present
	if len(pod.Volumes) > 0 {
		volumeNames := make(map[string]bool)
		for j, volume := range pod.Volumes {
			v.validateVolume(index, j, volume, volumeNames)
		}
	}

	// Validate environment variables
	if len(pod.Vars) > 0 {
		envVarNames := make(map[string]bool)
		for _, env := range pod.Vars {
			if envVarNames[env.Key] {
				v.errors = append(v.errors, ValidationError{
					Field:   fmt.Sprintf("pods[%d].vars[%s]", index, env.Key),
					Message: fmt.Sprintf("duplicate environment variable: %s", env.Key),
				})
			}
			envVarNames[env.Key] = true
		}
	}
}

// validateVolume validates a volume configuration
func (v *Validator) validateVolume(podIndex, volumeIndex int, volume template.Volume, volumeNames map[string]bool) {
	if volume.Name == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].name", podIndex, volumeIndex),
			Message: "volume name is required",
		})
	} else if !isValidName(volume.Name) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].name", podIndex, volumeIndex),
			Message: "volume name must be lowercase alphanumeric with hyphens",
		})
	} else if volumeNames[volume.Name] {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].name", podIndex, volumeIndex),
			Message: fmt.Sprintf("duplicate volume name: %s", volume.Name),
		})
	}
	volumeNames[volume.Name] = true

	if volume.Path == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].path", podIndex, volumeIndex),
			Message: "volume path is required",
		})
	} else if !strings.HasPrefix(volume.Path, "/") {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].path", podIndex, volumeIndex),
			Message: fmt.Sprintf("volume path must start with '/': %s", volume.Path),
			Suggestions: []string{
				fmt.Sprintf("Change to '/%s'", strings.TrimPrefix(volume.Path, "/")),
			},
		})
	}

	if volume.Size == "" {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].size", podIndex, volumeIndex),
			Message: "volume size is required",
			Suggestions: []string{
				"Specify size with units (e.g., 1Gi, 500Mi)",
			},
		})
	} else if !isValidVolumeSize(volume.Size) {
		v.errors = append(v.errors, ValidationError{
			Field:   fmt.Sprintf("pods[%d].volumes[%d].size", podIndex, volumeIndex),
			Message: fmt.Sprintf("invalid volume size format: %s", volume.Size),
			Suggestions: []string{
				"Use format: <number><unit>",
				"Valid units: Ki, Mi, Gi, Ti",
				"Example: 1Gi, 500Mi",
			},
		})
	}
}

// Helper functions for validation

func isValidName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

func isValidPodName(name string) bool {
	return isValidName(name)
}

func isValidURL(url string) bool {
	// More comprehensive URL validation
	if strings.ContainsAny(url, " \t\n\r") {
		return false
	}

	// Must contain at least one dot and no consecutive dots
	if !strings.Contains(url, ".") || strings.Contains(url, "..") {
		return false
	}

	// Split into parts
	parts := strings.Split(url, ".")
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
		// Each part must be alphanumeric or hyphen
		for _, c := range part {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
				return false
			}
		}
		// Cannot start or end with hyphen
		if part[0] == '-' || part[len(part)-1] == '-' {
			return false
		}
	}
	return true
}

func isValidRegistryHost(host string) bool {
	// More comprehensive registry host validation
	if strings.ContainsAny(host, " \t\n\r") {
		return false
	}

	// Must contain at least one dot and no consecutive dots
	if !strings.Contains(host, ".") || strings.Contains(host, "..") {
		return false
	}

	// Split into parts
	parts := strings.Split(host, ".")
	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		// Each part must be alphanumeric or hyphen
		for _, c := range part {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') {
				return false
			}
		}
		// Cannot start or end with hyphen
		if part[0] == '-' || part[len(part)-1] == '-' {
			return false
		}
	}
	return true
}

func isValidProtocol(protocol string) bool {
	switch protocol {
	case "TCP", "UDP", "SCTP":
		return true
	default:
		return false
	}
}

func isValidVolumeSize(size string) bool {
	// More comprehensive volume size validation
	re := regexp.MustCompile(`^([1-9][0-9]*(?:\.[0-9]+)?|0\.[0-9]*[1-9][0-9]*)[KMGT]i$`)
	if !re.MatchString(size) {
		return false
	}

	// Extract numeric part and unit
	numStr := size[:len(size)-2]
	unit := size[len(size)-2:]

	// Parse numeric part
	num, err := parseFloat(numStr)
	if err != nil {
		return false
	}

	// Check size limits based on unit
	switch unit[0] {
	case 'K':
		return num <= 1024*1024 // Max 1024Ki
	case 'M':
		return num <= 1024*1024 // Max 1024Mi
	case 'G':
		return num <= 1024 // Max 1024Gi
	case 'T':
		return num <= 64 // Max 64Ti
	default:
		return false
	}
}

// parseFloat parses a string to float64, handling scientific notation
func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
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

// formatErrors formats all validation errors into a single error message
func (v *Validator) formatErrors() error {
	if len(v.errors) == 0 {
		return nil
	}

	// Pre-allocate buffer with estimated size
	var errMsg strings.Builder
	errMsg.Grow(len(v.errors) * 100) // Estimate 100 bytes per error

	errMsg.WriteString("Validation failed:\n")

	// Use map to group errors by category
	categories := map[string][]ValidationError{
		"application": make([]ValidationError, 0),
		"pods":        make([]ValidationError, 0),
		"volumes":     make([]ValidationError, 0),
		"vars":        make([]ValidationError, 0),
	}

	// Group errors by category
	for _, err := range v.errors {
		category := strings.Split(err.Field, ".")[0]
		categories[category] = append(categories[category], err)
	}

	// Print errors by category in a consistent order
	for _, category := range []string{"application", "pods", "volumes", "vars"} {
		if errors := categories[category]; len(errors) > 0 {
			errMsg.WriteString(fmt.Sprintf("\n%s:\n", strings.Title(category)))
			for _, err := range errors {
				errMsg.WriteString(fmt.Sprintf("  - %s: %s\n", err.Field, err.Message))
				for _, suggestion := range err.Suggestions {
					errMsg.WriteString(fmt.Sprintf("    • %s\n", suggestion))
				}
			}
		}
	}

	return fmt.Errorf(errMsg.String())
}

// ValidatePod validates a single pod configuration
func ValidatePod(pod template.Pod) error {
	validator := NewValidator(&template.NexlayerYAML{
		Application: template.ApplicationYAML{
			Name: "temp",
			Pods: []template.PodYAML{{
				Name:  pod.Name,
				Image: pod.Image,
			}},
		},
	})
	return validator.Validate()
}
