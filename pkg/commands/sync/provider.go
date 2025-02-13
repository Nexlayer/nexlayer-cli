// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sync

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/spf13/cobra"
)

type Provider struct{}

func NewProvider() registry.CommandProvider {
	return &Provider{}
}

// NewCommand creates a new sync command
func (p *Provider) NewCommand(deps *registry.CommandDependencies) *cobra.Command {
	return NewCommand(deps.APIClient)
}

func (p *Provider) Name() string {
	return "sync"
}

func (p *Provider) Description() string {
	return "Synchronize local configuration with Nexlayer"
}

func (p *Provider) Dependencies() []string {
	return nil
}

func (p *Provider) Commands() []*cobra.Command {
	return nil
}
