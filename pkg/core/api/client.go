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
)

// Client represents an API client for interacting with the Nexlayer API
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://api.nexlayer.com"
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

// StartDeployment starts a new deployment using a YAML file
func (c *Client) StartDeployment(ctx context.Context, appID string, yamlFile string) (*types.StartDeploymentResponse, error) {
	// Read YAML file
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/startUserDeployment/%s", c.baseURL, appID)
	resp, err := c.post(ctx, url, data)
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

	if !result.Success {
		return fmt.Errorf("failed to save custom domain: %s", result.Message)
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

	var result []types.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// GetDeploymentInfo gets detailed information about a deployment
func (c *Client) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error) {
	url := fmt.Sprintf("%s/getDeploymentInfo/%s/%s", c.baseURL, namespace, appID)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}

	var result types.DeploymentInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
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
