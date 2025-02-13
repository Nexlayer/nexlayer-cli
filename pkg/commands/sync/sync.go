// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sync

import (
	"context"
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

// NewCommand creates a new sync command
func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync local configuration with Nexlayer",
		Long: `Synchronize your local nexlayer.yaml configuration with Nexlayer.
This ensures your local configuration matches the deployed state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync(cmd.Context(), client)
		},
	}
	return cmd
}

func runSync(ctx context.Context, client *api.Client) error {
	_ = ctx
	_ = client
	// TODO: Implement sync functionality
	return fmt.Errorf("sync command not yet implemented")
}
