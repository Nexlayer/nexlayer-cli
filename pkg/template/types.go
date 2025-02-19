// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

// NexlayerYAML represents the root structure of a Nexlayer YAML file
type NexlayerYAML struct {
	Application Application `yaml:"application"`
}

// Application represents the application configuration
type Application struct {
	Name          string         `yaml:"name"`
	URL           string         `yaml:"url,omitempty"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
	Pods          []Pod          `yaml:"pods"`
}

// RegistryLogin represents container registry authentication
type RegistryLogin struct {
	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Pod represents a pod configuration
type Pod struct {
	Name         string            `yaml:"name"`
	Type         string            `yaml:"type"`
	Path         string            `yaml:"path,omitempty"`
	Image        string            `yaml:"image,omitempty"`
	ServicePorts []ServicePort     `yaml:"servicePorts,omitempty"`
	Vars         []EnvVar          `yaml:"vars,omitempty"`
	Volumes      []Volume          `yaml:"volumes,omitempty"`
	Secrets      []Secret          `yaml:"secrets,omitempty"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

// ServicePort represents a service port configuration
type ServicePort struct {
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
	Protocol   string `yaml:"protocol,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Volume represents a volume configuration
type Volume struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Size     string `yaml:"size,omitempty"`
	Type     string `yaml:"type,omitempty"`
	ReadOnly bool   `yaml:"readOnly,omitempty"`
}

// Secret represents a secret configuration
type Secret struct {
	Name  string `yaml:"name"`
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
