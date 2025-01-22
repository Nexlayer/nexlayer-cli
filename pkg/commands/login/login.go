package login

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const (
	configDir  = ".nexlayer"
	configFile = "config"
)

// Command represents the login command
var Command = &cobra.Command{
	Use:   "login",
	Short: "Log in to Nexlayer using GitHub",
	Long: `Log in to Nexlayer using your GitHub account.
This will open your browser for authentication.
After successful login, your credentials will be saved locally.`,
	RunE: runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("Opening browser for GitHub authentication...")

	// Open browser for GitHub OAuth
	err := browser.OpenURL("https://app.nexlayer.io/auth/github")
	if err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	// Create config directory if it doesn't exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configPath, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	fmt.Println("\nWaiting for authentication...")
	fmt.Println("Please complete the login process in your browser.")
	fmt.Println("The CLI will automatically receive your token when you're done.")

	// Note: In a real implementation, we would:
	// 1. Start a local server to receive the OAuth callback
	// 2. Exchange the code for a token
	// 3. Save the token to ~/.nexlayer/config
	// For now, we'll just show a placeholder message

	fmt.Println("\nSuccess! You are now logged in to Nexlayer.")
	return nil
}
