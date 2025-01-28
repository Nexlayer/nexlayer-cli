package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

type mockClient struct {
	deploymentInfo types.DeploymentInfo
	err           error
}

func (m *mockClient) StartDeployment(ctx context.Context, appID string, configPath string) (*types.StartDeploymentResponse, error) {
	return &types.StartDeploymentResponse{
		Namespace: "test-ns",
		URL:      "https://test-ns.nexlayer.com",
	}, m.err
}

func (m *mockClient) SaveCustomDomain(ctx context.Context, appID string, domain string) error {
	return m.err
}

func (m *mockClient) GetDeployments(ctx context.Context, appID string) ([]types.Deployment, error) {
	return []types.Deployment{
		{
			Namespace:        "test-ns",
			TemplateName:     "python",
			TemplateID:       "123",
			DeploymentStatus: "running",
		},
	}, m.err
}

func (m *mockClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error) {
	return &m.deploymentInfo, m.err
}

func TestNewCommand(t *testing.T) {
	mock := &mockClient{
		deploymentInfo: types.DeploymentInfo{
			Namespace:        "test-ns",
			TemplateName:     "python",
			TemplateID:       "123",
			DeploymentStatus: "running",
		},
	}

	cmd := NewCommand(mock)
	assert.NotNil(t, cmd)
	assert.Equal(t, "app", cmd.Use)
}

func TestRunCommand(t *testing.T) {
	mock := &mockClient{
		deploymentInfo: types.DeploymentInfo{
			Namespace:        "test-ns",
			TemplateName:     "python",
			TemplateID:       "123",
			DeploymentStatus: "running",
		},
	}

	cmd := NewCommand(mock)
	cmd.SetArgs([]string{"info", "--app", "test-app", "--namespace", "test-ns"})

	err := cmd.Execute()
	assert.NoError(t, err)
}
