package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDir  = ".nexlayer"
	tokenFile  = "token.json"
	defaultToken = "default-token" // This will be replaced by actual authentication
)

type TokenConfig struct {
	Token string `json:"token"`
}

// GetToken returns the authentication token
func GetToken() (string, error) {
	// TODO: Implement proper token management
	// For now, return the default token
	return defaultToken, nil
}

// SaveToken saves the authentication token
func SaveToken(token string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, configDir)
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	tokenPath := filepath.Join(configPath, tokenFile)
	config := TokenConfig{Token: token}
	
	data, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal token config: %w", err)
	}

	if err := os.WriteFile(tokenPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}
