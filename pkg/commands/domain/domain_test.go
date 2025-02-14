package domain_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
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
	return args.Get(0).(*schema.APIResponse[schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) GetDeployments(ctx context.Context, appID string) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx, appID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) SendFeedback(ctx context.Context, text string) error {
	args := m.Called(ctx, text)
	return args.Error(0)
}

func (m *mockAPIClient) GetLogs(ctx context.Context, deploymentID string, containerName string, follow bool, tail int) ([]string, error) {
	args := m.Called(ctx, deploymentID, containerName, follow, tail)
	return args.Get(0).([]string), args.Error(1)
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

func (m *mockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx)
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	args := m.Called(ctx, appID, configPath)
	return args.Get(0).(*schema.APIResponse[schema.DeploymentResponse]), args.Error(1)
}

func TestNewDomainCommand(t *testing.T) {
	client := &commands.MockAPIClient{}
	cmd := domain.NewDomainCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "domain", cmd.Use)
	assert.Equal(t, "Manage custom domains", cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Check that set subcommand exists
	setCmd, _, err := cmd.Find([]string{"set"})
	assert.NoError(t, err)
	assert.NotNil(t, setCmd)
	assert.Equal(t, "set <applicationID>", setCmd.Use)
}

func TestSetDomain(t *testing.T) {
	client := &mockAPIClient{}
	cmd := domain.NewDomainCommand(client)

	tests := []struct {
		name    string
		args    []string
		domain  string
		wantErr bool
	}{
		{
			name:    "set custom domain",
			args:    []string{"set", "myapp"},
			domain:  "example.com",
			wantErr: false,
		},
		{
			name:    "missing domain flag",
			args:    []string{"set", "myapp"},
			domain:  "",
			wantErr: true,
		},
		{
			name:    "missing application ID",
			args:    []string{"set"},
			domain:  "example.com",
			wantErr: true,
		},
		{
			name:    "invalid domain with spaces",
			args:    []string{"set", "myapp"},
			domain:  "my domain.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetArgs(tt.args)
			if tt.domain != "" {
				cmd.Flags().Set("domain", tt.domain)
			}

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, buf.String(), "Custom domain")
			assert.Contains(t, buf.String(), "saved")
		})
	}
}

func TestValidateDomain(t *testing.T) {
	err := domain.ValidateDomain("example.com")
	assert.NoError(t, err)

	err = domain.ValidateDomain("")
	assert.Error(t, err)
}
