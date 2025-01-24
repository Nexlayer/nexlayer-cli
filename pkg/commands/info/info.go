package info

import (
	"context"
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
	"github.com/spf13/cobra"
)

// NewInfoCmd creates a new info command
func NewInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [app]",
		Short: "Get information about an application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(vars.APIEndpoint)
			client.SetToken(vars.Token)

			app, err := client.GetAppInfo(context.Background(), args[0])
			if err != nil {
				return fmt.Errorf("failed to get application info: %w", err)
			}

			fmt.Printf("Application: %s\n", app.Name)
			fmt.Printf("ID: %s\n", app.ID)
			fmt.Printf("Status: %s\n", app.Status)
			fmt.Printf("Created: %s\n", app.CreatedAt.Format("2006-01-02 15:04:05"))

			return nil
		},
	}

	return cmd
}
