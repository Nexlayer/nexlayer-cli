package service

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure service settings",
	Long: `Configure settings for a service in your Nexlayer deployment.
Examples:
  nexlayer service configure --app my-app --service frontend --env API_URL=https://api.example.com
  nexlayer service configure --app my-app --service backend --env "DB_URL=postgres://localhost:5432/db" --env "REDIS_URL=redis://localhost:6379"`,
	RunE: runConfigure,
}

func init() {
	configureCmd.Flags().StringVar(&vars.AppName, "app", "", "Application name")
	configureCmd.Flags().StringVar(&vars.ServiceName, "service", "", "Service name")
	configureCmd.Flags().StringSliceVar(&vars.EnvVars, "env", []string{}, "Environment variables (KEY=VALUE)")
	configureCmd.Flags().StringVar(&vars.APIURL, "api-url", "https://api.nexlayer.io", "API URL")

	configureCmd.MarkFlagRequired("app")
	configureCmd.MarkFlagRequired("service")
}

func runConfigure(cmd *cobra.Command, args []string) error {
	client := api.NewClient(vars.APIURL)

	fmt.Printf("⚙️  Configuring service %s in app %s...\n", vars.ServiceName, vars.AppName)

	if err := client.Configure(vars.AppName, vars.ServiceName, vars.EnvVars); err != nil {
		return fmt.Errorf("failed to configure service: %w", err)
	}

	fmt.Printf("✅ Successfully configured service %s\n", vars.ServiceName)
	return nil
}
