package domain

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

// validateDomain checks if a domain name is valid
func validateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain name cannot be empty")
	}

	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return fmt.Errorf("domain name cannot start or end with a dot")
	}

	if strings.Contains(domain, "..") {
		return fmt.Errorf("domain name cannot contain consecutive dots")
	}

	if _, err := net.LookupHost(domain); err != nil {
		return fmt.Errorf("invalid domain name or domain does not exist: %w", err)
	}

	return nil
}

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
	var applicationID string

	cmd := &cobra.Command{
		Use:   "add [domain]",
		Short: "Add a custom domain to an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			if err := validateDomain(domain); err != nil {
				return fmt.Errorf("domain validation failed: %w", err)
			}

			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			err := client.SaveCustomDomain(context.Background(), applicationID, domain)
			if err != nil {
				return fmt.Errorf("failed to add custom domain: %w", err)
			}

			fmt.Printf("Added custom domain %s to application %s\n", domain, applicationID)
			return nil
		},
	}

	cmd.Flags().StringVar(&applicationID, "app-id", "", "Application ID")
	if err := cmd.MarkFlagRequired("app-id"); err != nil {
		panic(fmt.Sprintf("failed to mark app-id flag as required: %v", err))
	}

	return cmd
}

func newListCmd() *cobra.Command {
	var applicationID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List custom domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			domains, err := client.GetDomains(applicationID)
			if err != nil {
				return fmt.Errorf("failed to list domains: %w", err)
			}

			if len(domains) == 0 {
				fmt.Printf("No custom domains found for application %s\n", applicationID)
				return nil
			}

			fmt.Printf("Custom domains for application %s:\n", applicationID)
			fmt.Printf("%-30s %-15s\n", "DOMAIN", "STATUS")
			for _, domain := range domains {
				fmt.Printf("%-30s %-15s\n", domain.Domain, domain.Status)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&applicationID, "app-id", "", "Application ID")
	if err := cmd.MarkFlagRequired("app-id"); err != nil {
		panic(fmt.Sprintf("failed to mark app-id flag as required: %v", err))
	}

	return cmd
}

func newRemoveCmd() *cobra.Command {
	var applicationID string

	cmd := &cobra.Command{
		Use:   "remove [domain]",
		Short: "Remove a custom domain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			if err := client.RemoveDomain(applicationID, domain); err != nil {
				return fmt.Errorf("failed to remove domain: %w", err)
			}

			fmt.Printf("Removed custom domain %s from application %s\n", domain, applicationID)
			return nil
		},
	}

	cmd.Flags().StringVar(&applicationID, "app-id", "", "Application ID")
	if err := cmd.MarkFlagRequired("app-id"); err != nil {
		panic(fmt.Sprintf("failed to mark app-id flag as required: %v", err))
	}

	return cmd
}
