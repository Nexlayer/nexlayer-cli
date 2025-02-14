// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package login

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

// NewLoginCommand creates a new login command
func NewLoginCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to Nexlayer",
		Long:  "Log in to your Nexlayer account to access deployment features",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("login flow not yet implemented")
		},
	}

	return cmd
}
