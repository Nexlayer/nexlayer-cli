// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package info

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
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

	// Section styles
	sectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffff")).
			Bold(true)
)

// NewInfoCommand creates a new info command
func NewInfoCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <namespace> <applicationID>",
		Short: "Get detailed deployment information",
		Long: `Retrieve detailed information about a specific deployment.

The command shows:
  • Deployment status and health
  • Pod statuses and readiness
  • Resource usage and limits
  • Environment variables
  • Volume mounts
  • Network configuration

Arguments:
  namespace      The deployment namespace (required)
  applicationID  The application ID (required)

Examples:
  nexlayer info my-namespace my-app
  nexlayer info production api-backend`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]
			appID := args[1]

			// Validate namespace
			if namespace == "" {
				return fmt.Errorf("namespace is required")
			}
			namespace = strings.TrimSpace(namespace)
			if namespace == "" {
				return fmt.Errorf("namespace cannot be only whitespace")
			}
			if strings.Contains(namespace, "/") {
				return fmt.Errorf("namespace cannot contain slashes")
			}

			// Validate applicationID
			if appID == "" {
				return fmt.Errorf("applicationID is required")
			}
			appID = strings.TrimSpace(appID)
			if appID == "" {
				return fmt.Errorf("applicationID cannot be only whitespace")
			}
			if strings.Contains(appID, "/") {
				return fmt.Errorf("applicationID cannot contain slashes")
			}

			// Show progress
			fmt.Fprintf(cmd.OutOrStdout(), "📊 Fetching deployment information...\n\n")

			// Get deployment info using namespace
			resp, err := client.GetDeploymentInfo(cmd.Context(), namespace)
			if err != nil {
				return fmt.Errorf("failed to get deployment info: %w", err)
			}

			// Check JSON output flag
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return json.NewEncoder(os.Stdout).Encode(resp)
			}

			// Print deployment overview
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", sectionStyle.Render("Deployment Overview"))
			fmt.Fprintf(cmd.OutOrStdout(), "Status:       %s\n", formatStatus(resp.Data.Status))
			fmt.Fprintf(cmd.OutOrStdout(), "URL:          %s\n", resp.Data.URL)
			if resp.Data.CustomDomain != "" {
				fmt.Fprintf(cmd.OutOrStdout(), "Domain:       %s\n", resp.Data.CustomDomain)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Version:      %s\n", resp.Data.Version)
			fmt.Fprintf(cmd.OutOrStdout(), "Last Updated: %s\n", formatTime(resp.Data.LastUpdated))

			// Print pod statuses
			if len(resp.Data.PodStatuses) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", sectionStyle.Render("Pod Status"))
				table := ui.NewTable()
				table.AddHeader("NAME", "STATUS", "READY", "RESTARTS", "AGE")
				for _, pod := range resp.Data.PodStatuses {
					table.AddRow(
						pod.Name,
						formatStatus(pod.Status),
						formatBool(pod.Ready),
						fmt.Sprintf("%d", pod.Restarts),
						formatAge(pod.CreatedAt),
					)
				}
				table.Render()
			}

			// Print next steps based on status
			fmt.Fprintf(cmd.OutOrStdout(), "\n%s\n", sectionStyle.Render("Available Commands"))
			switch resp.Data.Status {
			case "running":
				fmt.Fprintf(cmd.OutOrStdout(), "• View logs:        nexlayer logs %s %s\n", namespace, appID)
				if resp.Data.CustomDomain == "" {
					fmt.Fprintf(cmd.OutOrStdout(), "• Set domain:       nexlayer domain set %s --domain example.com\n", appID)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "• Update config:    nexlayer deploy\n")
			case "failed":
				fmt.Fprintf(cmd.OutOrStdout(), "• View logs:        nexlayer logs %s %s\n", namespace, appID)
				fmt.Fprintf(cmd.OutOrStdout(), "• Check events:     nexlayer events %s %s\n", namespace, appID)
				fmt.Fprintf(cmd.OutOrStdout(), "• Redeploy:        nexlayer deploy\n")
			case "pending":
				fmt.Fprintf(cmd.OutOrStdout(), "• View logs:        nexlayer logs %s %s\n", namespace, appID)
				fmt.Fprintf(cmd.OutOrStdout(), "• Check events:     nexlayer events %s %s\n", namespace, appID)
				fmt.Fprintf(cmd.OutOrStdout(), "• Cancel deploy:    nexlayer cancel %s %s\n", namespace, appID)
			}

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

// formatBool returns a colored boolean string
func formatBool(b bool) string {
	if b {
		return runningStyle.Render("Yes")
	}
	return failedStyle.Render("No")
}

// formatTime formats a time.Time into a human-readable string
func formatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("2006-01-02 15:04:05")
}

// formatAge returns a human-readable duration since the given time
func formatAge(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}

	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}
