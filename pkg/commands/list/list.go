// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package list

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	// Status styles
	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)

	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00")).
			Bold(true)

	failedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000")).
			Bold(true)
)

// NewListCommand creates a new list command
func NewListCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [applicationID]",
		Short: "List your Nexlayer deployments",
		Long: `List all your deployments on the Nexlayer platform.

The command shows:
  • Deployment status and health
  • Application URLs
  • Version information
  • Last update time

Examples:
  nexlayer list                    # List all deployments
  nexlayer list my-app            # List deployments for specific application`,
		Aliases: []string{"ls"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Show progress
			fmt.Fprintf(cmd.OutOrStdout(), "📋 Fetching your deployments...\n\n")

			var resp *schema.APIResponse[[]schema.Deployment]
			var err error

			// Check if an application ID was provided
			if len(args) > 0 {
				appID := args[0]
				// If we have an appID, use GetDeployments which filters by appID
				// This requires a cast to APIClientForCommands since GetDeployments is not in the main APIClient interface
				if clientWithCommands, ok := client.(api.APIClientForCommands); ok {
					resp, err = clientWithCommands.GetDeployments(cmd.Context(), appID)
					if err != nil {
						return fmt.Errorf("failed to get deployments for application %s: %w", appID, err)
					}
				} else {
					// Fallback to ListDeployments if the client doesn't implement APIClientForCommands
					fmt.Fprintf(cmd.OutOrStdout(), "Warning: Filtering by application ID not supported. Showing all deployments.\n\n")
					resp, err = client.ListDeployments(cmd.Context())
					if err != nil {
						return fmt.Errorf("failed to get deployments: %w", err)
					}
				}
			} else {
				// Get all deployments
				resp, err = client.ListDeployments(cmd.Context())
				if err != nil {
					return fmt.Errorf("failed to get deployments: %w", err)
				}
			}

			// Check JSON output flag
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return json.NewEncoder(os.Stdout).Encode(resp)
			}

			// Print human-readable table
			if len(resp.Data) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No deployments found. Use 'nexlayer deploy' to deploy your first application.")
				return nil
			}

			// Print table
			table := ui.NewTable()
			table.AddHeader("STATUS", "URL", "VERSION", "LAST UPDATED")
			for _, d := range resp.Data {
				url := d.URL
				if d.CustomDomain != "" {
					url = fmt.Sprintf("%s (custom domain: %s)", d.URL, d.CustomDomain)
				}
				table.AddRow(
					formatStatus(d.Status),
					url,
					d.Version,
					formatTime(d.LastUpdated),
				)
			}
			table.Render()

			// Print help text
			fmt.Fprintf(cmd.OutOrStdout(), "\nℹ️  Available Commands:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "• View details:    nexlayer info <namespace> <appID>\n")
			fmt.Fprintf(cmd.OutOrStdout(), "• View logs:       nexlayer logs <namespace> <appID>\n")
			fmt.Fprintf(cmd.OutOrStdout(), "• Set domain:      nexlayer domain set <appID> --domain example.com\n")
			fmt.Fprintf(cmd.OutOrStdout(), "• Update config:   nexlayer deploy\n")

			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}

// formatStatus returns a colored status string
func formatStatus(status string) string {
	switch status {
	case "running":
		return runningStyle.Render(status)
	case "pending":
		return pendingStyle.Render(status)
	case "failed":
		return failedStyle.Render(status)
	default:
		return status
	}
}

// formatTime formats a time.Time into a human-readable string
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}
