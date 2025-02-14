// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
)

// Provider represents an AI code assistant provider.
type Provider interface {
	// GenerateText generates text using the AI provider.
	GenerateText(ctx context.Context, prompt string) (string, error)
}

// DefaultProvider is a simple implementation of the Provider interface.
type DefaultProvider struct{}

// GenerateText generates text using the default provider.
func (p *DefaultProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	// Return a basic template matching the new Nexlayer schema.
	return `application:
  name: "generated-app"
  pods:
    - name: "app"
      image: "nginx:latest"
      servicePorts: [80]`, nil
}

// NewDefaultProvider creates a new default provider.
func NewDefaultProvider() Provider {
	return &DefaultProvider{}
}

// GetPreferredProvider returns the preferred AI provider based on environment and capabilities.
func GetPreferredProvider(ctx context.Context, cap Capability) *AIProvider {
	// Check for Windsurf Editor
	if os.Getenv("WINDSURF_EDITOR_ACTIVE") == "true" {
		return &AIProvider{
			Name:         "Windsurf Editor",
			Description:  "AI-powered code editor",
			EnvVarKey:    "WINDSURF_EDITOR_ACTIVE",
			Capabilities: CapCodeGeneration | CapCodeCompletion | CapDeploymentAssistance,
		}
	}

	// Add more provider checks here if needed

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
