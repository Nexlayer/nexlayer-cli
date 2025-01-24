package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ConfigureCmd represents the configure command
var ConfigureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure a service",
	Long: `Configure a service in your application.
This command is currently not implemented.`,
	RunE: runConfigure,
}

func init() {
	ConfigureCmd.Flags().StringVar(&applicationID, "app", "", "Application ID")
	if err := ConfigureCmd.MarkFlagRequired("app"); err != nil {
		panic(fmt.Sprintf("failed to mark app flag as required: %v", err))
	}
}

func runConfigure(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("service configuration is not yet implemented")
}
