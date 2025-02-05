package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/compose"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/feedback"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/observability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// logger is the global structured logger instance
	logger *observability.Logger
	// configOnce ensures thread-safe lazy loading of config
	configOnce sync.Once
	// rootCmd is the primary cobra command
	rootCmd *cobra.Command
	// jsonOutput toggles JSON-formatted error output
	jsonOutput bool
)

func init() {
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
	cmd.AddCommand(ai.NewAICommand())
	cmd.AddCommand(compose.NewCommand())
	cmd.AddCommand(deploy.NewCommand(apiClient))
	cmd.AddCommand(domain.NewCommand(apiClient))
	cmd.AddCommand(feedback.NewCommand(apiClient))
	cmd.AddCommand(list.NewCommand(apiClient))
	cmd.AddCommand(status.NewCommand(apiClient))

	return cmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		reportError(err)
		os.Exit(1)
	}
}

// reportError handles error output with either JSON or structured logging.
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
