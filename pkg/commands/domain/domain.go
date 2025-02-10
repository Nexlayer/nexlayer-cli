// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package domain

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "Manage custom domains",
		Long:  `Add or remove custom domains for your applications.`,
	}

	cmd.AddCommand(newAddCommand(client))
	return cmd
}

func newAddCommand(client *api.Client) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "add [domain]",
		Short: "Add a custom domain",
		Long: `Add a custom domain to your application.
		
Example:
  nexlayer domain add example.com --app myapp`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddDomain(cmd, client, appID, args[0])
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID")
	cmd.MarkFlagRequired("app")

	return cmd
}

// validateDomain checks if a domain name is valid
func validateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}
	// Basic domain validation
	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return fmt.Errorf("invalid domain format: %s", domain)
	}
	return nil
}

func runAddDomain(cmd *cobra.Command, client *api.Client, appID string, domain string) error {
	if err := validateDomain(domain); err != nil {
		return err
	}
	cmd.Println(ui.RenderTitleWithBorder("Adding Custom Domain"))

	err := client.SaveCustomDomain(cmd.Context(), appID, domain)
	if err != nil {
		return fmt.Errorf("failed to add custom domain: %w", err)
	}

	cmd.Printf("Successfully added custom domain: %s\n", domain)
	return nil
}
