// Package commands contains the CLI commands for the application.
package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// InitCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new project",
	Long: `Initialize a new project with the given name.
Example: nexlayer init my-app`,
	Args: cobra.ExactArgs(1),
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	projectName := args[0]

	// Create project directory
	if err := os.MkdirAll(projectName, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create templates directory
	templatesDir := filepath.Join(projectName, "templates")
	if err := os.MkdirAll(templatesDir, 0o755); err != nil {
		return fmt.Errorf("failed to create templates directory: %w", err)
	}

	// Create default template file
	templatePath := filepath.Join(templatesDir, "default.yaml")
	defaultTemplate := []byte(`name: default
version: 1.0.0
description: Default template
components:
  - name: app
    type: container
    image: nginx:latest
    ports:
      - 80:80
`)
	if err := os.WriteFile(templatePath, defaultTemplate, 0o644); err != nil {
		return fmt.Errorf("failed to create default template: %w", err)
	}

	fmt.Printf("Created new project: %s\n", projectName)
	fmt.Printf("  - Created %s/\n", projectName)
	fmt.Printf("  - Created %s/templates/\n", projectName)
	fmt.Printf("  - Created %s/templates/default.yaml\n", projectName)
	fmt.Println("\nNext steps:")
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  nexlayer deploy default")

	return nil
}
