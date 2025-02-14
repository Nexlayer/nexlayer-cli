package list_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
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
	// Create mock client and command for each test case in the test function

	tests := []struct {
		name        string
		args        []string
		setupMock   func(*mockAPIClient)
		wantJSON    bool
		wantErr     bool
		wantOutput  string
	}{
		{
			name: "list all deployments",
			args: []string{},
			setupMock: func(m *mockAPIClient) {
				m.On("GetDeployments", mock.Anything, "").Return(&schema.APIResponse[[]schema.Deployment]{
					Data: []schema.Deployment{
						{
							Status:      "running",
							URL:         "https://app1.nexlayer.dev",
							Version:     "v1.0.0",
							LastUpdated: time.Now(),
						},
						{
							Status:      "stopped",
							URL:         "https://app2.nexlayer.dev",
							Version:     "v1.1.0",
							LastUpdated: time.Now().Add(-24 * time.Hour),
						},
					},
				}, nil)
			},
			wantJSON:   false,
			wantErr:    false,
			wantOutput: `Status\s+URL\s+Version\s+Last Updated`,
		},
		{
			name: "list deployments for specific app",
			args: []string{"myapp"},
			setupMock: func(m *mockAPIClient) {
				m.On("GetDeployments", mock.Anything, "myapp").Return(&schema.APIResponse[[]schema.Deployment]{
					Data: []schema.Deployment{
						{
							Status:      "running",
							URL:         "https://myapp.nexlayer.dev",
							Version:     "v1.0.0",
							LastUpdated: time.Now(),
						},
					},
				}, nil)
			},
			wantJSON:   false,
			wantErr:    false,
			wantOutput: `running\s+https://myapp.nexlayer.dev`,
		},
		{
			name: "list all deployments as JSON",
			args: []string{"--json"},
			setupMock: func(m *mockAPIClient) {
				m.On("GetDeployments", mock.Anything, "").Return(&schema.APIResponse[[]schema.Deployment]{
					Data: []schema.Deployment{
						{
							Status:      "running",
							URL:         "https://app1.nexlayer.dev",
							Version:     "v1.0.0",
							LastUpdated: time.Now(),
						},
					},
				}, nil)
			},
			wantJSON:   true,
			wantErr:    false,
			wantOutput: `{"deployments":[{"status":"running","url":"https://app1.nexlayer.dev"}]}`,
		},
		{
			name: "handle API error",
			args: []string{},
			setupMock: func(m *mockAPIClient) {
				m.On("GetDeployments", mock.Anything, "").Return(nil, assert.AnError)
			},
			wantJSON:   false,
			wantErr:    true,
			wantOutput: "Error: Failed to get deployments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new mock client for each test
			mockClient := &mockAPIClient{}
			tt.setupMock(mockClient)

			// Create new command with mock client
			cmd := NewListCommand(mockClient)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetArgs(tt.args)

			// Set JSON flag if needed
			if tt.wantJSON {
				cmd.Flags().Set("json", "true")
			}

			// Execute command
			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, buf.String(), tt.wantOutput)
				return
			}

			// Verify no error and output matches
			assert.NoError(t, err)
			assert.NotEmpty(t, buf.String())

			// Verify output format and content
			if tt.wantJSON {
				assert.Contains(t, buf.String(), "{")
				assert.JSONEq(t, tt.wantOutput, buf.String())
			} else {
				assert.Regexp(t, tt.wantOutput, buf.String())
			}

			// Verify all mock expectations were met
			mockClient.AssertExpectations(t)
		})
	}
}
