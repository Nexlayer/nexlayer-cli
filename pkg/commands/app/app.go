package app

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

var Cmd *cobra.Command

func init() {
	client := api.NewClient("https://api.nexlayer.com")
	Cmd = NewCommand(client)
}

// NewCommand creates a new app command
func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Manage your applications",
		Long:  `View and manage your Nexlayer applications.`,
	}

	cmd.AddCommand(
		newInfoCommand(client),
	)

	return cmd
}

func newInfoCommand(client *api.Client) *cobra.Command {
	var appID string
	var namespace string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get application info",
		Long: `Get detailed information about an application deployment.
		
Example:
  nexlayer app info --app myapp --namespace ecstatic-frog`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(cmd, client, appID, namespace)
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Deployment namespace")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func runInfo(cmd *cobra.Command, client *api.Client, appID string, namespace string) error {
	cmd.Println(ui.RenderTitleWithBorder("Application Info"))

	// Get deployment info
	deployment, err := client.GetDeploymentInfo(cmd.Context(), namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get application info: %w", err)
	}

	// Display info
	cmd.Printf("Application ID: %s\n", appID)
	cmd.Printf("Namespace:      %s\n", deployment.Namespace)
	cmd.Printf("Template:       %s\n", deployment.TemplateName)
	cmd.Printf("Status:         %s\n", deployment.DeploymentStatus)

	return nil
}
