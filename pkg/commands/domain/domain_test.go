package domain

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()
	assert.Equal(t, "domain", cmd.Use)
	assert.Equal(t, "Manage custom domains", cmd.Short)
	assert.True(t, len(cmd.Commands()) > 0, "Should have subcommands")
}

func TestAddCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		wantErr  bool
		errMsg   string
		wantText string
	}{
		{
			name:    "No args",
			args:    []string{},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:    "No app flag",
			args:    []string{"example.com"},
			wantErr: true,
			errMsg:  "required flag(s) \"app\" not set",
		},
		{
			name: "Valid args",
			args: []string{"example.com"},
			flags: map[string]string{
				"app": "test-app",
			},
			wantErr:  false,
			wantText: "Added custom domain example.com to application test-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newAddCmd()
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
			name:    "No args",
			args:    []string{},
			wantErr: true,
			errMsg:  "accepts 1 arg(s), received 0",
		},
		{
			name:     "Valid args",
			args:     []string{"test-app"},
			wantErr:  false,
			wantText: "DOMAIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newListCmd()
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

func TestRemoveCmd(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		errMsg   string
		wantText string
	}{
		{
			name:    "No args",
			args:    []string{},
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 0",
		},
		{
			name:    "Missing domain",
			args:    []string{"test-app"},
			wantErr: true,
			errMsg:  "accepts 2 arg(s), received 1",
		},
		{
			name:    "Valid args",
			args:    []string{"test-app", "example.com"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newRemoveCmd()
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
