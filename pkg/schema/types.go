// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

// NexlayerYAML represents a complete Nexlayer application template
type NexlayerYAML struct {
	Application Application `yaml:"application" validate:"required"`
}

// Application represents a Nexlayer application configuration
type Application struct {
	Name          string         `yaml:"name" validate:"required,appname"`
	URL           string         `yaml:"url,omitempty" validate:"omitempty,url"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
	Pods          []Pod          `yaml:"pods" validate:"required,min=1,dive"`
}

// RegistryLogin represents private registry authentication
type RegistryLogin struct {
	Registry            string `yaml:"registry" validate:"required,hostname"`
	Username            string `yaml:"username" validate:"required"`
	PersonalAccessToken string `yaml:"personalAccessToken" validate:"required"`
}

// Pod represents a container in the deployment
type Pod struct {
	Name         string            `yaml:"name" validate:"required,podname"`
	Type         string            `yaml:"type,omitempty" validate:"omitempty,podtype"`
	Path         string            `yaml:"path,omitempty" validate:"omitempty,startswith=/"`
	Image        string            `yaml:"image" validate:"required,image"`
	Volumes      []Volume          `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets      []Secret          `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Vars         []EnvVar          `yaml:"vars,omitempty" validate:"omitempty,dive"`
	ServicePorts []int             `yaml:"servicePorts" validate:"required,dive,gt=0,lt=65536"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name      string `yaml:"name" validate:"required,volumename"`
	Size      string `yaml:"size" validate:"required,volumesize"`
	MountPath string `yaml:"mountPath" validate:"required,startswith=/"`
}

// Secret represents encrypted credentials or config files
type Secret struct {
	Name      string `yaml:"name" validate:"required,secretname"`
	Data      string `yaml:"data" validate:"required"`
	MountPath string `yaml:"mountPath" validate:"required,startswith=/"`
	FileName  string `yaml:"fileName" validate:"required,filename"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}
