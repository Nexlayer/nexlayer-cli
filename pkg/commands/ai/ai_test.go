package ai

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "ai [subcommand]", cmd.Use)
	assert.Equal(t, "AI-powered features for Nexlayer", cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Check that subcommands exist
	generateCmd, _, err := cmd.Find([]string{"generate"})
	assert.NoError(t, err)
	assert.NotNil(t, generateCmd)
	assert.Equal(t, "generate <app-name>", generateCmd.Use)

	detectCmd, _, err := cmd.Find([]string{"detect"})
	assert.NoError(t, err)
	assert.NotNil(t, detectCmd)
	assert.Equal(t, "detect", detectCmd.Use)
}

func TestGenerateCommand(t *testing.T) {
	cmd := NewCommand()
	generateCmd, _, _ := cmd.Find([]string{"generate"})

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "generate template",
			args:    []string{"myapp"},
			wantErr: false,
		},
		{
			name:    "missing app name",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			generateCmd.SetOut(buf)
			generateCmd.SetArgs(tt.args)

			err := generateCmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, buf.String(), "Successfully generated")

			// Clean up generated file
			if len(tt.args) > 0 {
				os.Remove("nexlayer.yaml")
			}
		})
	}
}

func TestDetectCommand(t *testing.T) {
	cmd := NewCommand()
	detectCmd, _, _ := cmd.Find([]string{"detect"})

	buf := new(bytes.Buffer)
	detectCmd.SetOut(buf)

	err := detectCmd.Execute()
	assert.NoError(t, err)
	output := buf.String()

	// Should either detect an AI assistant or show "No AI assistants detected"
	assert.True(t, 
		assert.Contains(t, output, "No AI assistants detected") ||
		assert.Contains(t, output, "Detected AI assistant"),
	)
}

func TestGetPreferredProvider(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		capability  Capability
		wantNil    bool
	}{
		{
			name: "windsurf editor available",
			envVars: map[string]string{
				"WINDSURF_EDITOR_ACTIVE": "true",
			},
			capability: CapDeploymentAssistance,
			wantNil:    false,
		},
		{
			name: "no providers available",
			envVars: map[string]string{},
			capability: CapDeploymentAssistance,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			provider := GetPreferredProvider(context.Background(), tt.capability)
			if tt.wantNil {
				assert.Nil(t, provider)
			} else {
				assert.NotNil(t, provider)
				assert.True(t, provider.Capabilities&tt.capability != 0)
			}
		})
	}
}
