// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package domain

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
	"github.com/spf13/cobra"
)

// NewDomainCommand creates a new domain command group
func NewDomainCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage custom domains for your applications",
		Long: `Configure and manage custom domains for your Nexlayer applications.

Key Features:
  • Map custom domains to your applications
  • Automatic SSL certificate provisioning
  • DNS validation and health checks
  • Zero-downtime domain updates`,
	}

	// Add set subcommand
	cmd.AddCommand(newSetCommand(client))

	return cmd
}

// newSetCommand creates the set subcommand
func newSetCommand(client api.APIClient) *cobra.Command {
	var customDomain string

	cmd := &cobra.Command{
		Use:   "set <applicationID>",
		Short: "Configure custom domain for an application",
		Long: `Configure a custom domain for your Nexlayer application.

The domain will be automatically configured with:
  • SSL/TLS certificate provisioning
  • DNS validation
  • Health monitoring
  • Zero-downtime updates

Examples:
  nexlayer domain set my-app --domain example.com
  nexlayer domain set api-backend --domain api.mycompany.com`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			applicationID := args[0]

			// Validate domain
			if err := ValidateDomain(customDomain); err != nil {
				return err
			}

			// Show progress
			fmt.Fprintf(cmd.OutOrStdout(), "🔄 Configuring domain %s for application %s...\n", customDomain, applicationID)

			// Call API to save custom domain
			if _, err := client.SaveCustomDomain(cmd.Context(), applicationID, customDomain); err != nil {
				return fmt.Errorf("failed to save custom domain: %w", err)
			}

			// Get deployment info to verify it exists
			deployInfo, err := client.GetDeploymentInfo(cmd.Context(), applicationID)
			if err != nil {
				return fmt.Errorf("failed to get deployment info: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\n✨ Custom domain configured successfully!\n")
			fmt.Fprintf(cmd.OutOrStdout(), "\nNext Steps:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "1. Add the following DNS record to your domain:\n")
			fmt.Fprintf(cmd.OutOrStdout(), "   CNAME %s -> %s\n", customDomain, deployInfo.Data.URL)
			fmt.Fprintf(cmd.OutOrStdout(), "2. Wait for DNS propagation (may take up to 24 hours)\n")
			fmt.Fprintf(cmd.OutOrStdout(), "3. Your domain will be automatically validated and SSL certificate provisioned\n")

			return nil
		},
	}

	cmd.Flags().StringVar(&customDomain, "domain", "", "Custom domain to configure (required)")
	cmd.MarkFlagRequired("domain")

	return cmd
}

// ValidateDomain checks if a domain name is valid using the centralized validation system
func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	// Create a minimal YAML with just the domain to validate
	yaml := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: "temp",
			URL:  domain,
		},
	}

	validator := schema.NewValidator(true)
	errors := validator.ValidateYAML(yaml)

	if len(errors) > 0 {
		// Return the first error with suggestions
		err := errors[0]
		return fmt.Errorf("invalid domain name: %s\n\nSuggestions:\n• %s",
			err.Message,
			strings.Join(err.Suggestions, "\n• "))
	}

	return nil
}
