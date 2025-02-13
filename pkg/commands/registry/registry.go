// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package registry

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

// CommandDependencies holds common dependencies for commands
type CommandDependencies struct {
	APIClient *api.Client
}

// CommandProvider defines an interface for command providers
type CommandProvider interface {
	NewCommand(deps *CommandDependencies) *cobra.Command
}

// Provider provides command registration
type Provider struct {
	providers map[string]CommandProvider
}

// NewProvider creates a new command provider
func NewProvider() *Provider {
	return &Provider{
		providers: make(map[string]CommandProvider),
	}
}

// Register adds a command provider
func (p *Provider) Register(name string, provider CommandProvider) {
	p.providers[name] = provider
}

// GetCommand returns a command by name
func (p *Provider) GetCommand(name string, deps *CommandDependencies) *cobra.Command {
	if provider, ok := p.providers[name]; ok {
		return provider.NewCommand(deps)
	}
	return nil
}

// GetCommands returns all registered commands
func (p *Provider) GetCommands(deps *CommandDependencies) []*cobra.Command {
	var commands []*cobra.Command
	for _, provider := range p.providers {
		if cmd := provider.NewCommand(deps); cmd != nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}
