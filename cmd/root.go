package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/compose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// logger is the global structured logger instance
	logger *zap.Logger
	// configOnce ensures thread-safe lazy loading of config
	configOnce sync.Once
	// rootCmd is the primary cobra command
	rootCmd = NewRootCommand()
	// jsonOutput toggles JSON-formatted error output
	jsonOutput bool
)

// initLogger creates a production-ready zap logger
func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
}

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy AI applications with ease",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Initialize logger once
			if logger == nil {
				initLogger()
			}

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

	cmd.AddCommand(initcmd.NewCommand())
	cmd.AddCommand(ai.NewCommand())
	cmd.AddCommand(compose.NewCommand())

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
		logger.Error("Command execution error", zap.Error(err))
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
				logger.Error("Error reading config file", zap.Error(err))
			} else {
				logger.Info("No config file found; using defaults")
			}
		} else {
			logger.Info("Configuration loaded", 
				zap.String("file", viper.ConfigFileUsed()))
		}
	})
}
