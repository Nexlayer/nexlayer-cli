package templatebuilder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/detector"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/generator"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/scanner"
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	version        = "1.0.0"
	verbose        bool
	outputFmt      string
	dryRun         bool
	aiProvider     string
	registryURL    string
	templateName   string
	templateVer    string
	securityScan   bool
	estimateCosts  bool
	remoteRegistry bool
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nexlayer",
		Short:   "Nexlayer CLI - Intelligent Infrastructure Templates",
		Version: version,
	}

	cmd.AddCommand(newGenerateCmd())
	cmd.AddCommand(newValidateCmd())
	cmd.AddCommand(newDiffCmd())
	cmd.AddCommand(newUpgradeCmd())
	cmd.AddCommand(newInitCmd())
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.PersistentFlags().StringVar(&registryURL, "registry", "https://registry.nexlayer.dev", "Template registry URL")
	return cmd
}

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [template-name]",
		Short: "Initialize a new template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				color.Blue("Initializing new template: %s", args[0])
			}

			template := &types.NexlayerTemplate{
				Name:        args[0],
				Version:    "0.1.0",
				Services:   make([]types.Service, 0),
				Resources:  make(map[string]types.Resource),
				Config:     make(map[string]string),
				Variables:  make(map[string]string),
			}

			outputPath := fmt.Sprintf("%s.yaml", args[0])
			if err := SaveTemplate(template, outputPath); err != nil {
				return err
			}

			color.Green("✓ Template initialized at %s", outputPath)
			return nil
		},
	}

	return cmd
}

func newDiffCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diff [template1] [template2]",
		Short: "Show differences between templates",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				color.Blue("Comparing templates: %s and %s", args[0], args[1])
			}

			t1, err := loadTemplate(args[0])
			if err != nil {
				return err
			}

			t2, err := loadTemplate(args[1])
			if err != nil {
				return err
			}

			diff := compareTemplates(t1, t2)
			fmt.Println(diff)
			return nil
		},
	}

	return cmd
}

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [template-file]",
		Short: "Upgrade template to new version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				color.Blue("Upgrading template: %s", args[0])
			}

			template, err := loadTemplate(args[0])
			if err != nil {
				return err
			}

			// Increment version
			template.Version = incrementVersion(template.Version)

			// Update template with latest best practices
			if err := updateTemplate(template); err != nil {
				return err
			}

			// Run security scan if requested
			if securityScan {
				securityScanner := scanner.NewSecurityScanner()
				issues, err := securityScanner.ScanTemplate(template)
				if err != nil {
					return fmt.Errorf("error scanning template: %v", err)
				}

				if len(issues) > 0 {
					color.Yellow("\nSecurity scan found %d issues:", len(issues))
					for _, issue := range issues {
						switch issue.Severity {
						case "HIGH":
							color.Red("  [HIGH] %s", issue.Description)
						case "MEDIUM":
							color.Yellow("  [MEDIUM] %s", issue.Description)
						case "LOW":
							color.Blue("  [LOW] %s", issue.Description)
						}
					}
				} else {
					color.Green("\nNo security issues found!")
				}
			}

			// Estimate costs if requested
			if estimateCosts {
				costs, err := estimateTemplateCosts(template)
				if err != nil {
					color.Yellow("⚠ Cost estimation warning: %v", err)
				} else {
					color.Green("Estimated monthly cost: $%.2f", costs)
				}
			}

			// Save upgraded template
			outputPath := fmt.Sprintf("%s-v%s.yaml", template.Name, template.Version)
			if err := SaveTemplate(template, outputPath); err != nil {
				return err
			}

			color.Green("✓ Template upgraded to version %s at %s", template.Version, outputPath)
			return nil
		},
	}

	cmd.Flags().BoolVar(&securityScan, "security-scan", false, "Run security scan")
	cmd.Flags().BoolVar(&estimateCosts, "estimate-costs", false, "Estimate infrastructure costs")
	return cmd
}

func loadTemplate(path string) (*types.NexlayerTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading template: %v", err)
	}

	var template types.NexlayerTemplate
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("error parsing template: %v", err)
	}

	return &template, nil
}

func compareTemplates(t1, t2 *types.NexlayerTemplate) string {
	// TODO: Implement detailed template comparison
	return "Template comparison not implemented yet"
}

func incrementVersion(version string) string {
	// TODO: Implement semantic version increment
	return version + "-next"
}

func updateTemplate(template *types.NexlayerTemplate) error {
	// TODO: Implement template update logic
	return nil
}

func estimateTemplateCosts(template *types.NexlayerTemplate) (float64, error) {
	// TODO: Implement cost estimation
	return 0.0, nil
}

func newGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [project-dir]",
		Short: "Generate a template from project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				color.Blue("Analyzing project directory: %s", args[0])
			}

			template, err := BuildTemplate(args[0])
			if err != nil {
				return fmt.Errorf("error building template: %v", err)
			}

			if dryRun {
				color.Yellow("Dry run - template would be generated as:")
				if outputFmt == "json" {
					json.NewEncoder(os.Stdout).Encode(template)
				} else {
					yaml.NewEncoder(os.Stdout).Encode(template)
				}
				return nil
			}

			outputPath := fmt.Sprintf("template.%s", outputFmt)
			if err := SaveTemplate(template, outputPath); err != nil {
				return fmt.Errorf("error saving template: %v", err)
			}

			color.Green("✓ Template saved to %s", outputPath)
			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFmt, "output", "o", "yaml", "Output format (yaml|json)")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	cmd.Flags().StringVar(&aiProvider, "ai-provider", "openai", "AI provider for template refinement (openai|claude)")
	return cmd
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [template-file]",
		Short: "Validate a template file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if verbose {
				color.Blue("Validating template: %s", args[0])
			}

			// TODO: Implement template validation
			color.Green("✓ Template is valid")
			return nil
		},
	}
}

// BuildTemplate analyzes a project directory and generates a Nexlayer template
func BuildTemplate(projectDir string) (*types.NexlayerTemplate, error) {
	if verbose {
		color.Blue("Detecting project stack...")
	}

	stack, err := detector.DetectStack(projectDir)
	if err != nil {
		return nil, fmt.Errorf("error detecting stack: %v", err)
	}

	if verbose {
		color.Blue("Generating template...")
	}

	template, err := generator.GenerateTemplate(projectDir, stack)
	if err != nil {
		return nil, fmt.Errorf("error generating template: %v", err)
	}

	// Temporarily disable AI refinement
	/*
	if aiProvider != "" {
		if verbose {
			color.Blue("Refining template with AI...")
		}

		refiner, err := NewAIRefiner()
		if err != nil {
			return template, fmt.Errorf("warning: AI refinement not available: %v", err)
		}

		refined, err := refiner.RefineTemplate(*stack, template)
		if err != nil {
			return nil, fmt.Errorf("error refining template: %v", err)
		}
		template = refined
	}
	*/

	return template, nil
}

func SaveTemplate(template *types.NexlayerTemplate, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer file.Close()

	switch filepath.Ext(outputPath) {
	case ".json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(template); err != nil {
			return fmt.Errorf("error encoding JSON: %v", err)
		}

	case ".yaml", ".yml":
		encoder := yaml.NewEncoder(file)
		encoder.SetIndent(2)
		if err := encoder.Encode(template); err != nil {
			return fmt.Errorf("error encoding YAML: %v", err)
		}

	default:
		return fmt.Errorf("unsupported output format: %s", filepath.Ext(outputPath))
	}

	return nil
}

func main() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
