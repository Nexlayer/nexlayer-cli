// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package initcmd

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/spf13/cobra"
)

// Provider struct represents the command provider
type Provider struct{}

// NewProvider creates a new provider instance
func NewProvider() registry.CommandProvider {
	return &Provider{}
}

// Name returns the provider's name
func (p *Provider) Name() string {
	return "init"
}

// Description returns the provider's description
func (p *Provider) Description() string {
	return "Provides commands for initializing new Nexlayer projects"
}

// Dependencies returns an empty dependency list
func (p *Provider) Dependencies() []string {
	return nil
}

// Commands returns the list of available commands
func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		NewCommand(),
	}
}
