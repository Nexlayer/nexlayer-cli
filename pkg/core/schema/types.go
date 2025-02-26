// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package template provides centralized schema management and template processing for Nexlayer YAML configurations.
package schema

import (
	"fmt"
	"time"
)

// NexlayerYAML represents a complete Nexlayer application template
type NexlayerYAML struct {
	Application Application `yaml:"application" validate:"required"`
}

// Application represents a Nexlayer application configuration
type Application struct {
	Name          string         `yaml:"name" validate:"required,podname"`
	URL           string         `yaml:"url,omitempty" validate:"omitempty,url"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
	Pods          []Pod          `yaml:"pods" validate:"required,min=1,dive"`
}

// RegistryLogin represents private registry authentication
type RegistryLogin struct {
	Registry            string `yaml:"registry" validate:"required"`
	Username            string `yaml:"username" validate:"required"`
	PersonalAccessToken string `yaml:"personalAccessToken" validate:"required"`
}

// Pod represents a container in the deployment
type Pod struct {
	Name         string            `yaml:"name" validate:"required,podname"`
	Type         string            `yaml:"type,omitempty" validate:"omitempty"`
	Path         string            `yaml:"path,omitempty" validate:"omitempty,startswith=/"`
	Image        string            `yaml:"image" validate:"required,image"`
	Entrypoint   string            `yaml:"entrypoint,omitempty" validate:"omitempty"`
	Command      string            `yaml:"command,omitempty" validate:"omitempty"`
	Volumes      []Volume          `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets      []Secret          `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Vars         []EnvVar          `yaml:"vars,omitempty" validate:"omitempty,dive"`
	ServicePorts []ServicePort     `yaml:"servicePorts" validate:"required,min=1,dive"`
	Annotations  map[string]string `yaml:"annotations,omitempty" validate:"omitempty"`
}

// UnmarshalYAML implements custom unmarshaling for Pod to handle environment variables
func (p *Pod) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Define a temporary type without the custom unmarshaling
	type tempPod Pod

	// First, unmarshal into a map to check if vars is a map
	var podMap map[string]interface{}
	if err := unmarshal(&podMap); err != nil {
		return err
	}

	// Create a temporary pod to unmarshal into
	var tmp tempPod

	// Unmarshal everything except vars
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	// Copy all fields from tmp to p
	*p = Pod(tmp)

	// Check if vars exists and is a map
	if varsInterface, ok := podMap["vars"]; ok {
		if varsMap, ok := varsInterface.(map[string]interface{}); ok {
			// Convert map to EnvVar slice
			for k, v := range varsMap {
				strValue := fmt.Sprintf("%v", v)
				p.Vars = append(p.Vars, EnvVar{
					Key:   k,
					Value: strValue,
				})
			}
			return nil
		}
	}

	return nil
}

// ServicePort represents a service port configuration
type ServicePort struct {
	Name       string `yaml:"name" validate:"required"`
	Port       int    `yaml:"port" validate:"required,min=1,max=65535"`
	TargetPort int    `yaml:"targetPort" validate:"required,min=1,max=65535"`
	Protocol   string `yaml:"protocol,omitempty" validate:"omitempty,oneof=TCP UDP"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name     string `yaml:"name" validate:"required,alphanum"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
	Size     string `yaml:"size,omitempty" validate:"omitempty,volumesize"`
	Type     string `yaml:"type,omitempty" validate:"omitempty"`
	ReadOnly bool   `yaml:"readOnly,omitempty"`
}

// Secret represents encrypted credentials or config files
type Secret struct {
	Name     string `yaml:"name" validate:"required,alphanum"`
	Data     string `yaml:"data" validate:"required"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
	FileName string `yaml:"fileName" validate:"required,filename"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}

// ProjectType represents the detected type of project
type ProjectType string

const (
	// Base project types
	TypeUnknown   ProjectType = "unknown"
	TypeNextjs    ProjectType = "nextjs"
	TypeReact     ProjectType = "react"
	TypeNode      ProjectType = "node"
	TypePython    ProjectType = "python"
	TypeGo        ProjectType = "go"
	TypeDockerRaw ProjectType = "docker"

	// AI/LLM project types
	TypeLangchainNextjs ProjectType = "langchain-nextjs"
	TypeOpenAINode      ProjectType = "openai-node"
	TypeLlamaPython     ProjectType = "llama-py"

	// Full-stack project types
	TypeMERN ProjectType = "mern" // MongoDB + Express + React + Node.js
	TypePERN ProjectType = "pern" // PostgreSQL + Express + React + Node.js
	TypeMEAN ProjectType = "mean" // MongoDB + Express + Angular + Node.js
)

// ProjectInfo contains detected information about a project
type ProjectInfo struct {
	Type         ProjectType       `json:"type"`
	Name         string            `json:"name"`
	Version      string            `json:"version,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Scripts      map[string]string `json:"scripts,omitempty"`
	Port         int               `json:"port,omitempty"`
	HasDocker    bool              `json:"has_docker"`
	LLMProvider  string            `json:"llm_provider,omitempty"` // AI-powered IDE
	LLMModel     string            `json:"llm_model,omitempty"`    // LLM Model being used
	ImageTag     string            `json:"image_tag,omitempty"`    // Docker image tag
}

// DeploymentStatus represents the current state of a deployment
type DeploymentStatus struct {
	Namespace    string      `json:"namespace"`
	TemplateID   string      `json:"templateId"`
	TemplateName string      `json:"templateName"`
	Status       string      `json:"status"`
	URL          string      `json:"url"`
	CustomDomain string      `json:"customDomain"`
	Version      string      `json:"version"`
	CreatedAt    time.Time   `json:"createdAt"`
	LastUpdated  time.Time   `json:"lastUpdated"`
	PodStatuses  []PodStatus `json:"podStatuses"`
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
