package info

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

// NewInfoCmd creates a new info command
func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get information about applications and deployments",
		Long:  `Get detailed information about applications and their deployments.`,
	}

	cmd.AddCommand(newAppInfoCmd())
	cmd.AddCommand(newDeploymentInfoCmd())

	return cmd
}

// newAppInfoCmd creates a command to get application info
func newAppInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "app [application-id]",
		Short: "Get information about an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			applicationID := args[0]
			app, err := client.GetAppInfo(context.Background(), applicationID)
			if err != nil {
				return fmt.Errorf("failed to get application info: %w", err)
			}

			fmt.Printf("Application: %s\n", app.Name)
			fmt.Printf("ID: %s\n", app.ID)
			fmt.Printf("Status: %s\n", app.Status)
			fmt.Printf("Created: %s\n", app.CreatedAt.Format("2006-01-02 15:04:05"))

			return nil
		},
	}
}

// newDeploymentInfoCmd creates a command to get deployment info
func newDeploymentInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deployment [namespace] [application-id]",
		Short: "Get detailed information about a deployment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			namespace := args[0]
			applicationID := args[1]

			deployment, err := client.GetDeploymentInfo(context.Background(), namespace, applicationID)
			if err != nil {
				return fmt.Errorf("failed to get deployment info: %w", err)
			}

			fmt.Printf("Deployment Information:\n")
			fmt.Printf("ID: %s\n", deployment.ID)
			fmt.Printf("Application ID: %s\n", deployment.ApplicationID)
			fmt.Printf("Status: %s\n", deployment.Status)
			fmt.Printf("Namespace: %s\n", deployment.Namespace)
			fmt.Printf("Created: %s\n", deployment.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Printf("Updated: %s\n", deployment.UpdatedAt.Format("2006-01-02 15:04:05"))
			if deployment.Config != "" {
				fmt.Printf("\nConfiguration:\n%s\n", deployment.Config)
			}

			return nil
		},
	}
}
