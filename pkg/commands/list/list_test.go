package list_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/stretchr/testify/assert"
)

func TestNewListCommand(t *testing.T) {
	client := &commands.MockAPIClient{}
	cmd := list.NewListCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "list [applicationID]", cmd.Use)
	assert.Equal(t, "List deployments", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestListDeployments(t *testing.T) {
	client := &commands.MockAPIClient{}

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
			setup: func() {
				client.ListDeploymentsFunc = func(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
					return successResp, nil
				}
			},
		},
		{
			name:     "list deployments as JSON",
			args:     []string{"--json"},
			wantJSON: true,
			wantErr:  false,
			setup: func() {
				client.ListDeploymentsFunc = func(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
					return successResp, nil
				}
			},
		},
		{
			name:     "handle API error",
			args:     []string{},
			wantJSON: false,
			wantErr:  true,
			errMsg:   "failed to get deployments",
			setup: func() {
				client.ListDeploymentsFunc = func(ctx context.Context) (*schema.APIResponse[[]schema.Deployment], error) {
					return nil, fmt.Errorf("failed to get deployments")
				}
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
				_ = cmd.Flags().Set("json", "true")
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
