// types.go
// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.

package components

// Pod represents a deployment pod configuration as defined in Nexlayer YAML templates.
type Pod struct {
	Type       string   `yaml:"type"`              // Component type (e.g., frontend, backend, database, nginx, llm).
	Name       string   `yaml:"name"`              // Unique name for the pod.
	Image      string   `yaml:"image"`             // Full Docker image URL including registry and tag.
	Command    []string `yaml:"command,omitempty"` // Optional command to run in the container.
	Vars       []EnvVar `yaml:"vars,omitempty"`    // Environment variables.
	Ports      []Port   `yaml:"ports,omitempty"`   // Port mappings.
	ExposeOn80 bool     `yaml:"exposeOn80"`        // Whether to expose this pod on port 80.
}

// ComponentDetector is an interface for detecting and configuring components.
type ComponentDetector interface {
	DetectAndConfigure(pod Pod) (DetectedComponent, error)
}

// ComponentConfig holds the default configuration for a specific component type.
type ComponentConfig struct {
	Image        string       `yaml:"image"`                  // Default Docker image.
	Ports        []Port       `yaml:"ports,omitempty"`        // Default ports to expose.
	Environment  []EnvVar     `yaml:"environment,omitempty"`  // Default environment variables.
	Command      []string     `yaml:"command,omitempty"`      // Default command.
	HealthCheck  *Healthcheck `yaml:"healthCheck,omitempty"`  // Health check configuration.
	Dependencies []string     `yaml:"dependencies,omitempty"` // Component dependencies.
	Volumes      []Volume     `yaml:"volumes,omitempty"`      // Volume mounts.
	RequiredVars []string     `yaml:"requiredVars,omitempty"` // Required environment variables.
}

// Port represents a port mapping.
type Port struct {
	Container int    // Port inside the container.
	Host      int    // Port on the host.
	Protocol  string // Protocol (e.g., tcp or udp).
	Name      string // Descriptive port name.
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Key      string // Variable name.
	Value    string // Variable value.
	Required bool   // Whether this variable is required.
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

// Volume represents a volume mount configuration.
type Volume struct {
	Source     string // Host path or volume name.
	Target     string // Container path.
	Type       string // Type of mount: bind, volume, or tmpfs.
	Persistent bool   // Whether the volume should persist.
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
