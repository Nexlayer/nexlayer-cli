package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

const (
	authServerPort = 3000
	callbackPath   = "/auth/callback"
)

// LoginCmd represents the login command
var LoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Nexlayer",
	Long:  `Authenticate with Nexlayer using your GitHub account`,
	RunE:  runLogin,
}

func runLogin(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Opening browser to login to Nexlayer...")
	
	// Start local server to receive callback
	tokenChan := make(chan string)
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", authServerPort),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == callbackPath {
				token := r.URL.Query().Get("token")
				if token != "" {
					tokenChan <- token
					fmt.Fprintf(w, "<h1>Successfully logged in!</h1><p>You can close this window and return to your terminal.</p>")
				}
			}
		}),
	}

	// Start server in background
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	// Open browser to start auth flow
	authURL := fmt.Sprintf("https://app.nexlayer.io/cli-auth?port=%d", authServerPort)
	if err := browser.OpenURL(authURL); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	// Wait for token with timeout
	select {
	case token := <-tokenChan:
		// Store token securely
		if err := saveToken(token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
		fmt.Println("✅ Successfully logged in to Nexlayer!")
		return nil
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("login timed out")
	}
}

func saveToken(token string) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	nexlayerDir := fmt.Sprintf("%s/nexlayer", configDir)
	if err := os.MkdirAll(nexlayerDir, 0700); err != nil {
		return err
	}

	config := struct {
		Token     string    `json:"token"`
		CreatedAt time.Time `json:"created_at"`
	}{
		Token:     token,
		CreatedAt: time.Now(),
	}

	configFile := fmt.Sprintf("%s/config.json", nexlayerDir)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return err
	}

	// Also set environment variable for immediate use
	return os.Setenv("NEXLAYER_AUTH_TOKEN", token)
}
