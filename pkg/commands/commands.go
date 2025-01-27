package commands

import (
	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
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

Need help? Use 'nexlayer debug' for deployment assistance.`,
	}

	return cmd
}
