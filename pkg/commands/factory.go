package commands

import (
	"context"
	"sync"
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
	"github.com/Nexlayer/nexlayer-cli/pkg/di"
	"github.com/Nexlayer/nexlayer-cli/pkg/plugins"
)

// Factory creates and configures commands.
type Factory struct {
	container *di.Container
	registry  *registry.Registry
	plugins   *plugins.Manager
}

// NewFactory creates a new command factory.
// It sets up command dependencies, the plugin manager, and the command registry.
func NewFactory(container *di.Container) *Factory {
	// Prepare dependencies for commands and plugins
	deps := &registry.CommandDependencies{
		APIClient: container.GetAPIClient(),
		Logger:    container.GetLogger(),
		UIManager: container.GetUIManager(),
	}

	// Initialize plugin manager with its dependencies.
	pluginDeps := &plugins.PluginDependencies{
		APIClient: deps.APIClient,
		Logger:    deps.Logger,
		UIManager: deps.UIManager,
	}
	pluginManager := plugins.NewManager(pluginDeps, container.GetConfig().PluginsDir)

	// Create a command registry using the dependencies.
	reg := registry.NewRegistry(deps)

	return &Factory{
		container: container,
		registry:  reg,
		plugins:   pluginManager,
	}
}

// CreateRootCommand creates the root command and attaches all subcommands.
// It uses concurrent plugin loading to reduce startup latency.
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
			// Create a context with a unique trace ID.
			traceID := time.Now().Format("20060102150405")
			ctx := context.WithValue(cmd.Context(), "trace_id", traceID)
			cmd.SetContext(ctx)

			// Cache the logger and metrics to avoid repeated container lookups.
			logger := f.container.GetLogger()
			logger.Info(ctx, "Executing command: %s %v", cmd.Name(), args)

			// Record the command execution metric.
			metrics := f.container.GetMetricsCollector()
			metrics.Counter("command_executions_total", 1, map[string]string{
				"command": cmd.Name(),
			})
		},
	}

	// Start concurrent plugin loading while registering built-in commands.
	var wg sync.WaitGroup
	var pluginLoadErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Load plugins concurrently. An empty string ("") implies using the default plugin directory.
		pluginLoadErr = f.plugins.LoadPluginsFromDir("")
	}()

	// Register built-in (core) command providers.
	f.registerBuiltinCommands()

	// Wait for plugin loading to finish.
	wg.Wait()
	if pluginLoadErr != nil {
		// Log the error but continue: built-in commands will still work.
		f.container.GetLogger().Error(nil, "Failed to load plugins: %v", pluginLoadErr)
	}

	// Add all commands from both the registry and loaded plugins.
	cmd.AddCommand(f.getAllCommands()...)

	return cmd
}

// registerBuiltinCommands registers all built-in command providers into the registry.
func (f *Factory) registerBuiltinCommands() {
	// List of built-in command providers.
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

	// Register each provider and log errors if registration fails.
	for _, p := range providers {
		if err := f.registry.Register(p); err != nil {
			f.container.GetLogger().Error(nil, "Failed to register command provider %s: %v", p.Name(), err)
		}
	}
}

// getAllCommands aggregates commands from the registry and plugins.
// The slice is pre-allocated to avoid multiple reallocations.
func (f *Factory) getAllCommands() []*cobra.Command {
	regCmds := f.registry.GetCommands()
	pluginCmds := f.plugins.GetCommands()

	total := len(regCmds) + len(pluginCmds)
	commands := make([]*cobra.Command, 0, total)
	commands = append(commands, regCmds...)
	commands = append(commands, pluginCmds...)

	return commands
}
