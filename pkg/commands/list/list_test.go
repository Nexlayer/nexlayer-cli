package list_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
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

func (m *mockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	args := m.Called(ctx, appID, configPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[schema.DeploymentResponse]), args.Error(1)
}

func (m *mockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error) {
	args := m.Called(ctx, namespace, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[schema.Deployment]), args.Error(1)
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

func (m *mockAPIClient) GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error) {
	args := m.Called(ctx, namespace, appID, follow, tail)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockAPIClient) SendFeedback(ctx context.Context, text string) error {
	args := m.Called(ctx, text)
	return args.Error(0)
}

func NewListCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "list [applicationID]",
		Short: "List deployments",
		Long:  "List all deployments associated with the specified application ID.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command execution logic
			return nil
		},
	}
}

func TestNewListCommand(t *testing.T) {
	client := &mockAPIClient{}
	cmd := list.NewListCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "list [applicationID]", cmd.Use)
	assert.Equal(t, "List deployments", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestListDeployments(t *testing.T) {
	client := &mockAPIClient{}

	// Setup default mock response
	successResp := &schema.APIResponse[[]schema.Deployment]{
		Message: "Success",
		Data: []schema.Deployment{
			{
				Namespace: "test",
				Status:    "Running",
				URL:       "https://test.nexlayer.dev",
				Version:   "v1.0.0",
			},
		},
	}
	client.On("ListDeployments", mock.Anything).Return(successResp, nil)

	tests := []struct {
		name     string
		args     []string
		wantJSON bool
		wantErr  bool
		errMsg   string
		setup    func()
	}{
		{
			name:     "list all deployments",
			args:     []string{},
			wantJSON: false,
			wantErr:  false,
		},
		{
			name:     "list deployments as JSON",
			args:     []string{"--json"},
			wantJSON: true,
			wantErr:  false,
		},
		{
			name:     "handle API error",
			args:     []string{},
			wantJSON: false,
			wantErr:  true,
			errMsg:   "failed to get deployments",
			setup: func() {
				client.On("ListDeployments", mock.Anything).Return(nil, fmt.Errorf("failed to get deployments")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			cmd := list.NewListCommand(client)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)

			if tt.wantJSON {
				cmd.Flags().Bool("json", true, "")
			}

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			assert.NoError(t, err)
			output := buf.String()

			if tt.wantJSON {
				assert.Contains(t, output, "{")
				assert.Contains(t, output, "}")
				var resp map[string]interface{}
				err = json.Unmarshal([]byte(output), &resp)
				assert.NoError(t, err)
			} else {
				assert.Contains(t, output, "STATUS")
				assert.Contains(t, output, "URL")
				assert.Contains(t, output, "VERSION")
				assert.Contains(t, output, "Running")
				assert.Contains(t, output, "https://test.nexlayer.dev")
			}
		})
	}
}
