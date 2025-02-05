// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package debug

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
		Use:   "debug",
		Short: "Debug deployment issues",
		Long: `Debug issues with your deployments by checking their status and configuration.
		
Example:
  nexlayer debug --app myapp --namespace ecstatic-frog`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDebug(cmd, client, appID, namespace)
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Deployment namespace")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func runDebug(cmd *cobra.Command, client *api.Client, appID string, namespace string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployment Debug"))

	// Get deployment info
	deployment, err := client.GetDeploymentInfo(cmd.Context(), namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	// Display debug info
	cmd.Printf("Deployment Status:\n")
	cmd.Printf("  Namespace:      %s\n", deployment.Namespace)
	cmd.Printf("  Template:       %s\n", deployment.TemplateName)
	cmd.Printf("  Template ID:    %s\n", deployment.TemplateID)
	cmd.Printf("  Status:         %s\n", deployment.DeploymentStatus)

	// Provide debug suggestions
	cmd.Printf("\nDebug Suggestions:\n")
	switch deployment.DeploymentStatus {
	case "pending":
		cmd.Println("- Deployment is still being created. Please wait a few minutes.")
	case "failed":
		cmd.Println("- Check your deployment configuration in the YAML file")
		cmd.Println("- Ensure all required environment variables are set")
		cmd.Println("- Verify your application code builds successfully")
	case "running":
		cmd.Println("- Deployment is running normally")
		cmd.Println("- Use 'nexlayer status' for more detailed information")
	default:
		cmd.Printf("- Unknown status: %s\n", deployment.DeploymentStatus)
		cmd.Println("- Please contact support for assistance")
	}

	return nil
}
