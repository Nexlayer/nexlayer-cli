// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
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
// This interface provides methods for name, description, version, additional commands,
// initialization with dependencies, and execution.
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

// PluginDependencies contains the dependencies available to plugins.
type PluginDependencies struct {
	APIClient api.APIClient
	Logger    *observability.Logger
	UIManager ui.Manager
}

// Manager is responsible for loading, initializing, and managing plugins.
type Manager struct {
	mu            sync.RWMutex
	plugins       map[string]Plugin   // Registered plugins by name.
	loadedPlugins map[string]bool     // Tracks already loaded plugin files.
	dependencies  *PluginDependencies // Shared dependencies for plugins.
	pluginsDir    string              // Default directory to load plugins from.
}

// NewManager creates and returns a new Manager instance.
func NewManager(deps *PluginDependencies, pluginsDir string) *Manager {
	return &Manager{
		plugins:       make(map[string]Plugin),
		loadedPlugins: make(map[string]bool),
		dependencies:  deps,
		pluginsDir:    pluginsDir,
	}
}

// LoadPlugin loads and initializes a plugin from the specified path.
// It ensures the plugin is not loaded more than once.
func (m *Manager) LoadPlugin(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.loadedPlugins[path] {
		// Plugin already loaded; nothing to do.
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
		return fmt.Errorf("plugin %s does not implement the Plugin interface", path)
	}

	// Initialize the plugin with the provided dependencies.
	if err := p.Init(m.dependencies); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", path, err)
	}

	// Save the plugin by name and mark the file as loaded.
	m.plugins[p.Name()] = p
	m.loadedPlugins[path] = true
	return nil
}

// LoadPluginsFromDir loads all plugins (files ending in .so) from the specified directory.
// If no directory is provided, the default plugins directory is used.
func (m *Manager) LoadPluginsFromDir(dir string) error {
	if dir == "" {
		dir = m.pluginsDir
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".so" {
			path := filepath.Join(dir, entry.Name())
			if err := m.LoadPlugin(path); err != nil {
				// Log error but continue loading other plugins.
				m.dependencies.Logger.Error(context.Background(), "Failed to load plugin %s: %v", path, err)
			}
		}
	}

	return nil
}

// GetPlugin returns the plugin instance by name.
func (m *Manager) GetPlugin(name string) (Plugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.plugins[name]
	return p, ok
}

// ListPlugins returns a map of plugin names and their versions.
func (m *Manager) ListPlugins() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make(map[string]string, len(m.plugins))
	for name, p := range m.plugins {
		result[name] = p.Version()
	}
	return result
}

// GetCommands aggregates and returns all CLI commands provided by loaded plugins.
func (m *Manager) GetCommands() []*cobra.Command {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var cmds []*cobra.Command
	for _, p := range m.plugins {
		cmds = append(cmds, p.Commands()...)
	}
	return cmds
}

// RunPlugin executes the specified plugin with the given options.
func (m *Manager) RunPlugin(name string, opts map[string]interface{}) error {
	m.mu.RLock()
	p, ok := m.plugins[name]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("plugin %s not found", name)
	}
	return p.Run(opts)
}
