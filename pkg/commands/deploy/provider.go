// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package deploy

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/spf13/cobra"
)

type Provider struct{}

func NewProvider() registry.CommandProvider {
	return &Provider{}
}

// NewCommand creates a new deploy command
func (p *Provider) NewCommand(deps *registry.CommandDependencies) *cobra.Command {
	return NewCommand(deps.APIClient)
}

func (p *Provider) Name() string {
	return "deploy"
}

func (p *Provider) Description() string {
	return "Provides commands for deploying applications to Nexlayer"
}

func (p *Provider) Dependencies() []string {
	// Deploy command has no dependencies on other commands
	return nil
}

func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		NewCommand(deps.APIClient),
	}
}
