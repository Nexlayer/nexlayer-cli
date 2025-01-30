package config

import (
	"sync"

	"github.com/spf13/viper"
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
	return viper.GetString(key)
}

// GetStringMap gets a string map from configuration
func (m *Manager) GetStringMap(key string) map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]string)
	for k, v := range viper.GetStringMap(key) {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}

// Set sets a configuration value
func (m *Manager) Set(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	viper.Set(key, value)
}

// Save saves the configuration to disk
func (m *Manager) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return viper.WriteConfig()
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
	return m.GetString("namespace")
}

// SetDefaultNamespace sets the default namespace
func (m *Manager) SetDefaultNamespace(namespace string) error {
	m.Set("namespace", namespace)
	return m.Save()
}
