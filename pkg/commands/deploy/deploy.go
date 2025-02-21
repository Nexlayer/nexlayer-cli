// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deploy

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

	return "", fmt.Errorf("no deployment file found in current directory\nExpected one of: %v\nCreate a deployment file or specify one with --file", possibleFiles)
}

// NewCommand creates a new deploy command
func NewCommand(apiClient api.APIClient) *cobra.Command {
	var yamlFile string

	cmd := &cobra.Command{
		Use:   "deploy [applicationID]",
		Short: "Deploy an application",
		Long: `Deploy an application using a deployment YAML file.

The deployment file should be named 'deployment.yaml' or 'nexlayer.yaml' in the current directory.
You can also specify a custom file path using the --file flag.

Arguments:
  applicationID     Optional application ID. If not provided, will use Nexlayer profile.
  --file, -f       Path to deployment YAML file (optional)

Example:
  nexlayer deploy                    # Deploy using deployment.yaml in current directory
  nexlayer deploy myapp             # Deploy specific application
  nexlayer deploy -f custom.yaml    # Deploy using custom file`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no file specified, try to find one
			if yamlFile == "" {
				file, err := findDeploymentFile()
				if err != nil {
					return err
				}
				yamlFile = file
				fmt.Printf("Using deployment file: %s\n", yamlFile)
			}

			// Get app ID if provided
			appID := ""
			if len(args) > 0 {
				appID = args[0]
			}

			return runDeploy(apiClient, yamlFile, appID)
		},
	}

	cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "Path to deployment YAML file")
	return cmd
}

func runDeploy(client api.APIClient, yamlFile string, appID string) error {
	ui.RenderTitleWithBorder("Deploying Application")

	// Read and parse the YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read deployment file: %w", err)
	}

	var config template.NexlayerYAML
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return fmt.Errorf("failed to parse deployment file: %w\nEnsure the file is valid YAML and follows the Nexlayer schema", err)
	}

	// Validate the configuration
	validator := NewValidator(&config)
	if err := validator.Validate(); err != nil {
		ui.RenderError("Validation failed")
		fmt.Println(err)
		return fmt.Errorf("deployment aborted due to validation errors")
	}

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

	// Wait a moment for deployment to initialize
	time.Sleep(2 * time.Second)

	// Get deployment info to show additional details
	info, err := client.GetDeploymentInfo(context.Background(), resp.Data.Namespace, appID)
	if err != nil {
		ui.RenderWarning("Could not fetch deployment details. The deployment is still in progress.")
		fmt.Printf("You can check the status later using: nexlayer info %s\n", appID)
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

// ValidateDeployConfig validates a deployment configuration
// This function is exported for use by other packages
func ValidateDeployConfig(yamlConfig *template.NexlayerYAML) error {
	validator := NewValidator(yamlConfig)
	return validator.Validate()
}
