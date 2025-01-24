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
	var app string
	VisualizeCmd.Flags().StringVar(&app, "app", "", "Application name")
	if err := VisualizeCmd.MarkFlagRequired("app"); err != nil {
		panic(fmt.Sprintf("failed to mark app flag as required: %v", err))
	}
}

func runVisualize(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("service visualization is not yet implemented")
}
