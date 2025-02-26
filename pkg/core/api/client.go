// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package api provides a client for interacting with the Nexlayer API.
// Generated from OpenAPI schema version 1.0.0

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
	"strings"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

// ClientAPI is an interface that abstracts the methods required for API interactions.
type ClientAPI interface {
	StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error)
	SendFeedback(ctx context.Context, text string) error
	SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error)
	ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error)
	GetDeploymentInfo(ctx context.Context, namespace string) (*schema.APIResponse[schema.Deployment], error)
	GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error)
}

// APIClient defines the interface for interacting with the Nexlayer API.
// The Nexlayer API provides endpoints for deploying applications, managing deployments,
// sending feedback, and handling custom domains. Designed for easy integration into
// CI/CD pipelines and automated deployment systems.

type APIClient interface {
	// StartDeployment starts a new deployment using a YAML configuration file.
	// The YAML file should be provided as binary data using the 'text/x-yaml' content type.
	// Endpoint: POST /startUserDeployment
	StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error)

	// SendFeedback submits feedback to Nexlayer regarding deployment or application experience.
	// Endpoint: POST /feedback
	SendFeedback(ctx context.Context, text string) error

	// SaveCustomDomain associates a custom domain with a specific application deployment.
	// Endpoint: POST /saveCustomDomain/{applicationID}
	SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error)

	// ListDeployments retrieves all deployments.
	// Endpoint: GET /listDeployments
	ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error)

	// GetDeploymentInfo retrieves detailed information about a specific deployment.
	// Endpoint: GET /getDeploymentInfo/{namespace}
	GetDeploymentInfo(ctx context.Context, namespace string) (*schema.APIResponse[schema.Deployment], error)

	// GetLogs retrieves logs for a specific deployment.
	// If follow is true, streams logs in real-time.
	// tail specifies the number of lines to return from the end of the logs.
	GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error)
}

// APIClientForCommands interface is used for API client operations used in commands.
type APIClientForCommands interface {
	GetDeploymentInfo(ctx context.Context, namespace string) (*schema.APIResponse[schema.Deployment], error)
	GetDeployments(ctx context.Context, appID string) (*schema.APIResponse[[]schema.Deployment], error)
	SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error)
	// Add other methods as needed
}

// Client represents an API client for interacting with the Nexlayer API.
// The Nexlayer API enables rapid deployment of full-stack AI-powered applications
// by providing a simple template-based interface that abstracts away deployment complexity.
type Client struct {
	baseURL    string       // Base URL of the Nexlayer API
	httpClient *http.Client // HTTP client for making API requests
	token      string       // Authentication token for API requests
}

// Ensure Client implements APIClientForCommands
var _ APIClientForCommands = (*Client)(nil)

// Ensure Client implements ClientAPI
var _ ClientAPI = (*Client)(nil)

// handleAPIError processes API error responses and returns a formatted error
func (c *Client) handleAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var errResp schema.APIError
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	return fmt.Errorf("API error (status %d): %s", resp.StatusCode, errResp.Message)
}

// NewClient creates a new Nexlayer API client.
// If baseURL is empty, defaults to the staging environment at app.staging.nexlayer.io.
// ListDeployments retrieves all deployments.
// Endpoint: GET /listDeployments
func (c *Client) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	url := fmt.Sprintf("%s/listDeployments", c.baseURL)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp)
	}

	var result schema.APIResponse[[]schema.Deployment]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode deployments response: %w", err)
	}

	return &result, nil
}

// GetLogs retrieves logs for a specific deployment
func (c *Client) GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error) {
	url := fmt.Sprintf("%s/getDeploymentLogs/%s/%s?follow=%v&tail=%d", c.baseURL, namespace, appID, follow, tail)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	defer resp.Body.Close()

	var logs []string
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, fmt.Errorf("failed to decode logs response: %w", err)
	}

	return logs, nil
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://app.staging.nexlayer.io"
	}

	transport := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: strings.Contains(baseURL, "staging")},
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout:   120 * time.Second,
			Transport: transport,
		},
	}
}

// SetToken sets the authentication token for the client
func (c *Client) SetToken(token string) {
	c.token = token
}

