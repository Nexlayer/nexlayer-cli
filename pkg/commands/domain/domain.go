// Package domain contains the CLI commands for the Nexlayer CLI.
package domain

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

// NewCommand creates a new domain command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage custom domains",
		Long:  `Add, list, and remove custom domains for your applications.`,
	}

	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newRemoveCmd())

	return cmd
}

func newAddCmd() *cobra.Command {
	var appName string

	cmd := &cobra.Command{
		Use:   "add [domain]",
		Short: "Add a custom domain to an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			err := client.SaveCustomDomain(context.Background(), appName, domain)
			if err != nil {
				return fmt.Errorf("failed to add custom domain: %w", err)
			}

			fmt.Printf("Added custom domain %s to application %s\n", domain, appName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&appName, "app", "a", "", "Application name")
	cmd.MarkFlagRequired("app")

	return cmd
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list [app]",
		Short: "List custom domains",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			client := api.NewClient(vars.APIURL)
			domains, err := client.GetDomains(appName)
			if err != nil {
				return err
			}
			
			cmd.Printf("%-30s %-15s\n", "DOMAIN", "STATUS")
			for _, domain := range domains {
				cmd.Printf("%-30s %-15s\n", domain.Domain, domain.Status)
			}
			return nil
		},
	}
}

func newRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [app] [domain]",
		Short: "Remove a custom domain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			domain := args[1]
			client := api.NewClient(vars.APIURL)
			return client.RemoveDomain(appName, domain)
		},
	}
}
