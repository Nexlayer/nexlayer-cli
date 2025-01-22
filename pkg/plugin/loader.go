package plugin

// Formatted with gofmt -s
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Nexlayer/nexlayer-cli/pkg/worker"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Plugin represents a loaded plugin with its metadata and execution path
type Plugin struct {
	Path     string
	Metadata struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Usage       string `json:"usage"`
		ExecPath    string `json:"exec_path,omitempty"`
		Checksum    string `json:"checksum,omitempty"`
	}
}
type Loader struct {
	pluginsDir string
	plugins    map[string]*Plugin
	mu         sync.RWMutex
	pool       *worker.Pool
}

func NewLoader(pluginsDir string) *Loader {
	return &Loader{
		pluginsDir: pluginsDir,
		plugins:    make(map[string]*Plugin),
		pool: worker.NewPool(worker.PoolConfig{
			MinWorkers:    2,
			MaxWorkers:    runtime.GOMAXPROCS(0),
			QueueSize:     100,
			ScaleInterval: time.Second * 5,
			IdleTimeout:   time.Minute,
		}),
	}
}
func (l *Loader) LoadPlugins(ctx context.Context) error {
	l.pool.Start()
	defer l.pool.Stop()
	entries, err := os.ReadDir(l.pluginsDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pluginPath := filepath.Join(l.pluginsDir, entry.Name())
		l.pool.Submit(func(ctx context.Context) error {
			return l.loadPlugin(ctx, pluginPath)
		})
	}
	// Collect results
	var loadErrors []error
	for err := range l.pool.Results() {
		if err != nil {
			loadErrors = append(loadErrors, err)
		}
	}
	if len(loadErrors) > 0 {
		return fmt.Errorf("failed to load some plugins: %v", loadErrors)
	}
	return nil
}
func (l *Loader) loadPlugin(ctx context.Context, pluginPath string) error {
	execPath := filepath.Join(pluginPath, "plugin")
	cmd := exec.CommandContext(ctx, execPath, "--describe")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get plugin metadata: %w", err)
	}
	plugin := &Plugin{
		Path: pluginPath,
	}
	if err := json.Unmarshal(output, &plugin.Metadata); err != nil {
		return fmt.Errorf("failed to parse plugin metadata: %w", err)
	}
	// Set the execution path
	plugin.Metadata.ExecPath = execPath
	l.mu.Lock()
	l.plugins[plugin.Metadata.Name] = plugin
	l.mu.Unlock()
	return nil
}
func (l *Loader) GetPlugin(name string) (*Plugin, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	plugin, exists := l.plugins[name]
	return plugin, exists
}
func (l *Loader) ListPlugins() []Plugin {
	l.mu.RLock()
	defer l.mu.RUnlock()
	plugins := make([]Plugin, 0, len(l.plugins))
	for _, p := range l.plugins {
		plugins = append(plugins, *p)
	}
	return plugins
}
