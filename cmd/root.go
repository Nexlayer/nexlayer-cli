package cmd

import (
	"fmt"
	"os"

	ci "github.com/Nexlayer/nexlayer-cli/pkg/commands/ci"
	deploy "github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	domain "github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	info "github.com/Nexlayer/nexlayer-cli/pkg/commands/info"
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/init"
	list "github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	login "github.com/Nexlayer/nexlayer-cli/pkg/commands/login"
	plugin "github.com/Nexlayer/nexlayer-cli/pkg/commands/plugin"
	scale "github.com/Nexlayer/nexlayer-cli/pkg/commands/scale"
	service "github.com/Nexlayer/nexlayer-cli/pkg/commands/service"
	status "github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	wizard "github.com/Nexlayer/nexlayer-cli/pkg/commands/wizard"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nexlayer",
	Short: "Nexlayer CLI",
	Long: `Nexlayer CLI is a command-line interface for managing your Nexlayer deployments.
It provides commands for:
- Service management (deploy, configure, scale)
- CI/CD pipeline management
- Infrastructure provisioning
- Monitoring and logging`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(service.ServiceCmd)
	rootCmd.AddCommand(ci.CICmd)
	rootCmd.AddCommand(deploy.DeployCmd)
	rootCmd.AddCommand(domain.DomainCmd)
	rootCmd.AddCommand(info.InfoCmd)
	rootCmd.AddCommand(initcmd.InitCmd)
	rootCmd.AddCommand(list.ListCmd)
	rootCmd.AddCommand(login.LoginCmd)
	rootCmd.AddCommand(plugin.PluginCmd)
	rootCmd.AddCommand(scale.ScaleCmd)
	rootCmd.AddCommand(status.StatusCmd)
	rootCmd.AddCommand(wizard.WizardCmd)
}
