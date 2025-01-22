package commands

// Formatted with gofmt -s
import (
	"github.com/Nexlayer/nexlayer-cli/pkg/tui"
	"github.com/spf13/cobra"
)

var WizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Start the interactive deployment wizard",
	Long:  "Start an interactive wizard to guide you through deploying your application",
	RunE:  runWizard,
}

func runWizard(cmd *cobra.Command, args []string) error {
	wizard := tui.NewDeploymentWizard()
	return wizard.Run()
}
