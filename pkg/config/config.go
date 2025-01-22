package config

import (
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	APIEndpoints map[string]string
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	return &Config{
		APIEndpoints: map[string]string{
			"staging":    "https://app.staging.nexlayer.io",
			"production": "https://app.nexlayer.io",
		},
	}
}

// GetAPIEndpoint returns the API endpoint for the given environment
func (c *Config) GetAPIEndpoint(env string) string {
	if endpoint, ok := c.APIEndpoints[env]; ok {
		return endpoint
	}
	return c.APIEndpoints["staging"] // default to staging
}

// GetConfigDir returns the configuration directory
func GetConfigDir() string {
	configDir := os.Getenv("NEXLAYER_CONFIG_DIR")
	if configDir == "" {
		configDir = filepath.Join(os.Getenv("HOME"), ".nexlayer")
	}
	return configDir
}
