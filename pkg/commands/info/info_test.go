package info

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoCommand(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:        "No args",
			args:        []string{},
			wantErr:     true,
			errContains: "requires at least 1 arg(s)",
		},
		{
			name:        "Invalid app name",
			args:        []string{""},
			wantErr:     true,
			errContains: "application name is required",
		},
		{
			name:       "Valid app name",
			args:       []string{"test-app"},
			wantErr:    false,
			wantOutput: "Application Information",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewInfoCmd()
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

func TestNewInfoCmd(t *testing.T) {
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
			name:     "Valid app name",
			args:     []string{"test-app"},
			wantErr:  false,
			wantText: "Application: test-app",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewInfoCmd()
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
