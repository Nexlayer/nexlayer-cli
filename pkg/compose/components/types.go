// types.go
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package components

// Pod represents a deployment pod configuration as defined in Nexlayer YAML templates v2.
type Pod struct {
	Name         string    `yaml:"name"`              // Pod name (lowercase alphanumeric, '-', '.')
	Path         string    `yaml:"path,omitempty"`     // Route path for frontend (e.g., "/")
	Image        string    `yaml:"image"`             // Docker image path (supports <% REGISTRY %>)
	Volumes      []Volume  `yaml:"volumes,omitempty"`  // List of persistent volumes
	Secrets      []Secret  `yaml:"secrets,omitempty"`  // List of secrets
	Vars         []EnvVar  `yaml:"vars,omitempty"`    // Environment variables
	ServicePorts []int     `yaml:"servicePorts"`      // List of ports to expose
	Command      []string  `yaml:"command,omitempty"` // Optional command to run
}

// ComponentDetector is an interface for detecting and configuring components.
type ComponentDetector interface {
	DetectAndConfigure(pod Pod) (DetectedComponent, error)
}

// ComponentConfig holds the default configuration for a specific component type.
type ComponentConfig struct {
	Image        string       `yaml:"image"`                  // Docker image path (supports <% REGISTRY %>)
	ServicePorts []int       `yaml:"servicePorts"`           // List of ports to expose
	Environment  []EnvVar     `yaml:"environment,omitempty"`  // Environment variables
	Command      []string     `yaml:"command,omitempty"`      // Default command
	HealthCheck  *Healthcheck `yaml:"healthCheck,omitempty"`  // Health check configuration
	Volumes      []Volume     `yaml:"volumes,omitempty"`      // Persistent volumes
	Secrets      []Secret     `yaml:"secrets,omitempty"`      // Secrets configuration
	Path         string       `yaml:"path,omitempty"`         // Route path for frontend
}

// Port represents a port mapping (deprecated in v2, use ServicePorts instead).
type Port struct {
	Container int    // Port inside the container.
	Host      int    // Port on the host.
	Protocol  string // Protocol (e.g., tcp or udp).
	Name      string // Descriptive port name.
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Key   string `yaml:"key"`   // Variable name
	Value string `yaml:"value"` // Variable value (supports template variables)
}

// Healthcheck represents a Docker healthcheck configuration.
type Healthcheck struct {
	Command             []string `yaml:"command"`                       // Healthcheck command.
	Interval            string   `yaml:"interval,omitempty"`            // Time between health checks.
	Timeout             string   `yaml:"timeout,omitempty"`             // Timeout for a health check.
	Retries             int      `yaml:"retries,omitempty"`             // Number of retries before failing.
	InitialDelaySeconds int      `yaml:"initialDelaySeconds,omitempty"` // Initial delay before starting health checks.
	PeriodSeconds       int      `yaml:"periodSeconds,omitempty"`       // Frequency of health checks.
}

// Volume represents a persistent volume configuration.
type Volume struct {
	Name      string `yaml:"name"`      // Volume name
	Size      string `yaml:"size"`      // Volume size (e.g., "1Gi")
	MountPath string `yaml:"mountPath"` // Container mount path
}

// Secret represents a secret configuration.
type Secret struct {
	Name      string `yaml:"name"`      // Secret name
	Data      string `yaml:"data"`      // Base64/raw secret value
	MountPath string `yaml:"mountPath"` // Secret mount directory
	FileName  string `yaml:"fileName"`  // Secret file name
}

// DetectedComponent represents a component detected from the project directory.
type DetectedComponent struct {
	Type          string          // Detected component type (e.g., postgres, redis).
	Category      string          // Component category (e.g., database, cache).
	Config        ComponentConfig // The default configuration for the component.
	Dependencies  []string        // List of dependencies.
	IsStateful    bool            // Whether this component maintains state.
	RequiresSetup bool            // Whether the component requires initialization.
}
