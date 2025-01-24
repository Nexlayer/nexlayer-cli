package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands"
)

var (
	cfgFile string
	verbose bool
	useAI   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nexlayer",
	Short: "Nexlayer CLI - Deploy and manage your applications",
	Long: `Nexlayer CLI helps you deploy and manage your applications.
	
Use the wizard command to get started with an interactive setup:
  nexlayer wizard [--ai]  # Use --ai for AI-powered recommendations`,
	Example: `  # Initialize a new project
  nexlayer init my-app

  # Deploy using the AI-powered wizard (recommended)
  nexlayer wizard

  # Deploy with AI optimization
  nexlayer deploy -f stack.yaml --ai

  # Add a custom domain
  nexlayer domain add my-app example.com`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := cmd.Help(); err != nil {
				fmt.Printf("Error displaying help: %v\n", err)
			}
			return
		}
	},
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
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nexlayer.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&useAI, "ai", false, "Enable AI-powered recommendations (requires OPENAI_API_KEY)")

	// Register all commands
	commands.RegisterCommands(rootCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".nexlayer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".nexlayer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
