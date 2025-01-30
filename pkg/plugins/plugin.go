package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// Plugin represents a Nexlayer plugin with enhanced capabilities
type Plugin interface {
	// Name returns the name of the plugin
	Name() string

	// Description returns a description of what the plugin does
	Description() string

	// Version returns the plugin version
	Version() string

	// Commands returns any additional CLI commands provided by the plugin
	Commands() []*cobra.Command

	// Init initializes the plugin with dependencies
	Init(deps *PluginDependencies) error

	// Run executes the plugin with the given options
	Run(opts map[string]interface{}) error
}

// PluginDependencies contains all dependencies available to plugins
type PluginDependencies struct {
	APIClient        api.APIClient
	Logger           *observability.Logger
	UIManager        ui.Manager
	MetricsCollector *observability.MetricsCollector
}

// Manager handles plugin loading, initialization and execution
type Manager struct {
	mu            sync.RWMutex
	plugins       map[string]Plugin
	dependencies  *PluginDependencies
	pluginsDir    string
	loadedPlugins map[string]bool
}

// NewManager creates a new plugin manager
func NewManager(deps *PluginDependencies, pluginsDir string) *Manager {
	return &Manager{
		plugins:       make(map[string]Plugin),
		dependencies:  deps,
		pluginsDir:    pluginsDir,
		loadedPlugins: make(map[string]bool),
	}
}

// LoadPlugin loads a plugin from the given path
func (m *Manager) LoadPlugin(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if plugin is already loaded
	if m.loadedPlugins[path] {
		return nil
	}

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

	// Initialize plugin with dependencies
	if err := plugin.Init(m.dependencies); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", path, err)
	}

	// Store the plugin
	m.plugins[plugin.Name()] = plugin
	m.loadedPlugins[path] = true

	return nil
}

// LoadPluginsFromDir loads all plugins from the given directory
func (m *Manager) LoadPluginsFromDir(dir string) error {
	if dir == "" {
		dir = m.pluginsDir
	}

	// Walk through the plugins directory
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only load .so files
		if !info.IsDir() && filepath.Ext(path) == ".so" {
			if err := m.LoadPlugin(path); err != nil {
				// Log error but continue loading other plugins
				if m.dependencies != nil && m.dependencies.Logger != nil {
					m.dependencies.Logger.Error(nil, "Failed to load plugin %s: %v", path, err)
				}
			}
		}
		return nil
	})
}

// GetPlugin returns a plugin by name
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugin, ok := m.plugins[name]
	return plugin, ok
}

// ListPlugins returns a list of loaded plugin names and versions
func (m *Manager) ListPlugins() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	plugins := make(map[string]string)
	for name, plugin := range m.plugins {
		plugins[name] = plugin.Version()
	}
	return plugins
}

// GetCommands returns all commands from all loaded plugins
func (m *Manager) GetCommands() []*cobra.Command {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var commands []*cobra.Command
	for _, p := range m.plugins {
		commands = append(commands, p.Commands()...)
	}
	return commands
}

// RunPlugin runs a plugin by name with the given options
func (m *Manager) RunPlugin(name string, opts map[string]interface{}) error {
	m.mu.RLock()
	plugin, ok := m.plugins[name]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}

	return plugin.Run(opts)
}
