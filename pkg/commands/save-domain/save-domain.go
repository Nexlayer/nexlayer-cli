// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package savedomain

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

func NewCommand(client *api.Client) *cobra.Command {
	var (
		appID string
		domain string
	)

	cmd := &cobra.Command{
		Use:   "save-domain [domain]",
		Short: "Save a custom domain for your application",
		Long: `Save a custom domain for your application.

Endpoint: POST /saveCustomDomain/{applicationID}

Arguments:
  --app      Application ID to save domain for
  domain     Domain name to save (e.g., example.com)

Example:
  nexlayer save-domain example.com --app myapp`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain = args[0]
			return runAddDomain(cmd, client, appID, domain)
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

	ui.RenderTitleWithBorder("Adding Custom Domain")

	// Save custom domain
	_, err := client.SaveCustomDomain(cmd.Context(), appID, domain)
	if err != nil {
		return fmt.Errorf("failed to save custom domain: %w", err)
	}

	// Print success message
	ui.RenderSuccess(fmt.Sprintf("Custom domain '%s' saved successfully", domain))
	return nil
}
