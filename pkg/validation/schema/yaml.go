// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

// NexlayerYAML represents the structure of a Nexlayer deployment template.
type NexlayerYAML struct {
	Application Application `yaml:"application" validate:"required"`
}

// Application represents the application-level configuration.
type Application struct {
	Name          string         `yaml:"name" validate:"required,podname"`
	URL           string         `yaml:"url,omitempty" validate:"omitempty,url"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
	Pods          []Pod          `yaml:"pods" validate:"required,dive"`
}

// RegistryLogin contains authentication details for private registries.
type RegistryLogin struct {
	Registry            string `yaml:"registry" validate:"required"`
	Username            string `yaml:"username" validate:"required"`
	PersonalAccessToken string `yaml:"personalAccessToken" validate:"required"`
}

// Pod represents a pod configuration.
type Pod struct {
	Name         string    `yaml:"name" validate:"required,podname"`
	Path         string    `yaml:"path,omitempty"`
	Image        string    `yaml:"image" validate:"required,image"`
	Volumes      []Volume  `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets      []Secret  `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Vars         []VarPair `yaml:"vars,omitempty" validate:"omitempty,dive"`
	ServicePorts []int     `yaml:"servicePorts,omitempty" validate:"omitempty,dive,gt=0,lt=65536"`
}

// Volume represents a storage volume configuration.
type Volume struct {
	Name      string `yaml:"name" validate:"required,filename"`
	Size      string `yaml:"size" validate:"required,volumesize"`
	MountPath string `yaml:"mountPath" validate:"required"`
}

// Secret represents a secret file configuration.
type Secret struct {
	Name      string `yaml:"name" validate:"required,filename"`
	Data      string `yaml:"data" validate:"required"`
	MountPath string `yaml:"mountPath" validate:"required"`
	FileName  string `yaml:"fileName" validate:"required,filename"`
}

// VarPair represents an environment variable key-value pair.
type VarPair struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}
