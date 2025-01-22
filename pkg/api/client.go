package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	baseURL       = "https://app.nexlayer.io" // From OpenAPI spec
)

// Client represents a Nexlayer API client
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new Nexlayer API client
func NewClient(baseURL string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip certificate verification for staging
		},
	}

	url := baseURL
	if url == "" {
		url = baseURL
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   defaultTimeout,
			Transport: tr,
		},
		baseURL: url,
	}
}

// StartUserDeploymentResponse matches OpenAPI spec
type StartUserDeploymentResponse struct {
	Message   string `json:"message"`
	Namespace string `json:"namespace"`
	URL       string `json:"url"`
}

// StartUserDeployment starts a new deployment with a YAML configuration
func (c *Client) StartUserDeployment(applicationID string, yamlContent []byte) (*StartUserDeploymentResponse, error) {
	url := fmt.Sprintf("%s/startUserDeployment/%s", c.baseURL, applicationID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(yamlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/x-yaml")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deployment failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result StartUserDeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SaveCustomDomainRequest matches OpenAPI spec
type SaveCustomDomainRequest struct {
	Domain string `json:"domain"`
}

// SaveCustomDomainResponse matches OpenAPI spec
type SaveCustomDomainResponse struct {
	Message string `json:"message"`
}

// SaveCustomDomain saves a custom domain for a deployment
func (c *Client) SaveCustomDomain(applicationID string, domain string) (*SaveCustomDomainResponse, error) {
	url := fmt.Sprintf("%s/saveCustomDomain/%s", c.baseURL, applicationID)

	payload := SaveCustomDomainRequest{
		Domain: domain,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to set custom domain: %s", string(body))
	}

	var result SaveCustomDomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetDeploymentsResponse matches OpenAPI spec
type GetDeploymentsResponse struct {
	Deployments []struct {
		Namespace        string `json:"namespace"`
		TemplateID       string `json:"templateID"`
		TemplateName     string `json:"templateName"`
		DeploymentStatus string `json:"deploymentStatus"`
	} `json:"deployments"`
}

// GetDeployments retrieves all deployments for a given application
func (c *Client) GetDeployments(applicationID string) (*GetDeploymentsResponse, error) {
	url := fmt.Sprintf("%s/getDeployments/%s", c.baseURL, applicationID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get deployments: %s", string(body))
	}

	var result GetDeploymentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetDeploymentInfoResponse matches OpenAPI spec
type GetDeploymentInfoResponse struct {
	Deployment struct {
		Namespace        string `json:"namespace"`
		TemplateID       string `json:"templateID"`
		TemplateName     string `json:"templateName"`
		DeploymentStatus string `json:"deploymentStatus"`
	} `json:"deployment"`
}

// GetDeploymentInfo retrieves detailed information about a specific deployment
func (c *Client) GetDeploymentInfo(namespace, applicationID string) (*GetDeploymentInfoResponse, error) {
	url := fmt.Sprintf("%s/getDeploymentInfo/%s/%s", c.baseURL, namespace, applicationID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get deployment info: %s", string(body))
	}

	var result GetDeploymentInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetUserDeployments retrieves all deployments for a user
func (c *Client) GetUserDeployments(sessionID string) (*GetDeploymentsResponse, error) {
	url := fmt.Sprintf("%s/deployments", c.baseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add session ID to headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionID))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result GetDeploymentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ScaleDeployment scales a deployment to the specified number of replicas
func (c *Client) ScaleDeployment(namespace, sessionID string, replicas int) error {
	url := fmt.Sprintf("%s/api/v1/deployments/%s/scale", c.baseURL, namespace)

	payload := struct {
		Replicas int `json:"replicas"`
	}{
		Replicas: replicas,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionID))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to scale deployment: %s", string(body))
	}

	return nil
}

// SetCustomDomain sets a custom domain for a deployment
func (c *Client) SetCustomDomain(namespace, sessionID, domain string) error {
	url := fmt.Sprintf("%s/api/v1/deployments/%s/domain", c.baseURL, namespace)

	payload := SaveCustomDomainRequest{
		Domain: domain,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionID))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to set custom domain: %s", string(body))
	}

	return nil
}

// GetAISuggestions gets AI-powered suggestions for a query
func (c *Client) GetAISuggestions(query, sessionID string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/ai/suggest", c.baseURL)

	payload := struct {
		Query string `json:"query"`
	}{
		Query: query,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", sessionID))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get AI suggestions: %s", string(body))
	}

	var result struct {
		Suggestions []string `json:"suggestions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Suggestions, nil
}

// ServiceConfig represents the configuration for a service
type ServiceConfig struct {
	AppName     string   `json:"appName"`
	ServiceName string   `json:"serviceName"`
	EnvVars     []string `json:"envVars"`
}

// Deploy deploys a service
func (c *Client) Deploy(token, appName, serviceName string, envVars []string) error {
	url := fmt.Sprintf("%s/api/v1/deploy", c.baseURL)

	config := ServiceConfig{
		AppName:     appName,
		ServiceName: serviceName,
		EnvVars:     envVars,
	}

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return nil
}

// Configure configures a service
func (c *Client) Configure(appName, serviceName string, envVars []string) error {
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	url := fmt.Sprintf("%s/api/v1/services/configure", c.baseURL)

	config := ServiceConfig{
		AppName:     appName,
		ServiceName: serviceName,
		EnvVars:     envVars,
	}

	body, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	return nil
}

// UpdateServiceConfig updates the configuration for a service
func (c *Client) UpdateServiceConfig(appName, service string, envVars map[string]string, token string) error {
	// TODO: Implement API call to update service configuration
	return nil
}

// GetServiceConnections retrieves the service connections for an application
func (c *Client) GetServiceConnections(appName, token string) ([]ServiceConnection, error) {
	// TODO: Implement API call to get service connections
	return []ServiceConnection{
		{
			From:        "frontend",
			To:          "backend",
			Description: "HTTP/REST",
		},
		{
			From:        "backend",
			To:          "database",
			Description: "PostgreSQL",
		},
	}, nil
}

// ServiceConnection represents a connection between two services
type ServiceConnection struct {
	From        string
	To          string
	Description string
}
