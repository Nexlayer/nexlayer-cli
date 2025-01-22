// Package commands contains the CLI commands for the Nexlayer CLI.
package commands

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

var (
	domainName string
)

func init() {
	DomainCmd.Flags().StringVarP(&domainName, "domain", "d", "", "Custom domain name")
	_ = DomainCmd.MarkFlagRequired("domain")
}

// DomainCmd represents the domain command
var DomainCmd = &cobra.Command{
	Use:   "domain [namespace]",
	Short: "Set custom domain",
	Long: `Set a custom domain for a deployment.
Example: nexlayer domain my-app --domain example.com`,
	Args: cobra.ExactArgs(1),
	RunE: runDomain,
}

func runDomain(cmd *cobra.Command, args []string) error {
	namespace := args[0]

	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create client
	client := api.NewClient("https://app.staging.nexlayer.io")
	err := client.SetCustomDomain(namespace, sessionID, domainName)
	if err != nil {
		return fmt.Errorf("failed to set custom domain: %w", err)
	}

	fmt.Printf("Successfully set custom domain %s for %s\n", domainName, namespace)
	return nil
}
