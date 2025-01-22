package list

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	applicationID string
)

// Command represents the list command
var Command = &cobra.Command{
	Use:   "list",
	Short: "List deployments",
	Long: `List all deployments for an application.
Example:
  nexlayer-cli list --app my-app`,
	RunE: runList,
}

func init() {
	Command.Flags().StringVar(&applicationID, "app", "", "Application ID")
	Command.MarkFlagRequired("app")
}

func runList(cmd *cobra.Command, args []string) error {
	// Create API client
	client, err := api.NewClient("https://app.nexlayer.io")
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Get deployments
	fmt.Printf("Fetching deployments for application %s...\n", applicationID)
	resp, err := client.GetDeployments(applicationID)
	if err != nil {
		return fmt.Errorf("failed to get deployments: %w", err)
	}

	if len(resp.Deployments) == 0 {
		fmt.Println("\nNo deployments found.")
		return nil
	}

	fmt.Printf("\nDeployments:\n")
	for _, d := range resp.Deployments {
		fmt.Printf("\nNamespace: %s\n", d.Namespace)
		fmt.Printf("Application ID: %s\n", d.ApplicationID)
		fmt.Printf("Status: %s\n", d.DeploymentStatus)
		fmt.Println("---")
	}

	return nil
}
