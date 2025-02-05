// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package plugin

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/Nexlayer/nexlayer-cli/pkg/plugins"
	"github.com/spf13/cobra"
)

type Provider struct {
	pluginManager *plugins.Manager
}

func NewProvider(pluginManager *plugins.Manager) registry.CommandProvider {
	return &Provider{
		pluginManager: pluginManager,
	}
}

func (p *Provider) Name() string {
	return "plugin"
}

func (p *Provider) Description() string {
	return "Provides commands for managing Nexlayer plugins"
}

func (p *Provider) Dependencies() []string {
	return nil
}

func (p *Provider) Commands(deps *registry.CommandDependencies) []*cobra.Command {
	pluginCmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage Nexlayer plugins",
		Long: `Manage Nexlayer plugins:
- List installed plugins
- Install new plugins
- Remove plugins
- Run plugin commands`,
	}

	// List plugins command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			plugins := p.pluginManager.ListPlugins()
			if len(plugins) == 0 {
				fmt.Println("No plugins installed")
				return nil
			}

			fmt.Println("Installed plugins:")
			for name, version := range plugins {
				fmt.Printf("  %s (v%s)\n", name, version)
			}
			return nil
		},
	}

	// Install plugin command
	installCmd := &cobra.Command{
		Use:   "install [plugin-path]",
		Short: "Install a plugin from a path",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return p.pluginManager.LoadPlugin(args[0])
		},
	}

	// Run plugin command
	runCmd := &cobra.Command{
		Use:   "run [plugin-name] [args...]",
		Short: "Run a plugin command",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := make(map[string]interface{})
			if len(args) > 1 {
				opts["args"] = args[1:]
			}
			return p.pluginManager.RunPlugin(args[0], opts)
		},
	}

	pluginCmd.AddCommand(listCmd, installCmd, runCmd)
	return []*cobra.Command{pluginCmd}
}
