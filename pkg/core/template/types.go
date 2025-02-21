// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"fmt"
)

// NexlayerYAML represents the root structure of a Nexlayer deployment template
type NexlayerYAML struct {
	Application Application `yaml:"application" validate:"required"`
}

// Application represents the application configuration
type Application struct {
	Name          string         `yaml:"name" validate:"required,name"`
	URL           string         `yaml:"url,omitempty" validate:"omitempty,url"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
	Pods          []Pod          `yaml:"pods" validate:"required,min=1,dive"`
}

// RegistryLogin represents container registry authentication
type RegistryLogin struct {
	Registry            string `yaml:"registry" validate:"required"`
	Username            string `yaml:"username" validate:"required"`
	PersonalAccessToken string `yaml:"personalAccessToken" validate:"required"`
}

// Pod represents a container in the deployment
type Pod struct {
	Name         string            `yaml:"name" validate:"required,name"`
	Type         string            `yaml:"type" validate:"required,oneof=frontend backend nextjs react node python go raw"`
	Path         string            `yaml:"path,omitempty" validate:"omitempty,startswith=/"`
	Image        string            `yaml:"image" validate:"required,image"`
	Command      string            `yaml:"command,omitempty"`
	Entrypoint   string            `yaml:"entrypoint,omitempty"`
	ServicePorts []ServicePort     `yaml:"servicePorts" validate:"required,min=1,dive"`
	Vars         []EnvVar          `yaml:"vars,omitempty" validate:"omitempty,dive"`
	Volumes      []Volume          `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets      []Secret          `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

// ServicePort represents a service port configuration
type ServicePort struct {
	Name       string `yaml:"name,omitempty" validate:"omitempty"`
	Port       int    `yaml:"port" validate:"required,min=1,max=65535"`
	TargetPort int    `yaml:"targetPort,omitempty" validate:"omitempty,min=1,max=65535"`
	Protocol   string `yaml:"protocol,omitempty" validate:"omitempty,oneof=TCP UDP"`
}

// UnmarshalYAML implements custom unmarshaling for ServicePort to support both formats
func (sp *ServicePort) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try simple format (just port number)
	var port int
	if err := unmarshal(&port); err == nil {
		sp.Port = port
		sp.TargetPort = port
		sp.Name = fmt.Sprintf("port-%d", port)
		sp.Protocol = ProtocolTCP
		return nil
	}

	// Try full format
	type fullServicePort ServicePort
	var full fullServicePort
	if err := unmarshal(&full); err != nil {
		return err
	}

	sp.Name = full.Name
	sp.Port = full.Port
	sp.TargetPort = full.TargetPort
	if sp.TargetPort == 0 {
		sp.TargetPort = sp.Port
	}
	sp.Protocol = full.Protocol
	if sp.Protocol == "" {
		sp.Protocol = ProtocolTCP
	}
	if sp.Name == "" {
		sp.Name = fmt.Sprintf("port-%d", sp.Port)
	}

	return nil
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name     string `yaml:"name" validate:"required,name"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
	Size     string `yaml:"size,omitempty" validate:"omitempty,volumesize"`
	Type     string `yaml:"type,omitempty" validate:"omitempty,oneof=persistent ephemeral"`
	ReadOnly bool   `yaml:"readOnly,omitempty"`
}

// Secret represents a secret configuration
type Secret struct {
	Name     string `yaml:"name" validate:"required,name"`
	Data     string `yaml:"data" validate:"required"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
	FileName string `yaml:"fileName" validate:"required"`
}
