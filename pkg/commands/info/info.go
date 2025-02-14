// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package info

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// NewInfoCommand creates a new info command
func NewInfoCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <namespace> <applicationID>",
		Short: "Get deployment info",
		Long:  "Retrieve detailed information about a specific deployment.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			namespace := args[0]
			appID := args[1]

			// Show progress
			spinner := ui.NewSpinner("Fetching deployment info...")
			spinner.Start()
			defer spinner.Stop()

			// Get deployment info
			resp, err := client.GetDeploymentInfo(cmd.Context(), namespace, appID)
			if err != nil {
				return fmt.Errorf("failed to get deployment info: %w", err)
			}

			// Check JSON output flag
			if jsonOutput, _ := cmd.Flags().GetBool("json"); jsonOutput {
				return json.NewEncoder(os.Stdout).Encode(resp)
			}

			// Print human-readable output
			fmt.Printf("Status: %s\n", resp.Data.Status)
			fmt.Printf("URL: %s\n", resp.Data.URL)
			fmt.Printf("Version: %s\n", resp.Data.Version)
			fmt.Printf("Last Updated: %s\n", resp.Data.LastUpdated.Format(time.RFC3339))

			if len(resp.Data.PodStatuses) > 0 {
				fmt.Println("\nPods:")
				for _, pod := range resp.Data.PodStatuses {
					fmt.Printf("  - %s: %s (Ready: %v)\n", pod.Name, pod.Status, pod.Ready)
				}
			}

			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")
	return cmd
}
