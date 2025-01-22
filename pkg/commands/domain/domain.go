// Package domain contains the CLI commands for the Nexlayer CLI.
package domain

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	applicationID string
	domain        string
)

// Command represents the domain command
var Command = &cobra.Command{
	Use:   "domain",
	Short: "Configure custom domain for an application",
	Long: `Configure a custom domain for your application.
Example:
  nexlayer-cli domain --app my-app --domain example.com`,
	RunE: runDomain,
}

func init() {
	Command.Flags().StringVar(&applicationID, "app", "", "Application ID")
	Command.Flags().StringVar(&domain, "domain", "", "Custom domain (e.g., example.com)")
	Command.MarkFlagRequired("app")
	Command.MarkFlagRequired("domain")
}

func runDomain(cmd *cobra.Command, args []string) error {
	// Create API client
	client, err := api.NewClient("https://app.nexlayer.io")
	if err != nil {
		return fmt.Errorf("failed to create API client: %w", err)
	}

	// Save custom domain
	fmt.Printf("Configuring custom domain for application %s...\n", applicationID)
	resp, err := client.SaveCustomDomain(applicationID, domain)
	if err != nil {
		return fmt.Errorf("failed to configure custom domain: %w", err)
	}

	fmt.Printf("\nCustom domain configured successfully!\n")
	fmt.Printf("Message: %s\n", resp.Message)

	return nil
}
