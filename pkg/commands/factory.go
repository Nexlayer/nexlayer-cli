package commands

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/debug"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/plugin"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/di"
	"github.com/Nexlayer/nexlayer-cli/pkg/plugins"
)

// Factory creates and configures commands
type Factory struct {
	container *di.Container
	registry  *registry.Registry
	plugins   *plugins.Manager
}

// NewFactory creates a new command factory
func NewFactory(container *di.Container) *Factory {
	// Create dependencies for commands and plugins
	deps := &registry.CommandDependencies{
		APIClient:        container.GetAPIClient(),
		Logger:           container.GetLogger(),
		UIManager:        container.GetUIManager(),
		MetricsCollector: container.GetMetricsCollector(),
	}

	// Create plugin manager
	pluginDeps := &plugins.PluginDependencies{
		APIClient:        deps.APIClient,
		Logger:           deps.Logger,
		UIManager:        deps.UIManager,
		MetricsCollector: deps.MetricsCollector,
	}
	pluginManager := plugins.NewManager(pluginDeps, container.GetConfig().PluginsDir)

	// Create command registry
	reg := registry.NewRegistry(deps)

	return &Factory{
		container: container,
		registry:  reg,
		plugins:   pluginManager,
	}
}

// CreateRootCommand creates the root command with all subcommands
func (f *Factory) CreateRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy applications to Nexlayer",
		Long: `Nexlayer CLI helps you deploy and manage your applications on Nexlayer.
	
Key features:
- Easy application deployment
- Custom domain management
- Deployment status monitoring
- Deployment assistance
- Plugin system for extensibility

Need help? Use 'nexlayer debug' for deployment assistance.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize context with trace ID
			ctx := context.WithValue(cmd.Context(), "trace_id", time.Now().Format("20060102150405"))
			cmd.SetContext(ctx)

			// Log command execution
			logger := f.container.GetLogger()
			logger.Info(ctx, "Executing command: %s %v", cmd.Name(), args)

			// Record metric
			metrics := f.container.GetMetricsCollector()
			metrics.Counter("command_executions_total", 1, map[string]string{
				"command": cmd.Name(),
			})
		},
	}

	// Register built-in command providers
	f.registerBuiltinCommands()

	// Load plugins
	if err := f.plugins.LoadPluginsFromDir(""); err != nil {
		f.container.GetLogger().Error(nil, "Failed to load plugins: %v", err)
	}

	// Add all commands from registry and plugins
	cmd.AddCommand(f.getAllCommands()...)

	return cmd
}

// registerBuiltinCommands registers all built-in command providers
func (f *Factory) registerBuiltinCommands() {
	// Create command dependencies
	deps := &registry.CommandDependencies{
		APIClient:        f.container.GetAPIClient(),
		Logger:           f.container.GetLogger(),
		UIManager:        f.container.GetUIManager(),
		MetricsCollector: f.container.GetMetricsCollector(),
	}

	// Register core command providers
	providers := []registry.CommandProvider{
		deploy.NewProvider(),
		domain.NewProvider(),
		list.NewProvider(),
		status.NewProvider(),
		debug.NewProvider(),
		initcmd.NewProvider(),
		ai.NewProvider(),
		plugin.NewProvider(f.plugins),
	}

	// Register each provider
	for _, p := range providers {
		if err := f.registry.Register(p); err != nil {
			f.container.GetLogger().Error(nil, "Failed to register command provider %s: %v", p.Name(), err)
		}
	}
}

// getAllCommands returns all commands from both the registry and plugins
func (f *Factory) getAllCommands() []*cobra.Command {
	var commands []*cobra.Command

	// Get commands from registry
	commands = append(commands, f.registry.GetCommands()...)

	// Get commands from plugins
	commands = append(commands, f.plugins.GetCommands()...)

	return commands
}
