// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/templates"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewCommand creates a new init command that supports both existing projects
// and new project creation from templates.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Nexlayer project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine project name from argument or use current directory name
			var projectName string
			if len(args) > 0 {
				projectName = args[0]
			} else {
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("unable to determine current directory. Please ensure you have proper permissions and try again. Error: %w", err)
				}
				projectName = filepath.Base(dir)
			}

			// Get current directory
			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			// Check if directory is empty
			files, err := os.ReadDir(dir)
			if err != nil {
				return fmt.Errorf("cannot read contents of directory '%s'. Please verify: (1) you have read permissions, (2) the directory is not locked, (3) the path is valid. Error: %w", dir, err)
			}

			// If directory is empty or only contains hidden files, offer templates
			hasVisibleFiles := false
			for _, file := range files {
				if !file.IsDir() && !isHiddenFile(file.Name()) {
					hasVisibleFiles = true
					break
				}
			}

			if !hasVisibleFiles {
				// Mode 2: Start from Scratch
				pterm.Info.Println("📦 No project detected. Would you like to start with a template?")

				// Create template selection list
				items := templates.GetTemplateItems()
				l := list.New(items, list.NewDefaultDelegate(), 0, 0)
				l.Title = "Select a template"
				l.SetShowStatusBar(false)
				l.SetFilteringEnabled(false)
				l.Styles.Title = lipgloss.NewStyle().MarginLeft(2)

				// Show template selection
				selected := l.SelectedItem()
				if selected == nil {
					return fmt.Errorf("no template was selected. Please choose a template from the list to continue")
				}

				// Create project from template
				templateName := selected.(templates.TemplateItem).Name
				pterm.Info.Printf("✨ Selected: %s\n", templateName)
				pterm.Info.Println("🚀 Generating starter files...")

				if err := templates.CreateProject(projectName, templateName); err != nil {
					return fmt.Errorf("failed to create project from template '%s'. Please ensure: (1) template name is valid, (2) you have write permissions, (3) template files are accessible. Error: %w", templateName, err)
				}
			}

			// Mode 1: Detect & Generate
			progress, _ := pterm.DefaultProgressbar.WithTotal(100).Start()
			progress.Title = "Analyzing project"

			// Check if nexlayer.yaml already exists
			configFile := "nexlayer.yaml"
			if _, err := os.Stat(configFile); err == nil {
				// File exists, create backup
				backupFile := configFile + ".backup"
				if err := os.Rename(configFile, backupFile); err != nil {
					return fmt.Errorf("failed to create backup of existing nexlayer.yaml. Please ensure: (1) you have write permissions, (2) original file is not in use, (3) sufficient disk space available. Error: %w", err)
				}
				cmd.Printf("Backed up existing %s to %s\n", configFile, backupFile)
			}

			// Generate the template using AI
			progress.Title = "Analyzing project and generating template"
			stackType, components := ai.DetectStack(dir)
			req := ai.TemplateRequest{
				ProjectName: projectName,
				TemplateType: stackType,
				RequiredFields: map[string]interface{}{
					"components": components,
				},
			}
			yamlStr, err := ai.GenerateTemplate(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("failed to generate template: %w", err)
			}
			progress.Add(90)

			// Write the generated template to 'nexlayer.yaml'
			if err := os.WriteFile("nexlayer.yaml", []byte(yamlStr), 0644); err != nil {
				return fmt.Errorf("failed to write template: %w", err)
			}
			progress.Add(10)

			// Display success message
			progress.Stop()
			pterm.Success.Printf("Created nexlayer.yaml for %s\n", projectName)
			fmt.Println("To deploy your application, run: nexlayer deploy")
			return nil
		},
	}

	return cmd
}

// isHiddenFile returns true if the file name starts with a dot
func isHiddenFile(name string) bool {
	return len(name) > 0 && name[0] == '.'
}
