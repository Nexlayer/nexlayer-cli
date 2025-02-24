// Package vars provides centralized configuration management for the Nexlayer CLI.
// Deprecated: Use pkg/config instead.
package vars

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Configuration constants
const (
	// DefaultAPIURL is the default Nexlayer API endpoint
	DefaultAPIURL = "https://api.nexlayer.dev"

	// ConfigFileName is the name of the configuration file
	ConfigFileName = "config.yaml"

	// CacheDir is the directory for caching CLI data
	CacheDir = ".nexlayer"

	// DefaultPort is the default port for local development
	DefaultPort = 3000
)

// Config holds the CLI configuration
type Config struct {
	// API configuration
	API struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token"`
	} `yaml:"api"`

	// Project configuration
	Project struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
		Domain    string `yaml:"domain"`
	} `yaml:"project"`

	// Registry configuration
	Registry struct {
		Type     string `yaml:"type"`     // Container registry type (ghcr, dockerhub, gcr, ecr, artifactory, gitlab)
		URL      string `yaml:"url"`      // Registry URL
		Username string `yaml:"username"` // Registry username
		Region   string `yaml:"region"`   // Registry region (for ECR)
		Project  string `yaml:"project"`  // Registry project ID (for GCR)
	} `yaml:"registry"`

	// Build configuration
	Build struct {
		Context string `yaml:"context"` // Docker build context path
		Tag     string `yaml:"tag"`     // Docker image tag
	} `yaml:"build"`

	// Environment variables
	Env map[string]string `yaml:"env"`
}

// GetConfigDir returns the platform-specific configuration directory
func GetConfigDir() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	configDir := filepath.Join(userConfigDir, "nexlayer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// GetCacheDir returns the platform-specific cache directory
func GetCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user cache directory: %w", err)
	}

	cacheDir := filepath.Join(userCacheDir, "nexlayer")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cacheDir, nil
}

// GetDefaultShell returns the default shell for the current platform
func GetDefaultShell() string {
	if runtime.GOOS == "windows" {
		return "cmd.exe"
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	return shell
}

// GetAPIURL returns the configured API URL or the default
func GetAPIURL() string {
	if url := os.Getenv("NEXLAYER_API_URL"); url != "" {
		return url
	}
	return DefaultAPIURL
}

// GetToken returns the authentication token from environment or config
func GetToken() string {
	return os.Getenv("NEXLAYER_TOKEN")
}

// GetProjectNamespace returns the current project namespace
func GetProjectNamespace() string {
	if ns := os.Getenv("NEXLAYER_NAMESPACE"); ns != "" {
		return ns
	}
	return "default"
}

// GetRegistryConfig returns the container registry configuration
func GetRegistryConfig() (string, string, string) {
	regType := os.Getenv("NEXLAYER_REGISTRY_TYPE")
	regURL := os.Getenv("NEXLAYER_REGISTRY_URL")
	regUser := os.Getenv("NEXLAYER_REGISTRY_USER")
	return regType, regURL, regUser
}

// ValidateConfig checks if all required configuration is present
func ValidateConfig() error {
	missing := make([]string, 0)

	// Check required environment variables
	required := []struct {
		name  string
		value string
	}{
		{"NEXLAYER_TOKEN", GetToken()},
	}

	for _, req := range required {
		if req.value == "" {
			missing = append(missing, req.name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	return nil
}

// GetLogLevel returns the configured log level
func GetLogLevel() string {
	if level := os.Getenv("NEXLAYER_LOG_LEVEL"); level != "" {
		return level
	}
	return "info"
}

// IsDevelopment returns true if running in development mode
func IsDevelopment() bool {
	return os.Getenv("NEXLAYER_ENV") == "development"
}

// IsDebug returns true if debug mode is enabled
func IsDebug() bool {
	return os.Getenv("NEXLAYER_DEBUG") == "true"
}

// GetUserAgent returns the CLI user agent string
func GetUserAgent() string {
	return fmt.Sprintf("nexlayer-cli/%s (%s; %s)",
		os.Getenv("NEXLAYER_VERSION"),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
