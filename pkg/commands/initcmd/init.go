// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package initcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

// NewCommand creates a new init command that uses AI to detect the project stack
// and generate an appropriate deployment template.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Nexlayer project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine project name from argument or use current directory name.
			var projectName string
			if len(args) > 0 {
				projectName = args[0]
			} else {
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				projectName = filepath.Base(dir)
			}

			// Create a progress bar for user feedback.
			progress, _ := pterm.DefaultProgressbar.WithTotal(100).Start()
			progress.Title = "Analyzing project"

			// Get current directory for analysis (if needed).
			dir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			// Generate the template using AI.
			progress.Title = "Analyzing project and generating template"
			yamlStr, err := ai.GenerateYAML(projectName, dir, nil)
			if err != nil {
				return fmt.Errorf("failed to generate template: %w", err)
			}
			progress.Add(90)

			// Write the generated template to 'nexlayer.yaml'.
			if err := os.WriteFile("nexlayer.yaml", []byte(yamlStr), 0644); err != nil {
				return fmt.Errorf("failed to write template: %w", err)
			}
			progress.Add(10)

			// Display success message.
			progress.Stop()
			pterm.Success.Printf("Created nexlayer.yaml for %s\n", projectName)
			fmt.Println("To deploy your application, run: nexlayer deploy")
			return nil
		},
	}

	return cmd
}
