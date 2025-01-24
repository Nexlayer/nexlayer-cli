package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

// NewListCmd creates a new list command
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
	}

	cmd.AddCommand(listDeploymentsCmd())
	return cmd
}

// listDeploymentsCmd creates a command to list deployments
func listDeploymentsCmd() *cobra.Command {
	var appName string

	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "List deployments for an application",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			deployments, err := client.GetDeployments(context.Background(), appName)
			if err != nil {
				return fmt.Errorf("failed to list deployments: %w", err)
			}

			if len(deployments) == 0 {
				fmt.Println("No deployments found")
				return nil
			}

			fmt.Printf("Deployments for application %s:\n", appName)
			for _, d := range deployments {
				fmt.Printf("- ID: %s\n", d.ID)
				fmt.Printf("  Status: %s\n", d.Status)
				fmt.Printf("  Created: %s\n", d.CreatedAt.Format("2006-01-02 15:04:05"))
				fmt.Printf("  Updated: %s\n", d.UpdatedAt.Format("2006-01-02 15:04:05"))
				fmt.Println()
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	if err := cmd.MarkFlagRequired("app"); err != nil {
		panic(fmt.Sprintf("failed to mark app flag as required: %v", err))
	}

	return cmd
}
