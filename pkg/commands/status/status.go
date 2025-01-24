// Package status contains the CLI commands for the Nexlayer CLI.
package status

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applicationID string

// Command represents the status command
var Command = &cobra.Command{
	Use:   "status",
	Short: "Check deployment status",
	Long: `Check the status of a deployment.
This command is currently not implemented.`,
	RunE: runStatus,
}

func init() {
	Command.Flags().StringVar(&applicationID, "app", "", "Application ID")
	if err := Command.MarkFlagRequired("app"); err != nil {
		panic(fmt.Sprintf("failed to mark app flag as required: %v", err))
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("status check is not yet implemented")
}
