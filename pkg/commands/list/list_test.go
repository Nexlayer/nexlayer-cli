package list

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListCmd(t *testing.T) {
	cmd := NewListCmd()
	assert.Equal(t, "list", cmd.Use)
	assert.Equal(t, "List resources", cmd.Short)
	assert.True(t, len(cmd.Commands()) > 0, "Should have subcommands")
}

func TestListDeploymentsCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		wantErr  bool
		errMsg   string
		wantText string
	}{
		{
			name:    "No app flag",
			args:    []string{},
			wantErr: true,
			errMsg:  "required flag(s) \"app\" not set",
		},
		{
			name: "Empty app name",
			args: []string{},
			flags: map[string]string{
				"app": "",
			},
			wantErr: true,
			errMsg:  "app name is required",
		},
		{
			name: "Valid app name",
			args: []string{},
			flags: map[string]string{
				"app": "test-app",
			},
			wantErr:  false,
			wantText: "Deployments for application test-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := listDeploymentsCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetArgs(tt.args)
			cmd.SetErr(b)

			for flag, value := range tt.flags {
				err := cmd.Flags().Set(flag, value)
				assert.NoError(t, err)
			}

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				if tt.wantText != "" {
					assert.Contains(t, b.String(), tt.wantText)
				}
			}
		})
	}
}
