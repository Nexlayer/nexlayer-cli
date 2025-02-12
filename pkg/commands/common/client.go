// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package common

import (
	"context"
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

// CommandClient abstracts the API client for command usage
type CommandClient interface {
	// GetDeploymentInfo retrieves detailed information about a specific deployment
	GetDeploymentInfo(ctx context.Context, namespace, appID string) (*schema.APIResponse[schema.Deployment], error)

	// ListDeployments retrieves all deployments
	ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error)

	// SaveCustomDomain associates a custom domain with a deployment
	SaveCustomDomain(ctx context.Context, appID, domain string) (*schema.APIResponse[struct{}], error)

	// StartDeployment starts a new deployment from a YAML file
	StartDeployment(ctx context.Context, appID, yamlFile string) (*schema.APIResponse[schema.DeploymentResponse], error)

	// GetLogs retrieves logs for a specific deployment
	GetLogs(ctx context.Context, namespace, appID string, follow bool, tail int) ([]string, error)

	// SendFeedback sends user feedback
	SendFeedback(ctx context.Context, text string) error
}

// clientAdapter implements CommandClient by wrapping the API client
type clientAdapter struct {
	api api.APIClient
}

// NewCommandClient creates a new CommandClient
func NewCommandClient(apiClient api.APIClient) CommandClient {
	return &clientAdapter{api: apiClient}
}

func (c *clientAdapter) GetDeploymentInfo(ctx context.Context, namespace, appID string) (*schema.APIResponse[schema.Deployment], error) {
	resp, err := c.api.GetDeploymentInfo(ctx, namespace, appID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment info: %w", err)
	}
	return resp, nil
}

func (c *clientAdapter) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	resp, err := c.api.ListDeployments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}
	return resp, nil
}

func (c *clientAdapter) SaveCustomDomain(ctx context.Context, appID, domain string) (*schema.APIResponse[struct{}], error) {
	resp, err := c.api.SaveCustomDomain(ctx, appID, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to save custom domain: %w", err)
	}
	return resp, nil
}

func (c *clientAdapter) StartDeployment(ctx context.Context, appID, yamlFile string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	resp, err := c.api.StartDeployment(ctx, appID, yamlFile)
	if err != nil {
		return nil, fmt.Errorf("failed to start deployment: %w", err)
	}
	return resp, nil
}

func (c *clientAdapter) GetLogs(ctx context.Context, namespace, appID string, follow bool, tail int) ([]string, error) {
	logs, err := c.api.GetLogs(ctx, namespace, appID, follow, tail)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs: %w", err)
	}
	return logs, nil
}

func (c *clientAdapter) SendFeedback(ctx context.Context, text string) error {
	err := c.api.SendFeedback(ctx, text)
	if err != nil {
		return fmt.Errorf("failed to send feedback: %w", err)
	}
	return nil
}
