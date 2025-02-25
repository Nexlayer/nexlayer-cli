// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// DefaultConfigName is the default name for the configuration file
const DefaultConfigName = "config"

// DefaultConfigType is the default type for the configuration file
const DefaultConfigType = "yaml"

// DefaultConfigDir is the default directory for the configuration file
const DefaultConfigDir = ".config/nexlayer"

// providerKey is the key used to store the configuration provider in the context
var providerKey = struct{}{}

// Manager handles configuration loading and access
type Manager struct {
	provider Provider
}

// NewManager creates a new configuration manager with the given provider
func NewManager(provider Provider) *Manager {
	return &Manager{
		provider: provider,
	}
}

// DefaultManager creates a new configuration manager with the default Viper provider
func DefaultManager() *Manager {
	return NewManager(NewViperProvider())
}

// Provider returns the configuration provider
func (m *Manager) Provider() Provider {
	return m.provider
}

// WithContext returns a new context with the configuration provider
func (m *Manager) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, providerKey, m.provider)
}

// FromContext retrieves the configuration provider from the context
func FromContext(ctx context.Context) (Provider, bool) {
	provider, ok := ctx.Value(providerKey).(Provider)
	return provider, ok
}

// MustFromContext retrieves the configuration provider from the context or panics
func MustFromContext(ctx context.Context) Provider {
	provider, ok := FromContext(ctx)
	if !ok {
		panic("config provider not found in context")
	}
	return provider
}

// InitConfig initializes the configuration with default values and paths
func (m *Manager) InitConfig(configPath string) error {
	p, ok := m.provider.(*ViperProvider)
	if !ok {
		return fmt.Errorf("provider does not support initialization")
	}

	// Set defaults for configuration
	p.SetDefault("nexlayer.api_url", "https://api.nexlayer.io")
	p.SetDefault("nexlayer.port", 8080)

	// If a config file is provided, use it
	if configPath != "" {
		p.SetConfigFile(configPath)
	} else {
		// Otherwise, look in default locations
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting user home directory: %w", err)
		}

		// Add default config paths
		p.AddConfigPath(".")
		p.AddConfigPath(filepath.Join(homeDir, DefaultConfigDir))
		p.SetConfigName(DefaultConfigName)
		p.SetConfigType(DefaultConfigType)
	}

	// Enable environment variable overrides
	p.AutomaticEnv()
	p.SetEnvPrefix("NEXLAYER")

	// Try to read the config file
	if err := p.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	return nil
}

// SaveConfig saves the current configuration to disk
func (m *Manager) SaveConfig() error {
	p, ok := m.provider.(*ViperProvider)
	if !ok {
		return fmt.Errorf("provider does not support saving")
	}

	return p.WriteConfig()
}

// GetConfigDir returns the directory where the configuration file is located
func (m *Manager) GetConfigDir() (string, error) {
	p, ok := m.provider.(*ViperProvider)
	if !ok {
		return "", fmt.Errorf("provider does not support config file path")
	}

	configFile := p.ConfigFileUsed()
	if configFile == "" {
		return "", fmt.Errorf("no config file used")
	}

	return filepath.Dir(configFile), nil
}

// GetAPIURL returns the API URL from the configuration
func (m *Manager) GetAPIURL() string {
	return m.provider.GetString("nexlayer.api_url")
}

// SetAPIURL sets the API URL in the configuration
func (m *Manager) SetAPIURL(url string) {
	m.provider.Set("nexlayer.api_url", url)
}

// GetToken returns the authentication token from the configuration
func (m *Manager) GetToken() string {
	return m.provider.GetString("nexlayer.token")
}

// SetToken sets the authentication token in the configuration
func (m *Manager) SetToken(token string) {
	m.provider.Set("nexlayer.token", token)
}

// GetDefaultNamespace returns the default namespace from the configuration
func (m *Manager) GetDefaultNamespace() string {
	return m.provider.GetString("nexlayer.default_namespace")
}

// SetDefaultNamespace sets the default namespace in the configuration
func (m *Manager) SetDefaultNamespace(namespace string) {
	m.provider.Set("nexlayer.default_namespace", namespace)
}
