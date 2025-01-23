package types

import "time"

// DeploymentResponse represents the response from the deployment API
type DeploymentResponse struct {
	Message   string `json:"message"`
	Namespace string `json:"namespace"`
	URL       string `json:"url"`
}

// GetDeploymentsResponse represents the response from the get deployments API
type GetDeploymentsResponse struct {
	Deployments []DeploymentInfo `json:"deployments"`
}

// DeploymentInfo represents information about a deployment
type DeploymentInfo struct {
	Namespace        string `json:"namespace"`
	TemplateID       string `json:"templateID"`
	TemplateName     string `json:"templateName"`
	DeploymentStatus string `json:"deploymentStatus"`
}

// SaveCustomDomainRequest represents the request to save a custom domain
type SaveCustomDomainRequest struct {
	Domain string `json:"domain"`
}

// SaveCustomDomainResponse represents the response from saving a custom domain
type SaveCustomDomainResponse struct {
	Message string `json:"message"`
}

// Application represents a Nexlayer application
type Application struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateApplicationResponse represents the response from creating an application
type CreateApplicationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// StartDeploymentRequest represents the request to start a deployment
type StartDeploymentRequest struct {
	YAML string `json:"yaml"`
}

// Config represents the client configuration
type Config struct {
	Token string `json:"token"`
}
