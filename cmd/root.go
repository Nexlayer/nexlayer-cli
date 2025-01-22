package cmd

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nexlayer",
	Short: "Nexlayer CLI - Deploy your applications with ease",
	Long: `Nexlayer CLI is a tool for deploying and managing your applications.
It provides a simple interface for deploying applications to Nexlayer's platform.`,
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
	// Add subcommands
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.DeployCmd)
	rootCmd.AddCommand(commands.LoginCmd)
	rootCmd.AddCommand(commands.StatusCmd)
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.InfoCmd)
	rootCmd.AddCommand(commands.DomainCmd)
	rootCmd.AddCommand(commands.ScaleCmd)
	rootCmd.AddCommand(commands.PluginCmd)
	rootCmd.AddCommand(commands.AISuggestCmd)
	rootCmd.AddCommand(commands.CICmd)
}
