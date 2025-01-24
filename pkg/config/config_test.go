package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()
	assert.NotNil(t, cfg)
	assert.NotNil(t, cfg.APIEndpoints)
	assert.Equal(t, "https://app.nexlayer.io", cfg.APIEndpoints["production"])
	assert.Equal(t, "https://app.staging.nexlayer.io", cfg.APIEndpoints["staging"])
	assert.Equal(t, "https://app.staging.nexlayer.io", cfg.APIEndpoints["default"])
}

func TestGetAPIEndpoint(t *testing.T) {
	cfg := GetConfig()

	tests := []struct {
		name string
		env  string
		want string
	}{
		{
			name: "Production environment",
			env:  "production",
			want: "https://app.nexlayer.io",
		},
		{
			name: "Staging environment",
			env:  "staging",
			want: "https://app.staging.nexlayer.io",
		},
		{
			name: "Default environment",
			env:  "default",
			want: "https://app.staging.nexlayer.io",
		},
		{
			name: "Unknown environment",
			env:  "unknown",
			want: "https://app.staging.nexlayer.io", // Should return default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.GetAPIEndpoint(tt.env)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetConfigDir(t *testing.T) {
	tests := []struct {
		name           string
		configDirEnv   string
		homeEnv        string
		expectedPrefix string
	}{
		{
			name:           "Custom config directory",
			configDirEnv:   "/custom/config/dir",
			homeEnv:        "/home/user",
			expectedPrefix: "/custom/config/dir",
		},
		{
			name:           "Default config directory",
			configDirEnv:   "",
			homeEnv:        "/home/user",
			expectedPrefix: "/home/user/.nexlayer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original environment
			oldConfigDir := os.Getenv("NEXLAYER_CONFIG_DIR")
			oldHome := os.Getenv("HOME")
			defer func() {
				os.Setenv("NEXLAYER_CONFIG_DIR", oldConfigDir)
				os.Setenv("HOME", oldHome)
			}()

			// Set test environment
			os.Setenv("NEXLAYER_CONFIG_DIR", tt.configDirEnv)
			os.Setenv("HOME", tt.homeEnv)

			got := GetConfigDir()
			assert.Equal(t, tt.expectedPrefix, got)
		})
	}
}
