package ci

import (
	"github.com/spf13/cobra"
)

// CICmd represents the ci command
var CICmd = &cobra.Command{
	Use:   "ci",
	Short: "CI/CD pipeline management",
	Long: `Manage CI/CD pipelines for your Nexlayer deployments.
Examples:
  nexlayer ci setup github-actions
  nexlayer ci customize github-actions --build-context ./src --image-tag v1.0.0`,
}

func init() {
	// Add subcommands
	CICmd.AddCommand(customizeCmd)
	CICmd.AddCommand(setupCmd)
}
