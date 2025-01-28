package commands

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/debug"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/init"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/wizard"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/di"
)

// Factory creates and configures commands
type Factory struct {
	container *di.Container
}

// NewFactory creates a new command factory
func NewFactory(container *di.Container) *Factory {
	return &Factory{
		container: container,
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

	// Get dependencies
	client := f.container.GetAPIClient()

	// Add subcommands
	cmd.AddCommand(
		initcmd.NewCommand(),
		f.createDeployCommand(client),
		f.createDomainCommand(client),
		f.createListCommand(client),
		f.createStatusCommand(client),
		f.createWizardCommand(client),
		f.createDebugCommand(client),
	)

	return cmd
}

func (f *Factory) createDeployCommand(client api.APIClient) *cobra.Command {
	if apiClient, ok := client.(*api.Client); ok {
		return deploy.NewCommand(apiClient)
	}
	// This should never happen as we control the DI container
	panic("invalid API client type")
}

func (f *Factory) createDomainCommand(client api.APIClient) *cobra.Command {
	if apiClient, ok := client.(*api.Client); ok {
		return domain.NewCommand(apiClient)
	}
	panic("invalid API client type")
}

func (f *Factory) createListCommand(client api.APIClient) *cobra.Command {
	if apiClient, ok := client.(*api.Client); ok {
		return list.NewCommand(apiClient)
	}
	panic("invalid API client type")
}

func (f *Factory) createStatusCommand(client api.APIClient) *cobra.Command {
	if apiClient, ok := client.(*api.Client); ok {
		return status.NewCommand(apiClient)
	}
	panic("invalid API client type")
}

func (f *Factory) createWizardCommand(client api.APIClient) *cobra.Command {
	if apiClient, ok := client.(*api.Client); ok {
		return wizard.NewCommand(apiClient)
	}
	panic("invalid API client type")
}

func (f *Factory) createDebugCommand(client api.APIClient) *cobra.Command {
	if apiClient, ok := client.(*api.Client); ok {
		return debug.NewCommand(apiClient)
	}
	panic("invalid API client type")
}
