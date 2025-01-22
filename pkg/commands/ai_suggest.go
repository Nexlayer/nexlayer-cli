package commands

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/spf13/cobra"
)

// AISuggestCmd represents the ai-suggest command
var AISuggestCmd = &cobra.Command{
	Use:   "ai-suggest [query]",
	Short: "Get AI suggestions",
	Long: `Get AI-powered suggestions for your deployment.
Example: nexlayer ai-suggest "optimize my nodejs app"`,
	Args: cobra.ExactArgs(1),
	RunE: runAISuggest,
}

func runAISuggest(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Get session ID from environment
	sessionID := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if sessionID == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create client
	client := api.NewClient("https://app.staging.nexlayer.io")
	suggestions, err := client.GetAISuggestions(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get AI suggestions: %w", err)
	}

	fmt.Printf("AI Suggestions for: %s\n\n", query)
	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion)
	}

	return nil
}
