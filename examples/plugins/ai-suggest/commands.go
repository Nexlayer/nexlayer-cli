package main

import (
	"github.com/spf13/cobra"
)

var (
	modelFlag        string
	docsPathFlag     string
	templatesPathFlag string
)

var rootCmd = &cobra.Command{
	Use:   "ai-suggest",
	Short: "Get AI-powered suggestions for your Nexlayer applications",
	Long: `AI-powered suggestions for optimizing your Nexlayer applications.
Supports both OpenAI (GPT-4) and Anthropic (Claude) models.

Environment variables:
  OPENAI_API_KEY    - Required for GPT-4
  ANTHROPIC_API_KEY - Required for Claude`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create AI client with documentation paths
		client, err := NewAIClient(modelFlag, docsPathFlag, templatesPathFlag)
		if err != nil {
			return err
		}

		// Run the mods UI
		return RunModsUI(client)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&modelFlag, "model", "m", "openai", "AI model to use (openai or claude)")
	rootCmd.PersistentFlags().StringVar(&docsPathFlag, "docs", "", "Path to documentation directory")
	rootCmd.PersistentFlags().StringVar(&templatesPathFlag, "templates", "", "Path to templates directory")
}

// Execute starts the application
func Execute() error {
	return rootCmd.Execute()
}
