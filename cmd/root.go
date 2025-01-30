package cmd

import (
	"fmt"
	"os"

	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:          "nexlayer",
	Short:        "Nexlayer CLI",
	Long:         "Nexlayer CLI - Deploy AI applications with ease",
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add init command
	rootCmd.AddCommand(initcmd.NewCommand())
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
