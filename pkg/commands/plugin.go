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
	Use:   "plugin [command]",
	Short: "Manage plugins",
	Long: `Manage Nexlayer plugins.
Example: nexlayer plugin install my-plugin`,
	Args: cobra.MinimumNArgs(1),
	RunE: runPlugin,
}

func runPlugin(cmd *cobra.Command, args []string) error {
	command := args[0]

	switch command {
	case "install":
		if len(args) < 2 {
			return fmt.Errorf("plugin name is required")
		}
		return installPlugin(args[1])
	case "list":
		return listPlugins()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func installPlugin(name string) error {
	pluginDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// Clone plugin repository
	repoURL := fmt.Sprintf("https://github.com/nexlayer/plugin-%s.git", name)
	cmd := exec.Command("git", "clone", repoURL, filepath.Join(pluginDir, name))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone plugin repository: %w", err)
	}

	fmt.Printf("Successfully installed plugin: %s\n", name)
	return nil
}

func listPlugins() error {
	pluginDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
	entries, err := os.ReadDir(pluginDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No plugins installed")
			return nil
		}
		return fmt.Errorf("failed to list plugins: %w", err)
	}

	fmt.Println("Installed plugins:")
	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("  - %s\n", entry.Name())
		}
	}

	return nil
}
