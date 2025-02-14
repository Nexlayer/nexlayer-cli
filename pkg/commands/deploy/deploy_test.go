package deploy

import (
	"bytes"
	"context"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/stretchr/testify/assert"
)

type mockAPIClient struct{}

func (m *mockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	// Mock different responses based on whether appID is provided
	if appID == "" {
		return &schema.APIResponse[schema.DeploymentResponse]{
			Data: schema.DeploymentResponse{
				Namespace: "profile-namespace",
				URL:       "https://profile-app.nexlayer.dev",
			},
		}, nil
	}
	return &schema.APIResponse[schema.DeploymentResponse]{
		Data: schema.DeploymentResponse{
			Namespace: "test-namespace",
			URL:       "https://test.nexlayer.dev",
		},
	}, nil
}

func (m *mockAPIClient) SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error) {
	return &schema.APIResponse[struct{}]{Data: struct{}{}}, nil
}

func (m *mockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	return &schema.APIResponse[[]schema.Deployment]{
		Data: []schema.Deployment{},
	}, nil
}

func (m *mockAPIClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error) {
	return &schema.APIResponse[schema.Deployment]{
		Data: schema.Deployment{
			Status:  "running",
			URL:     "https://test.nexlayer.dev",
			Version: "v1.0.0",
			PodStatuses: []schema.PodStatus{
				{
					Name:   "web",
					Status: "Running",
					Ready:  true,
				},
			},
		},
	}, nil
}

func (m *mockAPIClient) GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error) {
	return []string{}, nil
}

func (m *mockAPIClient) SendFeedback(ctx context.Context, text string) error {
	return nil
}

func TestNewCommand(t *testing.T) {
	client := &commands.MockAPIClient{}
	cmd := NewCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "deploy [applicationID]", cmd.Use)
	assert.Contains(t, cmd.Short, "Deploy")

	// Test required flags
	configFlag := cmd.Flags().Lookup("config")
	assert.NotNil(t, configFlag)

	// Verify flag is required by attempting to execute without it
	cmd.SetArgs([]string{"myapp"})
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required flag(s) \"config\" not set")

	// Test deployment execution
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"myapp", "--config", "nexlayer.yaml"})

	// Mock the API response
	client.StartDeploymentFunc = func(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
		return &schema.APIResponse[schema.DeploymentResponse]{
			Data: schema.DeploymentResponse{
				Namespace: "test",
				URL:       "https://myapp.nexlayer.dev",
			},
		}, nil
	}

	err = cmd.Execute()
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "Deployment started")
}

func TestValidateDeployConfig(t *testing.T) {
	tests := []struct {
		name    string
		yaml    *schema.NexlayerYAML
		wantErr bool
	}{
		{
			name: "valid deployment",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "test-app",
					Pods: []schema.Pod{
						{
							Name:  "web",
							Image: "nginx:latest",
							Ports: []schema.Port{
								{
									ContainerPort: 80,
									ServicePort:   80,
									Name:          "web",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing required fields",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "",
					Pods: []schema.Pod{},
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
