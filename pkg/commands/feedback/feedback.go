package feedback

import (
	"context"
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

// NewCommand creates a new feedback command
func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedback [text]",
		Short: "Send feedback to Nexlayer",
		Long: `Send feedback to help us improve Nexlayer.
Your feedback is valuable and will be used to enhance the service.

Example:
  nexlayer feedback "Great service! The deployment was super fast."`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFeedback(cmd.Context(), client, args[0])
		},
	}

	return cmd
}

func runFeedback(ctx context.Context, client *api.Client, text string) error {
	if err := client.SendFeedback(ctx, text); err != nil {
		return fmt.Errorf("failed to send feedback: %w", err)
	}

	fmt.Println("Thank you for your feedback!")
	return nil
}
