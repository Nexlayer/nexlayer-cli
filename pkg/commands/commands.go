package commands

import (
	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
)

// NewRootCommand creates a new root command with all subcommands
func NewRootCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy applications to Nexlayer",
		Long: `Nexlayer CLI helps you deploy and manage your applications on Nexlayer.
		
Key features:
- Easy application deployment
- Container registry management
- Resource monitoring
- Deployment assistance
- AI-powered features

Need help? Use 'nexlayer debug' for deployment assistance.`,
	}

	// Add subcommands
	cmd.AddCommand(deploy.NewCommand(client))
	cmd.AddCommand(status.NewCommand(client))
	cmd.AddCommand(domain.NewCommand(client))
	cmd.AddCommand(ai.NewCommand(client))

	return cmd
}
