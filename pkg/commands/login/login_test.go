package login_test

import (
	"bytes"
	"context"
	"testing"

	logincmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
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

func (m *mockAPIClient) SendFeedback(ctx context.Context, text string) error {
	args := m.Called(ctx, text)
	return args.Error(0)
}

func (m *mockAPIClient) GetLogs(ctx context.Context, deploymentID string, containerName string, follow bool, tail int) ([]string, error) {
	args := m.Called(ctx, deploymentID, containerName, follow, tail)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) StartDeployment(ctx context.Context, appID string, configPath string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	args := m.Called(ctx, appID, configPath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[schema.DeploymentResponse]), args.Error(1)
}

func NewCommand(client api.APIClient) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Log in to Nexlayer",
		Long:  "Authenticate with the Nexlayer platform using your credentials.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Command execution logic
			return nil
		},
	}
}

func TestNewCommand(t *testing.T) {
	client := new(mockAPIClient)
	cmd := logincmd.NewLoginCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "login", cmd.Use)
	assert.Equal(t, "Log in to Nexlayer", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestLoginCommand(t *testing.T) {
	client := new(mockAPIClient)
	cmd := logincmd.NewLoginCommand(client)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)

	// Execute the command
	err := cmd.Execute()
	assert.Error(t, err)
	assert.Equal(t, "login flow not yet implemented", err.Error())
}

func TestLoginCommandStructure(t *testing.T) {
	client := new(mockAPIClient)
	cmd := logincmd.NewLoginCommand(client)

	assert.NotNil(t, cmd)
	assert.Equal(t, "login", cmd.Use)
	assert.Equal(t, "Log in to Nexlayer", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}
