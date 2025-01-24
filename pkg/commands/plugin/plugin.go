// Package commands contains the CLI commands for the Nexlayer application.
package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// PluginCmd represents the plugin command
var PluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage Nexlayer plugins",
	Long:  "Manage Nexlayer plugins - install, list, and remove plugins",
}

func init() {
	PluginCmd.AddCommand(newInstallCommand())
	PluginCmd.AddCommand(newListCommand())
	PluginCmd.AddCommand(newRemoveCommand())
}

func newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install [plugin]",
		Short: "Install a plugin",
		Long:  "Install a plugin from the Nexlayer plugin registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("requires plugin name")
			}
			pluginDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
			if err := os.MkdirAll(pluginDir, 0755); err != nil {
				return fmt.Errorf("failed to create plugin directory: %w", err)
			}

			// Clone plugin repository
			repoURL := fmt.Sprintf("https://github.com/nexlayer/plugin-%s.git", args[0])
			cmdInstall := exec.Command("git", "clone", repoURL, filepath.Join(pluginDir, args[0]))
			if err := cmdInstall.Run(); err != nil {
				return fmt.Errorf("failed to clone plugin repository: %w", err)
			}

			fmt.Printf("Successfully installed plugin: %s\n", args[0])
			return nil
		},
	}
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		Long:  "List all installed Nexlayer plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
			entries, err := os.ReadDir(pluginDir)
			if err != nil {
				if os.IsNotExist(err) {
					cmd.Println("No plugins installed")
					return nil
				}
				return fmt.Errorf("failed to list plugins: %w", err)
			}

			if len(entries) == 0 {
				cmd.Println("No plugins installed")
				return nil
			}

			cmd.Println("Installed plugins:")
			for _, entry := range entries {
				if entry.IsDir() {
					cmd.Printf("  %s\n", entry.Name())
				}
			}

			return nil
		},
	}
}

func newRemoveCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "remove [plugin]",
		Short: "Remove a plugin",
		Long:  "Remove an installed Nexlayer plugin",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("requires plugin name")
			}
			return nil
		},
	}
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage Nexlayer plugins",
		Long:  "Manage Nexlayer plugins - install, list, and remove plugins",
	}

	cmd.AddCommand(newInstallCommand())
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newRemoveCommand())

	return cmd
}
