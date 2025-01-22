package commands

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	debug       bool
	environment string
)

func init() {
	DeployCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	DeployCmd.Flags().StringVarP(&environment, "env", "e", "staging", "Environment to deploy to (staging or production)")
}

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy [template]",
	Short: "Deploy a template",
	Long: `Deploy a template to Nexlayer.
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

	// Read template file
	yamlFile := fmt.Sprintf("templates/%s.yaml", templateName)
	yamlContent, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Get API endpoint based on environment
	var env string
	switch environment {
	case "staging":
		env = "staging"
	case "production":
		env = "production"
	default:
		return fmt.Errorf("invalid environment: %s (must be staging or production)", environment)
	}

	cfg := config.GetConfig()
	baseURL := cfg.GetAPIEndpoint(env)

	if debug {
		fmt.Printf("Debug: Using API endpoint %s\n", baseURL)
	}

	client := api.NewClient(baseURL)
	resp, err := client.StartUserDeployment(sessionID, yamlContent)
	if err != nil {
		return fmt.Errorf("failed to start deployment: %w", err)
	}

	fmt.Printf("Starting deployment...\n")
	fmt.Printf("Template: %s\n", templateName)
	fmt.Printf("Namespace: %s\n", resp.Namespace)
	fmt.Printf("URL: %s\n", resp.URL)
	fmt.Printf("Message: %s\n", resp.Message)
	fmt.Printf("Deployment started successfully!\n")

	return nil
}
