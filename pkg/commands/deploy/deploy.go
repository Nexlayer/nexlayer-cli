// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
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
func NewCommand(client api.APIClient) *cobra.Command {
	var yamlFile string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application to Nexlayer",
		Long: `Deploy an application to Nexlayer using a YAML configuration file.

Endpoint: POST /startUserDeployment/{applicationID?}

Arguments:
  --app      Optional: Application ID to deploy to. If not provided, uses Nexlayer profile
  --config   Path to YAML configuration file

The YAML file must follow the Nexlayer schema v2 format with required fields:
  application:
    name: string      # Unique deployment name
    url: string      # Optional permanent domain
    pods:            # List of pod configurations
      - name: string   # Pod name (lowercase alphanumeric)
        image: string  # Fully qualified image path
        servicePorts: []

If no config file is specified, it will look for one of these files:
- deployment.yaml
- deployment.yml
- nexlayer.yaml
- nexlayer.yml

Response will include:
- Deployment status message
- Generated namespace
- Application URL

Examples:
  # Deploy with application ID
  nexlayer deploy --app my-app-123 --config deploy.yaml

  # Deploy using Nexlayer profile
  nexlayer deploy --config deploy.yaml`,
		Args: cobra.ExactArgs(0),
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

			return runDeploy(cmd, yamlFile, appID)
		},
	}

	cmd.Flags().StringVarP(&yamlFile, "config", "f", "", "Path to YAML configuration file")
	// Make app flag optional
	var appID string
	cmd.Flags().StringVar(&appID, "app", "", "Application ID (optional)")
	
	// Mark config flag as required
	cmd.MarkFlagRequired("config")

	return cmd
}

func runDeploy(cmd *cobra.Command, yamlFile string, appID string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deploying Application"))

	// Create API client
	apiClient := api.NewClient("")

	// Start deployment
	if appID == "" {
		cmd.Println("No application ID provided, using Nexlayer profile")
	}
	resp, err := apiClient.StartDeployment(context.Background(), appID, yamlFile)
	if err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}

	cmd.Printf("\nDeployment started successfully!\n")
	cmd.Printf("URL: %s\n", resp.URL)
	cmd.Printf("Namespace: %s\n", resp.Namespace)
	cmd.Printf("Status: %s\n", resp.Message)

	return nil
}
