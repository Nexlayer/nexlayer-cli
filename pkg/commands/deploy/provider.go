package deploy

import (
	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
)

type Provider struct{}

func NewProvider() registry.CommandProvider {
	return &Provider{}
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
