package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/compose"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = NewRootCommand()

// Global flag for JSON output
var jsonOutput bool

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nexlayer",
		Short: "Nexlayer CLI - Deploy AI applications with ease",
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
		if jsonOutput {
			// Convert error to JSON format
			jsonErr := map[string]interface{}{
				"error_context": map[string]interface{}{
					"type":    "CommandError",
					"message": err.Error(),
				},
			}
			if jsonBytes, err := json.Marshal(jsonErr); err == nil {
				fmt.Println(string(jsonBytes))
			}
		} else {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Find config file
	viper.AddConfigPath("$HOME/.config/nexlayer")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Error reading config file: %s\n", err)
		}
	}
}
