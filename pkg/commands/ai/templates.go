package ai

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Template represents a Nexlayer YAML template.
type Template struct {
	Application Application `yaml:"application"`
}

// Application represents the top-level application configuration.
type Application struct {
	Name          string         `yaml:"name"`
	URL           string         `yaml:"url,omitempty"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
	Pods          []Pod          `yaml:"pods"`
	Entrypoint    string         `yaml:"entrypoint,omitempty"`
	Command       string         `yaml:"command,omitempty"`
}

// RegistryLogin represents private registry authentication.
type RegistryLogin struct {
	Registry            string `yaml:"registry"`
	Username            string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// Pod represents a single pod configuration.
type Pod struct {
	Name         string        `yaml:"name"`
	Path         string        `yaml:"path,omitempty"`
	Image        string        `yaml:"image"`
	Volumes      []Volume      `yaml:"volumes,omitempty"`
	Secrets      []Secret      `yaml:"secrets,omitempty"`
	Vars         []EnvVar      `yaml:"vars,omitempty"`
	ServicePorts []ServicePort `yaml:"servicePorts"`
}

// EnvVar represents an environment variable for a pod.
type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Volume represents persistent storage configuration.
type Volume struct {
	Name      string `yaml:"name"`
	Size      string `yaml:"size"`
	MountPath string `yaml:"mountPath"`
}

// Secret represents secret configuration.
type Secret struct {
	Name      string `yaml:"name"`
	Data      string `yaml:"data"`
	MountPath string `yaml:"mountPath"`
	FileName  string `yaml:"fileName"`
}

// ServicePort represents port configuration.
// It supports both a simple integer shorthand and a detailed mapping.
type ServicePort struct {
	Port       int `yaml:"port,omitempty"`
	TargetPort int `yaml:"targetPort,omitempty"`
}

// UnmarshalYAML implements custom unmarshaling for ServicePort so that it supports
// both an integer shorthand (e.g., "3000") and a detailed configuration (e.g., { port: 80, targetPort: 8080 }).
func (sp *ServicePort) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML node is a scalar, treat it as a simple port.
	if value.Kind == yaml.ScalarNode {
		var port int
		if err := value.Decode(&port); err != nil {
			return fmt.Errorf("failed to decode service port: %w", err)
		}
		sp.Port = port
		sp.TargetPort = port
		return nil
	}

	// Otherwise, decode as a mapping.
	type plain ServicePort
	var plainSP plain
	if err := value.Decode(&plainSP); err != nil {
		return fmt.Errorf("failed to decode detailed service port: %w", err)
	}
	*sp = ServicePort(plainSP)
	return nil
}
