package plugin

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
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
