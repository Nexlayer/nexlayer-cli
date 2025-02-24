// Package schema provides centralized schema management for Nexlayer YAML configurations.
package schema

import "time"

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
