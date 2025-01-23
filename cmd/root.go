package cmd

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/app"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nexlayer-cli",
	Short: "A CLI tool for managing Nexlayer deployments",
	Long: `Nexlayer CLI is a command-line tool for managing your Nexlayer deployments.
It provides commands for deploying applications, managing domains, and viewing deployment information.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Add commands
	rootCmd.AddCommand(
		login.Command,
	)
	rootCmd.AddCommand(
		deploy.Command,
	)
	rootCmd.AddCommand(
		domain.Command,
	)
	rootCmd.AddCommand(
		info.Command,
	)
	rootCmd.AddCommand(
		list.Command,
	)
	rootCmd.AddCommand(
		app.Cmd,
	)
}
