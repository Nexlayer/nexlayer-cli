// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/feedback"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	syncCmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/sync"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/watch"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Package cmd provides the command-line interface for the Nexlayer CLI.

var (
	// logger is the global structured logger instance
	// It is used for logging messages with different severity levels.
	logger *observability.Logger
	// configOnce ensures thread-safe lazy loading of config
	// It is used to initialize configuration only once.
	configOnce sync.Once
	// rootCmd is the primary cobra command
	// It serves as the entry point for all CLI commands.
	rootCmd *cobra.Command
	// jsonOutput toggles JSON-formatted error output
	// It is used to output errors in a structured JSON format.
	jsonOutput bool
)

// init initializes the logger and sets default configuration values.
// It creates the root command and registers all subcommands.
func init() {
	// Enable colors for Windows
	os.Setenv("TERM", "xterm-256color")

	// Initialize the logger first with JSON mode and rotation settings.
	logger = observability.NewLogger(
		observability.INFO,
		observability.WithJSON(),
		observability.WithRotation(50, 7), // 50MB max size, 7 days retention
	)

	// Set default configuration values.
	viper.SetDefault("nexlayer.api_url", "https://app.staging.nexlayer.io")

	// Create the root command.
	rootCmd = NewRootCommand()
}

// NewRootCommand creates and returns the root command for the CLI.
// It sets up the API client, adds global flags, and registers subcommands.
// Returns the configured cobra.Command instance.
func NewRootCommand() *cobra.Command {
	// Retrieve API URL from configuration (allows override via config/env)
	apiURL := viper.GetString("nexlayer.api_url")
	apiClient := api.NewClient(apiURL)

	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy AI applications with ease",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Load configuration only when needed.
			if cmd.Name() != "help" && cmd.Name() != "version" {
				lazyInitConfig()
			}

			// Set a background context.
			cmd.SetContext(context.Background())
		},
	}

	// Add a global JSON flag for structured error output.
	cmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output errors in JSON format")

	// Register all commands.
	cmd.AddCommand(initcmd.NewCommand())
	cmd.AddCommand(ai.NewCommand())

	cmd.AddCommand(deploy.NewCommand(apiClient))
	cmd.AddCommand(domain.NewCommand(apiClient))
	cmd.AddCommand(feedback.NewCommand(apiClient))
	cmd.AddCommand(list.NewCommand(apiClient))
	cmd.AddCommand(syncCmd.NewCommand())
	cmd.AddCommand(status.NewCommand(apiClient))
	cmd.AddCommand(watch.NewWatchCommand())

	return cmd
}

// Execute runs the root command.
// It handles errors by reporting them and exiting the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		reportError(err)
		os.Exit(1)
	}
}

// reportError handles error output with either JSON or structured logging.
// It formats errors based on the jsonOutput flag and logs them.
func reportError(err error) {
	if jsonOutput {
		jsonErr := map[string]interface{}{
			"error_context": map[string]interface{}{
				"type":    "CommandError",
				"message": err.Error(),
			},
		}
		if jsonBytes, jerr := json.Marshal(jsonErr); jerr == nil {
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println(err)
		}
	} else {
		// Log the error using structured logging and print the error.
		logger.Error(context.Background(), "Command execution error: %v", err)
		fmt.Println(err)
	}
}

// lazyInitConfig loads configuration files and environment variables.
// It searches for config files in predefined locations and enables env var overrides.
func lazyInitConfig() {
	configOnce.Do(func() {
		// Search for configuration files in multiple locations.
		viper.AddConfigPath("$HOME/.config/nexlayer")
		viper.AddConfigPath(".") // Also look in the current directory.
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// Enable environment variable overrides.
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				logger.Error(context.Background(), "Error reading config file: %v", err)
			} else {
				logger.Info(context.Background(), "No config file found; using defaults")
			}
		} else {
			logger.Info(context.Background(), "Configuration loaded from %s", viper.ConfigFileUsed())
		}
	})
}
