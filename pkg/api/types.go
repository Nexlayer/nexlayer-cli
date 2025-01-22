package api

// DeploymentResponse represents the response from starting a deployment
type DeploymentResponse struct {
	Message   string `json:"message"`
	URL       string `json:"url"`
	Namespace string `json:"namespace"`
}

// CustomDomainResponse represents the response from saving a custom domain
type CustomDomainResponse struct {
	Message string `json:"message"`
}

// DeploymentInfo represents information about a deployment
type DeploymentInfo struct {
	Namespace        string `json:"namespace"`
	ApplicationID    string `json:"applicationId"`
	DeploymentStatus string `json:"deploymentStatus"`
}

// DeploymentsResponse represents the response from getting all deployments
type DeploymentsResponse struct {
	Deployments []DeploymentInfo `json:"deployments"`
}

// DeploymentInfoResponse represents the response from getting deployment info
type DeploymentInfoResponse struct {
	Deployment DeploymentInfo `json:"deployment"`
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

// SaveCustomDomainRequest represents the request to save a custom domain
type SaveCustomDomainRequest struct {
	Domain string `json:"domain"`
}
