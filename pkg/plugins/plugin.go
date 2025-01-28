package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
)

// Plugin represents a Nexlayer plugin
type Plugin interface {
	// Name returns the name of the plugin
	Name() string
	
	// Description returns a description of what the plugin does
	Description() string
	
	// Run executes the plugin with the given options
	Run(opts map[string]interface{}) error
}

// Manager handles plugin loading and execution
type Manager struct {
	plugins map[string]Plugin
}

// NewManager creates a new plugin manager
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

// LoadPlugin loads a plugin from the given path
func (m *Manager) LoadPlugin(path string) error {
	// Open the plugin
	plug, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// Look up the Plugin symbol
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'Plugin' symbol: %w", path, err)
	}

	// Assert that the symbol is a Plugin
	plugin, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}

	// Store the plugin
	m.plugins[plugin.Name()] = plugin
	return nil
}

// LoadPluginsFromDir loads all plugins from the given directory
func (m *Manager) LoadPluginsFromDir(dir string) error {
	// Get plugin directory
	pluginDir := os.ExpandEnv(dir)
	if pluginDir == "" {
		pluginDir = filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
	}

	// Create plugin directory if it doesn't exist
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Walk the plugin directory
	return filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-.so files
		if info.IsDir() || filepath.Ext(path) != ".so" {
			return nil
		}

		// Load the plugin
		return m.LoadPlugin(path)
	})
}

// GetPlugin returns a plugin by name
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	plugin, ok := m.plugins[name]
	return plugin, ok
}

// ListPlugins returns a list of loaded plugin names
func (m *Manager) ListPlugins() []string {
	var names []string
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

// RunPlugin runs a plugin by name with the given options
func (m *Manager) RunPlugin(name string, opts map[string]interface{}) error {
	plugin, ok := m.GetPlugin(name)
	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}
	return plugin.Run(opts)
}
