package commands

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

// InfoCmd represents the info command
var InfoCmd = &cobra.Command{
	Use:   "info [namespace]",
	Short: "Get deployment info",
	Long: `Get detailed information about a deployment.
Example: nexlayer info my-app`,
	Args: cobra.ExactArgs(1),
	RunE: runInfo,
}

func runInfo(cmd *cobra.Command, args []string) error {
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
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	fmt.Printf("Deployment Information:\n")
	fmt.Printf("  Namespace: %s\n", resp.Deployment.Namespace)
	fmt.Printf("  Template: %s\n", resp.Deployment.TemplateName)
	fmt.Printf("  Status: %s\n", resp.Deployment.DeploymentStatus)
	fmt.Printf("  Template ID: %s\n", resp.Deployment.TemplateID)

	return nil
}
