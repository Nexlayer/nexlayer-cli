package wizard

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

type DeploymentConfig struct {
	AppID     string            `yaml:"appId"`
	Template  string            `yaml:"template"`
	Variables map[string]string `yaml:"variables,omitempty"`
}

func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive deployment wizard",
		Long: `Create a new deployment using an interactive wizard.
		
The wizard will guide you through:
1. Selecting a template
2. Configuring environment variables
3. Creating a deployment configuration file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWizard(cmd)
		},
	}

	return cmd
}

func runWizard(cmd *cobra.Command) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployment Wizard"))

	// Get application ID
	var appID string
	cmd.Print("Enter application ID: ")
	fmt.Scanln(&appID)

	// Get template name
	var template string
	cmd.Print("Enter template name (e.g., 'node', 'python'): ")
	fmt.Scanln(&template)

	// Create config
	config := DeploymentConfig{
		AppID:    appID,
		Template: template,
		Variables: map[string]string{
			"PORT": "8080",
		},
	}

	// Save config to file
	configPath := "deployment.yaml"
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
	cmd.Printf("  nexlayer deploy --app %s --file %s\n", appID, configPath)

	return nil
}
