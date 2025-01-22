package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
)

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy [templateId]",
	Short: "Deploy an application",
	Long: `Deploy an application using either a predefined template or a custom YAML configuration.
Example: nexlayer deploy hello-world`,
	Args: cobra.ExactArgs(1),
	RunE: runDeploy,
}

func runDeploy(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	
	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Read the template YAML file from examples
	templatePath := filepath.Join("examples", "plugins", "template-builder", "template-builder-nexlayer-template.yaml")
	yamlContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	fmt.Printf("Deploying template %s...\n", templateName)

	// Create client with staging URL
	client := api.NewClient("https://app.staging.nexlayer.io")
	resp, err := client.StartDeployment(sessionID, yamlContent)
	if err != nil {
		return fmt.Errorf("failed to deploy: %w", err)
	}

	fmt.Printf("✓ Deployment successful!\n")
	fmt.Printf("  Namespace: %s\n", resp.Namespace)
	fmt.Printf("  URL: %s\n", resp.URL)
	fmt.Printf("  Message: %s\n", resp.Message)

	return nil
}
