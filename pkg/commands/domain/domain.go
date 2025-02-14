// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package domain

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// NewDomainCommand creates a new domain command group
func NewDomainCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage custom domains",
		Long:  "Configure and manage custom domains for your applications",
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
		Short: "Configure custom domain",
		Long: `Configure a custom domain for your application.

Endpoint: POST /saveCustomDomain/{applicationID}

Arguments:
  applicationID   Application ID to save domain for
  --domain        Domain name to save (e.g., example.com)

Example:
  nexlayer domain set myapp --domain example.com`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			applicationID := args[0]

			// Validate domain
			if err := ValidateDomain(customDomain); err != nil {
				return err
			}

			// Show progress
			spinner := ui.NewSpinner(fmt.Sprintf("Saving domain %s...", customDomain))
			spinner.Start()
			defer spinner.Stop()

			// Call POST /saveCustomDomain/{applicationID}
			_, err := client.SaveCustomDomain(cmd.Context(), applicationID, customDomain)
			if err != nil {
				return fmt.Errorf("failed to save custom domain: %w", err)
			}

			fmt.Printf("✓ Custom domain %s saved for application %s\n", customDomain, applicationID)
			return nil
		},
	}

	cmd.Flags().StringVar(&customDomain, "domain", "", "Custom domain to set (required)")
	cmd.MarkFlagRequired("domain")

	return cmd
}

// ValidateDomain checks if a domain name is valid
func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}

	if strings.Contains(domain, " ") {
		return fmt.Errorf("domain cannot contain spaces")
	}

	return nil
}
