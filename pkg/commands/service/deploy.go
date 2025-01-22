package service

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a service",
	Long: `Deploy a service to your Nexlayer deployment.
Examples:
  nexlayer service deploy --app my-app --service frontend --env API_URL=https://api.example.com
  nexlayer service deploy --app my-app --service backend --env "DB_URL=postgres://localhost:5432/db" --env "REDIS_URL=redis://localhost:6379"`,
	RunE: runDeploy,
}

func init() {
	deployCmd.Flags().StringVar(&vars.AppName, "app", "", "Application name")
	deployCmd.Flags().StringVar(&vars.ServiceName, "service", "", "Service name")
	deployCmd.Flags().StringSliceVar(&vars.EnvVars, "env", []string{}, "Environment variables (KEY=VALUE)")
	deployCmd.Flags().StringVar(&vars.APIURL, "api-url", "https://api.nexlayer.io", "API URL")

	deployCmd.MarkFlagRequired("app")
	deployCmd.MarkFlagRequired("service")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	// Get auth token
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	client := api.NewClient(vars.APIURL)

	fmt.Printf("🚀 Deploying service %s in app %s...\n", vars.ServiceName, vars.AppName)

	if err := client.Deploy(token, vars.AppName, vars.ServiceName, vars.EnvVars); err != nil {
		return fmt.Errorf("failed to deploy service: %w", err)
	}

	fmt.Printf("✅ Successfully deployed service %s\n", vars.ServiceName)
	return nil
}
