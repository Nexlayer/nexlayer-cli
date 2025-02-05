// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package config

import (
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	APIEndpoints map[string]string
	PluginsDir   string
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	configDir := GetConfigDir()
	return &Config{
		APIEndpoints: map[string]string{
			"production": "https://app.nexlayer.io",
			"staging":    "https://app.staging.nexlayer.io",
			"default":    "https://app.staging.nexlayer.io",
		},
		PluginsDir: filepath.Join(configDir, "plugins"),
	}
}

// GetAPIEndpoint returns the API endpoint for the given environment
func (c *Config) GetAPIEndpoint(env string) string {
	if endpoint, ok := c.APIEndpoints[env]; ok {
		return endpoint
	}
	return c.APIEndpoints["default"] // default to staging
}

// GetConfigDir returns the configuration directory
func GetConfigDir() string {
	configDir := os.Getenv("NEXLAYER_CONFIG_DIR")
	if configDir == "" {
		configDir = filepath.Join(os.Getenv("HOME"), ".nexlayer")
	}
	return configDir
}

// GetPluginsDir returns the plugins directory
func (c *Config) GetPluginsDir() string {
	if c.PluginsDir == "" {
		return filepath.Join(GetConfigDir(), "plugins")
	}
	return c.PluginsDir
}
