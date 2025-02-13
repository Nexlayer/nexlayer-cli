// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package types defines the data structures used in the Nexlayer API.
// Generated from OpenAPI schema version 1.0.0

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
	Registry            string `yaml:"registry" validate:"required,hostname"`
	Username            string `yaml:"username" validate:"required"`
	PersonalAccessToken string `yaml:"personalAccessToken" validate:"required"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name      string `yaml:"name" validate:"required,alphanum"`
	Size      string `yaml:"size" validate:"required,volumesize"`
	MountPath string `yaml:"mountPath" validate:"required,startswith=/"`
}

// Secret represents encrypted credentials or config files
type Secret struct {
	Name      string `yaml:"name" validate:"required,alphanum"`
	Data      string `yaml:"data" validate:"required"`
	MountPath string `yaml:"mountPath" validate:"required,startswith=/"`
	FileName  string `yaml:"fileName" validate:"required,filename"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}

// Pod represents a pod configuration in the template
type Pod struct {
	Name         string   `yaml:"name" validate:"required,podname"`
	Path         string   `yaml:"path,omitempty" validate:"omitempty,startswith=/"`
	Image        string   `yaml:"image" validate:"required,image"`
	Volumes      []Volume `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets      []Secret `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Vars         []EnvVar `yaml:"vars,omitempty" validate:"omitempty,dive"`
	ServicePorts []int    `yaml:"servicePorts" validate:"required,dive,gt=0,lt=65536"`
}

// NexlayerYAML represents the structure of a Nexlayer deployment template
// NexlayerYAML represents the structure of a Nexlayer deployment template
// Application represents a Nexlayer application configuration
type Application struct {
	Name          string         `yaml:"name" validate:"required"`
	URL           string         `yaml:"url,omitempty" validate:"omitempty,url"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
	Pods          []Pod          `yaml:"pods" validate:"required,dive"`
}

// NexlayerYAML represents the structure of a Nexlayer deployment template
type NexlayerYAML struct {
	Application Application `yaml:"application" validate:"required"`
}

// StartDeploymentResponse represents the response from starting a deployment
// StartDeploymentResponse represents the response from the /startUserDeployment endpoint
type StartDeploymentResponse struct {
	// Deployment status message (e.g., "Deployment started successfully")
	Message string `json:"message" validate:"required" example:"Deployment started successfully"`

	// Generated namespace for the deployment (e.g., "fantastic-fox")
	Namespace string `json:"namespace" validate:"required" example:"fantastic-fox"`

	// Application URL where the deployment can be accessed
	// Format: https://{namespace}-{app-name}.alpha.nexlayer.ai
	URL string `json:"url" validate:"required,url" example:"https://fantastic-fox-my-mern-app.alpha.nexlayer.ai"`
}

// SaveCustomDomainRequest represents the request body for saving a custom domain
type SaveCustomDomainRequest struct {
	Domain string `json:"domain" example:"mydomain.com"`
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
