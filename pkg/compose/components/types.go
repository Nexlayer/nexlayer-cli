// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package components

// Pod represents a deployment pod configuration
type Pod struct {
	Type       string   `yaml:"type"`
	Name       string   `yaml:"name"`
	Image      string   `yaml:"image"` // Full image URL including registry and tag
	Command    []string `yaml:"command,omitempty"`
	Vars       []EnvVar `yaml:"vars,omitempty"`
	Ports      []Port   `yaml:"ports,omitempty"`
	ExposeOn80 bool     `yaml:"exposeOn80"`
}

// ComponentDetector interface for detecting and configuring components
type ComponentDetector interface {
	DetectAndConfigure(pod Pod) (DetectedComponent, error)
}

// ComponentConfig holds the configuration for a specific component type
type ComponentConfig struct {
	Image        string       // Default Docker image
	Ports        []Port       // Default ports to expose
	Environment  []EnvVar     // Default environment variables
	Command      []string     // Default command to run
	Healthcheck  *Healthcheck // Health check configuration
	Dependencies []string     // Other components this depends on
	Volumes      []Volume     // Default volume mounts
	RequiredVars []string     // Required environment variables
}

// Port represents a port mapping
type Port struct {
	Container int    // Container port
	Host      int    // Host port
	Protocol  string // Protocol (tcp/udp)
	Name      string // Port name/description
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key      string
	Value    string
	Required bool
}

// Healthcheck represents a Docker healthcheck configuration
type Healthcheck struct {
	Test     []string
	Interval string
	Timeout  string
	Retries  int
}

// Volume represents a volume mount
type Volume struct {
	Source     string // Host path or volume name
	Target     string // Container path
	Type       string // bind/volume/tmpfs
	Persistent bool   // Whether this volume should persist
}

// DetectedComponent represents a detected component with its configuration
type DetectedComponent struct {
	Type          string          // Component type (e.g., postgres, redis)
	Category      string          // Component category (e.g., database, cache)
	Config        ComponentConfig // Component configuration
	Dependencies  []string        // Dependencies on other components
	IsStateful    bool            // Whether this component maintains state
	RequiresSetup bool            // Whether this component needs initialization
}
