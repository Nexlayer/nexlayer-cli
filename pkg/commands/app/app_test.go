package app

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

type mockClient struct {
	getDeploymentInfoFunc func(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error)
}

func (m *mockClient) GetDeploymentInfo(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error) {
	return m.getDeploymentInfoFunc(ctx, namespace, appID)
}

func TestAppCommand(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		mockSetup func(*mockClient)
		wantErr   bool
		wantText  string
	}{
		{
			name: "get app info",
			args: []string{"info", "--app", "testapp", "--namespace", "test-ns"},
			mockSetup: func(m *mockClient) {
				m.getDeploymentInfoFunc = func(ctx context.Context, namespace string, appID string) (*types.DeploymentInfo, error) {
					return &types.DeploymentInfo{
						Namespace:        "test-ns",
						TemplateName:     "python",
						TemplateID:       "123",
						DeploymentStatus: "running",
					}, nil
				}
			},
			wantErr:  false,
			wantText: "Application ID: testapp\nNamespace:      test-ns\nTemplate:       python\nStatus:         running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockClient{}
			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			cmd := NewCommand(mock)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, buf.String(), tt.wantText)
		})
	}
}
