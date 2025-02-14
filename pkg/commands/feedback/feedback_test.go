package feedback_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/feedback"
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

func (m *mockAPIClient) GetLogs(ctx context.Context, namespace string, appID string, follow bool, tail int) ([]string, error) {
	args := m.Called(ctx, namespace, appID, follow, tail)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockAPIClient) ListDeployments(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
	args := m.Called(ctx)
	return args.Get(0).(*schema.APIResponse[[]schema.Deployment]), args.Error(1)
}

func (m *mockAPIClient) StartDeployment(ctx context.Context, param1 string, param2 string) (*schema.APIResponse[schema.DeploymentResponse], error) {
	args := m.Called(ctx, param1, param2)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*schema.APIResponse[schema.DeploymentResponse]), args.Error(1)
}

func TestNewFeedbackCommand(t *testing.T) {
	tests := []struct {
		name           string
		expectedUse    string
		expectedShort  string
		shouldHaveLong bool
		flagName       string
	}{
		{
			name:           "creates command with correct properties",
			expectedUse:    "feedback",
			expectedShort:  "Send feedback",
			shouldHaveLong: true,
			flagName:       "message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := new(mockAPIClient)
			cmd := feedback.NewFeedbackCommand(client)

			assert.NotNil(t, cmd)
			assert.Equal(t, tt.expectedUse, cmd.Use)
			assert.Equal(t, tt.expectedShort, cmd.Short)
			if tt.shouldHaveLong {
				assert.NotEmpty(t, cmd.Long)
			}

			// Verify send subcommand exists
			sendCmd, _, err := cmd.Find([]string{"send"})
			assert.NoError(t, err)
			assert.NotNil(t, sendCmd)

			// Verify required flags on send subcommand
			flag := sendCmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag)
			assert.True(t, sendCmd.MarkFlagRequired(tt.flagName) == nil)
		})
	}
}

func TestSendFeedback(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		mockError     error
		expectedError string
		expectedMsg   string
	}{
		{
			name:          "successful feedback submission",
			message:       "Great product!",
			mockError:     nil,
			expectedError: "",
			expectedMsg:   "Thank you",
		},
		{
			name:          "missing message",
			message:       "",
			mockError:     nil,
			expectedError: "required flag(s) \"message\" not set",
			expectedMsg:   "",
		},
		{
			name:          "api error",
			message:       "Test feedback",
			mockError:     errors.New("api error"),
			expectedError: "api error",
			expectedMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			client := new(mockAPIClient)
			if tt.message != "" {
				client.On("SendFeedback", mock.Anything, tt.message).Return(tt.mockError)
			}

			// Setup command
			cmd := feedback.NewFeedbackCommand(client)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)

			// Set up command arguments
			args := []string{"send"}
			if tt.message != "" {
				args = append(args, "--message", tt.message)
			}
			cmd.SetArgs(args)

			// Execute command
			err := cmd.Execute()

			// Verify results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Contains(t, buf.String(), tt.expectedMsg)
			}

			// Verify all mock expectations were met
			client.AssertExpectations(t)
		})
	}
}
