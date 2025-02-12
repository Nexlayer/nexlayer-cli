// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/factory"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/feedback"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	syncCmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/sync"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/validate"
	versionCmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/version"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/watch"
	"github.com/Nexlayer/nexlayer-cli/pkg/version"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/errors"
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

	// Initialize the logger
	logger = observability.NewLogger(
		observability.INFO,
		observability.WithJSON(),
		observability.WithRotation(10, 5),
	)

	// Create command factory with middlewares
	cmdFactory := factory.NewCommandFactory(logger)
	cmdFactory.AddMiddleware(factory.RecoveryMiddleware(logger))
	cmdFactory.AddMiddleware(factory.ErrorHandlingMiddleware())
	cmdFactory.AddMiddleware(factory.LoggingMiddleware(logger))
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
		RunE: func(cmd *cobra.Command, args []string) error {
			showVersion, _ := cmd.Flags().GetBool("version")
			if showVersion {
				fmt.Printf("Nexlayer CLI version %s\n", version.GetVersion())
				return nil
			}
			return cmd.Help()
		},
	}

	// Add global flags
	cmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output errors in JSON format")
	cmd.Flags().Bool("version", false, "Print version information")

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
	cmd.AddCommand(validate.NewCommand())
	cmd.AddCommand(versionCmd.NewCommand())

	return cmd
}

// Execute runs the root command.
// It handles errors by reporting them and exiting the application.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Check if it's our custom error type
		if nexErr, ok := err.(*errors.Error); ok {
			// Log with stack trace for internal errors
			if nexErr.Type == errors.ErrorTypeInternal {
				logger.Error(context.Background(), "Internal error occurred: %s [file=%s, line=%d]", 
					nexErr.Error(), nexErr.File, nexErr.Line)
			}
			// For user errors, just show the message
			if nexErr.Type == errors.ErrorTypeUser {
				fmt.Fprintln(os.Stderr, nexErr.Message)
				os.Exit(1)
			}
		}
		reportError(err)
		os.Exit(1)
	}
}

// reportError handles error output with either JSON or structured logging.
// It formats errors based on the jsonOutput flag and logs them.
func reportError(err error) {
	if jsonOutput {
		// For structured errors, include all available context
		if nexErr, ok := err.(*errors.Error); ok {
			output := map[string]interface{}{
				"error":   nexErr.Message,
				"type":    nexErr.Type.String(),
				"file":    nexErr.File,
				"line":    nexErr.Line,
				"context": map[string]interface{}{},
			}
			if nexErr.Cause != nil {
				output["cause"] = nexErr.Cause.Error()
			}
			if len(nexErr.Stack) > 0 {
				output["stack"] = nexErr.Stack
			}
			json.NewEncoder(os.Stderr).Encode(output)
			return
		}
		// Fall back to simple error message for non-structured errors
		json.NewEncoder(os.Stderr).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	if logger != nil {
		// Use structured logging for our custom error type
		if nexErr, ok := err.(*errors.Error); ok {
			fields := map[string]interface{}{
				"type": nexErr.Type.String(),
				"file": nexErr.File,
				"line": nexErr.Line,
			}
			if nexErr.Cause != nil {
				fields["cause"] = nexErr.Cause.Error()
			}
			logger.Error(context.Background(), "Command error: %s [type=%s, file=%s, line=%d]", 
				nexErr.Message, nexErr.Type.String(), nexErr.File, nexErr.Line)
			return
		}
		// Fall back to simple error logging
		logger.Error(context.Background(), "Command error: %v", err)
		return
	}

	fmt.Fprintln(os.Stderr, err)
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
