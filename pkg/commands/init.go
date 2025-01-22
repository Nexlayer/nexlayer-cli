package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Nexlayer CLI configuration",
	Long: `Initialize Nexlayer CLI configuration and set up your environment.
This command will help you set up your authentication token and configure
basic settings for the Nexlayer CLI.`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Welcome to Nexlayer CLI!")
	fmt.Println("\nLet's get you set up with everything you need.")
	
	// Create config directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".nexlayer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if auth token is already set
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		fmt.Println("\n📝 To get started, you'll need an authentication token.")
		fmt.Println("You can get your token by visiting: https://app.nexlayer.io/settings/tokens")
		fmt.Println("\nOnce you have your token, set it in your environment:")
		fmt.Println("\nexport NEXLAYER_AUTH_TOKEN=your_token_here")
		
		// Add to common shell config files
		fmt.Println("\nPro tip: Add this to your shell configuration file (~/.bashrc, ~/.zshrc, etc.)")
		fmt.Println("to make it permanent.")
	} else {
		fmt.Println("\n✅ Authentication token is already set!")
	}

	fmt.Println("\n🎉 Next steps:")
	fmt.Println("1. Run 'nexlayer wizard' to start the interactive deployment wizard")
	fmt.Println("2. Or use 'nexlayer deploy' to deploy an existing template")
	fmt.Println("\n📚 For more information:")
	fmt.Println("- Run 'nexlayer --help' to see all available commands")
	fmt.Println("- Visit our documentation at https://docs.nexlayer.io")
	fmt.Println("- Join our community at https://community.nexlayer.io")

	return nil
}
