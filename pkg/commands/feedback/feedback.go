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
		Short: "Send feedback",
		Long:  "Send feedback about your experience with Nexlayer",
	}

	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "Send feedback",
		Long: `Send feedback about Nexlayer.

Endpoint: POST /feedback

Arguments:
  --message        Your feedback message

Example:
  nexlayer feedback send --message "Great product!"`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			msg, _ := cmd.Flags().GetString("message")
			return runFeedback(cmd, cmd.Context(), client, msg)
		},
	}

	sendCmd.Flags().String("message", "", "Feedback message (required)")
	sendCmd.MarkFlagRequired("message")

	cmd.AddCommand(sendCmd)
	return cmd
}

func runFeedback(cmd *cobra.Command, ctx context.Context, client api.APIClient, text string) error {
	if err := client.SendFeedback(ctx, text); err != nil {
		return fmt.Errorf("failed to send feedback: %w", err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Thank you for your feedback!")
	return nil
}
