package deploy

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a service",
	Long: `Deploy a service to your Nexlayer deployment.
Examples:
  nexlayer deploy --app my-app --service frontend
  nexlayer deploy --app my-app --service backend --env "DB_URL=postgres://localhost:5432/db"`,
	RunE: runDeploy,
}

func init() {
	DeployCmd.Flags().StringVar(&vars.AppName, "app", "", "Application name")
	DeployCmd.Flags().StringVar(&vars.ServiceName, "service", "", "Service name")
	DeployCmd.Flags().StringSliceVar(&vars.EnvVars, "env", []string{}, "Environment variables (KEY=VALUE)")

	DeployCmd.MarkFlagRequired("app")
	DeployCmd.MarkFlagRequired("service")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	// Get auth token from environment
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create API client
	client := api.NewClient("https://api.nexlayer.io")

	fmt.Printf("🚀 Deploying service %s in app %s...\n", vars.ServiceName, vars.AppName)

	// Deploy service
	if err := client.Deploy(token, vars.AppName, vars.ServiceName, vars.EnvVars); err != nil {
		return fmt.Errorf("failed to deploy service: %w", err)
	}

	fmt.Printf("✅ Successfully deployed service %s\n", vars.ServiceName)
	return nil
}
