package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/browser"
)

const (
	configDir = ".nexlayer"
	tokenFile = "token.json"
	sessionID = "R99BLvxPgvaXW" // Test session ID
)

type TokenConfig struct {
	Token string `json:"token"`
}

// GetToken returns the authentication token
func GetToken() (string, error) {
	// For testing, return the session ID
	return sessionID, nil
}

// GetAuthURL returns the authentication URL
func GetAuthURL() string {
	if os.Getenv("NEXLAYER_TEST_MODE") == "true" {
		return "https://app.staging.nexlayer.io/auth/github"
	}
	return "https://app.nexlayer.io/auth/github"
}

// Login opens the authentication URL in the browser
func Login() error {
	url := GetAuthURL()
	if err := browser.OpenURL(url); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}
	return nil
}

// SaveToken saves the authentication token
func SaveToken(token string) error {
	// Trim whitespace and validate
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	// Get config directory path
	configPath := os.Getenv("NEXLAYER_CONFIG_DIR")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, configDir)
	}

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	config := TokenConfig{
		Token: token,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal token config: %w", err)
	}

	tokenPath := filepath.Join(configPath, tokenFile)
	if err := os.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// LoadToken loads the saved authentication token
func LoadToken() (string, error) {
	// Get config directory path
	configPath := os.Getenv("NEXLAYER_CONFIG_DIR")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, configDir)
	}

	// Read token file
	tokenPath := filepath.Join(configPath, tokenFile)
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("not logged in")
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}

	// Parse token config
	var config TokenConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse token config: %w", err)
	}

	token := strings.TrimSpace(config.Token)
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	return token, nil
}
