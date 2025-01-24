package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPluginCmd(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:       "No args",
			args:       []string{},
			wantErr:    false,
			wantOutput: "Usage:",
		},
		{
			name:        "Invalid command",
			args:        []string{"invalid"},
			wantErr:     true,
			errContains: "unknown command",
		},
		{
			name:        "Install without name",
			args:        []string{"install"},
			wantErr:     true,
			errContains: "requires plugin name",
		},
		{
			name:       "List plugins",
			args:       []string{"list"},
			wantErr:    false,
			wantOutput: "No plugins installed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			cmd.SetArgs(tt.args)

			output := new(bytes.Buffer)
			cmd.SetOut(output)

			err := cmd.Execute()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Contains(t, output.String(), tt.wantOutput)
			}
		})
	}
}
