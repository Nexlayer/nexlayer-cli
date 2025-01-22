package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// AISuggestCmd represents the ai-suggest command
var AISuggestCmd = &cobra.Command{
	Use:   "ai-suggest",
	Short: "Get AI-powered suggestions",
	Long: `Get AI-powered suggestions for your application.
This command is currently not implemented.`,
	RunE: runAISuggest,
}

func init() {
	AISuggestCmd.Flags().String("app", "", "Application ID")
	AISuggestCmd.MarkFlagRequired("app")
}

func runAISuggest(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("AI suggestions are not yet implemented")
}
