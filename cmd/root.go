package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
)

var rootCmd = &cobra.Command{
	Use:   "nexlayer",
	Short: "Nexlayer CLI - Deploy and manage your applications",
	Long: `Nexlayer CLI is a command-line tool for deploying and managing your applications.
It provides a simple interface to interact with the Nexlayer platform.

Example usage:
  nexlayer deploy hello-world    # Deploy a template
  nexlayer list                  # List all deployments
  nexlayer info my-namespace     # Get deployment information
  nexlayer domain mydomain.com   # Set a custom domain`,
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
	rootCmd.AddCommand(commands.InitCmd)
	rootCmd.AddCommand(commands.DeployCmd)
	rootCmd.AddCommand(commands.ListCmd)
	rootCmd.AddCommand(commands.StatusCmd)
	rootCmd.AddCommand(commands.DomainCmd)
	rootCmd.AddCommand(commands.InfoCmd)
	rootCmd.AddCommand(commands.ScaleCmd)
	rootCmd.AddCommand(commands.WizardCmd)
}
