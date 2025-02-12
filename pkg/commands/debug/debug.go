// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package debug

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/common"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

func NewCommand(client common.CommandClient) *cobra.Command {
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

func runDebug(cmd *cobra.Command, client common.CommandClient, appID string, namespace string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployment Debug"))

	// Get deployment info
	resp, err := client.GetDeploymentInfo(cmd.Context(), namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	deployment := resp.Data

	// Display debug info
	cmd.Printf("Deployment Status:\n")
	cmd.Printf("  Namespace:      %s\n", deployment.Namespace)
	cmd.Printf("  Template:       %s\n", deployment.TemplateName)
	cmd.Printf("  Template ID:    %s\n", deployment.TemplateID)
	cmd.Printf("  Status:         %s\n", deployment.Status)
	cmd.Printf("  URL:            %s\n", deployment.URL)
	cmd.Printf("  Custom Domain:  %s\n", deployment.CustomDomain)
	cmd.Printf("  Version:        %s\n", deployment.Version)
	cmd.Printf("  Created:        %s\n", deployment.CreatedAt.Format(time.RFC3339))
	cmd.Printf("  Last Updated:   %s\n", deployment.LastUpdated.Format(time.RFC3339))

	// Display pod statuses
	if len(deployment.PodStatuses) > 0 {
		cmd.Printf("\nPod Statuses:\n")
		for _, pod := range deployment.PodStatuses {
			cmd.Printf("\n  Pod: %s (%s)\n", pod.Name, pod.Type)
			cmd.Printf("    Status:    %s\n", pod.Status)
			cmd.Printf("    Ready:     %v\n", pod.Ready)
			cmd.Printf("    Restarts:  %d\n", pod.Restarts)
			cmd.Printf("    Image:     %s\n", pod.Image)
			cmd.Printf("    Created:   %s\n", pod.CreatedAt.Format(time.RFC3339))
		}
	}
	// Provide debug suggestions
	cmd.Printf("\nDebug Suggestions:\n")
	switch deployment.Status {
	case "Pending":
		cmd.Println("- Deployment is still being created. Please wait a few minutes.")
	case "Failed":
		cmd.Println("- Check your deployment configuration in the YAML file")
		cmd.Println("- Ensure all required environment variables are set")
		cmd.Println("- Verify your application code builds successfully")
	case "Running":
		cmd.Println("- Deployment is running normally")
		cmd.Println("- Use 'nexlayer status' for more detailed information")
	default:
		cmd.Printf("- Unknown status: %s\n", deployment.Status)
		cmd.Println("- Please contact support for assistance")
	}

	return nil
}
