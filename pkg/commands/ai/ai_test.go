package ai

import (
	"bytes"
	"context"
	"os"
	"strings"
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
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		want    string
	}{
		{
			name:    "generate template",
			args:    []string{"generate", "myapp"},
			wantErr: false,
			want:    "Successfully generated",
		},
		{
			name:    "missing app name",
			args:    []string{"generate"},
			wantErr: true,
			want:    "requires exactly 1 arg(s)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCommand()
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tt.want)
				}
				return
			}

			assert.NoError(t, err)
			assert.Contains(t, buf.String(), tt.want)
		})
	}
}

func TestDetectCommand(t *testing.T) {
	cmd := NewCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"detect"})

	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.True(t, strings.Contains(output, "No AI assistants detected") ||
		strings.Contains(output, "Detected AI assistant"))
}

func TestGetPreferredProvider(t *testing.T) {
	tests := []struct {
		name     string
		envVars  map[string]string
		wantName string
		wantNil  bool
	}{
		{
			name: "windsurf editor available",
			envVars: map[string]string{
				"WINDSURF_EDITOR_ACTIVE": "true",
			},
			wantName: "Windsurf Editor",
			wantNil:  false,
		},
		{
			name:     "no providers available",
			envVars:  map[string]string{},
			wantName: "",
			wantNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original env
			origEnv := make(map[string]string)
			for k := range tt.envVars {
				origEnv[k] = os.Getenv(k)
			}

			// Set test env
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			// Restore env after test
			defer func() {
				for k, v := range origEnv {
					os.Setenv(k, v)
				}
			}()

			provider := GetPreferredProvider(context.Background(), CapDeploymentAssistance)
			if tt.wantNil {
				assert.Nil(t, provider)
			} else {
				assert.NotNil(t, provider)
				assert.Equal(t, tt.wantName, provider.Name)
			}
		})
	}
}
