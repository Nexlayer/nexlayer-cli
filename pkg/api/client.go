package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/Nexlayer/nexlayer-cli/pkg/auth"
	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
)

// Client represents the Nexlayer API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	mockMode   bool
}

// NewClient creates a new Nexlayer API client
func NewClient(baseURL string) (*Client, error) {
	// Skip TLS verification for staging environment
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second,
		},
		mockMode: false, // Disable mock mode for real testing
	}, nil
}

// doRequest makes an HTTP request with proper authentication
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	if c.mockMode {
		return c.getMockResponse(method, path, body)
	}

	var buf io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		buf = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, c.baseURL+path, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token, err := auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %s", string(data))
	}

	return data, nil
}

// getMockResponse returns mock data for testing
func (c *Client) getMockResponse(method, path string, body interface{}) ([]byte, error) {
	switch {
	case path == "/api/v1/applications" && method == "GET":
		return []byte(`{
			"applications": [
				{
					"id": "app-123",
					"name": "test-app",
					"created_at": "2025-01-23T01:22:27-05:00"
				},
				{
					"id": "app-124",
					"name": "todo-mern-app",
					"created_at": "2025-01-23T01:23:27-05:00"
				}
			]
		}`), nil
	case path == "/api/v1/applications" && method == "POST":
		name := ""
		if m, ok := body.(map[string]string); ok {
			name = m["name"]
		}
		return []byte(fmt.Sprintf(`{
			"id": "app-%s",
			"name": "%s"
		}`, time.Now().Format("20060102150405"), name)), nil
	case path == "/startUserDeployment/app-123" && method == "POST":
		return []byte(`{"message": "Deployment started", "namespace": "default", "url": "https://example.com"}`), nil
	case path == "/saveCustomDomain/app-123" && method == "POST":
		return []byte(`{"message": "Custom domain saved successfully"}`), nil
	case path == "/getDeployments/app-123" && method == "GET":
		return []byte(`{
			"deployments": [
				{
					"namespace": "default",
					"templateID": "123",
					"templateName": "MERN Todo",
					"deploymentStatus": "running"
				}
			]
		}`), nil
	case path == "/getDeploymentInfo/default/app-123" && method == "GET":
		return []byte(`{
			"deployment": {
				"namespace": "default",
				"templateID": "123",
				"templateName": "MERN Todo",
				"deploymentStatus": "running"
			}
		}`), nil
	default:
		return nil, fmt.Errorf("mock: unknown endpoint %s %s", method, path)
	}
}

// StartUserDeployment starts a new deployment
func (c *Client) StartUserDeployment(applicationID, yamlPath string) (*types.DeploymentResponse, error) {
	// Read YAML file
	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/startUserDeployment/%s", c.baseURL, applicationID), bytes.NewBuffer(yamlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	token, err := auth.GetToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "text/x-yaml")
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deployment failed: %s", string(body))
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected content type: %s. Response: %s", contentType, string(body))
	}

	// Parse response
	var deployResp types.DeploymentResponse
	if err := json.Unmarshal(body, &deployResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w. Response body: %s", err, string(body))
	}

	return &deployResp, nil
}

// SaveCustomDomain configures a custom domain for an application
func (c *Client) SaveCustomDomain(applicationID, domain string) (*types.SaveCustomDomainResponse, error) {
	// Create request body
	reqBody := types.SaveCustomDomainRequest{
		Domain: domain,
	}

	// Make request
	respBody, err := c.doRequest(http.MethodPost, fmt.Sprintf("/saveCustomDomain/%s", applicationID), reqBody)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp types.SaveCustomDomainResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetDeployments gets all deployments for an application
func (c *Client) GetDeployments(applicationID string) (*types.GetDeploymentsResponse, error) {
	// Make request
	respBody, err := c.doRequest(http.MethodGet, fmt.Sprintf("/getDeployments/%s", applicationID), nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp types.GetDeploymentsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetDeploymentInfo gets detailed information about a deployment
func (c *Client) GetDeploymentInfo(namespace, applicationID string) (*types.DeploymentInfo, error) {
	// Make request
	respBody, err := c.doRequest(http.MethodGet, fmt.Sprintf("/getDeploymentInfo/%s/%s", namespace, applicationID), nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp struct {
		Deployment types.DeploymentInfo `json:"deployment"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp.Deployment, nil
}

// CreateApplication creates a new application
func (c *Client) CreateApplication(name string) (*types.CreateApplicationResponse, error) {
	payload := map[string]string{
		"name": name,
	}

	data, err := c.doRequest("POST", "/api/v1/applications", payload)
	if err != nil {
		return nil, err
	}

	var resp types.CreateApplicationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// ListApplications returns all applications for the authenticated user
func (c *Client) ListApplications() ([]types.Application, error) {
	data, err := c.doRequest("GET", "/api/v1/applications", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Applications []types.Application `json:"applications"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Applications, nil
}
