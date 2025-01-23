// Package scale contains the CLI commands for the nexlayer CLI tool.
package scale

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	applicationID string
	replicas      int
)

// Command represents the scale command
var Command = &cobra.Command{
	Use:   "scale",
	Short: "Scale a deployment",
	Long: `Scale a deployment by specifying the number of replicas.
This command is currently not implemented.`,
	Args: cobra.NoArgs,
	RunE: runScale,
}

func init() {
	Command.Flags().StringVar(&applicationID, "app", "", "Application ID")
	Command.Flags().IntVar(&replicas, "replicas", 1, "Number of replicas")
	Command.MarkFlagRequired("app")
}

func runScale(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("scaling is not yet implemented")
}
