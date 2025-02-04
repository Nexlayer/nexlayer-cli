package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
	"gopkg.in/yaml.v3"
)

// APIClient defines the interface for interacting with the Nexlayer API.
// The Nexlayer API allows deploying full-stack AI-powered applications using templates.
// Each template defines the application stack, pods, and environment variables,
// abstracting away the complexity of deployment.
type APIClient interface {
	// StartDeployment starts a new deployment using a YAML template file.
	// The template must follow the Nexlayer template structure (see Templates docs)
	// containing application.template section and pods array.
	// If appID is empty, creates a new deployment from the template.
	StartDeployment(ctx context.Context, appID string, configPath string) (*types.StartDeploymentResponse, error)

	// SaveCustomDomain saves a custom domain for an application.
	// The domain will be associated with the specified application ID.
	SaveCustomDomain(ctx context.Context, appID string, domain string) error

	// GetDeployments retrieves all deployments for a given application.
	GetDeployments(ctx context.Context, appID string) ([]types.Deployment, error)

	// GetDeploymentInfo retrieves detailed information about a specific deployment.
	// Requires both the namespace and application ID to uniquely identify the deployment.
	GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error)
}

// Client represents an API client for interacting with the Nexlayer API.
// The Nexlayer API enables rapid deployment of full-stack AI-powered applications
// by providing a simple template-based interface that abstracts away deployment complexity.
type Client struct {
	baseURL    string      // Base URL of the Nexlayer API
	httpClient *http.Client // HTTP client for making API requests
	token      string      // Authentication token for API requests
}

// NewClient creates a new Nexlayer API client.
// If baseURL is empty, defaults to the staging environment at app.staging.nexlayer.io.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://app.staging.nexlayer.io"
	}
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// SetToken sets the authentication token for the client
func (c *Client) SetToken(token string) {
	c.token = token
}

// StartDeployment starts a new deployment using a YAML template file.
// The template must follow the Nexlayer template structure, including:
// - Required fields: application.template.name, deploymentName, registryLogin
// - Pods array with valid pod types (frontend, backend, database, nginx, llm)
// - Each pod must have: type, name, tag, and optional vars
// If appID is empty, creates a new deployment from the template.
func (c *Client) StartDeployment(ctx context.Context, appID string, yamlFile string) (*types.StartDeploymentResponse, error) {
	// Read and validate YAML file
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Validate template structure (basic check)
	var template types.NexlayerYAML
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("invalid template format: %w", err)
	}

	// Validate required fields
	if template.Application.Template.Name == "" ||
		template.Application.Template.DeploymentName == "" ||
		template.Application.Template.RegistryLogin.Registry == "" {
		return nil, fmt.Errorf("missing required fields in template")
	}

	// Make request
	url := fmt.Sprintf("%s/startUserDeployment", c.baseURL)
	if appID != "" {
		url = fmt.Sprintf("%s/%s", url, appID)
	}

	// Send as text/x-yaml content type
	resp, err := c.postYAML(ctx, url, data)
	if err != nil {
		return nil, err
	}

	var result types.StartDeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// SaveCustomDomain saves a custom domain for an application
func (c *Client) SaveCustomDomain(ctx context.Context, appID string, domain string) error {
	url := fmt.Sprintf("%s/saveCustomDomain/%s", c.baseURL, appID)
	data := map[string]string{"domain": domain}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.post(ctx, url, body)
	if err != nil {
		return err
	}

	var result types.SaveCustomDomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// GetDeployments gets all deployments for an application
func (c *Client) GetDeployments(ctx context.Context, appID string) ([]types.Deployment, error) {
	url := fmt.Sprintf("%s/getDeployments/%s", c.baseURL, appID)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}

	var result struct {
		Deployments []types.Deployment `json:"deployments"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Deployments, nil
}

// GetDeploymentInfo gets detailed information about a deployment
func (c *Client) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error) {
	url := fmt.Sprintf("%s/getDeploymentInfo/%s/%s", c.baseURL, namespace, appID)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}

	var result struct {
		Deployment types.DeploymentInfo `json:"deployment"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Deployment, nil
}

// Helper methods for making HTTP requests
func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// postYAML sends a POST request with YAML content type.
// The Nexlayer API expects deployment templates to be sent as text/x-yaml.
func (c *Client) postYAML(ctx context.Context, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/x-yaml")
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
