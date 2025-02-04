package logs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// NewCommand creates a new logs command.
func NewCommand(client *api.Client) *cobra.Command {
	var appID string
	var follow bool
	var tail int
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "logs [namespace]",
		Short: "View application logs",
		Long: `View logs from your application deployments.

Example:
  nexlayer logs my-namespace --app my-app --follow`,
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

// runLogs fetches and outputs logs. If jsonOutput is true, then errors and log lines are printed
// as JSON objects (which include additional metadata) so that both humans and AI agents can understand them.
func runLogs(cmd *cobra.Command, client *api.Client, namespace, appID string, follow bool, tail int, jsonOutput bool) error {
	// Render a title for human-friendly output.
	if !jsonOutput {
		cmd.Println(ui.RenderTitleWithBorder("Application Logs"))
	}

	// Retrieve logs from the API client.
	logs, err := client.GetLogs(cmd.Context(), namespace, appID, follow, tail)
	if err != nil {
		// If JSON output is requested, print a structured error log.
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
				// Fallback to plain text if JSON marshaling fails.
				fmt.Printf("failed to get logs: %v\n", err)
			}
			// Return nil so that the error is already output in JSON.
			return nil
		}
		// For human-friendly output, return the wrapped error.
		return fmt.Errorf("failed to get logs: %w", err)
	}

	// Print each log line.
	if jsonOutput {
		// For machine readability, output each log line as a JSON object.
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
				// Fallback: output plain text if JSON marshaling fails.
				fmt.Println(line)
			}
		}
	} else {
		// Human-friendly output: simply print each log line.
		for _, line := range logs {
			fmt.Println(line)
		}
	}

	return nil
}
