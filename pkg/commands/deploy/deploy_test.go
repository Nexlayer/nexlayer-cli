package deploy

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDeployCmd(t *testing.T) {
	// Create a temporary YAML file for testing
	yamlContent := []byte("name: test-app\nversion: 1.0.0")
	tmpfile, err := os.CreateTemp("", "test*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(yamlContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		args     []string
		flags    map[string]string
		wantErr  bool
		errMsg   string
		wantText string
	}{
		{
			name:    "No flags",
			args:    []string{},
			wantErr: true,
			errMsg:  "required flag(s) \"app\", \"file\" not set",
		},
		{
			name: "Missing app flag",
			args: []string{},
			flags: map[string]string{
				"file": tmpfile.Name(),
			},
			wantErr: true,
			errMsg:  "required flag(s) \"app\" not set",
		},
		{
			name: "Missing file flag",
			args: []string{},
			flags: map[string]string{
				"app": "test-app",
			},
			wantErr: true,
			errMsg:  "required flag(s) \"file\" not set",
		},
		{
			name: "Invalid YAML file",
			args: []string{},
			flags: map[string]string{
				"app":  "test-app",
				"file": "nonexistent.yaml",
			},
			wantErr: true,
			errMsg:  "failed to read YAML file",
		},
		{
			name: "Valid flags",
			args: []string{},
			flags: map[string]string{
				"app":  "test-app",
				"file": tmpfile.Name(),
			},
			wantErr:  false,
			wantText: "Started deployment",
		},
		{
			name: "With AI flag",
			args: []string{},
			flags: map[string]string{
				"app":  "test-app",
				"file": tmpfile.Name(),
				"ai":   "true",
			},
			wantErr:  false,
			wantText: "Started deployment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewDeployCmd()
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
