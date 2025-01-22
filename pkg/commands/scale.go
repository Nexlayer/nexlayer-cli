// Package commands contains the CLI commands for the nexlayer CLI tool.
package commands

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	replicas int
)

func init() {
	ScaleCmd.Flags().IntVarP(&replicas, "replicas", "r", 1, "Number of replicas")
}

// ScaleCmd represents the scale command
var ScaleCmd = &cobra.Command{
	Use:   "scale [namespace]",
	Short: "Scale a deployment",
	Long: `Scale a deployment to the specified number of replicas.
Example: nexlayer scale my-app --replicas 3`,
	Args: cobra.ExactArgs(1),
	RunE: runScale,
}

func runScale(cmd *cobra.Command, args []string) error {
	namespace := args[0]

	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create client
	client := api.NewClient("https://app.staging.nexlayer.io")
	err := client.ScaleDeployment(namespace, sessionID, replicas)
	if err != nil {
		return fmt.Errorf("failed to scale deployment: %w", err)
	}

	fmt.Printf("Successfully scaled %s to %d replicas\n", namespace, replicas)
	return nil
}
