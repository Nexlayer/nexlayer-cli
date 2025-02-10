// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sync

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewCommand creates a new sync command that detects project changes
// and updates the nexlayer.yaml file accordingly.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync nexlayer.yaml with project changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if nexlayer.yaml exists
			configFile := "nexlayer.yaml"
			if _, err := os.Stat(configFile); os.IsNotExist(err) {
				return fmt.Errorf("nexlayer.yaml not found. Run 'nexlayer init' first")
			}

			// Create a progress bar for user feedback
			progress, _ := pterm.DefaultProgressbar.WithTotal(100).Start()
			progress.Title = "Scanning project changes"

			// Get current directory for analysis
			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			// Read existing nexlayer.yaml
			existingYAML, err := os.ReadFile(configFile)
			if err != nil {
				return fmt.Errorf("failed to read existing nexlayer.yaml: %w", err)
			}

			// Generate new YAML based on current project state
			progress.Title = "Analyzing project and generating updated template"
			stackType, components := ai.DetectStack(dir)
			req := ai.TemplateRequest{
				ProjectName: "", // Empty since we're updating existing
				TemplateType: stackType,
				RequiredFields: map[string]interface{}{
					"components": components,
				},
			}
			newYAML, err := ai.GenerateTemplate(cmd.Context(), req)
			if err != nil {
				return fmt.Errorf("failed to generate template: %w", err)
			}
			progress.Add(90)

			// If there are no changes, don't update the file
			if string(existingYAML) == newYAML {
				progress.Stop()
				pterm.Info.Println("No changes detected in project structure")
				return nil
			}

			// Confirm with user before updating
			confirm, err := pterm.DefaultInteractiveConfirm.Show("Would you like to update nexlayer.yaml?")
			if err != nil {
				return fmt.Errorf("failed to get user confirmation: %w", err)
			}
			if !confirm {
				progress.Stop()
				return nil
			}

			// Create backup of existing file
			backupFile := configFile + ".backup"
			if err := os.Rename(configFile, backupFile); err != nil {
				return fmt.Errorf("failed to backup existing config: %w", err)
			}

			// Write the updated template
			if err := os.WriteFile(configFile, []byte(newYAML), 0644); err != nil {
				return fmt.Errorf("failed to write template: %w", err)
			}
			progress.Add(10)

			// Display success message
			progress.Stop()
			pterm.Success.Printf("Updated nexlayer.yaml\n")
			pterm.Info.Printf("Backup saved as %s\n", backupFile)
			return nil
		},
	}

	return cmd
}
