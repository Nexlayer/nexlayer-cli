package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appName string

// Command represents the service command
var Command = &cobra.Command{
	Use:   "service [name]",
	Short: "Manage services",
	Long: `Manage services for your applications.
This command is currently not implemented.`,
	Args: cobra.ExactArgs(1),
	RunE: runService,
}

func init() {
	Command.Flags().StringVar(&appName, "app", "", "Application name")
	if err := Command.MarkFlagRequired("app"); err != nil {
		panic(fmt.Sprintf("failed to mark app flag as required: %v", err))
	}
}

func runService(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("service management is not yet implemented")
}
