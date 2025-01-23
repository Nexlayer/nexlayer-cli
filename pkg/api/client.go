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
	"github.com/Nexlayer/nexlayer-cli/pkg/auth"
)

const (
	configDir  = ".nexlayer"
	configFile = "config"
)

// Client represents the Nexlayer API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	mockMode   bool
}

// Config represents the configuration file structure
type Config struct {
	Token string `json:"token"`
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
		mockMode: true, // Enable mock mode for testing
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
		return []byte(`{"id": "dep-123", "status": "running"}`), nil
	case path == "/saveCustomDomain/app-123" && method == "POST":
		return []byte(`{"domain": "example.com", "status": "active"}`), nil
	case path == "/getDeployments/app-123" && method == "GET":
		return []byte(`{"deployments": [{"id": "dep-123", "status": "running"}]}`), nil
	case path == "/getDeploymentInfo/default/app-123" && method == "GET":
		return []byte(`{"id": "dep-123", "status": "running", "logs": "example logs"}`), nil
	default:
		return nil, fmt.Errorf("mock: unknown endpoint %s %s", method, path)
	}
}

// StartUserDeployment starts a new deployment
func (c *Client) StartUserDeployment(applicationID, yamlPath string) (*DeploymentResponse, error) {
	// Read YAML file
	yaml, err := os.ReadFile(yamlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Create request body
	reqBody := StartDeploymentRequest{
		YAML: string(yaml),
	}

	// Make request
	respBody, err := c.doRequest(http.MethodPost, fmt.Sprintf("/startUserDeployment/%s", applicationID), reqBody)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp DeploymentResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// SaveCustomDomain configures a custom domain for an application
func (c *Client) SaveCustomDomain(applicationID, domain string) (*CustomDomainResponse, error) {
	// Create request body
	reqBody := SaveCustomDomainRequest{
		Domain: domain,
	}

	// Make request
	respBody, err := c.doRequest(http.MethodPost, fmt.Sprintf("/saveCustomDomain/%s", applicationID), reqBody)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp CustomDomainResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetDeployments gets all deployments for an application
func (c *Client) GetDeployments(applicationID string) (*DeploymentsResponse, error) {
	// Make request
	respBody, err := c.doRequest(http.MethodGet, fmt.Sprintf("/getDeployments/%s", applicationID), nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp DeploymentsResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// GetDeploymentInfo gets detailed information about a deployment
func (c *Client) GetDeploymentInfo(namespace, applicationID string) (*DeploymentInfoResponse, error) {
	// Make request
	respBody, err := c.doRequest(http.MethodGet, fmt.Sprintf("/getDeploymentInfo/%s/%s", namespace, applicationID), nil)
	if err != nil {
		return nil, err
	}

	// Parse response
	var resp DeploymentInfoResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
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

// CreateApplication creates a new application
func (c *Client) CreateApplication(name string) (*CreateApplicationResponse, error) {
	payload := map[string]string{
		"name": name,
	}

	data, err := c.doRequest("POST", "/api/v1/applications", payload)
	if err != nil {
		return nil, err
	}

	var resp CreateApplicationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &resp, nil
}

// ListApplications returns all applications for the authenticated user
func (c *Client) ListApplications() ([]Application, error) {
	data, err := c.doRequest("GET", "/api/v1/applications", nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Applications []Application `json:"applications"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return resp.Applications, nil
}
