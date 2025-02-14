// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/spf13/cobra"
)

// Capability represents an AI provider's capabilities.
type Capability int

const (
	CapCodeGeneration Capability = 1 << iota
	CapCodeCompletion
	CapDeploymentAssistance
	CapErrorDiagnosis
)

// AIProvider represents an AI code assistant provider.
type AIProvider struct {
	Name         string
	Description  string
	EnvVarKey    string
	Capabilities Capability
}

// GenerateText generates text using the AI provider.
func (p *AIProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement provider-specific text generation.
	// For now, return a basic template matching the new Nexlayer schema.
	return `application:
  name: "generated-app"
  pods:
    - name: "app"
      image: "nginx:latest"
      servicePorts: [80]`, nil
}

// providerCache caches the detected provider.
type providerCache struct {
	sync.RWMutex
	provider   *AIProvider
	expiration time.Time
}

const cachePeriod = 5 * time.Minute

var (
	// cache is the global provider cache.
	cache = &providerCache{}

	// Predefined AI providers.
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

	ZedEditor = AIProvider{
		Name:         "Zed Editor",
		Description:  "Zed's built-in AI assistant",
		EnvVarKey:    "ZED_AI_ACTIVE",
		Capabilities: CapCodeGeneration | CapCodeCompletion | CapDeploymentAssistance,
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

	// AllProviders is a list of all available AI providers.
	AllProviders = []AIProvider{
		WindsurfEditor,
		GitHubCopilot,
		ZedEditor,
		CursorAI,
		VSCodeAI,
	}
)

// GetPreferredProvider returns the first configured AI provider with the required capabilities.
func GetPreferredProvider(ctx context.Context, requiredCaps Capability) *AIProvider {
	// Check cache first.
	cache.RLock()
	if time.Now().Before(cache.expiration) {
		provider := cache.provider
		cache.RUnlock()
		return provider
	}
	cache.RUnlock()

	// Define provider priority order.
	priorityProviders := []AIProvider{
		WindsurfEditor, // Windsurf has the most capabilities.
		ZedEditor,      // Zed has deployment assistance.
		CursorAI,       // Cursor is widely used.
		GitHubCopilot,  // Copilot is also popular.
		VSCodeAI,       // VSCode is common.
	}

	// Try each provider in priority order.
	for _, provider := range priorityProviders {
		if os.Getenv(provider.EnvVarKey) != "" && provider.Capabilities&requiredCaps == requiredCaps {
			// Log which provider we're using.
			if os.Getenv("NEXLAYER_DEBUG") != "" {
				fmt.Printf("Using AI provider: %s (%s)\n", provider.Name, provider.Description)
			}

			// Cache the result.
			cache.Lock()
			cache.provider = &provider
			cache.expiration = time.Now().Add(cachePeriod)
			cache.Unlock()
			return &provider
		}
	}

	// No provider found.
	if os.Getenv("NEXLAYER_DEBUG") != "" {
		fmt.Println("No AI provider found with required capabilities, using fallback template")
	}
	return nil
}

// CommandProvider implements the registry.CommandProvider interface.
type CommandProvider struct{}

// NewProvider creates a new AI command provider.
func NewProvider() *CommandProvider {
	return &CommandProvider{}
}

// Name returns the unique name of this command provider.
func (p *CommandProvider) Name() string {
	return "ai"
}

// Description returns a description of what commands this provider offers.
func (p *CommandProvider) Description() string {
	return "AI-powered features for generating and optimizing deployment templates"
}

// Dependencies returns a list of other provider names that this provider depends on.
func (p *CommandProvider) Dependencies() []string {
	return nil
}

// Commands returns the AI-related commands.
func (p *CommandProvider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		NewCommand(),
	}
}
