// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validate

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewCommand creates a new validate command
func NewCommand() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate a Nexlayer YAML configuration file",
		Long: `Validate a Nexlayer YAML configuration file for correctness.

This command performs comprehensive validation of your configuration file, checking:
- Required fields are present
- Field values are in correct format
- Resource names follow Nexlayer conventions
- Volume sizes are properly formatted
- Registry credentials are complete
- Service ports are valid

If no file is specified, it will look for:
- deployment.yaml
- deployment.yml
- nexlayer.yaml
- nexlayer.yml

Examples:
  # Validate file in current directory
  nexlayer validate

  # Validate specific file
  nexlayer validate custom-config.yaml

  # Validate with strict mode
  nexlayer validate --strict`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var yamlFile string

			// If file specified as argument, use it
			if len(args) > 0 {
				yamlFile = args[0]
			} else {
				// Otherwise look for default files
				possibleFiles := []string{
					"deployment.yaml",
					"deployment.yml",
					"nexlayer.yaml",
					"nexlayer.yml",
				}

				for _, file := range possibleFiles {
					if _, err := os.Stat(file); err == nil {
						yamlFile = file
						cmd.Printf("Using configuration file: %s\n", yamlFile)
						break
					}
				}

				if yamlFile == "" {
					return fmt.Errorf("no configuration file found. Expected one of: %v", possibleFiles)
				}
			}

			// Read and parse YAML file
			data, err := os.ReadFile(yamlFile)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			var config schema.NexlayerYAML
			if err := yaml.Unmarshal(data, &config); err != nil {
				return fmt.Errorf("failed to parse YAML: %w", err)
			}

			// Create validator and validate
			validator := validation.NewValidator(strict)
			errors := validator.ValidateYAML(&config)

			if len(errors) > 0 {
				cmd.Println(ui.RenderError("Validation Failed"))
				for _, err := range errors {
					cmd.Printf("❌ %s\n", err.Error())
				}
				return fmt.Errorf("%d validation error(s) found", len(errors))
			}

			cmd.Println(ui.RenderSuccess("Validation Passed"))
			cmd.Println("✅ Configuration file is valid")
			return nil
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "Enable strict validation mode")

	return cmd
}
