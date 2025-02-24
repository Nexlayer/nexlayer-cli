// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package feedback

import (
	"context"
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

// NewFeedbackCommand creates a new feedback command
func NewFeedbackCommand(client api.ClientAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feedback",
		Short: "Send feedback about your Nexlayer experience",
		Long: `Send feedback about your experience with the Nexlayer platform.
Your feedback helps us improve the platform and build better features.

Examples:
  • Report bugs or issues
  • Suggest new features
  • Share your success stories
  • Request improvements`,
	}

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send feedback to Nexlayer team",
		Long: `Send feedback about your experience with the Nexlayer platform.

The feedback will be sent directly to our development team and product managers.
We read every piece of feedback and use it to improve the platform.

Examples:
  nexlayer feedback send --message "Love the new AI deployment features!"
  nexlayer feedback send --message "Would like to see support for custom domains"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, _ := cmd.Flags().GetString("message")
			return runFeedback(cmd, cmd.Context(), client, msg)
		},
	}

	sendCmd.Flags().String("message", "", "Your feedback message (required)")
	sendCmd.MarkFlagRequired("message")

	cmd.AddCommand(sendCmd)
	return cmd
}

func runFeedback(cmd *cobra.Command, ctx context.Context, client api.APIClient, text string) error {
	fmt.Fprintln(cmd.OutOrStdout(), "📝 Sending feedback to Nexlayer team...")

	if err := client.SendFeedback(ctx, text); err != nil {
		return fmt.Errorf("failed to send feedback: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "\n✨ Thank you for your feedback!")
	fmt.Fprintln(cmd.OutOrStdout(), "Your input helps us improve the Nexlayer platform.")
	return nil
}
