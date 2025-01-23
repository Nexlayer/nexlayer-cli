package app

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

var (
	appName string
	Cmd     *cobra.Command
)

func init() {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new application",
		Long:  `Create a new application in Nexlayer.`,
		RunE:  runCreate,
	}
	createCmd.Flags().StringVarP(&appName, "name", "n", "", "Name of the application (required)")
	createCmd.MarkFlagRequired("name")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all applications",
		Long:  `List all applications in your Nexlayer account.`,
		RunE:  runList,
	}

	Cmd = &cobra.Command{
		Use:   "app",
		Short: "Manage applications",
		Long:  `Create, list, and manage your Nexlayer applications.`,
	}

	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(listCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(vars.APIURL)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	fmt.Printf("Creating application '%s'...\n", appName)
	resp, err := client.CreateApplication(appName)
	if err != nil {
		return fmt.Errorf("failed to create application: %w", err)
	}

	fmt.Printf("\nApplication created successfully!\n")
	fmt.Printf("Application ID: %s\n", resp.ID)
	fmt.Printf("Name: %s\n", resp.Name)
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(vars.APIURL)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	fmt.Println("Fetching applications...")
	apps, err := client.ListApplications()
	if err != nil {
		return fmt.Errorf("failed to list applications: %w", err)
	}

	if len(apps) == 0 {
		fmt.Println("\nNo applications found.")
		return nil
	}

	fmt.Println("\nApplications:")
	for _, app := range apps {
		fmt.Printf("\nID: %s\n", app.ID)
		fmt.Printf("Name: %s\n", app.Name)
		fmt.Printf("Created: %s\n", app.CreatedAt)
		fmt.Println("---")
	}

	return nil
}
