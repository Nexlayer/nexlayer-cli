package service

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VisualizeCmd represents the visualize command
var VisualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize service dependencies",
	Long: `Visualize the connections between services in your application.
This command is currently not implemented.`,
	RunE: runVisualize,
}

func init() {
	VisualizeCmd.Flags().StringVar(&AppName, "app", "", "Application name")
	VisualizeCmd.MarkFlagRequired("app")
}

func runVisualize(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("service visualization is not yet implemented")
}
