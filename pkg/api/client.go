package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	configDir  = ".nexlayer"
	configFile = "config"
)

// Client represents the Nexlayer API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// Config represents the configuration file structure
type Config struct {
	Token string `json:"token"`
}

// NewClient creates a new Nexlayer API client
func NewClient(baseURL string) (*Client, error) {
	client := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}

	// Load token from config file
	if err := client.loadToken(); err != nil {
		return nil, fmt.Errorf("failed to load token: %w", err)
	}

	return client, nil
}

// loadToken loads the Bearer token from the config file
func (c *Client) loadToken() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configDir, configFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w. Please run 'nexlayer-cli login' first", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.Token == "" {
		return fmt.Errorf("no token found in config file. Please run 'nexlayer-cli login' first")
	}

	c.token = config.Token
	return nil
}

// doRequest makes an HTTP request with proper authentication
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.baseURL, path), reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle response status
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: please run 'nexlayer-cli login' to reauthenticate")
	}
	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, fmt.Errorf("request failed: %s (code: %s)", errResp.Message, errResp.Code)
	}

	return respBody, nil
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
