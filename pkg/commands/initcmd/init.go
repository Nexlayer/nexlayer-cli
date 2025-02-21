// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package initcmd

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
)

// NewCommand initializes a new Nexlayer project
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Nexlayer project",
		Long:  "Initialize a new Nexlayer project by creating a nexlayer.yaml file in the current directory.",
		Args:  cobra.NoArgs,
		RunE:  runInitCommand,
	}
	return cmd
}

// runInitCommand handles the execution of the init command
func runInitCommand(cmd *cobra.Command, args []string) error {
	// Display welcome message
	ui.RenderWelcome("Welcome to Nexlayer CLI!\nLet's set up your project configuration.")

	// Start progress bar
	progress, err := pterm.DefaultProgressbar.WithTotal(100).Start()
	if err != nil {
		return fmt.Errorf("failed to start progress bar: %w", err)
	}
	defer progress.Stop()

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Create detector registry
	registry := detection.NewRegistry()
	progress.Add(20)

	// Detect project type
	progress.UpdateTitle("Analyzing project...")
	info, err := registry.DetectProject(cmd.Context(), cwd)
	if err != nil {
		return fmt.Errorf("failed to detect project: %w", err)
	}
	progress.Add(20)

	// Create template generator
	progress.UpdateTitle("Generating Nexlayer configuration...")
	generator := template.NewGenerator()

	// Generate template
	tmpl, err := generator.GenerateFromProjectInfo(info.Name, string(info.Type), info.Port)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}
	progress.Add(20)

	// Add database if needed
	if hasDatabase(info) {
		if err := generator.AddPod(tmpl, template.PodTypePostgres, 0); err != nil {
			return fmt.Errorf("failed to add database: %w", err)
		}
	}
	progress.Add(20)

	// Marshal template to YAML
	yamlData, err := yaml.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	// Write YAML to file
	if err := writeYAMLToFile("nexlayer.yaml", string(yamlData)); err != nil {
		return err
	}
	progress.Add(20)

	// Validate YAML
	if err := validateGeneratedYAML(string(yamlData)); err != nil {
		return err
	}

	// Completion message
	progress.Stop()
	fmt.Fprintln(cmd.OutOrStdout(), "✨ Your Nexlayer project is ready!")
	fmt.Fprintln(cmd.OutOrStdout(), "\nDeploy with:")
	fmt.Fprintln(cmd.OutOrStdout(), "  nexlayer deploy")
	return nil
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

	if err := os.WriteFile(filename, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write YAML file: %w", err)
	}
	pterm.Info.Println("✅ Configuration written to nexlayer.yaml")
	return nil
}

// validateGeneratedYAML checks for YAML syntax errors
func validateGeneratedYAML(yamlContent string) error {
	validationErrors, err := validation.ValidateYAMLString(yamlContent)
	if err != nil {
		return fmt.Errorf("failed to validate YAML: %w", err)
	}

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
		return fmt.Errorf("validation failed")
	}

	pterm.Success.Println("✅ Validation passed")
	return nil
}

// hasDatabase checks if the project needs a database
func hasDatabase(info *types.ProjectInfo) bool {
	if info == nil || len(info.Dependencies) == 0 {
		return false
	}

	// Check dependencies for database-related packages
	dbPackages := map[string]bool{
		"pg":         true,
		"postgres":   true,
		"postgresql": true,
		"sequelize":  true,
		"typeorm":    true,
		"prisma":     true,
		"mongoose":   true,
		"mongodb":    true,
		"mysql":      true,
		"mysql2":     true,
		"sqlite3":    true,
		"redis":      true,
	}

	for dep := range info.Dependencies {
		if dbPackages[dep] {
			return true
		}
	}
	return false
}
