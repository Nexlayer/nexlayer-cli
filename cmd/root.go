package cmd

import (
	"fmt"
	"os"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/help"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "nexlayer",
	Short:   "A modern cloud application deployment tool",
	Long:    help.RootLongDesc,
	Example: help.RootExample,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
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

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nexlayer.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Version flag
	rootCmd.Flags().Bool("version", false, "display version information")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		// viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".nexlayer" (without extension).
		cfgFile = fmt.Sprintf("%s/.nexlayer.yaml", home)
	}

	// If a config file is found, read it in.
	if _, err := os.Stat(cfgFile); err == nil {
		fmt.Println("Using config file:", cfgFile)
	}
}
