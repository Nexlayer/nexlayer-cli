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

// newListCommand creates a command that wraps GET /getDeployments
func newListCommand(apiClient *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List all deployments",
		Long:    "List all deployments (GET /getDeployments)",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show progress
			spinner := ui.NewSpinner("Fetching deployments...")
			spinner.Start()
			defer spinner.Stop()

			// Get application ID from config
			appID := "default" // TODO: Get from config

			// Call GET /getDeployments
			resp, err := apiClient.GetDeployments(cmd.Context(), appID)
			if err != nil {
				return fmt.Errorf("failed to get deployments: %w", err)
			}

			// Check JSON output flag
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return json.NewEncoder(os.Stdout).Encode(resp)
			}

			// Print human-readable table
			if len(resp.Data) == 0 {
				fmt.Println("No deployments found")
				return nil
			}

			table := ui.NewTable()
			table.AddHeader("STATUS", "URL", "VERSION")
			for _, d := range resp.Data {
				table.AddRow(d.Status, d.URL, d.Version)
			}
			return table.Render()
		},
	}

	return cmd
}
