package version

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/version"
	"github.com/spf13/cobra"
)

// NewCommand creates and returns the version command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of Nexlayer CLI",
		Long:  `Display the version and build information for your Nexlayer CLI installation.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Nexlayer CLI version %s\n", version.GetVersion())
		},
	}

	return cmd
}
