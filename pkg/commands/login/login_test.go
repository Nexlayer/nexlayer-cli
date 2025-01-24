package login

import (
	"bytes"
	"os"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginCommand(t *testing.T) {
	// Set test mode
	oldTestMode := os.Getenv("NEXLAYER_TEST_MODE")
	oldConfigDir := os.Getenv("NEXLAYER_CONFIG_DIR")

	// Create temp config dir
	tmpDir := t.TempDir()
	os.Setenv("NEXLAYER_CONFIG_DIR", tmpDir)
	os.Setenv("NEXLAYER_TEST_MODE", "true")

	defer func() {
		os.Setenv("NEXLAYER_TEST_MODE", oldTestMode)
		os.Setenv("NEXLAYER_CONFIG_DIR", oldConfigDir)
	}()

	tests := []struct {
		name        string
		args        []string
		setup       func(t *testing.T)
		cleanup     func(t *testing.T)
		wantErr     bool
		errContains string
		wantOutput  string
	}{
		{
			name:        "No args",
			args:        []string{},
			wantErr:     true,
			errContains: "token is required",
		},
		{
			name:       "With token flag",
			args:       []string{"--token", "valid-token"},
			wantErr:    false,
			wantOutput: "Successfully logged in",
		},
		{
			name:        "Invalid token (empty)",
			args:        []string{"--token", ""},
			wantErr:     true,
			errContains: "token is required",
		},
		{
			name:        "Invalid token (whitespace)",
			args:        []string{"--token", "   "},
			wantErr:     true,
			errContains: "token is required",
		},
		{
			name:    "Short flag",
			args:    []string{"-t", "valid-token"},
			wantErr: false,
		},
		{
			name: "Config directory not writable",
			args: []string{"--token", "valid-token"},
			setup: func(t *testing.T) {
				// Create a subdirectory in the temp dir
				readOnlyDir := tmpDir + "/readonly"
				require.NoError(t, os.MkdirAll(readOnlyDir, 0755))

				// Make it read-only
				require.NoError(t, os.Chmod(readOnlyDir, 0500))

				// Set it as the config dir
				os.Setenv("NEXLAYER_CONFIG_DIR", readOnlyDir)
			},
			cleanup: func(t *testing.T) {
				// Restore write permissions for cleanup
				readOnlyDir := tmpDir + "/readonly"
				_ = os.Chmod(readOnlyDir, 0755)
				os.Setenv("NEXLAYER_CONFIG_DIR", tmpDir)
			},
			wantErr:     true,
			errContains: "failed to save token",
		},
		{
			name: "Already logged in",
			args: []string{"--token", "new-token"},
			setup: func(t *testing.T) {
				require.NoError(t, auth.SaveToken("existing-token"))
			},
			wantErr:    false,
			wantOutput: "Successfully logged in",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run setup if provided
			if tt.setup != nil {
				tt.setup(t)
			}

			// Cleanup after test
			if tt.cleanup != nil {
				defer tt.cleanup(t)
			}

			cmd := NewCommand()
			b := bytes.NewBufferString("")
			cmd.SetOut(b)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.wantOutput != "" {
					assert.Contains(t, b.String(), tt.wantOutput)
				}

				// Verify token was saved
				if len(tt.args) > 1 {
					savedToken, err := auth.LoadToken()
					assert.NoError(t, err)
					assert.Equal(t, tt.args[1], savedToken)
				}
			}
		})
	}
}

func TestNewCommand(t *testing.T) {
	// Set test mode
	oldTestMode := os.Getenv("NEXLAYER_TEST_MODE")
	os.Setenv("NEXLAYER_TEST_MODE", "true")
	defer os.Setenv("NEXLAYER_TEST_MODE", oldTestMode)

	t.Run("Login command", func(t *testing.T) {
		cmd := NewCommand()
		assert.Equal(t, "login", cmd.Use)
		assert.Equal(t, "Login to Nexlayer", cmd.Short)
		assert.Equal(t, "Login to Nexlayer using your API token", cmd.Long)
		assert.NotNil(t, cmd.RunE)

		// Test flag existence
		tokenFlag := cmd.Flag("token")
		assert.NotNil(t, tokenFlag)
		assert.Equal(t, "token", tokenFlag.Name)
		assert.Equal(t, "t", tokenFlag.Shorthand)
		assert.Equal(t, "API token", tokenFlag.Usage)
	})

	t.Run("Help text", func(t *testing.T) {
		cmd := NewCommand()
		assert.Contains(t, cmd.UsageString(), "--token")
		assert.Contains(t, cmd.UsageString(), "-t")
		assert.Contains(t, cmd.UsageString(), "API token")
	})
}
