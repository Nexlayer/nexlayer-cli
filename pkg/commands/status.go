// Package commands contains the CLI commands for the Nexlayer CLI.
package commands

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status [namespace]",
	Short: "Get deployment status",
	Long: `Get the current status of a deployment.
Example: nexlayer status my-app`,
	Args: cobra.ExactArgs(1),
	RunE: runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	namespace := args[0]

	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create client
	client := api.NewClient("https://app.staging.nexlayer.io")
	resp, err := client.GetDeploymentInfo(namespace, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get deployment status: %w", err)
	}

	fmt.Printf("Deployment Status:\n")
	fmt.Printf("  Namespace: %s\n", resp.Deployment.Namespace)
	fmt.Printf("  Template: %s\n", resp.Deployment.TemplateName)
	fmt.Printf("  Status: %s\n", resp.Deployment.DeploymentStatus)
	fmt.Printf("  Template ID: %s\n", resp.Deployment.TemplateID)

	return nil
}
