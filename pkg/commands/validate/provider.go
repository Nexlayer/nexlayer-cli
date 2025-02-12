// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validate

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/spf13/cobra"
)

type Provider struct{}

// NewProvider creates a new validate command provider
func NewProvider() registry.CommandProvider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "validate"
}

func (p *Provider) Description() string {
	return "Validate Nexlayer YAML configuration files"
}

func (p *Provider) Dependencies() []string {
	return nil
}

func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		NewCommand(),
	}
}
