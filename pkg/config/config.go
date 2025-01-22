package config

import (
	"os"
	"sync"
)

// Environment represents different deployment environments
type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
)

var (
	once     sync.Once
	instance *Config
)

// Config holds all configuration values
type Config struct {
	APIEndpoints map[Environment]string
}

// GetConfig returns a singleton instance of Config
func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			APIEndpoints: map[Environment]string{
				Production: getEnvOrDefault("NEXLAYER_API_URL", "https://app.nexlayer.io"),
				Staging:    getEnvOrDefault("NEXLAYER_STAGING_API_URL", "https://app.staging.nexlayer.io"),
			},
		}
	})
	return instance
}

// GetAPIEndpoint returns the appropriate API endpoint for the given environment
func (c *Config) GetAPIEndpoint(env Environment) string {
	if endpoint, ok := c.APIEndpoints[env]; ok {
		return endpoint
	}
	return c.APIEndpoints[Staging] // Default to staging if environment not found
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ValidateEnvironment checks if the provided environment string is valid
func ValidateEnvironment(env string) (Environment, bool) {
	switch Environment(env) {
	case Production, Staging:
		return Environment(env), true
	default:
		return "", false
	}
}
