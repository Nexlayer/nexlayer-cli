package plugins

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// Plugin is the interface that all Nexlayer plugins must implement.
type Plugin interface {
	// Name returns the plugin's name.
	Name() string
	// Description returns a description of what the plugin does.
	Description() string
	// Version returns the plugin version.
	Version() string
	// Commands returns any additional CLI commands provided by the plugin.
	Commands() []*cobra.Command
	// Init initializes the plugin with dependencies.
	Init(deps *PluginDependencies) error
	// Run executes the plugin with the given options.
	Run(opts map[string]interface{}) error
}

// PluginDependencies contains dependencies available to plugins.
type PluginDependencies struct {
	APIClient api.APIClient
	Logger    *observability.Logger
	UIManager ui.Manager
}

// Manager handles plugin loading, initialization, and execution.
type Manager struct {
	mu            sync.RWMutex
	plugins       map[string]Plugin
	loadedPlugins map[string]bool
	dependencies  *PluginDependencies
	pluginsDir    string
}

// NewManager creates a new plugin manager.
func NewManager(deps *PluginDependencies, pluginsDir string) *Manager {
	return &Manager{
		plugins:       make(map[string]Plugin),
		loadedPlugins: make(map[string]bool),
		dependencies:  deps,
		pluginsDir:    pluginsDir,
	}
}

// LoadPlugin loads a plugin from the specified path.
func (m *Manager) LoadPlugin(path string) error {
	// Use a write lock to prevent concurrent writes.
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.loadedPlugins[path] {
		// Plugin already loaded; skip.
		return nil
	}

	// Open the plugin file.
	plug, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// Look up the "Plugin" symbol.
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("plugin %s does not export 'Plugin' symbol: %w", path, err)
	}

	// Assert that the symbol implements the Plugin interface.
	p, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("plugin %s does not implement Plugin interface", path)
	}

	// Initialize the plugin with our dependencies.
	if err := p.Init(m.dependencies); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", path, err)
	}

	// Save the plugin.
	m.plugins[p.Name()] = p
	m.loadedPlugins[path] = true
	return nil
}

// LoadPluginsFromDir loads all plugins (files with .so extension) from the given directory.
func (m *Manager) LoadPluginsFromDir(dir string) error {
	if dir == "" {
		dir = m.pluginsDir
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	// Iterate over all directory entries.
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".so" {
			path := filepath.Join(dir, entry.Name())
			if err := m.LoadPlugin(path); err != nil {
				// Log the error but continue loading remaining plugins.
				if m.dependencies != nil && m.dependencies.Logger != nil {
					m.dependencies.Logger.Error(context.TODO(), "Failed to load plugin %s: %v", path, err)
				}
			}
		}
	}

	return nil
}

// GetPlugin returns a plugin by name.
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.plugins[name]
	return p, ok
}

// ListPlugins returns a map of loaded plugin names and their versions.
func (m *Manager) ListPlugins() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]string, len(m.plugins))
	for name, p := range m.plugins {
		result[name] = p.Version()
	}
	return result
}

// GetCommands aggregates and returns all commands provided by loaded plugins.
func (m *Manager) GetCommands() []*cobra.Command {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var cmds []*cobra.Command
	for _, p := range m.plugins {
		cmds = append(cmds, p.Commands()...)
	}
	return cmds
}

// RunPlugin executes a plugin by name with the provided options.
func (m *Manager) RunPlugin(name string, opts map[string]interface{}) error {
	m.mu.RLock()
	p, ok := m.plugins[name]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}
	return p.Run(opts)
}
