// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package info

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// NewInfoCommand creates a new info command
func NewInfoCommand(client api.ClientAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info <namespace> <applicationID>",
		Short:   "Get deployment info",
		Long:    "Get detailed information about a deployment for a specific application",
		Aliases: []string{"status"},
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]
			applicationID := args[1]

			// Show progress
			spinner := ui.NewSpinner("Fetching deployment info...")
			spinner.Start()
			defer spinner.Stop()

			// Call GET /getDeploymentInfo/{namespace}/{applicationID}
			resp, err := client.GetDeploymentInfo(cmd.Context(), namespace, applicationID)
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
			
			// Print additional details if available
			if !deployment.LastUpdated.IsZero() {
				fmt.Printf("🕒 Updated: %s\n", deployment.LastUpdated.Format("2006-01-02 15:04:05"))
			}

			return nil
		},
	}

	return cmd
}
