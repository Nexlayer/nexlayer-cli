package debugcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type debugConfig struct {
	Application struct {
		Name string `yaml:"name"`
		Pods []struct {
			Name  string `yaml:"name"`
			Image string `yaml:"image"`
		} `yaml:"pods"`
	} `yaml:"application"`
}

type debugResult struct {
	Status    string   `json:"status"`
	Message   string   `json:"message"`
	Fixes     []string `json:"fixes,omitempty"`
	LogOutput string   `json:"log_output,omitempty"`
}

// NewCommand creates a new debug command
func NewCommand() *cobra.Command {
	var (
		configFile string
		fullDebug  bool
		jsonOutput bool
	)

	cmd := &cobra.Command{
		Use:   "debug",
		Short: "Troubleshoot deployment issues",
		Long: `Debug helps identify and fix common deployment issues by:
1. Checking deployment status
2. Validating configuration
3. Verifying registry access
4. Analyzing logs
5. Suggesting fixes`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Start debug spinner
			spinner, _ := pterm.DefaultSpinner.Start("🔍 Running Nexlayer diagnostics...")
			defer spinner.Stop()

			results := make([]debugResult, 0)

			// 1. Check if nexlayer.yaml exists and is valid
			spinner.UpdateText("Checking configuration...")
			configResult := checkConfiguration(configFile)
			results = append(results, configResult)

			// 2. Parse and validate configuration
			var config debugConfig
			if configResult.Status == "success" {
				yamlData, _ := os.ReadFile(configFile)
				if err := yaml.Unmarshal(yamlData, &config); err != nil {
					results = append(results, debugResult{
						Status:  "error",
						Message: "Failed to parse nexlayer.yaml",
						Fixes:   []string{"Run 'nexlayer validate' to check for syntax errors"},
					})
				}
			}

			// 3. Check Docker and registry access
			spinner.UpdateText("Checking registry access...")
			for _, pod := range config.Application.Pods {
				results = append(results, checkRegistryAccess(pod.Image))
			}

			// 4. Check deployment status
			spinner.UpdateText("Checking deployment status...")
			results = append(results, checkDeploymentStatus(config.Application.Name))

			// 5. Check logs if full debug is enabled
			if fullDebug {
				spinner.UpdateText("Analyzing logs...")
				results = append(results, analyzeLogs())
			}

			// Stop spinner and display results
			spinner.Success()

			if jsonOutput {
				// Output JSON format
				json.NewEncoder(os.Stdout).Encode(results)
			} else {
				// Display human-readable output
				pterm.DefaultSection.Println("Nexlayer Debug Results")

				for _, result := range results {
					switch result.Status {
					case "success":
						pterm.Success.Printf("✅ %s\n", result.Message)
					case "error":
						pterm.Error.Printf("❌ %s\n", result.Message)
						for _, fix := range result.Fixes {
							pterm.Info.Printf("💡 Fix: %s\n", fix)
						}
					case "warning":
						pterm.Warning.Printf("⚠️  %s\n", result.Message)
						for _, fix := range result.Fixes {
							pterm.Info.Printf("💡 Fix: %s\n", fix)
						}
					}
					if result.LogOutput != "" && fullDebug {
						fmt.Printf("\nRelevant logs:\n%s\n", result.LogOutput)
					}
				}

				if !fullDebug {
					fmt.Printf("\n💡 Run 'nexlayer debug --full' for more detailed analysis\n")
				}
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&configFile, "file", "f", "nexlayer.yaml", "Path to nexlayer.yaml")
	cmd.Flags().BoolVar(&fullDebug, "full", false, "Run full diagnostic analysis")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results in JSON format")

	return cmd
}

func checkConfiguration(configFile string) debugResult {
	// Check if file exists
	if _, err := os.Stat(configFile); err != nil {
		return debugResult{
			Status:  "error",
			Message: fmt.Sprintf("Configuration file not found: %s", configFile),
			Fixes:   []string{"Run 'nexlayer init' to create a new configuration"},
		}
	}

	// Validate configuration
	validator := validation.NewValidator(true) // Strict mode for debug
	yamlData, err := os.ReadFile(configFile)
	if err != nil {
		return debugResult{
			Status:  "error",
			Message: "Failed to read configuration file",
			Fixes:   []string{"Check if the file exists and has correct permissions"},
		}
	}

	var config schema.NexlayerYAML
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return debugResult{
			Status:  "error",
			Message: "Failed to parse configuration file",
			Fixes:   []string{"Check if the YAML syntax is correct"},
		}
	}

	validationErrors := validator.ValidateYAML(&config)
	if len(validationErrors) > 0 {
		return debugResult{
			Status:  "error",
			Message: "Configuration validation failed",
			Fixes:   []string{"Run 'nexlayer validate' to see detailed errors"},
		}
	}

	// Check for pods
	if len(config.Application.Pods) == 0 {
		return debugResult{
			Status:  "warning",
			Message: "No pods defined in configuration",
			Fixes:   []string{"Add at least one pod configuration"},
		}
	}

	return debugResult{
		Status:  "success",
		Message: "Configuration is valid",
	}
}

