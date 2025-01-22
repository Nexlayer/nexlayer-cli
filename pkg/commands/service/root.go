package service

import (
	"github.com/spf13/cobra"
)

// ServiceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage services",
	Long: `Manage services in your Nexlayer deployment.
Includes commands for:
- Deploying services
- Configuring services
- Visualizing service connections
- Managing service dependencies`,
}

func init() {
	ServiceCmd.AddCommand(deployCmd)
	ServiceCmd.AddCommand(configureCmd)
	ServiceCmd.AddCommand(visualizeCmd)
}
