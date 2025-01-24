package deploy

import (
	"context"
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// NewDeployCmd creates a new deploy command
func NewDeployCmd() *cobra.Command {
	var (
		yamlFile string
		appID    string
		useAI    bool
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application",
		Long: `Deploy an application using a YAML configuration file.
Use the --ai flag to get AI-powered suggestions for optimizing your deployment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Handle AI suggestions if enabled
			if useAI {
				aiClient, err := ai.NewClient()
				if err != nil {
					return fmt.Errorf("failed to initialize AI client: %w", err)
				}

				if err := aiClient.HandleAIFlag(cmd.Context(), "deploy", args); err != nil {
					return fmt.Errorf("failed to handle AI suggestions: %w", err)
				}
			}

			// Read YAML file
			yamlContent, err := os.ReadFile(yamlFile)
			if err != nil {
				return fmt.Errorf("failed to read YAML file: %w", err)
			}

			// Create API client
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			// Start deployment
			req := &types.DeployRequest{
				YAML:          string(yamlContent),
				ApplicationID: appID,
			}

			deployment, err := client.StartUserDeployment(context.Background(), appID, req)
			if err != nil {
				return fmt.Errorf("failed to start deployment: %w", err)
			}

			fmt.Printf("Started deployment %s for application %s\n", deployment.ID, deployment.ApplicationID)
			fmt.Printf("Status: %s\n", deployment.Status)
			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "YAML file containing deployment configuration")
	cmd.Flags().StringVar(&appID, "app", "", "Application ID")
	cmd.Flags().BoolVar(&useAI, "ai", false, "Enable AI-powered suggestions")

	// Mark required flags
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagRequired("app")

	return cmd
}
