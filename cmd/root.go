package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/compose"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
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
	// Initialize logger first with JSON mode and rotation
	logger = observability.NewLogger(
		observability.INFO,
		observability.WithJSON(),
		observability.WithRotation(50, 7), // 50MB max size, 7 days retention
	)
	// Then create root command
	rootCmd = NewRootCommand()
}

func NewRootCommand() *cobra.Command {
	// Create API client
	apiClient := api.NewClient("https://api.nexlayer.dev")

	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy AI applications with ease",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Only load config if needed
			if cmd.Name() != "help" && cmd.Name() != "version" {
				lazyInitConfig()
			}

			// Add timeout context for long-running commands
			ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
			defer cancel()
			cmd.SetContext(ctx)
		},
	}

	// Add global JSON flag
	cmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output errors in JSON format")

	// Add commands
	cmd.AddCommand(initcmd.NewCommand())
	cmd.AddCommand(ai.NewAICommand())
	cmd.AddCommand(compose.NewCommand())
	cmd.AddCommand(deploy.NewCommand(apiClient))

	return cmd
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		reportError(err)
		os.Exit(1)
	}
}

// reportError handles error output with either JSON or structured logging
func reportError(err error) {
	if jsonOutput {
		jsonErr := map[string]interface{}{
			"error_context": map[string]interface{}{
				"type":    "CommandError",
				"message": err.Error(),
			},
		}
		if jsonBytes, jsonErr := json.Marshal(jsonErr); jsonErr == nil {
			fmt.Println(string(jsonBytes))
		} else {
			fmt.Println(err)
		}
	} else {
		// Use structured logging with stack trace
		logger.Error(context.Background(), "Command execution error: %v", err)
		fmt.Println(err)
	}
}

func lazyInitConfig() {
	configOnce.Do(func() {
		// Search config in multiple locations
		viper.AddConfigPath("$HOME/.config/nexlayer")
		viper.AddConfigPath(".")  // Also look in current directory
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")

		// Enable environment variable override
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
