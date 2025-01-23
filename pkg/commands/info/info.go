package info

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

func NewInfoCmd() *cobra.Command {
	var namespace string
	var applicationID string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get detailed information about a deployment",
		Long:  `Get detailed information about a deployment, including its status and logs.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create API client
			client, err := api.NewClient("https://api.nexlayer.dev")
			if err != nil {
				return fmt.Errorf("failed to create API client: %w", err)
			}

			// Get deployment info
			resp, err := client.GetDeploymentInfo(namespace, applicationID)
			if err != nil {
				return fmt.Errorf("failed to get deployment info: %w", err)
			}

			// Print deployment info
			fmt.Printf("Deployment Information:\n")
			fmt.Printf("  Namespace: %s\n", resp.Namespace)
			fmt.Printf("  Template: %s (%s)\n", resp.TemplateName, resp.TemplateID)
			fmt.Printf("  Status: %s\n", resp.DeploymentStatus)

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace of the deployment")
	cmd.Flags().StringVarP(&applicationID, "app", "a", "", "Application ID")

	// Mark flags as required
	cmd.MarkFlagRequired("namespace")
	cmd.MarkFlagRequired("app")

	return cmd
}
