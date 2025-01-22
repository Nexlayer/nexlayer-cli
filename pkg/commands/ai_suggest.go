// Formatted with gofmt -s
package commands

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/ai"
	"github.com/spf13/cobra"
)

var AIClient ai.AIClient

// AISuggestCmd represents the ai:suggest command
var AISuggestCmd = &cobra.Command{
	Use:   "ai:suggest [prompt]",
	Short: "Send a prompt to your configured AI provider and receive a suggestion",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if AIClient == nil {
			fmt.Println("AI not configured. Set either OPENAI_API_KEY or CLAUDE_API_KEY environment variable.")
			return
		}

		prompt := strings.Join(args, " ")
		suggestion, err := AIClient.Suggest(prompt)
		if err != nil {
			fmt.Printf("AI Suggestion Error: %v\n", err)
			return
		}

		fmt.Printf("[%s/%s]: %s\n", AIClient.GetProvider(), AIClient.GetModel(), suggestion)
	},
}
