package info_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/stretchr/testify/assert"
)

func TestNewInfoCommand(t *testing.T) {
	client := &commands.MockAPIClient{}
	cmd := info.NewInfoCommand(client)
	assert.NotNil(t, cmd)
	assert.Equal(t, "info <namespace> <applicationID>", cmd.Use)
	assert.Equal(t, "Get deployment info", cmd.Short)
	assert.NotEmpty(t, cmd.Long)
}

func TestGetDeploymentInfo(t *testing.T) {
	client := &commands.MockAPIClient{}

	// Setup mock response
	deploymentResp := &schema.APIResponse[schema.Deployment]{
		Message: "Success",
		Data: schema.Deployment{
			Namespace:   "default",
			Status:      "Running",
			URL:         "https://myapp.nexlayer.dev",
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
	}

	// Setup mock expectations
	client.GetDeploymentInfoFunc = func(ctx context.Context, namespace string, appID string) (*schema.APIResponse[schema.Deployment], error) {
		return deploymentResp, nil
	}

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
			cmd := info.NewInfoCommand(client)
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)

			if tt.wantJSON {
				cmd.Flags().Bool("json", true, "")
			}

			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
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
				assert.Contains(t, output, "Status:")
				assert.Contains(t, output, "URL:")
				assert.Contains(t, output, "Version:")
			}
		})
	}
}
