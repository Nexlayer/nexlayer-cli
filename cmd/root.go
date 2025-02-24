// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/feedback"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/watch"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/errors"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Package cmd provides the command-line interface for the Nexlayer CLI.

var (
	// logger is the global structured logger instance for logging messages.
	logger *observability.Logger
	// configOnce ensures thread-safe lazy loading of configuration.
	configOnce sync.Once
	// rootCmd is the primary cobra command, the entry point for all CLI commands.
	rootCmd *cobra.Command
	// jsonOutput toggles JSON-formatted output for errors and responses.
	jsonOutput bool
)

// init initializes the logger, sets default config values, and creates the root command.
func init() {
	// Enable colors for Windows terminals.
	os.Setenv("TERM", "xterm-256color")

	// Initialize the logger with rotation settings.
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

// NewRootCommand creates and configures the root command for the CLI.
func NewRootCommand() *cobra.Command {
	// Retrieve API URL from configuration (overridable via config/env).
	apiURL := viper.GetString("nexlayer.api_url")
	apiClient := api.NewClient(apiURL)

	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy applications with ease",
		Long:  `Nexlayer CLI – Deploy Full-Stack Applications in Seconds ⚡️`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Load configuration only when needed.
			if cmd.Name() != "help" {
				lazyInitConfig()
			}

			// Set a background context.
			cmd.SetContext(context.Background())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return cmd.Help()
		},
	}

	// Add global flags
	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output response in JSON format")

	// Disable auto-generation of completion command
	cmd.CompletionOptions.DisableDefaultCmd = true

	// Register commands in desired order
	cmd.AddCommand(
		initcmd.NewCommand(),
		deploy.NewCommand(apiClient),
		list.NewListCommand(apiClient),
		info.NewInfoCommand(apiClient),
		domain.NewDomainCommand(apiClient),
		login.NewLoginCommand(apiClient),
		watch.NewCommand(apiClient),
		feedback.NewFeedbackCommand(apiClient),
	)

	// Disable suggestions and help command
	cmd.DisableSuggestions = true
	cmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// Set custom help template to control command order
	cmd.SetUsageTemplate(`Core Commands:
  init        Initialize a new project (auto-detects type)
  deploy      Deploy an application (uses nexlayer.yaml if present)
  list        List active deployments
  info        Get deployment details <namespace> <appID>
  domain      Manage custom domains
  login       Authenticate with Nexlayer
  watch       Monitor & auto-deploy changes
  feedback    Send CLI feedback

Flags:
  -h, --help         Show help for commands
      --preview      (Future) Show changes without applying them

Global Flags:
  --json          Output response in JSON format

For more details:
  {{.CommandPath}} [command] --help
`)

	return cmd
}

// Execute runs the root command and handles errors gracefully.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Handle custom errors with context and suggestions.
		if nexErr, ok := err.(*errors.Error); ok {
			if nexErr.Type == errors.ErrorTypeInternal {
				logger.Error(rootCmd.Context(), "Internal error in '%s': %s [file=%s, line=%d]",
					rootCmd.Name(), nexErr.Error(), nexErr.File, nexErr.Line)
				fmt.Fprintf(os.Stderr, "Internal error: %s\n", nexErr.Message)
				if nexErr.Cause != nil {
					fmt.Fprintf(os.Stderr, "Caused by: %v\n", nexErr.Cause)
				}
			} else if nexErr.Type == errors.ErrorTypeUser {
				fmt.Fprintf(os.Stderr, "Error: %s\n", nexErr.Message)
				if nexErr.Cause != nil {
					fmt.Fprintf(os.Stderr, "Details: %v\n", nexErr.Cause)
				}
			}
			os.Exit(1)
		}
		// Fallback for non-custom errors.
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// lazyInitConfig loads configuration files and environment variables.
func lazyInitConfig() {
	configOnce.Do(func() {
		// Use custom config path if provided via flag.
		if customConfig := viper.GetString("config"); customConfig != "" {
			viper.SetConfigFile(customConfig)
		} else {
			viper.AddConfigPath("$HOME/.config/nexlayer")
			viper.AddConfigPath(".") // Current directory fallback.
			viper.SetConfigName("config")
			viper.SetConfigType("yaml")
		}

		// Enable environment variable overrides.
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				logger.Info(context.Background(), "No config file found; using defaults")
			} else {
				logger.Error(context.Background(), "Error reading config file: %v", err)
				fmt.Fprintln(os.Stderr, "Configuration file is invalid. Please check the syntax and try again.")
			}
		} else {
			logger.Info(context.Background(), "Configuration loaded from %s", viper.ConfigFileUsed())
		}
	})
}

// CommandDependencies holds dependencies for commands.
type CommandDependencies struct {
	APIClient *api.Client
	Logger    *observability.Logger
}

// CommandRegistry manages command registration for scalability.
type CommandRegistry struct {
	commands []*cobra.Command
}

// Register adds a command to the registry.
func (r *CommandRegistry) Register(cmd *cobra.Command) {
	r.commands = append(r.commands, cmd)
}

// AddToRoot attaches all registered commands to the root command.
func (r *CommandRegistry) AddToRoot(root *cobra.Command) {
	for _, cmd := range r.commands {
		root.AddCommand(cmd)
	}
}
