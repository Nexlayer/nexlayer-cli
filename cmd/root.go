package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/app"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
	"github.com/spf13/cobra"
)

var (
	configFile string
	verbose    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nexlayer-cli",
	Short: "Nexlayer CLI tool",
	Long: `Nexlayer CLI is a command-line tool for managing your Nexlayer applications.
It provides commands for deploying applications, managing domains, and viewing deployment information.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Initialize config file path
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		os.Exit(1)
	}
	configFile = filepath.Join(home, ".nexlayer", "config")

	// Add commands
	rootCmd.AddCommand(login.Command)
	rootCmd.AddCommand(deploy.Command)
	rootCmd.AddCommand(domain.Command)
	rootCmd.AddCommand(info.NewInfoCmd())
	rootCmd.AddCommand(list.Command)
	rootCmd.AddCommand(app.Cmd)

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", configFile, "config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}
