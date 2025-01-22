// Formatted with gofmt -s
package cmd

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
	"github.com/Nexlayer/nexlayer-cli/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "nexlayer",
	Short:   "Nexlayer CLI - Deploy and manage full-stack applications",
	Version: version.Version,
	Long: `Nexlayer CLI helps you deploy and manage full-stack applications with ease.
Built for developers who value simplicity without sacrificing power.

Find more information at: https://docs.nexlayer.io`,
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
