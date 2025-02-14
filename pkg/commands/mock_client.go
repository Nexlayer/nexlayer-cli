package commands

import (
	"context"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

// MockAPIClient implements both api.APIClient and api.ClientAPI interfaces for testing
type MockAPIClient struct {
	// Mock fields for test control
	GetDeploymentsFunc    func(ctx context.Context) ([]schema.Deployment, error)
	SaveCustomDomainFunc  func(ctx context.Context, appID, domain string) (*schema.APIResponse[struct{}], error)
	GetLogsFunc           func(ctx context.Context, name string) ([]string, error)
	GetDeploymentInfoFunc func(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error)
	ListDeploymentsFunc   func(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error)
	StartDeploymentFunc   func(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error)
	SendFeedbackFunc      func(ctx context.Context, text string) error
}

// GetDeployments implements api.APIClient
func (m *MockAPIClient) GetDeployments(ctx context.Context) ([]schema.Deployment, error) {
	if m.GetDeploymentsFunc != nil {
		return m.GetDeploymentsFunc(ctx)
	}
	return nil, nil
}

// SaveCustomDomain implements api.APIClient
func (m *MockAPIClient) SaveCustomDomain(ctx context.Context, appID, domain string) (*schema.APIResponse[struct{}], error) {
	if m.SaveCustomDomainFunc != nil {
		return m.SaveCustomDomainFunc(ctx, appID, domain)
	}
	return nil, nil
}

// GetLogs implements api.APIClient
func (m *MockAPIClient) GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error) {
	if m.GetLogsFunc != nil {
		return m.GetLogsFunc(ctx, namespace)
	}
	return nil, nil
}

// GetDeploymentInfo implements api.APIClient
func (m *MockAPIClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error) {
	if m.GetDeploymentInfoFunc != nil {
		return m.GetDeploymentInfoFunc(ctx, namespace, appID)
	}
	return nil, nil
}

// ListDeployments implements api.APIClient
func (m *MockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	if m.ListDeploymentsFunc != nil {
		return m.ListDeploymentsFunc(ctx)
	}
	return nil, nil
}

// StartDeployment implements api.APIClient
func (m *MockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	if m.StartDeploymentFunc != nil {
		return m.StartDeploymentFunc(ctx, appID, configPath)
	}
	return nil, nil
}

// SendFeedback implements api.APIClient
func (m *MockAPIClient) SendFeedback(ctx context.Context, text string) error {
	if m.SendFeedbackFunc != nil {
		return m.SendFeedbackFunc(ctx, text)
	}
	return nil
}

// Ensure MockAPIClient implements both interfaces
var (
	_ api.APIClient = (*MockAPIClient)(nil)
	_ api.ClientAPI = (*MockAPIClient)(nil)
)
