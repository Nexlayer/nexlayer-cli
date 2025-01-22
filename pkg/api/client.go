package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	defaultTimeout = 30 * time.Second
	baseURL        = "https://app.nexlayer.io" // Production URL from OpenAPI spec
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

// StartDeployment initiates a new deployment
func (c *Client) StartDeployment(applicationID string, yamlContent []byte) (*StartDeploymentResponse, error) {
	url := fmt.Sprintf("%s/startUserDeployment/%s", c.baseURL, applicationID)

	fmt.Printf("DEBUG: Making request to %s\n", url)
	fmt.Printf("DEBUG: YAML content:\n%s\n", string(yamlContent))

	req, err := http.NewRequest("POST", url, bytes.NewReader(yamlContent))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers exactly as curl does
	req.Header.Set("Content-Type", "text/x-yaml")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "curl/8.7.1")

	// Debug headers
	fmt.Printf("DEBUG: Request headers:\n")
	for k, v := range req.Header {
		fmt.Printf("  %s: %s\n", k, v[0])
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("deployment failed (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var deployResp StartDeploymentResponse
	if err := json.Unmarshal(body, &deployResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &deployResp, nil
}

// GetDeployments retrieves all deployments for a given application
func (c *Client) GetDeployments(applicationID string) (*GetDeploymentsResponse, error) {
	url := fmt.Sprintf("%s/getDeployments/%s", c.baseURL, applicationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get deployments (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var deploymentsResp GetDeploymentsResponse
	if err := json.Unmarshal(body, &deploymentsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &deploymentsResp, nil
}

// GetDeploymentInfo retrieves detailed information about a specific deployment
func (c *Client) GetDeploymentInfo(namespace, applicationID string) (*GetDeploymentInfoResponse, error) {
	url := fmt.Sprintf("%s/getDeploymentInfo/%s/%s", c.baseURL, namespace, applicationID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get deployment info (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var deploymentInfoResp GetDeploymentInfoResponse
	if err := json.Unmarshal(body, &deploymentInfoResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &deploymentInfoResp, nil
}

// SaveCustomDomain saves a custom domain for a deployment
func (c *Client) SaveCustomDomain(applicationID string, domain string) (*SaveCustomDomainResponse, error) {
	url := fmt.Sprintf("%s/saveCustomDomain/%s", c.baseURL, applicationID)

	reqBody := SaveCustomDomainRequest{
		Domain: domain,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to save custom domain (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var saveResp SaveCustomDomainResponse
	if err := json.Unmarshal(body, &saveResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &saveResp, nil
}
