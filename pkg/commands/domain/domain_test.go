package domain_test

import (
	"bytes"
	"context"
	"fmt"
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
	return args.Get(0).(*schema.APIResponse[struct{}]), args.Error(1)
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

	// Setup mock responses
	successResp := &schema.APIResponse[struct{}]{
		Message: "Success",
		Data:    struct{}{},
	}

	// Setup mock expectations
	client.On("SaveCustomDomain", mock.Anything, "myapp", "example.com").Return(successResp, nil)
	client.On("SaveCustomDomain", mock.Anything, "myapp", "").Return(nil, fmt.Errorf("domain cannot be empty"))
	client.On("SaveCustomDomain", mock.Anything, "myapp", "invalid domain.com").Return(nil, fmt.Errorf("domain cannot contain spaces"))

	tests := []struct {
		name    string
		args    []string
		domain  string
		wantErr bool
		errMsg  string
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
			errMsg:  "domain cannot be empty",
		},
		{
			name:    "missing application ID",
			args:    []string{"set"},
			domain:  "example.com",
			wantErr: true,
			errMsg:  "accepts 1 arg(s)",
		},
		{
			name:    "invalid domain with spaces",
			args:    []string{"set", "myapp"},
			domain:  "invalid domain.com",
			wantErr: true,
			errMsg:  "domain cannot contain spaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := domain.NewDomainCommand(client)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)

			// Set command arguments including the domain flag if provided
			args := tt.args
			if tt.domain != "" {
				args = append(args, "--domain", tt.domain)
			}
			cmd.SetArgs(args)

			// Execute the command
			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			output := buf.String()
			assert.Contains(t, output, fmt.Sprintf("✓ Custom domain %s saved for application %s", tt.domain, tt.args[1]))
		})
	}
}

func TestValidateDomain(t *testing.T) {
	err := domain.ValidateDomain("example.com")
	assert.NoError(t, err)

	err = domain.ValidateDomain("")
	assert.Error(t, err)
}
