package status

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

func NewCommand(client *api.Client) *cobra.Command {
	var appID string
	var namespace string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Get deployment status",
		Long: `Get detailed status information about a deployment.
		
Example:
  nexlayer status --app myapp --namespace ecstatic-frog`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(cmd, client, namespace, appID)
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Deployment namespace")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func runStatus(cmd *cobra.Command, client *api.Client, namespace string, appID string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployment Status"))

	deployment, err := client.GetDeploymentInfo(cmd.Context(), namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get deployment status: %w", err)
	}

	// Display deployment info
	cmd.Printf("Namespace:     %s\n", deployment.Namespace)
	cmd.Printf("Template:      %s\n", deployment.TemplateName)
	cmd.Printf("Template ID:   %s\n", deployment.TemplateID)
	cmd.Printf("Status:        %s\n", deployment.DeploymentStatus)

	return nil
}
