package ci

import (
	"github.com/spf13/cobra"
)

// CICmd represents the ci command
var CICmd = &cobra.Command{
	Use:   "ci",
	Short: "Manage CI/CD workflows and Docker images",
	Long: `Manage continuous integration and delivery workflows, including:
- GitHub Actions workflow generation
- Docker image management
- Build and deployment monitoring`,
}

func init() {
	// Add subcommands
	CICmd.AddCommand(setupCmd)
	CICmd.AddCommand(customizeCmd)
	CICmd.AddCommand(imagesCmd)
}
