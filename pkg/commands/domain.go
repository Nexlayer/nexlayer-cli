package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/api"
)

// DomainCmd represents the domain command
var DomainCmd = &cobra.Command{
	Use:   "domain [domain]",
	Short: "Set a custom domain",
	Long: `Set a custom domain for your deployment.
Example: nexlayer domain mydomain.com`,
	Args: cobra.ExactArgs(1),
	RunE: runDomain,
}

func runDomain(cmd *cobra.Command, args []string) error {
	domain := args[0]
	
	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create client with staging URL
	client := api.NewClient("https://app.staging.nexlayer.io")
	resp, err := client.SaveCustomDomain(sessionID, domain)
	if err != nil {
		return fmt.Errorf("failed to save custom domain: %w", err)
	}

	fmt.Printf("✓ %s\n", resp.Message)
	return nil
}