// StartDeployment starts a new deployment using a YAML configuration file.
// Endpoint: POST /startUserDeployment
func (c *Client) StartDeployment(ctx context.Context, appID string, yamlFile string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	// Read and validate YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file: %w", err)
	}

	// Create request body
	reqBody := map[string]interface{}{
		"yamlContent": string(yamlData),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	var url string
	if appID != "" {
		// If appID is provided, include it in the URL
		url = fmt.Sprintf("%s/startUserDeployment/%s", c.baseURL, appID)
	} else {
		// If no appID, use base endpoint
		url = fmt.Sprintf("%s/startUserDeployment", c.baseURL)
	}

	// Debug: Print the URL we're requesting
	fmt.Printf("DEBUG: Starting deployment at URL: %s\n", url)

	resp, err := c.post(ctx, url, jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to start deployment: %w", err)
	}

	// Parse response
	var apiResp schema.APIResponse[schema.DeploymentResponse]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &apiResp, nil
}

// SaveCustomDomain associates a custom domain with a specific application deployment.
// Endpoint: POST /saveCustomDomain/{applicationID}
func (c *Client) SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error) {
	url := fmt.Sprintf("%s/saveCustomDomain/%s", c.baseURL, appID)

	req := struct {
		Domain string `json:"domain"`
	}{
		Domain: domain,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.post(ctx, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to save custom domain: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp)
	}

	var response schema.APIResponse[struct{}]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetDeployments retrieves all deployments associated with the specified application ID.
// Endpoint: GET /getDeployments/{applicationID}
func (c *Client) GetDeployments(ctx context.Context, appID string) (*schema.APIResponse[[]schema.Deployment], error) {
	url := fmt.Sprintf("%s/getDeployments/%s", c.baseURL, appID)
	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleAPIError(resp)
	}

	var response schema.APIResponse[[]schema.Deployment]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// GetDeploymentInfo retrieves detailed information about a specific deployment.
// Endpoint: GET /getDeploymentInfo/{namespace}
func (c *Client) GetDeploymentInfo(ctx context.Context, namespace string) (*schema.APIResponse[schema.Deployment], error) {
	// Ensure namespace is not empty
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required and cannot be empty")
	}

	// Remove any leading or trailing whitespace
	namespace = strings.TrimSpace(namespace)
	if namespace == "" {
		return nil, fmt.Errorf("namespace cannot be only whitespace")
	}

	// Additional validation to prevent double slashes in URL
	if strings.Contains(namespace, "/") {
		return nil, fmt.Errorf("namespace cannot contain slashes")
	}

	// Construct URL properly preserving the scheme
	url := fmt.Sprintf("%s/getDeploymentInfo/%s", c.baseURL, namespace)

	// Debug: Print the URL we're requesting
	fmt.Printf("DEBUG: Checking deployment status at URL: %s\n", url)

	resp, err := c.get(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment info: %w", err)
	}

	// Check for non-200 responses
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResp schema.APIResponse[schema.Deployment]
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &apiResp, nil
}

// Helper methods for making HTTP requests
func (c *Client) get(ctx context.Context, url string) (*http.Response, error) {
	// Check for double slashes in URL (except for http:// or https://)
	if strings.Contains(url, "//") &&
		!strings.Contains(url, "http://") &&
		!strings.Contains(url, "https://") {
		fmt.Printf("WARNING: URL contains double slashes: %s\n", url)
		// Fix the URL by replacing multiple slashes with a single slash
		fixedURL := strings.Replace(url, "//", "/", -1)
		// But preserve http:// or https://
		fixedURL = strings.Replace(fixedURL, "http:/", "http://", 1)
		fixedURL = strings.Replace(fixedURL, "https:/", "https://", 1)
		fmt.Printf("WARNING: Fixed URL: %s\n", fixedURL)
		url = fixedURL
	}

	fmt.Printf("GET Request URL: %s\n", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	fmt.Printf("Making GET request with headers: %v\n", req.Header)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	fmt.Printf("Response status: %s\n", resp.Status)
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) post(ctx context.Context, url string, body []byte) (*http.Response, error) {
	fmt.Printf("POST Request URL: %s\n", url)
	fmt.Printf("POST Request Body: %s\n", string(body))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
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

// SendFeedback sends user feedback to Nexlayer.
// The feedback text will be used to improve the service.
func (c *Client) SendFeedback(ctx context.Context, text string) error {
	url := fmt.Sprintf("%s/feedback", c.baseURL)
	fmt.Printf("Sending feedback to: %s\n", url)

	feedback := map[string]string{"text": text}
	body, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("failed to marshal feedback: %w", err)
	}

	resp, err := c.post(ctx, url, body)
	if err != nil {
		fmt.Printf("Error sending feedback: %v\n", err)
		return fmt.Errorf("failed to send feedback: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Feedback sent successfully\n")
	return nil
}
