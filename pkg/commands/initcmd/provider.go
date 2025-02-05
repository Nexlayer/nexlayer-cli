// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package initcmd

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/spf13/cobra"
)

type Provider struct{}

func NewProvider() registry.CommandProvider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "init"
}

func (p *Provider) Description() string {
	return "Provides commands for initializing new Nexlayer projects"
}

func (p *Provider) Dependencies() []string {
	return nil
}

func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		NewCommand(),
	}
}