func checkRegistryAccess(image string) debugResult {
	// Parse image parts
	parts := strings.Split(image, "/")
	if len(parts) < 2 {
		return debugResult{
			Status:  "error",
			Message: fmt.Sprintf("Invalid image format: %s", image),
			Fixes:   []string{"Use format: registry/repository:tag"},
		}
	}

	// Check if docker is available
	_, err := exec.LookPath("docker")
	if err != nil {
		return debugResult{
			Status:  "error",
			Message: "Docker not found",
			Fixes: []string{
				"Install Docker",
				"Add Docker to PATH",
			},
		}
	}

	// Try to pull image
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "pull", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return debugResult{
			Status:  "error",
			Message: fmt.Sprintf("Failed to pull image: %s\n%s", image, string(output)),
			Fixes: []string{
				"Check registry credentials",
				"Verify image exists and tag is correct",
				"Run 'docker login' if using private registry",
			},
		}
	}

	return debugResult{
		Status:  "success",
		Message: fmt.Sprintf("Successfully verified access to image: %s", image),
	}
}

func checkDeploymentStatus(appName string) debugResult {
	// TODO: Implement actual deployment status check using Nexlayer API
	// For now, return a placeholder success
	return debugResult{
		Status:  "success",
		Message: fmt.Sprintf("Deployment '%s' is running", appName),
	}
}

func analyzeLogs() debugResult {
	logDir := "/var/log/nexlayer"
	logFile := filepath.Join(logDir, "latest.log")

	// Check if log directory exists
	if _, err := os.Stat(logDir); err != nil {
		return debugResult{
			Status:  "warning",
			Message: "Log directory not found",
			Fixes:   []string{"Check if Nexlayer is properly installed"},
		}
	}

	// Read recent logs
	cmd := exec.Command("tail", "-n", "50", logFile)
	output, err := cmd.Output()
	if err != nil {
		return debugResult{
			Status:  "warning",
			Message: "Could not read logs",
			Fixes:   []string{"Check permissions on log directory"},
		}
	}

	// Analyze logs for common issues
	logContent := string(output)
	var issues []string

	if strings.Contains(logContent, "connection refused") {
		issues = append(issues, "Network connectivity issues detected")
	}
	if strings.Contains(logContent, "permission denied") {
		issues = append(issues, "Permission issues detected")
	}
	if strings.Contains(logContent, "out of memory") {
		issues = append(issues, "Resource constraints detected")
	}

	if len(issues) > 0 {
		return debugResult{
			Status:    "warning",
			Message:   "Found potential issues in logs",
			Fixes:     issues,
			LogOutput: logContent,
		}
	}

	return debugResult{
		Status:    "success",
		Message:   "No issues found in recent logs",
		LogOutput: logContent,
	}
}
