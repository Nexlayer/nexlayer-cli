package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/plugins"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plugin",
		Short: "Manage Nexlayer plugins",
		Long: `Manage Nexlayer plugins.
		
Available Commands:
  list    List installed plugins
  run     Run a plugin
  install Install a plugin from a .so file`,
	}

	cmd.AddCommand(
		newListCommand(),
		newRunCommand(),
		newInstallCommand(),
	)

	return cmd
}

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed plugins",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := plugins.NewManager()
			if err := manager.LoadPluginsFromDir(""); err != nil {
				return err
			}

			plugins := manager.ListPlugins()
			if len(plugins) == 0 {
				cmd.Println("No plugins installed")
				return nil
			}

			cmd.Println("Installed plugins:")
			for _, name := range plugins {
				plugin, _ := manager.GetPlugin(name)
				cmd.Printf("  %s - %s\n", name, plugin.Description())
			}
			return nil
		},
	}
}

func newRunCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run [plugin-name]",
		Short: "Run a plugin",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := plugins.NewManager()
			if err := manager.LoadPluginsFromDir(""); err != nil {
				return err
			}

			// Get options from flags
			opts := make(map[string]interface{})
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				opts[f.Name] = f.Value.String()
			})

			return manager.RunPlugin(args[0], opts)
		},
	}
}

func newInstallCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "install [plugin-file]",
		Short: "Install a plugin from a .so file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginFile := args[0]
			if filepath.Ext(pluginFile) != ".so" {
				return fmt.Errorf("plugin file must be a .so file")
			}

			// Get plugin directory
			pluginDir := filepath.Join(os.Getenv("HOME"), ".nexlayer", "plugins")
			if err := os.MkdirAll(pluginDir, 0755); err != nil {
				return fmt.Errorf("failed to create plugin directory: %w", err)
			}

			// Copy plugin to plugins directory
			dest := filepath.Join(pluginDir, filepath.Base(pluginFile))
			data, err := os.ReadFile(pluginFile)
			if err != nil {
				return fmt.Errorf("failed to read plugin file: %w", err)
			}

			if err := os.WriteFile(dest, data, 0644); err != nil {
				return fmt.Errorf("failed to install plugin: %w", err)
			}

			cmd.Printf("Plugin installed successfully to %s\n", dest)
			return nil
		},
	}
}
