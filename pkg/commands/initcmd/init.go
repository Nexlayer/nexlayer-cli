// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package initcmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewCommand initializes a new Nexlayer project with automatic detection
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Nexlayer project",
		Long:  "Initialize a new Nexlayer project with intelligent stack detection and configuration generation.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runInitCommand,
	}
	return cmd
}

// DetectIDE checks which AI-powered IDE is currently being used
func DetectIDE() string {
	// 1️⃣ Check environment variables
	if os.Getenv("CURSOR") != "" {
		return "Cursor"
	}
	if os.Getenv("WINDSURF") != "" {
		return "Windsurf"
	}
	if os.Getenv("VSCODE_GIT_IPC_HANDLE") != "" || os.Getenv("VSCODE_PID") != "" {
		return "VSCode"
	}
	if os.Getenv("ZED_ROOT") != "" {
		return "Zed"
	}
	if os.Getenv("AIDER_PROJECT") != "" {
		return "Aider"
	}

	// 2️⃣ Check running processes
	processes := []string{"cursor", "code", "windsurf", "zed", "aider"}
	for _, process := range processes {
		cmd := exec.Command("pgrep", "-x", process)
		if err := cmd.Run(); err == nil {
			return strings.Title(process)
		}
	}

	// 3️⃣ Check configuration files in the project directory
	configFiles := map[string]string{
		".cursor":          "Cursor",
		".vscode":          "VSCode",
		"windsurf.json":    "Windsurf",
		"zed-settings.json": "Zed",
		".aider":           "Aider",
	}

	for file, ide := range configFiles {
		if _, err := os.Stat(file); err == nil {
			return ide
		}
	}

	return "Unknown"
}

// runInitCommand handles the execution of the init command
func runInitCommand(cmd *cobra.Command, args []string) error {
	// Determine project name
	projectName := getProjectName(args)

	// Display welcome message
	fmt.Print(ui.RenderBox(ui.RenderWelcome()))

	// Start progress bar
	progress, err := pterm.DefaultProgressbar.WithTotal(100).Start()
	if err != nil {
		return fmt.Errorf("failed to start progress bar: %w", err)
	}
	defer progress.Stop()

	// Detect IDE being used
	ide := DetectIDE()
	if ide != "Unknown" {
		fmt.Printf("🖥️  Detected AI-powered IDE: %s\n", ide)
	}

	// Detect Project Type
	progress.UpdateTitle("Analyzing project...")
	info, err := detectProject(progress)
	if err != nil {
		return err
	}

	// Detect the LLM provider
	llmProvider := detection.DetectAIIDE()

	// Generate YAML based on detected project
	progress.UpdateTitle("Generating Nexlayer configuration...")
	yamlContent, err := generateProjectYAML(projectName, info, ide, llmProvider)
	if err != nil {
		return err
	}
	progress.Add(20)

	// Write YAML to file
	if err := writeYAMLToFile("nexlayer.yaml", yamlContent); err != nil {
		return err
	}
	progress.Add(20)

	// Validate YAML
	if err := validateGeneratedYAML(yamlContent); err != nil {
		return err
	}
	progress.Add(20)

	// Completion message
	progress.Stop()
	pterm.Success.Printf("\n✨ Your Nexlayer project is ready!\n")
	pterm.Info.Println("\nDeploy with:")
	fmt.Println("  nexlayer deploy")
	return nil
}

// getProjectName determines the project name from arguments or directory
func getProjectName(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	dir, err := os.Getwd()
	if err != nil {
		pterm.Error.Println("Unable to determine current directory. Defaulting to 'new-project'.")
		return "new-project"
	}
	return filepath.Base(dir)
}

// detectProject attempts to detect the project type
func detectProject(progress *pterm.ProgressbarPrinter) (*detection.ProjectInfo, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	detector := detection.NewDetectorRegistry()
	progress.Add(20)

	info, err := detector.DetectProject(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to detect project type: %w", err)
	}

	pterm.Success.Printf("✅ Found %s project\n", info.Type)
	return info, nil
}

// generateProjectYAML creates the Nexlayer configuration file
func generateProjectYAML(projectName string, info *detection.ProjectInfo, ide string, llmProvider string) (string, error) {
	// Try to generate YAML based on project detection
	yamlContent, err := detection.GenerateYAML(info)
	if err != nil {
		pterm.Warning.Println("⚠️  Using basic template - some features may need manual configuration")
		// Use a basic template following v1.2 schema
		yamlContent = fmt.Sprintf(`application:
  name: %s
  pods:
    - name: app
      type: frontend
      image: nginx:latest
      servicePorts:
        - 80
`, projectName)
	}

	// Validate the generated YAML
	err = validateGeneratedYAML(yamlContent)
	if err != nil {
		return "", fmt.Errorf("generated YAML validation failed: %w", err)
	}

	return yamlContent, nil
}

// writeYAMLToFile writes the YAML content to a file
func writeYAMLToFile(filename string, content string) error {
	if _, err := os.Stat(filename); err == nil {
		backupFile := filename + ".backup"
		if err := os.Rename(filename, backupFile); err != nil {
			return fmt.Errorf("failed to create backup of existing %s: %w", filename, err)
		}
		pterm.Info.Printf("📦 Backed up existing %s to %s\n", filename, backupFile)
	}

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}
	pterm.Info.Println("✅ Configuration written to nexlayer.yaml")
	return nil
}

// validateGeneratedYAML checks for YAML syntax errors
func validateGeneratedYAML(yamlContent string) error {
	var config schema.NexlayerYAML
	if err := yaml.Unmarshal([]byte(yamlContent), &config); err != nil {
		return fmt.Errorf("failed to parse generated YAML: %w", err)
	}

	validator := validation.NewValidator(false)
	validationErrors := validator.ValidateYAML(&config)
	if len(validationErrors) > 0 {
		pterm.Warning.Println("⚠️  Validation found issues:")

		for _, err := range validationErrors {
			fmt.Printf("❌ %s: %s\n", err.Field, err.Message)
			for _, suggestion := range err.Suggestions {
				fmt.Printf("   💡 %s\n", suggestion)
			}
		}

		fmt.Println("\n📝 Next Steps:")
		fmt.Println("1. Fix issues in nexlayer.yaml")
		fmt.Println("2. Run 'nexlayer validate'")
		fmt.Println("3. Once validation passes, run 'nexlayer deploy'")
	}
	pterm.Success.Println("✅ Validation passed")
	return nil
}
