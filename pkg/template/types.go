// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

// NexlayerYAML represents a complete Nexlayer application template
type NexlayerYAML struct {
	Application Application `yaml:"application" validate:"required"`
}

// Application represents a Nexlayer application configuration
type Application struct {
	Name          string         `yaml:"name" validate:"required,appname"` // lowercase, alphanumeric, '-', or '.'
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
	Name        string            `yaml:"name" validate:"required,podname"` // lowercase, alphanumeric, '-', or '.'
	Type        PodType           `yaml:"type,omitempty" validate:"omitempty,podtype"`
	Path        string            `yaml:"path,omitempty" validate:"omitempty,startswith=/"`
	Image       string            `yaml:"image" validate:"required,image"` // Full image URL including registry and tag
	Volumes     []Volume          `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets     []Secret          `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Vars        []EnvVar          `yaml:"vars,omitempty" validate:"omitempty,dive"`
	Ports       []Port            `yaml:"ports" validate:"required,dive"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// Port represents a port configuration
type Port struct {
	ContainerPort int    `yaml:"containerPort" validate:"required,gt=0,lt=65536"`
	ServicePort   int    `yaml:"servicePort" validate:"required,gt=0,lt=65536"`
	Name          string `yaml:"name" validate:"required,portname"` // web, api, db, etc
}

// Volume represents a persistent storage volume
type Volume struct {
	Name      string `yaml:"name" validate:"required,volumename"` // lowercase, alphanumeric, '-'
	Size      string `yaml:"size" validate:"required,volumesize"` // e.g., "1Gi", "500Mi"
	MountPath string `yaml:"mountPath" validate:"required,startswith=/"`
}

// Secret represents encrypted credentials or config files
type Secret struct {
	Name      string `yaml:"name" validate:"required,secretname"` // lowercase, alphanumeric, '-'
	Data      string `yaml:"data" validate:"required"`            // Raw or Base64-encoded secret content
	MountPath string `yaml:"mountPath" validate:"required,startswith=/"`
	FileName  string `yaml:"fileName" validate:"required,filename"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}

// PodType represents the type of a pod
type PodType string

const (
	// Frontend pod types
	Frontend PodType = "frontend"
	React    PodType = "react"
	NextJS   PodType = "nextjs"
	Vue      PodType = "vue"

	// Backend pod types
	Backend PodType = "backend"
	Express PodType = "express"
	Django  PodType = "django"
	FastAPI PodType = "fastapi"
	Node    PodType = "node"
	Python  PodType = "python"
	Golang  PodType = "golang"
	Java    PodType = "java"

	// Database pod types
	Database   PodType = "database"
	MongoDB    PodType = "mongodb"
	Postgres   PodType = "postgres"
	Redis      PodType = "redis"
	MySQL      PodType = "mysql"
	Clickhouse PodType = "clickhouse"

	// Message Queue types
	RabbitMQ PodType = "rabbitmq"
	Kafka    PodType = "kafka"

	// Storage types
	Minio   PodType = "minio"
	Elastic PodType = "elasticsearch"

	// Web Server types
	Nginx   PodType = "nginx"
	Traefik PodType = "traefik"

	// AI/ML pod types
	LLM      PodType = "llm"
	Ollama   PodType = "ollama"
	HFModel  PodType = "huggingface"
	VertexAI PodType = "vertexai"
	Jupyter  PodType = "jupyter"
)
