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
	apischema "github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Add at the top with other style variables
var (
	// ... existing styles ...
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffff"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00"))

	// Status styles
	runningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00"))

	pendingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00"))

	failedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))
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

// runDeploy handles the deployment process
func runDeploy(client api.APIClient, yamlFile string, appID string) error {
	ui.RenderTitleWithBorder("Deploying Application")

	// Read and parse the YAML file
	yamlData, err := os.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read deployment file: %w", err)
	}

	var config schema.NexlayerYAML
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return fmt.Errorf("failed to parse deployment file: %w\nEnsure the file is valid YAML and follows the Nexlayer schema", err)
	}

	// Validate the configuration
	validator := schema.NewValidator(true)
	if errors := validator.ValidateYAML(&config); len(errors) > 0 {
		ui.RenderError("Validation failed")
		for _, err := range errors {
			fmt.Println(err)
		}
		return fmt.Errorf("deployment aborted due to validation errors")
	}

	// Start deployment
	if appID == "" {
		fmt.Println("No application ID provided, using Nexlayer profile")
	}

	// Show deployment summary before proceeding
	fmt.Println("\n📋 Deployment Summary:")
	fmt.Printf("• Application: %s\n", config.Application.Name)
	fmt.Printf("• Pods: %d\n", len(config.Application.Pods))
	for _, pod := range config.Application.Pods {
		fmt.Printf("  - %s (%s)\n", pod.Name, pod.Image)
	}
	fmt.Println()

	resp, err := client.StartDeployment(context.Background(), appID, yamlFile)
	if err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}

	ui.RenderSuccess("Deployment started successfully!")
	fmt.Printf("🚀 URL: %s\n", resp.Data.URL)

	// Create context with timeout for status polling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Poll for deployment status with exponential backoff
	fmt.Println("\nWaiting for deployment to stabilize...")
	backoff := 2 * time.Second
	maxBackoff := 10 * time.Second
	spinner := ui.NewSpinner("Checking deployment status")
	spinner.Start()

	for {
		select {
		case <-ctx.Done():
			spinner.Stop()
			ui.RenderWarning("Deployment status check timed out after 5 minutes")
			fmt.Printf("The deployment is still in progress. Check status with: nexlayer info %s\n", appID)
			return nil
		case <-time.After(backoff):
			info, err := client.GetDeploymentInfo(ctx, resp.Data.Namespace, appID)
			if err != nil {
				// Only log the error and continue polling
				fmt.Fprintf(os.Stderr, "Error checking status: %v\n", err)
			} else {
				spinner.Stop()

				// Print deployment details
				ui.RenderTitleWithBorder("Deployment Details")
				fmt.Printf("✨ Status:  %s\n", info.Data.Status)
				fmt.Printf("🌐 URL:     %s\n", info.Data.URL)
				fmt.Printf("📚 Version: %s\n", info.Data.Version)

				if len(info.Data.PodStatuses) > 0 {
					fmt.Println("\n📦 Pods:")
					table := ui.NewTable()
					table.AddHeader("NAME", "STATUS", "READY", "RESTARTS")
					for _, pod := range info.Data.PodStatuses {
						table.AddRow(
							pod.Name,
							formatPodStatus(pod.Status),
							fmt.Sprintf("%v", pod.Ready),
							fmt.Sprintf("%d", pod.Restarts),
						)
					}
					table.Render()
				}

				// Check if deployment is complete
				if isDeploymentStable(info.Data) {
					if info.Data.Status == "running" {
						ui.RenderSuccess("\n✅ Deployment completed successfully!")
						printNextSteps(info.Data)
					} else {
						ui.RenderError("\n❌ Deployment failed")
						printTroubleshootingSteps(info.Data)
					}
					return nil
				}

				// Continue polling with updated spinner message
				spinner = ui.NewSpinner(fmt.Sprintf("Waiting for pods to be ready (%s)", info.Data.Status))
				spinner.Start()
			}

			// Increase backoff time exponentially, but cap it
			backoff = time.Duration(float64(backoff) * 1.5)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

// isDeploymentStable checks if the deployment has reached a stable state
func isDeploymentStable(deployment apischema.Deployment) bool {
	if deployment.Status == "running" || deployment.Status == "failed" {
		return true
	}

	// Check if all pods are ready
	for _, pod := range deployment.PodStatuses {
		if !pod.Ready {
			return false
		}
	}

	return true
}

// formatPodStatus returns a colored status string
func formatPodStatus(status string) string {
	switch status {
	case "running":
		return runningStyle.Render(status)
	case "pending":
		return pendingStyle.Render(status)
	case "failed":
		return failedStyle.Render(status)
	default:
		return status
	}
}

// printNextSteps prints helpful next steps after a successful deployment
func printNextSteps(deployment apischema.Deployment) {
	fmt.Println("\n📝 Next Steps:")
	fmt.Printf("1. Access your application at: %s\n", deployment.URL)
	if deployment.CustomDomain != "" {
		fmt.Printf("2. Custom domain configured: %s\n", deployment.CustomDomain)
	} else {
		fmt.Printf("2. Configure a custom domain: nexlayer domain set %s --domain your-domain.com\n", deployment.Namespace)
	}
	fmt.Printf("3. Monitor logs: nexlayer logs %s\n", deployment.Namespace)
	fmt.Printf("4. Check status: nexlayer info %s %s\n", deployment.Namespace, deployment.Namespace)
}

// printTroubleshootingSteps prints helpful debugging steps when deployment fails
func printTroubleshootingSteps(deployment apischema.Deployment) {
	fmt.Println("\n🔍 Troubleshooting Steps:")
	fmt.Printf("1. Check pod logs: nexlayer logs %s %s\n", deployment.Namespace, deployment.Namespace)
	fmt.Printf("2. View detailed status: nexlayer info %s %s --verbose\n", deployment.Namespace, deployment.Namespace)
	fmt.Println("3. Common issues:")
	fmt.Println("   - Image pull errors: Check image name and registry credentials")
	fmt.Println("   - Resource limits: Ensure pods have sufficient CPU/memory")
	fmt.Println("   - Port conflicts: Verify service port configurations")
	fmt.Println("4. For more help: https://docs.nexlayer.io/troubleshooting")
}

// ValidateDeployConfig validates a deployment configuration
// This function is exported for use by other packages
func ValidateDeployConfig(yamlConfig *schema.NexlayerYAML) error {
	validator := schema.NewValidator(true)
	errors := validator.ValidateYAML(yamlConfig)

	if len(errors) > 0 {
		report := schema.NewValidationReport()
		report.AddErrors(errors)
		return fmt.Errorf("validation failed:\n%v", report.String())
	}

	return nil
}
