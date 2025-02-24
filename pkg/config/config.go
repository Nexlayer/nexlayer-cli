// Package config provides centralized configuration management for the Nexlayer CLI.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"gopkg.in/yaml.v3"
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

var (
	instance *Config
	once     sync.Once
)

// Load loads the configuration from disk
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{}
		err = loadConfig()
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// loadConfig reads and parses the configuration file
func loadConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config if it doesn't exist
			instance = &Config{}
			return saveConfig()
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, instance); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	return nil
}

// Save saves the current configuration to disk
func Save() error {
	if instance == nil {
		return fmt.Errorf("configuration not loaded")
	}
	return saveConfig()
}

// saveConfig writes the configuration to disk
func saveConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(instance)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(configDir, ConfigFileName)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
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
	if instance != nil && instance.API.URL != "" {
		return instance.API.URL
	}
	return DefaultAPIURL
}

// GetToken returns the authentication token from environment or config
func GetToken() string {
	if token := os.Getenv("NEXLAYER_TOKEN"); token != "" {
		return token
	}
	if instance != nil {
		return instance.API.Token
	}
	return ""
}

// GetProjectNamespace returns the current project namespace
func GetProjectNamespace() string {
	if ns := os.Getenv("NEXLAYER_NAMESPACE"); ns != "" {
		return ns
	}
	if instance != nil && instance.Project.Namespace != "" {
		return instance.Project.Namespace
	}
	return "default"
}

// GetRegistryConfig returns the container registry configuration
func GetRegistryConfig() (string, string, string) {
	// Environment variables take precedence
	regType := os.Getenv("NEXLAYER_REGISTRY_TYPE")
	regURL := os.Getenv("NEXLAYER_REGISTRY_URL")
	regUser := os.Getenv("NEXLAYER_REGISTRY_USER")

	// Fall back to config file values
	if instance != nil {
		if regType == "" {
			regType = instance.Registry.Type
		}
		if regURL == "" {
			regURL = instance.Registry.URL
		}
		if regUser == "" {
			regUser = instance.Registry.Username
		}
	}

	return regType, regURL, regUser
}

// ValidateConfig checks if all required configuration is present
func ValidateConfig() error {
	if instance == nil {
		return fmt.Errorf("configuration not loaded")
	}

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
		return fmt.Errorf("missing required configuration: %v", missing)
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
