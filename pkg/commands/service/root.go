package service

import (
	"github.com/spf13/cobra"
)

// ServiceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Long: `Manage services in your Nexlayer application.
This command is currently not implemented.`,
}

func init() {
	// Add subcommands
	ServiceCmd.AddCommand(
		ConfigureCmd,
		DeployCmd,
		VisualizeCmd,
	)
}
