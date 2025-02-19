// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
	"context"
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

// findDeploymentFile looks for a deployment file in the current directory
func findDeploymentFile() (string, error) {
	// List of possible deployment file names
	possibleFiles := []string{
		"deployment.yaml",
		"deployment.yml",
		"nexlayer.yaml",
		"nexlayer.yml",
	}

	for _, file := range possibleFiles {
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}

	return "", fmt.Errorf("no deployment file found in current directory. Expected one of: %v", possibleFiles)
}

// NewCommand creates a new deploy command
func NewCommand(apiClient api.APIClient) *cobra.Command {
	var yamlFile string

	cmd := &cobra.Command{
		Use:   "deploy [applicationID]",
		Short: "Deploy an application",
		Long: `Deploy an application using a deployment YAML file.

Endpoint: POST /startUserDeployment/{applicationID?}

Arguments:
  applicationID     Optional application ID
  --file           Path to deployment YAML file

Example:
  nexlayer deploy myapp --file deployment.yaml`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no file specified, try to find one
			if yamlFile == "" {
				file, err := findDeploymentFile()
				if err != nil {
					return err
				}
				yamlFile = file
				cmd.Printf("Using deployment file: %s\n", yamlFile)
			}

			// Get app ID if provided
			appID, _ := cmd.Flags().GetString("app")

			return runDeploy(cmd, apiClient, yamlFile, appID)
		},
	}

	cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "Path to deployment YAML file")
	// Mark file flag as required
	cmd.MarkFlagRequired("file")

	return cmd
}

func runDeploy(_ *cobra.Command, client api.APIClient, yamlFile string, appID string) error {
	ui.RenderTitleWithBorder("Deploying Application")

	// Start deployment
	if appID == "" {
		fmt.Println("No application ID provided, using Nexlayer profile")
	}
	resp, err := client.StartDeployment(context.Background(), appID, yamlFile)
	if err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}

	// Print deployment info
	ui.RenderSuccess("Deployment started successfully!")
	fmt.Printf("🚀 URL: %s\n", resp.Data.URL)

	// Get deployment info to show additional details
	info, err := client.GetDeploymentInfo(context.Background(), resp.Data.Namespace, appID)
	if err != nil {
		ui.RenderError(fmt.Sprintf("Could not fetch deployment details: %v", err))
		return nil
	}

	// Print deployment details
	ui.RenderTitleWithBorder("Deployment Details")
	fmt.Printf("✨ Status:  %s\n", info.Data.Status)
	fmt.Printf("🌐 URL:     %s\n", info.Data.URL)
	fmt.Printf("📚 Version: %s\n", info.Data.Version)

	// Print pod statuses
	if len(info.Data.PodStatuses) > 0 {
		fmt.Println("\n📦 Pods:")
		table := ui.NewTable()
		table.AddHeader("NAME", "STATUS", "READY")
		for _, pod := range info.Data.PodStatuses {
			table.AddRow(pod.Name, pod.Status, fmt.Sprintf("%v", pod.Ready))
		}
		table.Render()
	}

	return nil
}
