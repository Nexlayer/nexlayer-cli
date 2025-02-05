// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package logs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// Package logs provides functionality to view application logs.

// NewCommand creates a new logs command.
// It sets up flags for application ID, follow mode, tail count, and JSON output.
// Returns a configured cobra.Command instance.
func NewCommand(client *api.Client) *cobra.Command {
	var appID string
	var follow bool
	var tail int
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "logs [namespace]",
		Short: "View application logs",
		Long: `View logs from your application deployments.

Examples:
  nexlayer logs my-namespace --app my-app --follow
  nexlayer logs my-namespace --app my-app --tail 200 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogs(cmd, client, args[0], appID, follow, tail, jsonOutput)
		},
	}

	// Define flags.
	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID (required)")
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().IntVarP(&tail, "tail", "t", 100, "Number of lines to show from the end of the logs")
	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output logs in JSON format")
	cmd.MarkFlagRequired("app")

	return cmd
}

// runLogs retrieves and prints logs.
// It uses the API client to fetch logs based on namespace, appID, and other flags.
// Logs are printed in either plain text or JSON format based on the jsonOutput flag.
// Returns an error if log retrieval fails.
func runLogs(cmd *cobra.Command, client *api.Client, namespace, appID string, follow bool, tail int, jsonOutput bool) error {
	// For human-friendly output, render a title.
	if !jsonOutput {
		cmd.Println(ui.RenderTitleWithBorder("Application Logs"))
	}

	// Retrieve logs using the API client.
	logs, err := client.GetLogs(cmd.Context(), namespace, appID, follow, tail)
	if err != nil {
		// Output structured error in JSON if requested.
		if jsonOutput {
			errorObj := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"error":     "failed to get logs",
				"details":   err.Error(),
				"namespace": namespace,
				"appID":     appID,
			}
			if jsonBytes, jerr := json.Marshal(errorObj); jerr == nil {
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Printf("failed to get logs: %v\n", err)
			}
			return nil
		}
		return fmt.Errorf("failed to get logs: %w", err)
	}

	// Print each log line.
	if jsonOutput {
		for _, line := range logs {
			logObj := map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339),
				"namespace": namespace,
				"appID":     appID,
				"line":      line,
			}
			if jsonBytes, err := json.Marshal(logObj); err == nil {
				fmt.Println(string(jsonBytes))
			} else {
				fmt.Println(line)
			}
		}
	} else {
		for _, line := range logs {
			fmt.Println(line)
		}
	}

	return nil
}
