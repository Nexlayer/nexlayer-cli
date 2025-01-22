package config

// Formatted with gofmt -s
import (
	"os"
	"sync"
)

// Environment represents different deployment environments
type Environment string

const (
	Production Environment = "production"
	Staging    Environment = "staging"
	// Default API endpoints
	productionAPI = "https://app.nexlayer.io"
	stagingAPI    = "https://app.staging.nexlayer.io"
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
				Production: productionAPI,
				Staging:    stagingAPI,
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

// ValidateEnvironment checks if the provided environment string is valid
func ValidateEnvironment(env string) (Environment, bool) {
	switch Environment(env) {
	case Production, Staging:
		return Environment(env), true
	default:
		return "", false
	}
}

// GetAuthToken returns the authentication token or an empty string if not set
func GetAuthToken() string {
	return os.Getenv("NEXLAYER_AUTH_TOKEN")
}
