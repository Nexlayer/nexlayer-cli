package login

import (
	"fmt"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/auth"
	"github.com/spf13/cobra"
)

var token string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to Nexlayer",
		Long:  "Login to Nexlayer using your API token",
		RunE:  runLogin,
	}

	cmd.Flags().StringVarP(&token, "token", "t", "", "API token")
	return cmd
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Trim whitespace and validate
	token = strings.TrimSpace(token)
	if token == "" {
		return fmt.Errorf("token is required")
	}

	if err := auth.SaveToken(token); err != nil {
		return fmt.Errorf("failed to save token: %w", err)
	}

	cmd.Printf("Successfully logged in\n")
	return nil
}
