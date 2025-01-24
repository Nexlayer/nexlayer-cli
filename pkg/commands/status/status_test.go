package status

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
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
			name: "With app flag",
			args: []string{},
			flags: map[string]string{
				"app": "test-app",
			},
			wantErr: true,
			errMsg:  "status check is not yet implemented",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := Command
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
