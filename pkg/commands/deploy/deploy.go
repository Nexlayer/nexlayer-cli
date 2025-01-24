package deploy

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

const (
	maxRetries = 3
	retryDelay = 2 * time.Second
)

// retryWithBackoff executes the given function with exponential backoff
func retryWithBackoff(ctx context.Context, operation func() error) error {
	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err != nil {
			lastErr = err
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(retryDelay * time.Duration(i+1)):
				continue
			}
		}
		return nil
	}
	return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}

// NewDeployCmd creates a new deploy command
func NewDeployCmd() *cobra.Command {
	var (
		yamlFile string
		applicationID string
		useAI    bool
	)

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application",
		Long: `Deploy an application using a YAML configuration file.
Use the --ai flag to get AI-powered suggestions for optimizing your deployment.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Handle AI suggestions if enabled
			if useAI {
				aiClient, err := ai.NewClient()
				if err != nil {
					return fmt.Errorf("failed to initialize AI client: %w", err)
				}

				if err := aiClient.HandleAIFlag(ctx, "deploy", args); err != nil {
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

			// Start deployment with retry logic
			var deployment *types.Deployment
			err = retryWithBackoff(ctx, func() error {
				req := &types.DeployRequest{
					YAML:          string(yamlContent),
					ApplicationID: applicationID,
				}

				var err error
				deployment, err = client.StartUserDeployment(ctx, applicationID, req)
				if err != nil {
					return fmt.Errorf("deployment attempt failed: %w", err)
				}
				return nil
			})

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
	cmd.Flags().StringVar(&applicationID, "app-id", "", "Application ID")
	cmd.Flags().BoolVar(&useAI, "ai", false, "Enable AI-powered suggestions")

	// Mark required flags
	if err := cmd.MarkFlagRequired("file"); err != nil {
		panic(fmt.Sprintf("failed to mark file flag as required: %v", err))
	}
	if err := cmd.MarkFlagRequired("app-id"); err != nil {
		panic(fmt.Sprintf("failed to mark app-id flag as required: %v", err))
	}

	return cmd
}
