package service

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Visualize service connections",
	Long: `Generate a visual diagram of service connections in your application.
Currently supports ASCII output to terminal.`,
	RunE: runVisualize,
}

func init() {
	visualizeCmd.Flags().StringVar(&vars.AppName, "app", "", "Application name")
	visualizeCmd.MarkFlagRequired("app")
}

func runVisualize(cmd *cobra.Command, args []string) error {
	// Get auth token
	token := os.Getenv("NEXLAYER_AUTH_TOKEN")
	if token == "" {
		return fmt.Errorf("NEXLAYER_AUTH_TOKEN environment variable is not set")
	}

	// Create API client
	client := api.NewClient("https://app.nexlayer.io")

	// Get service connections
	connections, err := client.GetServiceConnections(vars.AppName, token)
	if err != nil {
		return fmt.Errorf("failed to get service connections: %w", err)
	}

	return visualizeAscii(connections)
}

func visualizeAscii(connections []api.ServiceConnection) error {
	if len(connections) == 0 {
		fmt.Println("No service connections found")
		return nil
	}

	fmt.Println("Service Connections:")
	fmt.Println("-------------------")
	for _, conn := range connections {
		fmt.Printf("%s --> %s: %s\n", conn.From, conn.To, conn.Description)
	}
	return nil
}
