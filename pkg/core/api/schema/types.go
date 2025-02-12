// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import "time"

// APIResponse is a generic response type for all API responses
type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// DeploymentResponse represents the response from starting a deployment
type DeploymentResponse struct {
	Namespace string `json:"namespace"`
	URL      string `json:"url"`
}

// Deployment represents a deployment in the system
type Deployment struct {
	Namespace    string       `json:"namespace"`
	TemplateID   string       `json:"templateId"`
	TemplateName string       `json:"templateName"`
	Status       string       `json:"status"`
	URL          string       `json:"url"`
	CustomDomain string       `json:"customDomain"`
	Version      string       `json:"version"`
	CreatedAt    time.Time    `json:"createdAt"`
	LastUpdated  time.Time    `json:"lastUpdated"`
	PodStatuses  []PodStatus  `json:"podStatuses"`
}

// PodStatus represents the status of a pod in a deployment
type PodStatus struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Ready     bool      `json:"ready"`
	Restarts  int       `json:"restarts"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

// NexlayerYAML represents the structure of a Nexlayer deployment YAML file
type NexlayerYAML struct {
	Application Application `json:"application"`
}

// RegistryLogin represents registry authentication details
type RegistryLogin struct {
	Registry           string `json:"registry"`
	Username           string `json:"username"`
	PersonalAccessToken string `json:"personalAccessToken"`
}

// Application represents the root configuration of a Nexlayer application
type Application struct {
	Name          string         `json:"name"`
	URL           string         `json:"url,omitempty"`
	RegistryLogin *RegistryLogin `json:"registryLogin,omitempty"`
	Pods          []Pod          `json:"pods"`
}

// Port represents a port configuration in a pod
type Port struct {
	ContainerPort int    `json:"containerPort"`
	ServicePort   int    `json:"servicePort"`
	Name          string `json:"name"`
}

// Pod represents a container configuration in a Nexlayer application
type Pod struct {
	Name    string    `json:"name"`
	Path    string    `json:"path,omitempty"`
	Image   string    `json:"image"`
	Volumes []Volume  `json:"volumes,omitempty"`
	Secrets []Secret  `json:"secrets,omitempty"`
	Vars    []EnvVar  `json:"vars,omitempty"`
	Ports   []Port    `json:"ports"`
}

// Volume represents a persistent volume configuration
type Volume struct {
	Name      string `json:"name"`
	Size      string `json:"size"`
	MountPath string `json:"mountPath"`
}

// Secret represents a secret configuration
type Secret struct {
	Name      string `json:"name"`
	Data      string `json:"data"`
	MountPath string `json:"mountPath"`
	FileName  string `json:"fileName"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
