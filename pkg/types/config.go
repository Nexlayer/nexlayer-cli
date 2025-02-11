// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package types

// Config represents the application configuration
type Config struct {
	Application Application `yaml:"application"`
}

// Application represents the application configuration
type Application struct {
	Name          string       `yaml:"name"`
	URL           string       `yaml:"url,omitempty"`
	RegistryLogin *RegistryAuth `yaml:"registryLogin,omitempty"`
	Pods          []Pod        `yaml:"pods"`
}

// RegistryAuth represents registry authentication configuration
type RegistryAuth struct {
	Registry            string `yaml:"registry"`
	Username            string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// Pod represents a pod configuration
type Pod struct {
	Name         string    `yaml:"name"`
	Path         string    `yaml:"path,omitempty"`
	Image        string    `yaml:"image"`
	Volumes      []Volume  `yaml:"volumes,omitempty"`
	Secrets      []Secret  `yaml:"secrets,omitempty"`
	Vars         []EnvVar  `yaml:"vars,omitempty"`
	ServicePorts []int     `yaml:"servicePorts"`
	Command      []string  `yaml:"command,omitempty"`
}

// Volume represents a persistent volume configuration
type Volume struct {
	Name      string `yaml:"name"`
	Size      string `yaml:"size"`
	MountPath string `yaml:"mountPath"`
}

// Secret represents a secret configuration
type Secret struct {
	Name      string `yaml:"name"`
	Data      string `yaml:"data"`
	MountPath string `yaml:"mountPath"`
	FileName  string `yaml:"fileName"`
}

// EnvVar represents a key-value pair for environment variables
type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
