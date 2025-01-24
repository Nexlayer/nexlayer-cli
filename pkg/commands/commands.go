package commands

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/app"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	initCmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/init"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/registry"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/service"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/wizard"
	"github.com/spf13/cobra"
)

// RegisterCommands registers all commands with the root command
func RegisterCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		app.NewCommand(),
		deploy.NewDeployCmd(),
		domain.NewCommand(),
		info.NewInfoCmd(),
		initCmd.InitCmd,
		list.NewListCmd(),
		login.NewCommand(),
		registry.NewRegistryCmd(),
		service.ServiceCmd,
		wizard.NewWizardCmd(),
	)
}
