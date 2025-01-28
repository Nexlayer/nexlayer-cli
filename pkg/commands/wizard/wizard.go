package wizard

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

type Application struct {
	Template       string            `yaml:"template"`
	DeploymentName string            `yaml:"deploymentName"`
	Variables      map[string]string `yaml:"variables,omitempty"`
}

type DeploymentConfig struct {
	Application Application `yaml:"application"`
}

var validTemplates = map[string]bool{
	"langchain-nextjs":  true,
	"langchain-fastapi": true,
	"mern":             true,
	"pern":             true,
	"mean":             true,
}

func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive deployment wizard",
		Long: `Create a new deployment using an interactive wizard.
		
The wizard will guide you through:
1. Choosing a project name
2. Selecting a template
3. Creating a deployment configuration file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWizard(cmd)
		},
	}

	return cmd
}

func runWizard(cmd *cobra.Command) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployment Wizard"))

	// Get project name
	var projectName string
	cmd.Print("Enter project name: ")
	fmt.Scanln(&projectName)

	if projectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Get template name
	var template string
	cmd.Print("Enter template name (e.g., langchain-nextjs, langchain-fastapi): ")
	fmt.Scanln(&template)

	if !validTemplates[template] {
		return fmt.Errorf("invalid template: %s", template)
	}

	// Create config
	config := DeploymentConfig{
		Application: Application{
			Template:       template,
			DeploymentName: projectName,
			Variables: map[string]string{
				"PORT": "8080",
			},
		},
	}

	// Save config to file
	configPath := "nexlayer.yaml"
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	cmd.Printf("\nConfiguration saved to %s\n", configPath)
	cmd.Println("\nTo deploy your application, run:")
	cmd.Printf("  nexlayer deploy\n")

	return nil
}
