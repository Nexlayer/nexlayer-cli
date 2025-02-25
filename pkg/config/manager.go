// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package config provides backward compatibility with the old configuration system.
// New code should use the pkg/core/config package directly.
package config

import (
	"sync"

	coreconfig "github.com/Nexlayer/nexlayer-cli/pkg/core/config"
)

// Manager handles configuration management
type Manager struct {
	mu sync.RWMutex
}

// NewManager creates a new configuration manager
func NewManager() *Manager {
	return &Manager{}
}

// GetString gets a string value from configuration
func (m *Manager) GetString(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return coreconfig.GetConfigProvider().GetString(key)
}

// GetStringMap gets a string map from configuration
func (m *Manager) GetStringMap(key string) map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return coreconfig.GetConfigProvider().GetStringMapString(key)
}

// Set sets a configuration value
func (m *Manager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	coreconfig.GetConfigProvider().Set(key, value)
}

// Save saves the configuration to disk
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return coreconfig.SaveConfig()
}

// GetAPIEndpoint returns the API endpoint for the given environment
func (m *Manager) GetAPIEndpoint(env string) string {
	if env == "" {
		env = "production"
	}

	endpoints := m.GetStringMap("api.endpoints")
	if endpoint, ok := endpoints[env]; ok {
		return endpoint
	}

	// Default production endpoint
	return "https://api.nexlayer.com"
}

// GetDefaultNamespace returns the default namespace
func (m *Manager) GetDefaultNamespace() string {
	return coreconfig.GetDefaultNamespace()
}

// SetDefaultNamespace sets the default namespace
func (m *Manager) SetDefaultNamespace(namespace string) error {
	coreconfig.SetDefaultNamespace(namespace)
	return m.Save()
}
