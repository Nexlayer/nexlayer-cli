// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deployment

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// newInfoCommand creates a command that wraps GET /getDeploymentInfo
func newInfoCommand(apiClient *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info [deployment-id]",
		Short:   "Show deployment details",
		Long:    "Show detailed information about a specific deployment (GET /getDeploymentInfo)",
		Aliases: []string{"status"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deploymentID := args[0]

			// Show progress
			spinner := ui.NewSpinner("Fetching deployment info...")
			spinner.Start()
			defer spinner.Stop()

			// Get application ID from config
			appID := "default" // TODO: Get from config

			// Call GET /getDeploymentInfo
			resp, err := apiClient.GetDeploymentInfo(cmd.Context(), deploymentID, appID)
			if err != nil {
				return fmt.Errorf("failed to get deployment info: %w", err)
			}

			// Check JSON output flag
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return json.NewEncoder(os.Stdout).Encode(resp)
			}

			// Print human-readable output
			deployment := resp.Data
			fmt.Printf("✨ Status:  %s\n", deployment.Status)
			fmt.Printf("🌐 URL:     %s\n", deployment.URL)
			fmt.Printf("📚 Version: %s\n", deployment.Version)

			if len(deployment.PodStatuses) > 0 {
				fmt.Println("\n📦 Pods:")
				table := ui.NewTable()
				table.AddHeader("NAME", "STATUS", "READY")
				for _, pod := range deployment.PodStatuses {
					table.AddRow(pod.Name, pod.Status, fmt.Sprintf("%v", pod.Ready))
				}
				table.Render()
			}

			return nil
		},
	}

	return cmd
}
