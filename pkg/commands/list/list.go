// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package list

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List deployments",
		Long:  `List all deployments for an application.`,
	}

	cmd.AddCommand(newDeploymentsCommand(client))
	return cmd
}

func newDeploymentsCommand(client *api.Client) *cobra.Command {
	var appID string

	cmd := &cobra.Command{
		Use:   "deployments",
		Short: "List deployments",
		Long: `List all deployments for an application.
		
Example:
  nexlayer list deployments --app myapp`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListDeployments(cmd, client, appID)
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID")
	cmd.MarkFlagRequired("app")

	return cmd
}

func runListDeployments(cmd *cobra.Command, client *api.Client, appID string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployments"))

	resp, err := client.GetDeployments(cmd.Context(), appID)
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}
	deployments := resp.Data

	if len(deployments) == 0 {
		cmd.Println("No deployments found")
		return nil
	}

	// Prepare table data
	headers := []string{"Namespace", "Template", "Status"}
	rows := make([][]string, len(deployments))
	for i, d := range deployments {
		rows[i] = []string{
			d.Namespace,
			d.TemplateName,
			d.Status,
		}
	}

	cmd.Println(ui.RenderTable(headers, rows))
	return nil
}
