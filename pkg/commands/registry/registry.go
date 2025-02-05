package registry

import (
	"fmt"
	"sort"
	"sync"

	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// CommandDependencies contains dependencies that can be injected into commands
type CommandDependencies struct {
	APIClient api.APIClient
	Logger    *observability.Logger
	UIManager ui.Manager
}

// CommandProvider defines an interface for modules that provide commands
type CommandProvider interface {
	// Name returns the unique name of this command provider
	Name() string

	// Description returns a description of what commands this provider offers
	Description() string

	// Dependencies returns a list of other provider names that this provider depends on
	Dependencies() []string

	// Commands returns the commands provided by this module
	Commands(deps *CommandDependencies) []*cobra.Command
}

// Registry manages command providers and their dependencies
type Registry struct {
	mu        sync.RWMutex
	providers map[string]CommandProvider
	order     []string // Dependency-ordered list of provider names
	deps      *CommandDependencies
}

// NewRegistry creates a new command registry
func NewRegistry(deps *CommandDependencies) *Registry {
	return &Registry{
		providers: make(map[string]CommandProvider),
		deps:      deps,
	}
}

// Register adds a command provider to the registry
func (r *Registry) Register(provider CommandProvider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := provider.Name()
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("command provider %s already registered", name)
	}

	r.providers[name] = provider
	return r.updateDependencyOrder()
}

// updateDependencyOrder updates the order slice based on provider dependencies
func (r *Registry) updateDependencyOrder() error {
	// Reset order
	r.order = make([]string, 0, len(r.providers))
	visited := make(map[string]bool)
	temp := make(map[string]bool)

	var visit func(name string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("circular dependency detected involving %s", name)
		}
		if visited[name] {
			return nil
		}
		temp[name] = true

		provider := r.providers[name]
		for _, dep := range provider.Dependencies() {
			if _, exists := r.providers[dep]; !exists {
				return fmt.Errorf("provider %s depends on unregistered provider %s", name, dep)
			}
			if err := visit(dep); err != nil {
				return err
			}
		}

		temp[name] = false
		visited[name] = true
		r.order = append(r.order, name)
		return nil
	}

	// Visit all providers
	for name := range r.providers {
		if err := visit(name); err != nil {
			return err
		}
	}

	return nil
}

// GetCommands returns all registered commands in dependency order
func (r *Registry) GetCommands() []*cobra.Command {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var commands []*cobra.Command
	for _, name := range r.order {
		provider := r.providers[name]
		commands = append(commands, provider.Commands(r.deps)...)
	}
	return commands
}

// ListProviders returns information about registered providers
func (r *Registry) ListProviders() []ProviderInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	providers := make([]ProviderInfo, 0, len(r.providers))
	for _, name := range r.order {
		provider := r.providers[name]
		providers = append(providers, ProviderInfo{
			Name:         name,
			Description:  provider.Description(),
			Dependencies: provider.Dependencies(),
		})
	}

	// Sort by name for consistent output
	sort.Slice(providers, func(i, j int) bool {
		return providers[i].Name < providers[j].Name
	})

	return providers
}

// ProviderInfo contains information about a command provider
type ProviderInfo struct {
	Name         string
	Description  string
	Dependencies []string
}
