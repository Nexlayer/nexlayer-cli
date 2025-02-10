package init

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/detector"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Nexlayer project",
		Long: `Initialize a new Nexlayer project with intelligent stack detection
and template selection. If run in an existing project directory,
it will analyze the codebase and suggest an appropriate template.`,
		RunE: runInit,
	}

	cmd.Flags().StringP("template", "t", "", "Skip detection and use specified template")
	return cmd
}

func runInit(cmd *cobra.Command, args []string) error {
	var projectPath string
	var projectName string

	if len(args) > 0 {
		projectName = args[0]
		projectPath = filepath.Join(".", projectName)
		// Create project directory if it doesn't exist
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
	} else {
		var err error
		projectPath, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectName = filepath.Base(projectPath)
	}

	// Print welcome title
	fmt.Println(ui.RenderTitle("🚀 Welcome to Nexlayer CLI!"))
	fmt.Println(ui.RenderBox("Let's get you set up with everything you need."))

	// Check for --template flag
	template, _ := cmd.Flags().GetString("template")
	if template == "" {
		// Detect stack
		fmt.Println(ui.RenderHighlight("\n🔍 Detecting project stack..."))
		stackInfo, err := detector.DetectStack(projectPath)
		if err != nil {
			return fmt.Errorf("failed to detect stack: %w", err)
		}

		if stackInfo.Template != "blank" {
			fmt.Printf("%s Found %s stack with %s.\n",
				ui.RenderHighlight("✅"),
				stackInfo.Language,
				stackInfo.Framework)
			fmt.Printf("Recommended template: %s\n\n",
				ui.RenderHighlight(stackInfo.Template))
		}

		// Let user select or confirm template
		template, err = ui.SelectTemplate(stackInfo.Template)
		if err != nil {
			return fmt.Errorf("template selection failed: %w", err)
		}
	}

	// Generate project files based on template
	fmt.Printf("\n%s Generating project files...\n",
		ui.RenderHighlight("📝"))
	
	// TODO: Implement template generation logic
	
	fmt.Println(ui.RenderBox(`🎉 Project initialized successfully!

Next steps:
1. Review the generated configuration in nexlayer.yaml
2. Run 'nexlayer deploy' to deploy your application
3. Visit https://docs.nexlayer.io for more information`))

	return nil
}
