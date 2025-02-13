// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package deployment

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/spf13/cobra"
)

// NewCommand creates a command group that wraps the deployment-related endpoints:
// - GET /getDeployments
// - GET /getDeploymentInfo
func NewCommand(apiClient *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "Manage deployments",
		Long:  "Commands for managing and viewing deployment information using Nexlayer's REST API",
	}

	// Add global flags
	cmd.PersistentFlags().Bool("json", false, "Output in JSON format")

	// Add subcommands that map to REST endpoints
	cmd.AddCommand(
		newListCommand(apiClient),  // GET /getDeployments
		newInfoCommand(apiClient),  // GET /getDeploymentInfo
	)

	return cmd
}
