// Formatted with gofmt -s
package commands

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all deployments",
	Long: `List all your deployments and their current status.
Example: nexlayer list`,
	Args: cobra.NoArgs,
	RunE: runList,
}

func runList(cmd *cobra.Command, args []string) error {
	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create client with staging URL
	client := api.NewClient("https://app.staging.nexlayer.io")
	resp, err := client.GetDeployments(sessionID)
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(resp.Deployments) == 0 {
		fmt.Println("No deployments found")
		return nil
	}

	fmt.Println("Your deployments:")
	fmt.Println("----------------")
	for _, d := range resp.Deployments {
		fmt.Printf("Namespace: %s\n", d.Namespace)
		fmt.Printf("Template: %s (%s)\n", d.TemplateName, d.TemplateID)
		fmt.Printf("Status: %s\n", d.DeploymentStatus)
		fmt.Println("----------------")
	}

	return nil
}
