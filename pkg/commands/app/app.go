package app

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

var Cmd *cobra.Command

func init() {
	Cmd = NewCommand()
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Manage your applications",
		Long:  `Create, list, and manage your Nexlayer applications.`,
	}

	// Add subcommands
	cmd.AddCommand(CreateCmd())
	cmd.AddCommand(ListCmd())

	return cmd
}

// CreateCmd creates a new application
func CreateCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new application",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return fmt.Errorf("name is required")
			}

			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			req := &types.CreateAppRequest{
				Name: name,
			}

			app, err := client.CreateApplication(context.Background(), req)
			if err != nil {
				return fmt.Errorf("failed to create application: %w", err)
			}

			fmt.Printf("Created application %s\n", app.Name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Name of the application")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(fmt.Sprintf("failed to mark name flag as required: %v", err))
	}

	return cmd
}

// ListCmd lists all applications
func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			apps, err := client.ListApplications(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list applications: %w", err)
			}

			if len(apps) == 0 {
				fmt.Println("No applications found")
				return nil
			}

			fmt.Println("Applications:")
			for _, app := range apps {
				fmt.Printf("- %s (created: %s)\n", app.Name, app.CreatedAt.Format("2006-01-02 15:04:05"))
			}

			return nil
		},
	}

	return cmd
}
