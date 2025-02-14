// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package list

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// NewListCommand creates a new list command
func NewListCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list [applicationID]",
		Short:   "List deployments",
		Long:    "Retrieve a list of all deployments for a specific application",
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show progress
			spinner := ui.NewSpinner("Fetching deployments...")
			spinner.Start()
			defer spinner.Stop()

			// Get deployments
			resp, err := client.ListDeployments(cmd.Context())
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

			// Print table
			fmt.Printf("%-20s %-30s %-15s\n", "STATUS", "URL", "VERSION")
			fmt.Println(strings.Repeat("-", 65))
			for _, d := range resp.Data {
				fmt.Printf("%-20s %-30s %-15s\n", d.Status, d.URL, d.Version)
			}
			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}
