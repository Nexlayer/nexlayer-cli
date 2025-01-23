package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serviceName string
)

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a service",
	Long: `Deploy a service to your application.
This command is currently not implemented.`,
	RunE: runDeploy,
}

func init() {
	DeployCmd.Flags().StringVar(&serviceName, "service", "", "Service name")
	DeployCmd.MarkFlagRequired("service")
}

func runDeploy(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("service deployment is not yet implemented")
}
