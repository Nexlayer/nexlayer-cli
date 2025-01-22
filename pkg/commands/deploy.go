package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	environment string
	debug       bool
)

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy [templateId]",
	Short: "Deploy an application",
	Long: `Deploy an application using either a predefined template or a custom YAML configuration.
	
Examples:
  # Deploy using a template
  nexlayer deploy hello-world

  # Deploy to a specific environment
  nexlayer deploy hello-world --env production

  # Deploy with debug output
  nexlayer deploy hello-world --debug`,
	Args: cobra.ExactArgs(1),
	RunE: runDeploy,
}

func init() {
	DeployCmd.Flags().StringVarP(&environment, "env", "e", "staging", "Target environment (staging/production)")
	DeployCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	templateName := args[0]
	
	// Create a spinner for better UX
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Prefix = "Preparing deployment "
	s.Start()

	// Validate environment
	if environment != "staging" && environment != "production" {
		return errors.NewValidationError(
			"Invalid environment specified",
			nil,
			"Use --env staging or --env production",
			"Default environment is staging if not specified",
		)
	}

	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return errors.NewAuthError(
			"Authentication required",
			nil,
			"Visit https://app.nexlayer.io/settings/tokens to generate a token",
		)
	}

	// Read and validate the template YAML
	s.Suffix = " Reading template configuration"
	templatePath := filepath.Join("examples", "plugins", "template-builder", "template-builder-nexlayer-template.yaml")
	yamlContent, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return errors.NewConfigError(
			"Failed to read template file",
			err,
			fmt.Sprintf("Ensure the file exists at: %s", templatePath),
			"Check file permissions",
		)
	}

	// Validate YAML structure
	var templateConfig map[string]interface{}
	if err := yaml.Unmarshal(yamlContent, &templateConfig); err != nil {
		return errors.NewValidationError(
			"Invalid template YAML",
			err,
			"Check the YAML syntax",
			"Ensure all required fields are present",
			"Run 'nexlayer validate' to check template structure",
		)
	}

	// Update spinner for deployment
	s.Suffix = " Initiating deployment"
	
	// Create client with appropriate URL
	baseURL := "https://app.staging.nexlayer.io"
	if environment == "production" {
		baseURL = "https://app.nexlayer.io"
	}
	
	if debug {
		fmt.Printf("\n🔍 Debug: Using API endpoint %s\n", baseURL)
	}

	client := api.NewClient(baseURL)
	resp, err := client.StartDeployment(sessionID, yamlContent)
	if err != nil {
		return errors.NewDeploymentError(
			"Deployment failed",
			err,
			"Check your network connection",
			"Verify your authentication token",
			"Run with --debug flag for more information",
		)
	}

	// Stop spinner and show success message
	s.Stop()
	success := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("\n%s Deployment successful!\n", success("✓"))
	
	// Print deployment details
	fmt.Printf("\n📋 Deployment Details:\n")
	fmt.Printf("   • Namespace: %s\n", resp.Namespace)
	fmt.Printf("   • URL: %s\n", resp.URL)
	fmt.Printf("   • Environment: %s\n", environment)
	if resp.Message != "" {
		fmt.Printf("   • Message: %s\n", resp.Message)
	}

	fmt.Printf("\n💡 Next steps:\n")
	fmt.Printf("   • Monitor status: nexlayer status %s\n", resp.Namespace)
	fmt.Printf("   • View logs: nexlayer logs %s\n", resp.Namespace)
	fmt.Printf("   • Scale app: nexlayer scale %s --replicas 3\n", resp.Namespace)

	return nil
}
