package types

import "time"

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

// DeployRequest represents a deployment request
type DeployRequest struct {
	YAML          string `json:"yaml"`
	ApplicationID string `json:"application_id"`
}

// Deployment represents a deployment
type Deployment struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"applicationId"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// DeploymentInfo represents detailed deployment information
type DeploymentInfo struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"applicationId"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	Namespace     string    `json:"namespace"`
	Config        string    `json:"config"`
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

// Config represents the client configuration
type Config struct {
	Token string `json:"token"`
}
