package service

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	appName  string
	service  string
	envPairs []string
)

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure service settings",
	Long: `Configure settings for a service in your Nexlayer deployment.
Examples:
  nexlayer service configure --app my-app --service frontend --env API_URL=https://api.example.com
  nexlayer service configure --app my-app --service backend --env "DB_URL=postgres://localhost:5432/db" --env "REDIS_URL=redis://localhost:6379"`,
	RunE: runConfigure,
}

func init() {
	configureCmd.Flags().StringVar(&appName, "app", "", "Application name")
	configureCmd.Flags().StringVar(&service, "service", "", "Service name (e.g., frontend, backend)")
	configureCmd.Flags().StringArrayVar(&envPairs, "env", []string{}, "Environment variables in KEY=VALUE format")

	configureCmd.MarkFlagRequired("app")
	configureCmd.MarkFlagRequired("service")
}

func runConfigure(cmd *cobra.Command, args []string) error {
	// Parse environment variables
	envVars := make(map[string]string)
	for _, pair := range envPairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid environment variable format: %s (should be KEY=VALUE)", pair)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		envVars[key] = value
	}

	// Get auth token
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create API client
	client := api.NewClient("https://app.nexlayer.io")

	// Update service configuration
	err := client.UpdateServiceConfig(appName, service, envVars, token)
	if err != nil {
		return fmt.Errorf("failed to update service configuration: %w", err)
	}

	fmt.Printf("✅ Successfully updated configuration for %s service in %s\n", service, appName)
	for k, v := range envVars {
		maskedValue := maskSensitiveValue(v)
		fmt.Printf("  %s=%s\n", k, maskedValue)
	}

	return nil
}

// maskSensitiveValue masks potentially sensitive values
func maskSensitiveValue(value string) string {
	sensitiveKeys := []string{
		"password", "secret", "key", "token", "credential",
		"auth", "pwd", "pass",
	}

	// Check if the value is a URL
	if strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		return value
	}

	// Check if the value might be sensitive
	for _, key := range sensitiveKeys {
		if strings.Contains(strings.ToLower(value), key) {
			return "********"
		}
	}

	return value
}
