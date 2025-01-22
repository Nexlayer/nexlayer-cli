package service

import (
	"github.com/spf13/cobra"
)

// ServiceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage service configurations",
	Long: `Configure and manage services in your Nexlayer deployments.
Includes commands for:
- Configuring service settings
- Managing environment variables
- Visualizing service connections`,
}

func init() {
	// Add subcommands
	ServiceCmd.AddCommand(configureCmd)
	ServiceCmd.AddCommand(visualizeCmd)
}
