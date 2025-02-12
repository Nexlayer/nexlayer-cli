// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package app

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/common"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
)

var Cmd *cobra.Command

func init() {
	apiClient := api.NewClient("https://api.nexlayer.com")
	client := common.NewCommandClient(apiClient)
	Cmd = NewCommand(client)
}

// NewCommand creates a new app command
func NewCommand(client common.CommandClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Manage your applications",
		Long:  `View and manage your Nexlayer applications.`,
	}

	cmd.AddCommand(
		newInfoCommand(client),
	)

	return cmd
}

func newInfoCommand(client common.CommandClient) *cobra.Command {
	var appID string
	var namespace string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get application info",
		Long: `Get detailed information about an application deployment.
		
Example:
  nexlayer app info --app myapp --namespace ecstatic-frog`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInfo(cmd, client, appID, namespace)
		},
	}

	cmd.Flags().StringVarP(&appID, "app", "a", "", "Application ID")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Deployment namespace")
	cmd.MarkFlagRequired("app")
	cmd.MarkFlagRequired("namespace")

	return cmd
}

func runInfo(cmd *cobra.Command, client common.CommandClient, appID string, namespace string) error {
	info, err := client.GetDeploymentInfo(cmd.Context(), namespace, appID)
	if err != nil {
		return fmt.Errorf("failed to get deployment info: %w", err)
	}

	cmd.Printf("Application ID: %s\n", appID)
	cmd.Printf("Namespace:      %s\n", info.Data.Namespace)
	cmd.Printf("Template ID:    %s\n", info.Data.TemplateID)
	cmd.Printf("Template Name:  %s\n", info.Data.TemplateName)
	cmd.Printf("Status:         %s\n", info.Data.Status)
	cmd.Printf("URL:            %s\n", info.Data.URL)
	cmd.Printf("Custom Domain:  %s\n", info.Data.CustomDomain)
	cmd.Printf("Version:        %s\n", info.Data.Version)
	cmd.Printf("Created:        %s\n", info.Data.CreatedAt.Format(time.RFC3339))
	cmd.Printf("Last Updated:   %s\n", info.Data.LastUpdated.Format(time.RFC3339))

	// Display pod statuses
	if len(info.Data.PodStatuses) > 0 {
		cmd.Printf("\nPod Statuses:\n")
		for _, pod := range info.Data.PodStatuses {
			cmd.Printf("\n  Pod: %s (%s)\n", pod.Name, pod.Type)
			cmd.Printf("    Status:    %s\n", pod.Status)
			cmd.Printf("    Ready:     %v\n", pod.Ready)
			cmd.Printf("    Restarts:  %d\n", pod.Restarts)
			cmd.Printf("    Image:     %s\n", pod.Image)
			cmd.Printf("    Created:   %s\n", pod.CreatedAt.Format(time.RFC3339))
		}
	}
	return nil
}
