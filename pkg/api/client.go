package api

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
)

const (
	defaultBaseURL = "https://api.nexlayer.io/v1"
	defaultTimeout = 30 * time.Second
)

// Client represents the Nexlayer API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{},
	}

	// Skip SSL verification in test mode
	if os.Getenv("NEXLAYER_TEST_MODE") == "true" {
		tr.TLSClientConfig.InsecureSkipVerify = true
	}

	httpClient := &http.Client{
		Timeout:   defaultTimeout,
		Transport: tr,
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

// SetToken sets the authentication token
func (c *Client) SetToken(token string) {
	c.token = token
}

// CreateApplication creates a new application
func (c *Client) CreateApplication(ctx context.Context, req *types.CreateAppRequest) (*types.App, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, "/apps", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}
	defer resp.Body.Close()

	var app types.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &app, nil
}

// ListApplications lists all applications
func (c *Client) ListApplications(ctx context.Context) ([]types.App, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/apps", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list applications: %w", err)
	}
	defer resp.Body.Close()

	var apps []types.App
	if err := json.NewDecoder(resp.Body).Decode(&apps); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return apps, nil
}

// GetAppInfo gets information about an application
func (c *Client) GetAppInfo(ctx context.Context, appID string) (*types.App, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%s", appID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get app info: %w", err)
	}
	defer resp.Body.Close()

	var app types.App
	if err := json.NewDecoder(resp.Body).Decode(&app); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &app, nil
}

// StartUserDeployment starts a deployment for a user
func (c *Client) StartUserDeployment(ctx context.Context, appID string, req *types.DeployRequest) (*types.Deployment, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/apps/%s/deployments", appID), bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to start deployment: %w", err)
	}
	defer resp.Body.Close()

	var deployment types.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployment); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &deployment, nil
}

// GetDeployments gets all deployments for an application
func (c *Client) GetDeployments(ctx context.Context, appID string) ([]types.Deployment, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%s/deployments", appID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}
	defer resp.Body.Close()

	var deployments []types.Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return deployments, nil
}

// GetDeploymentInfo gets detailed information about a deployment
func (c *Client) GetDeploymentInfo(ctx context.Context, namespace, appID string) (*types.DeploymentInfo, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/deployments/%s/%s", namespace, appID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment info: %w", err)
	}
	defer resp.Body.Close()

	var info types.DeploymentInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &info, nil
}

// SaveCustomDomain saves a custom domain for an application
func (c *Client) SaveCustomDomain(ctx context.Context, appID string, domain string) error {
	data, err := json.Marshal(map[string]string{"domain": domain})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/apps/%s/domains", appID), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to save custom domain: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

// GetDomains gets all custom domains for an application
func (c *Client) GetDomains(appID string) ([]types.Domain, error) {
	ctx := context.Background()
	resp, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/apps/%s/domains", appID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get domains: %w", err)
	}
	defer resp.Body.Close()

	var domains []types.Domain
	if err := json.NewDecoder(resp.Body).Decode(&domains); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return domains, nil
}

// RemoveDomain removes a custom domain from an application
func (c *Client) RemoveDomain(appID, domain string) error {
	ctx := context.Background()
	resp, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/apps/%s/domains/%s", appID, domain), nil)
	if err != nil {
		return fmt.Errorf("failed to remove domain: %w", err)
	}
	defer resp.Body.Close()
	return nil
}

// doRequest makes an HTTP request with proper authentication and error handling
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		var errResp types.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("request failed with status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("request failed: %s", errResp.Message)
	}

	return resp, nil
}
