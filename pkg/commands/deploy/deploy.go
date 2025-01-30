package deploy

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

// NewCommand creates a new deploy command
func NewCommand(client *api.Client) *cobra.Command {
	var yamlFile string
	var appID string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application",
		Long: `Deploy an application using a YAML configuration file.
		
Example:
  nexlayer deploy --app myapp --file deploy.yaml`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(cmd, client, appID, yamlFile)
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID to deploy")
	cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "Path to YAML configuration file")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("file")

	return cmd
}

func runDeploy(cmd *cobra.Command, client *api.Client, appID string, yamlFile string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deploying Application"))

	// Start deployment
	resp, err := client.StartDeployment(cmd.Context(), appID, yamlFile)
	if err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}

	cmd.Printf("Deployment started successfully!\n")
	cmd.Printf("Namespace: %s\n", resp.Namespace)
	cmd.Printf("URL: %s\n", resp.URL)

	return nil
}
