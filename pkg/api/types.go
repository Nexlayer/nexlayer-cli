// Formatted with gofmt -s
package api

// StartDeploymentResponse represents the response from startUserDeployment endpoint
type StartDeploymentResponse struct {
	Message   string `json:"message"`
	Namespace string `json:"namespace"`
	URL       string `json:"url"`
}

// GetDeploymentsResponse represents the response from getDeployments endpoint
type GetDeploymentsResponse struct {
	Deployments []DeploymentInfo `json:"deployments"`
}

// GetDeploymentInfoResponse represents the response from getDeploymentInfo endpoint
type GetDeploymentInfoResponse struct {
	Deployment DeploymentInfo `json:"deployment"`
}

// DeploymentInfo represents information about a deployment
type DeploymentInfo struct {
	Namespace        string `json:"namespace"`
	TemplateID       string `json:"templateID"`
	TemplateName     string `json:"templateName"`
	DeploymentStatus string `json:"deploymentStatus"`
}

// SaveCustomDomainRequest represents the request body for saveCustomDomain endpoint
type SaveCustomDomainRequest struct {
	Domain string `json:"domain"`
}

// SaveCustomDomainResponse represents the response from saveCustomDomain endpoint
type SaveCustomDomainResponse struct {
	Message string `json:"message"`
}
