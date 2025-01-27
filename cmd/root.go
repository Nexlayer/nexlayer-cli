package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/debug"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/deploy"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/domain"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/list"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/status"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/wizard"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "nexlayer",
	Short: "Nexlayer CLI - Deploy applications to Nexlayer",
	Long: `Nexlayer CLI helps you deploy and manage your applications on Nexlayer.
	
Key features:
- Easy application deployment
- Custom domain management
- Deployment status monitoring
- Deployment assistance

Need help? Use 'nexlayer debug' for deployment assistance.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nexlayer.yaml)")

	client := api.NewClient("")

	rootCmd.AddCommand(
		deploy.NewCommand(client),
		domain.NewCommand(client),
		list.NewCommand(client),
		status.NewCommand(client),
		wizard.NewCommand(client),
		debug.NewCommand(client),
	)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".nexlayer")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
