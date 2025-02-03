package ai

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
)

// Capability represents an AI provider's capabilities
type Capability int

const (
	CapCodeGeneration Capability = 1 << iota
	CapCodeCompletion
	CapDeploymentAssistance
	CapErrorDiagnosis
)

// AIProvider represents an AI code assistant provider
type AIProvider struct {
	Name         string
	Description  string
	EnvVarKey    string
	Capabilities Capability
}

// providerCache caches the detected provider
type providerCache struct {
	sync.RWMutex
	provider    *AIProvider
	expiration  time.Time
}

const cachePeriod = 5 * time.Minute

var (
	// cache is the global provider cache
	cache = &providerCache{}

	// Predefined AI providers
	WindsurfEditor = AIProvider{
		Name:         "Windsurf Editor",
		Description:  "Built-in AI code assistant",
		EnvVarKey:    "WINDSURF_EDITOR_ACTIVE",
		Capabilities: CapCodeGeneration | CapCodeCompletion | CapDeploymentAssistance | CapErrorDiagnosis,
	}

	GitHubCopilot = AIProvider{
		Name:         "GitHub Copilot",
		Description:  "GitHub's AI pair programmer",
		EnvVarKey:    "GITHUB_COPILOT_ACTIVE",
		Capabilities: CapCodeGeneration | CapCodeCompletion,
	}

	CursorAI = AIProvider{
		Name:         "Cursor AI",
		Description:  "Cursor's AI code assistant",
		EnvVarKey:    "CURSOR_AI_ACTIVE",
		Capabilities: CapCodeGeneration | CapCodeCompletion,
	}

	VSCodeAI = AIProvider{
		Name:         "VS Code AI",
		Description:  "VS Code's AI assistant",
		EnvVarKey:    "VSCODE_AI_ACTIVE",
		Capabilities: CapCodeGeneration | CapCodeCompletion,
	}

	// AllProviders is a list of all available AI providers
	AllProviders = []AIProvider{
		WindsurfEditor,
		GitHubCopilot,
		CursorAI,
		VSCodeAI,
	}
)

// GetPreferredProvider returns the first configured AI provider with the required capabilities
func GetPreferredProvider(ctx context.Context, requiredCaps Capability) *AIProvider {
	// Check cache first
	cache.RLock()
	if time.Now().Before(cache.expiration) {
		provider := cache.provider
		cache.RUnlock()
		return provider
	}
	cache.RUnlock()

	for _, provider := range AllProviders {
		if os.Getenv(provider.EnvVarKey) != "" && provider.Capabilities&requiredCaps == requiredCaps {
			// Cache the result
			cache.Lock()
			cache.provider = &provider
			cache.expiration = time.Now().Add(cachePeriod)
			cache.Unlock()
			return &provider
		}
	}
	return nil
}

// Provider implements the registry.CommandProvider interface
type Provider struct{}

// NewProvider creates a new AI command provider
func NewProvider() *Provider {
	return &Provider{}
}

// Name returns the unique name of this command provider
func (p *Provider) Name() string {
	return "ai"
}

// Description returns a description of what commands this provider offers
func (p *Provider) Description() string {
	return "AI-powered features for generating and optimizing deployment templates"
}

// Dependencies returns a list of other provider names that this provider depends on
func (p *Provider) Dependencies() []string {
	return nil
}

// Commands returns the AI-related commands
func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		// TODO: Add AI-related commands here
	}
}
