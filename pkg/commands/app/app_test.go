package app

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		wantErr  bool
		errMsg   string
		wantText string
	}{
		{
			name:    "No name flag",
			args:    []string{},
			wantErr: true,
			errMsg:  "required flag(s) \"name\" not set",
		},
		{
			name: "Empty name",
			args: []string{},
			flags: map[string]string{
				"name": "",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "Valid name",
			args: []string{},
			flags: map[string]string{
				"name": "test-app",
			},
			wantErr:  false,
			wantText: "Created application test-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := CreateCmd()
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

func TestListCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		errMsg   string
		wantText string
	}{
		{
			name:     "List apps",
			args:     []string{},
			wantErr:  false,
			wantText: "Applications:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := ListCmd()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetArgs(tt.args)
			cmd.SetErr(b)

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

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()
	assert.Equal(t, "app", cmd.Use)
	assert.Equal(t, "Manage your applications", cmd.Short)
	assert.True(t, len(cmd.Commands()) > 0, "Should have subcommands")
}
