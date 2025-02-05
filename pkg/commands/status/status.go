// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package status

import (
	"context"
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewCommand creates a new status command
func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [namespace] [application-id]",
		Short: "Get deployment status",
		Long: `Get detailed status information about your deployments.
If no namespace and application ID are provided, lists all deployments.
If namespace and application ID are provided, shows detailed information about that specific deployment.`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 2 {
				// Get detailed info for specific deployment
				namespace := args[0]
				appID := args[1]
				return getDeploymentInfo(cmd.Context(), client, namespace, appID)
			}
			// List all deployments
			return listDeployments(cmd.Context(), client)
		},
	}

	return cmd
}

func getDeploymentInfo(ctx context.Context, client *api.Client, namespace, appID string) error {
	info, err := client.GetDeploymentInfo(ctx, namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	// Create a nice status display
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("\n%s\n\n", bold("Deployment Status"))

	// Add deployment details
	fmt.Printf("%s\n", bold("Details:"))
	fmt.Printf("  Namespace:     %s\n", info.Namespace)
	fmt.Printf("  Template:      %s (%s)\n", info.TemplateName, info.TemplateID)
	fmt.Printf("  Status:        %s\n", formatStatus(info.DeploymentStatus))

	// Add access URL
	fmt.Printf("\n%s\n", bold("Access:"))
	fmt.Printf("  URL: https://%s.%s\n", namespace, "nexlayer.io")

	return nil
}

func listDeployments(ctx context.Context, client *api.Client) error {
	deployments, err := client.GetDeployments(ctx, "")
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(deployments) == 0 {
		fmt.Println("No deployments found")
		return nil
	}

	// Print header
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("\n%s\n\n", bold("Your Deployments"))

	// Print table header
	headers := []string{"NAMESPACE", "TEMPLATE", "STATUS"}
	fmt.Printf("%-20s %-30s %s\n", headers[0], headers[1], headers[2])
	fmt.Printf("%-20s %-30s %s\n", strings.Repeat("-", len(headers[0])), strings.Repeat("-", len(headers[1])), strings.Repeat("-", len(headers[2])))

	// Print deployments
	for _, d := range deployments {
		fmt.Printf("%-20s %-30s %s\n",
			d.Namespace,
			fmt.Sprintf("%s (%s)", d.TemplateName, d.TemplateID),
			formatStatus(d.DeploymentStatus),
		)
	}
	fmt.Println()

	return nil
}

func formatStatus(status string) string {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	gray := color.New(color.FgHiBlack).SprintFunc()

	switch strings.ToLower(status) {
	case "running":
		return green("✓ Running")
	case "pending":
		return yellow("⟳ Pending")
	case "failed":
		return red("✗ Failed")
	default:
		return gray(status)
	}
}
