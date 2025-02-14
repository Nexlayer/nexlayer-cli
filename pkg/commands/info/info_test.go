package info_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockAPIClient is a mock implementation of api.APIClient
type mockAPIClient struct {
	mock.Mock
}

// Ensure mockAPIClient implements api.APIClient
var _ api.APIClient = (*mockAPIClient)(nil)

func (m *mockAPIClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error) {
	args := m.Called(ctx, namespace, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return &schema.APIResponse[schema.Deployment]{
		Data: schema.Deployment{
			Status:      "running",
			URL:         "https://test.nexlayer.dev",
			Version:     "v1.0.0",
			LastUpdated: time.Now(),
			PodStatuses: []schema.PodStatus{
				{
					Name:   "web",
					Status: "Running",
					Ready:  true,
				},
			},
		},
	}, args.Error(1)
}

func (m *mockAPIClient) GetDeployments(ctx context.Context, appID string) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) SaveCustomDomain(ctx context.Context, appID string, domain string) (*schema.APIResponse[struct{}], error) {
	args := m.Called(ctx, appID, domain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return &schema.APIResponse[struct{}]{
		Data: struct{}{},
	}, args.Error(1)
}

func (m *mockAPIClient) SendFeedback(ctx context.Context, text string) error {
	args := m.Called(ctx, text)
	return args.Error(0)
}

func (m *mockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx)
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) GetLogs(ctx context.Context, deploymentID string, containerName string, follow bool, tail int) ([]string, error) {
	args := m.Called(ctx, deploymentID, containerName, follow, tail)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	args := m.Called(ctx, appID, configPath)
	return args.Get(0).(*schema.APIResponse[schema.DeploymentResponse]), args.Error(1)
}

func NewInfoCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "info <namespace> <applicationID>",
		Short: "Get deployment info",
		Long:  "Retrieve detailed information about a specific deployment.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command execution logic
			return nil
		},
	}
}

func TestNewInfoCommand(t *testing.T) {
	client := &mockAPIClient{}
	cmd := NewInfoCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "info <namespace> <applicationID>", cmd.Use)
	assert.Equal(t, "Get deployment info", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestGetDeploymentInfo(t *testing.T) {
	client := &mockAPIClient{}
	cmd := NewInfoCommand(client)

	tests := []struct {
		name     string
		args     []string
		wantJSON bool
		wantErr  bool
	}{
		{
			name:     "get deployment info",
			args:     []string{"default", "myapp"},
			wantJSON: false,
			wantErr:  false,
		},
		{
			name:     "get deployment info as JSON",
			args:     []string{"default", "myapp", "--json"},
			wantJSON: true,
			wantErr:  false,
		},
		{
			name:     "missing arguments",
			args:     []string{"default"},
			wantJSON: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetArgs(tt.args)
			if tt.wantJSON {
				cmd.Flags().Set("json", "true")
			}

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, buf.String())

			if tt.wantJSON {
				assert.Contains(t, buf.String(), "{")
				assert.Contains(t, buf.String(), "}")
			} else {
				assert.Contains(t, buf.String(), "Status:")
				assert.Contains(t, buf.String(), "URL:")
				assert.Contains(t, buf.String(), "Version:")
			}
		})
	}
}
