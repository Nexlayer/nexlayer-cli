// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
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

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
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

	// GetLogs retrieves logs for a specific deployment.
	// If follow is true, streams logs in real-time.
	// tail specifies the number of lines to return from the end of the logs.
	GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error)

	// SendFeedback sends user feedback to Nexlayer.
	// The feedback text will be used to improve the service.
	SendFeedback(ctx context.Context, text string) error
}

// Client represents an API client for interacting with the Nexlayer API.
// The Nexlayer API enables rapid deployment of full-stack AI-powered applications
// by providing a simple template-based interface that abstracts away deployment complexity.
type Client struct {
	baseURL    string       // Base URL of the Nexlayer API
	httpClient *http.Client // HTTP client for making API requests
	token      string       // Authentication token for API requests
}

// handleAPIError processes API error responses and returns a formatted error
func (c *Client) handleAPIError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	var errResp types.ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	return fmt.Errorf("API error (status %d): %s (code: %s)", resp.StatusCode, errResp.Message, errResp.Code)
}

// NewClient creates a new Nexlayer API client.
// If baseURL is empty, defaults to the staging environment at app.staging.nexlayer.io.
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
			Timeout:   30 * time.Second,
			Transport: transport,
		},
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
	// Read YAML file
	data, err := os.ReadFile(yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read template file: %w", err)
	}

	// Make request
	url := fmt.Sprintf("%s/startUserDeployment", c.baseURL)
	if appID != "" {
		url = fmt.Sprintf("%s/%s", url, appID)
	}

	// Send as YAML
	resp, err := c.postYAML(ctx, url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Print status code and response body for debugging
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(respBody))

	// Check for error response
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error: %s - %s", resp.Status, string(respBody))
	}

	var result types.StartDeploymentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
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
	defer resp.Body.Close()

	var result types.GetDeploymentsResponse
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
	defer resp.Body.Close()

	var result types.GetDeploymentInfoResponse
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
