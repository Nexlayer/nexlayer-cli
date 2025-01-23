package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/errors"
)

// Config represents the CLI configuration
type Config struct {
	RegistryURL     string            `json:"registry_url"`
	DefaultTemplate string            `json:"default_template"`
	APIKeys         map[string]string `json:"api_keys"`
	Verbose         bool              `json:"verbose"`
	OutputFormat    string            `json:"output_format"`
}

// Manager handles configuration loading and saving
type Manager struct {
	config     *Config
	configPath string
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	return &Config{
		RegistryURL:     "https://registry.nexlayer.dev",
		DefaultTemplate: "default",
		APIKeys:         make(map[string]string),
		Verbose:         false,
		OutputFormat:    "yaml",
	}
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.NewError(errors.ErrConfigNotFound, "could not find home directory", err)
	}

	configDir := filepath.Join(homeDir, ".nexlayer")
	configPath := filepath.Join(configDir, "config.json")

	return &Manager{
		configPath: configPath,
	}, nil
}

// Load loads the configuration from disk
func (m *Manager) Load() error {
	if err := os.MkdirAll(filepath.Dir(m.configPath), 0755); err != nil {
		return errors.NewError(errors.ErrConfigNotFound, "could not create config directory", err)
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			m.config = DefaultConfig()
			return m.Save()
		}
		return errors.NewError(errors.ErrConfigNotFound, "could not read config file", err)
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return errors.NewError(errors.ErrConfigInvalid, "could not parse config file", err)
	}

	m.config = config
	return nil
}

// Save saves the configuration to disk
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return errors.NewError(errors.ErrConfigInvalid, "could not marshal config", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return errors.NewError(errors.ErrConfigInvalid, "could not write config file", err)
	}

	return nil
}

// Get returns the current configuration
func (m *Manager) Get() *Config {
	return m.config
}

// Set updates the configuration
func (m *Manager) Set(config *Config) {
	m.config = config
}

// SetAPIKey sets an API key in the configuration
func (m *Manager) SetAPIKey(service, key string) error {
	m.config.APIKeys[service] = key
	return m.Save()
}

// GetAPIKey gets an API key from the configuration
func (m *Manager) GetAPIKey(service string) string {
	return m.config.APIKeys[service]
}
