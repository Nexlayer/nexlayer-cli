package info

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	applicationID string
	namespace     string
)

// Command represents the info command
var Command = &cobra.Command{
	Use:   "info",
	Short: "Get deployment information",
	Long: `Get detailed information about a deployment.
Example:
  nexlayer-cli info --app my-app --namespace my-namespace`,
	RunE: runInfo,
}

func init() {
	Command.Flags().StringVar(&applicationID, "app", "", "Application ID")
	Command.Flags().StringVar(&namespace, "namespace", "", "Deployment namespace")
	Command.MarkFlagRequired("app")
	Command.MarkFlagRequired("namespace")
}

func runInfo(cmd *cobra.Command, args []string) error {
	// Create API client
	client, err := api.NewClient("https://app.nexlayer.io")
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Get deployment info
	fmt.Printf("Fetching deployment information for application %s...\n", applicationID)
	resp, err := client.GetDeploymentInfo(namespace, applicationID)
	if err != nil {
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	fmt.Printf("\nDeployment Information:\n")
	fmt.Printf("Namespace: %s\n", resp.Deployment.Namespace)
	fmt.Printf("Application ID: %s\n", resp.Deployment.ApplicationID)
	fmt.Printf("Status: %s\n", resp.Deployment.DeploymentStatus)

	return nil
}
