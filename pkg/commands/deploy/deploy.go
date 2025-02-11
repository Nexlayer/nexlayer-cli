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
		Short: "Deploy an application",
		Long: `Deploy an application using a YAML configuration file.

If no file is specified with -f flag, it will look for one of these files in the current directory:
- deployment.yaml
- deployment.yml
- nexlayer.yaml
- nexlayer.yml

Example:
  nexlayer deploy
  nexlayer deploy -f custom-deploy.yaml`,
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
	resp, err := apiClient.StartDeployment(context.Background(), appID, yamlFile)
	if err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}

	cmd.Printf("\nDeployment started successfully!\n")
	cmd.Printf("URL: %s\n", resp.URL)
	cmd.Printf("Namespace: %s\n", resp.Namespace)

	return nil
}
