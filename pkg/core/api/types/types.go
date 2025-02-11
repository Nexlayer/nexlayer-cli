// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package types

import "time"

// Info represents Nexlayer installation information
type Info struct {
	Version string `json:"version"`
	Build   string `json:"build"`
	API     string `json:"api"`
}

// Version represents version information
type Version struct {
	CLI string `json:"cli"`
	API string `json:"api"`
}

// App represents a Nexlayer application
type App struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	LastDeployAt time.Time `json:"lastDeployAt,omitempty"`
}

// CreateAppRequest represents a request to create a new application
type CreateAppRequest struct {
	Name string `json:"name"`
}

// RegistryLogin represents private registry authentication
type RegistryLogin struct {
	Registry           string `yaml:"registry"`
	Username           string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name      string `yaml:"name"`
	Size      string `yaml:"size"`
	MountPath string `yaml:"mountPath"`
}

// Secret represents encrypted credentials or config files
type Secret struct {
	Name      string `yaml:"name"`
	Data      string `yaml:"data"`
	MountPath string `yaml:"mountPath"`
	FileName  string `yaml:"fileName"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Pod represents a pod configuration in the template
type Pod struct {
	Name         string    `yaml:"name"`
	Path         string    `yaml:"path,omitempty"`
	Image        string    `yaml:"image"`
	Volumes      []Volume  `yaml:"volumes,omitempty"`
	Secrets      []Secret  `yaml:"secrets,omitempty"`
	Vars         []EnvVar  `yaml:"vars,omitempty"`
	ServicePorts []int     `yaml:"servicePorts,omitempty"`
}

// NexlayerYAML represents the structure of a Nexlayer deployment template
type NexlayerYAML struct {
	Application struct {
		Name         string       `yaml:"name"`
		URL          string       `yaml:"url,omitempty"`
		RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
		Pods         []Pod        `yaml:"pods"`
	} `yaml:"application"`
}

// StartDeploymentResponse represents the response from starting a deployment
type StartDeploymentResponse struct {
	Message   string `json:"message" example:"Deployment started successfully"`
	Namespace string `json:"namespace" example:"fantastic-fox"`
	URL       string `json:"url" example:"https://fantastic-fox-my-mern-app.alpha.nexlayer.ai"`
}

// SaveCustomDomainResponse represents the response from saving a custom domain
type SaveCustomDomainResponse struct {
	Message string `json:"message" example:"Custom domain saved successfully"`
}

// Deployment represents a deployment in the system
type Deployment struct {
	Namespace        string `json:"namespace" example:"ecstatic-frog"`
	TemplateID       string `json:"templateID" example:"0001"`
	TemplateName     string `json:"templateName" example:"K-d chat"`
	DeploymentStatus string `json:"deploymentStatus" example:"running"`
}

// DeploymentInfo represents detailed information about a deployment
type DeploymentInfo struct {
	Namespace        string `json:"namespace" example:"ecstatic-frog"`
	TemplateID       string `json:"templateID" example:"0001"`
	TemplateName     string `json:"templateName" example:"K-d chat"`
	DeploymentStatus string `json:"deploymentStatus" example:"running"`
}

// GetDeploymentsResponse represents the response from getting all deployments
type GetDeploymentsResponse struct {
	Deployments []Deployment `json:"deployments"`
	Pagination  *Pagination  `json:"pagination,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// GetDeploymentInfoResponse represents the response from getting deployment info
type GetDeploymentInfoResponse struct {
	Deployment DeploymentInfo `json:"deployment"`
}

// DeployResponse represents a deployment response
type DeployResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// Domain represents a custom domain
type Domain struct {
	Domain        string    `json:"domain"`
	ApplicationID string    `json:"applicationId"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	SSLEnabled    bool      `json:"sslEnabled"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// FeedbackRequest represents a user feedback submission
type FeedbackRequest struct {
	Text string `json:"text" example:"Sample text"`
}

// Config represents the client configuration
type Config struct {
	Token string `json:"token"`
}

// AppConfig represents the detected application configuration
type AppConfig struct {
	Name             string     `json:"name"`
	Type             string     `json:"type"`
	Container        *Container `json:"container"`
	Resources        *Resources `json:"resources"`
	Env              []string   `json:"env"`
	HasExistingImage bool       `json:"hasExistingImage"`
}

// Container represents container configuration
type Container struct {
	Command       string `json:"command,omitempty"`
	UseDockerfile bool   `json:"useDockerfile,omitempty"`
	Ports         []int  `json:"ports"`
}

// Resources represents compute resources
type Resources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
