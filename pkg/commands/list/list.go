package list

import (
	"fmt"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:   "list",
		Short: "List deployments",
		Long:  `List all deployments for an application.`,
		RunE:  runList,
	}
)

func init() {
	Command.Flags().StringVar(&vars.AppID, "app", "", "Application ID (optional)")
}

func runList(cmd *cobra.Command, args []string) error {
	client, err := api.NewClient(vars.APIURL)
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	if vars.AppID == "" {
		// List all applications first
		apps, err := client.ListApplications()
		if err != nil {
			return fmt.Errorf("failed to list applications: %w", err)
		}

		if len(apps) == 0 {
			fmt.Println("No applications found. Create one using 'nexlayer app create'")
			return nil
		}

		fmt.Println("Deployments across all applications:")
		for _, app := range apps {
			fmt.Printf("\nApplication: %s (%s)\n", app.Name, app.ID)
			fmt.Println("---")
			
			deployments, err := client.GetDeployments(app.ID)
			if err != nil {
				fmt.Printf("Error fetching deployments: %v\n", err)
				continue
			}

			if len(deployments.Deployments) == 0 {
				fmt.Println("No deployments found.")
				continue
			}

			for _, d := range deployments.Deployments {
				fmt.Printf("Namespace: %s\n", d.Namespace)
				fmt.Printf("Status: %s\n", d.DeploymentStatus)
				fmt.Println("---")
			}
		}
		return nil
	}

	// List deployments for specific app
	fmt.Printf("Fetching deployments for application %s...\n", vars.AppID)
	resp, err := client.GetDeployments(vars.AppID)
	if err != nil {
		return fmt.Errorf("failed to get deployments: %w", err)
	}

	if len(resp.Deployments) == 0 {
		fmt.Println("\nNo deployments found.")
		return nil
	}

	fmt.Printf("\nDeployments:\n")
	for _, d := range resp.Deployments {
		fmt.Printf("\nNamespace: %s\n", d.Namespace)
		fmt.Printf("Status: %s\n", d.DeploymentStatus)
		fmt.Println("---")
	}

	return nil
}
