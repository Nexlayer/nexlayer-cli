package deploy

import (
	"context"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
	"github.com/stretchr/testify/assert"
)

type mockAPIClient struct{}

func (m *mockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*types.StartDeploymentResponse, error) {
	return &types.StartDeploymentResponse{
		Message:   "Deployment started",
		Namespace: "test-namespace",
		URL:       "https://test.nexlayer.dev",
	}, nil
}

func (m *mockAPIClient) SaveCustomDomain(ctx context.Context, appID string, domain string) error {
	return nil
}

func (m *mockAPIClient) GetDeployments(ctx context.Context, appID string) ([]types.Deployment, error) {
	return []types.Deployment{}, nil
}

func (m *mockAPIClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error) {
	return &types.DeploymentInfo{
		Namespace:        namespace,
		TemplateID:       "test-id",
		TemplateName:     "test-app",
		DeploymentStatus: "running",
	}, nil
}

func (m *mockAPIClient) GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error) {
	return []string{}, nil
}

func (m *mockAPIClient) SendFeedback(ctx context.Context, text string) error {
	return nil
}

func TestNewCommand(t *testing.T) {
	client := &mockAPIClient{}
	cmd := NewCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "deploy", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestValidateDeployConfig(t *testing.T) {
	tests := []struct {
		name    string
		yaml    *types.NexlayerYAML
		wantErr bool
	}{
		{
			name: "valid deployment",
			yaml: &types.NexlayerYAML{
				Application: types.Application{
					Name: "test-app",
					Pods: []types.Pod{
						{
							Name:  "web",
							Type:  "react",
							Image: "nginx:latest",
							Path:  "/",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			yaml: &types.NexlayerYAML{
				Application: types.Application{
					Name: "",
					Pods: []types.Pod{},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDeployConfig(tt.yaml)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
