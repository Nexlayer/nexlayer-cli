package plugin

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/spf13/cobra"
)

// PluginMetadata represents the JSON structure that plugins must output when called with --describe
type PluginMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
	ExecPath    string `json:"execPath,omitempty"`
	Checksum    string `json:"checksum,omitempty"`
	Signature   string `json:"signature,omitempty"`
}

// Manager handles plugin discovery, validation, and registration
type Manager struct {
	pluginDir string
	plugins   map[string]*PluginMetadata
	mu        sync.RWMutex
}

// NewManager creates a new plugin manager instance
func NewManager() *Manager {
	pluginDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
	return &Manager{
		pluginDir: pluginDir,
		plugins:   make(map[string]*PluginMetadata),
	}
}

// calculateChecksum computes SHA-256 hash of the plugin binary
func (m *Manager) calculateChecksum(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash), nil
}

// validatePlugin performs security checks on the plugin
func (m *Manager) validatePlugin(meta *PluginMetadata, path string) error {
	// Calculate checksum of the binary
	checksum, err := m.calculateChecksum(path)
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	// If plugin provided a checksum, verify it matches
	if meta.Checksum != "" && meta.Checksum != checksum {
		return fmt.Errorf("checksum mismatch for plugin %s", meta.Name)
	}

	// TODO: Add signature verification when implemented
	return nil
}

// getPluginMetadata executes the plugin with --describe flag to get its metadata
func (m *Manager) getPluginMetadata(pluginPath string) (*PluginMetadata, error) {
	cmd := exec.Command(pluginPath, "--describe")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get plugin metadata: %w", err)
	}

	var meta PluginMetadata
	if err := json.Unmarshal(output, &meta); err != nil {
		return nil, fmt.Errorf("invalid plugin metadata: %w", err)
	}

	// Set the execution path if not provided
	if meta.ExecPath == "" {
		meta.ExecPath = pluginPath
	}

	return &meta, nil
}

// createPluginCommand creates a cobra.Command for the plugin
func (m *Manager) createPluginCommand(meta *PluginMetadata) *cobra.Command {
	return &cobra.Command{
		Use:   meta.Name,
		Short: meta.Description,
		Long:  fmt.Sprintf("%s\nVersion: %s\nUsage: %s", meta.Description, meta.Version, meta.Usage),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginCmd := exec.Command(meta.ExecPath, args...)
			pluginCmd.Stdin = os.Stdin
			pluginCmd.Stdout = os.Stdout
			pluginCmd.Stderr = os.Stderr
			return pluginCmd.Run()
		},
	}
}

// DiscoverAndLoad finds and loads all plugins in the plugin directory
func (m *Manager) DiscoverAndLoad(rootCmd *cobra.Command) error {
	// Ensure plugin directory exists
	if err := os.MkdirAll(m.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	files, err := os.ReadDir(m.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := filepath.Join(m.pluginDir, file.Name())

		// Check if file is executable
		info, err := file.Info()
		if err != nil {
			continue
		}
		if info.Mode()&0111 == 0 {
			continue // Skip non-executable files
		}

		// Get plugin metadata
		meta, err := m.getPluginMetadata(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to load plugin %s: %v\n", file.Name(), err)
			continue
		}

		// Validate plugin
		if err := m.validatePlugin(meta, path); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Plugin validation failed for %s: %v\n", file.Name(), err)
			continue
		}

		// Store plugin metadata
		m.plugins[meta.Name] = meta

		// Create and add the plugin command
		pluginCmd := m.createPluginCommand(meta)
		rootCmd.AddCommand(pluginCmd)
	}

	return nil
}
