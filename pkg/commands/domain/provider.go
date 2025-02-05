// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package domain

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

type Provider struct{}

func NewProvider() registry.CommandProvider {
	return &Provider{}
}

func (p *Provider) Name() string {
	return "domain"
}

func (p *Provider) Description() string {
	return "Provides commands for managing custom domains"
}

func (p *Provider) Dependencies() []string {
	return nil
}

func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	if apiClient, ok := deps.APIClient.(*api.Client); ok {
		return []*cobra.Command{
			NewCommand(apiClient),
		}
	}
	deps.Logger.Error(nil, "Invalid API client type for domain command")
	return nil
}
