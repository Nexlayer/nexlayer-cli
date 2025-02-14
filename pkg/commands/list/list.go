// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package list

import (
	"encoding/json"
	"fmt"
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
			fmt.Fprintln(cmd.ErrOrStderr(), "Debug: Creating spinner...")
			spinner := ui.NewSpinner("Fetching deployments...")
			if spinner == nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "Debug: Spinner is nil!")
			}
			spinner.Start()
			defer spinner.Stop()

			// Get deployments
			fmt.Fprintln(cmd.ErrOrStderr(), "Debug: Calling ListDeployments...")
			resp, err := client.ListDeployments(cmd.Context())
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Debug: ListDeployments error: %v\n", err)
				return fmt.Errorf("failed to get deployments: %w", err)
			}
			if resp == nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "Debug: Response is nil!")
			} else {
				fmt.Fprintf(cmd.ErrOrStderr(), "Debug: Got %d deployments\n", len(resp.Data))
			}

			// Check JSON output flag
			jsonOutput, err := cmd.Flags().GetBool("json")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Debug: Error getting json flag: %v\n", err)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Debug: JSON output flag: %v\n", jsonOutput)
			if jsonOutput {
				return json.NewEncoder(cmd.OutOrStdout()).Encode(resp)
			}

			// Print human-readable table
			if resp == nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "Debug: Cannot print table - response is nil")
				return fmt.Errorf("unexpected nil response from API")
			}
			if len(resp.Data) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No deployments found")
				return nil
			}

			// Print table
			fmt.Fprintln(cmd.ErrOrStderr(), "Debug: Printing table...")
			fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-30s %-15s\n", "STATUS", "URL", "VERSION")
			fmt.Fprintln(cmd.OutOrStdout(), strings.Repeat("-", 65))
			for i, d := range resp.Data {
				if d.Status == "" {
					fmt.Fprintf(cmd.ErrOrStderr(), "Debug: Deployment %d has empty status\n", i)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s %-30s %-15s\n", d.Status, d.URL, d.Version)
			}
			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}
