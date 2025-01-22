package deploy

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	applicationID string
	configFile    string
)

// Command represents the deploy command
var Command = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application",
	Long: `Deploy an application using a YAML configuration file.
Example:
  nexlayer-cli deploy --app my-app --config app.yaml`,
	RunE: runDeploy,
}

func init() {
	Command.Flags().StringVar(&applicationID, "app", "", "Application ID")
	Command.Flags().StringVar(&configFile, "config", "", "Path to YAML configuration file")
	Command.MarkFlagRequired("app")
	Command.MarkFlagRequired("config")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	// Create API client
	client, err := api.NewClient("https://app.nexlayer.io")
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Start deployment
	fmt.Printf("Starting deployment for application %s...\n", applicationID)
	resp, err := client.StartUserDeployment(applicationID, configFile)
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	fmt.Printf("\nDeployment started successfully!\n")
	fmt.Printf("Namespace: %s\n", resp.Namespace)
	fmt.Printf("URL: %s\n", resp.URL)
	fmt.Printf("Message: %s\n", resp.Message)

	return nil
}
