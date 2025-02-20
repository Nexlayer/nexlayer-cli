// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

// NexlayerYAML represents a complete Nexlayer application template
type NexlayerYAML struct {
	Application Application `yaml:"application"`
}

// Application represents a Nexlayer application configuration
type Application struct {
	Name          string         `yaml:"name"`
	URL           string         `yaml:"url,omitempty"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
	Pods          []Pod          `yaml:"pods"`
}

// RegistryLogin represents private registry authentication
type RegistryLogin struct {
	Registry            string `yaml:"registry"`
	Username            string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// Pod represents a container in the deployment
type Pod struct {
	Name         string            `yaml:"name"`
	Type         string            `yaml:"type,omitempty"`
	Path         string            `yaml:"path,omitempty"`
	Image        string            `yaml:"image"`
	Entrypoint   string            `yaml:"entrypoint,omitempty"`
	Command      string            `yaml:"command,omitempty"`
	Volumes      []Volume          `yaml:"volumes,omitempty"`
	Secrets      []Secret          `yaml:"secrets,omitempty"`
	Vars         []EnvVar          `yaml:"vars,omitempty"`
	ServicePorts []ServicePort     `yaml:"servicePorts"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

// ServicePort represents a service port configuration
type ServicePort struct {
	Name       string `yaml:"name"`
	Port       int    `yaml:"port"`
	TargetPort int    `yaml:"targetPort"`
	Protocol   string `yaml:"protocol,omitempty"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	Size     string `yaml:"size,omitempty"`
	Type     string `yaml:"type,omitempty"`
	ReadOnly bool   `yaml:"readOnly,omitempty"`
}

// Secret represents encrypted credentials or config files
type Secret struct {
	Name     string `yaml:"name"`
	Data     string `yaml:"data"`
	Path     string `yaml:"path"`
	FileName string `yaml:"fileName"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
