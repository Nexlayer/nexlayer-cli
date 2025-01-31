package ai

import (
	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
)

type Provider struct{}

func NewProvider() registry.CommandProvider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "ai"
}

func (p *Provider) Description() string {
	return "Provides AI-powered commands for deployment assistance"
}

func (p *Provider) Dependencies() []string {
	return nil
}

func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	return []*cobra.Command{
		NewCommand(),
	}
}
