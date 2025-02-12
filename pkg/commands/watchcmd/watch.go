package watchcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/validation"
	"github.com/fsnotify/fsnotify"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewCommand creates a new watch command that monitors nexlayer.yaml for changes
// and automatically syncs and redeploys when changes are detected.
func NewCommand() *cobra.Command {
	var (
		configFile string
		noSync    bool
		noDeploy  bool
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for changes and auto-sync/deploy",
		Long: `Watch for changes in nexlayer.yaml and automatically sync and redeploy.
		
When changes are detected:
1. Validates the new configuration
2. Syncs the configuration with project state
3. Triggers a redeployment if needed

Use --no-sync to disable auto-sync
Use --no-deploy to disable auto-deploy`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create watcher
			watcher, err := fsnotify.NewWatcher()
			if err != nil {
				return fmt.Errorf("failed to create file watcher: %w", err)
			}
			defer watcher.Close()

			// Get absolute path to config file
			absPath, err := filepath.Abs(configFile)
			if err != nil {
				return fmt.Errorf("failed to get absolute path for %s: %w", configFile, err)
			}

			// Verify config file exists
			if _, err := os.Stat(absPath); err != nil {
				return fmt.Errorf("nexlayer.yaml not found at %s. Run 'nexlayer init' first", absPath)
			}

			// Add file to watcher
			if err := watcher.Add(absPath); err != nil {
				return fmt.Errorf("failed to watch %s: %w", absPath, err)
			}

			// Create validator
			validator := validation.NewValidator(false) // Non-strict mode for watch

			// sync is a placeholder function for syncing changes.
			// This function should be implemented with the actual sync logic.
			sync := func() error {
				pterm.Info.Println("Syncing changes...")
				// TODO: Implement sync logic
				return nil
			}

			deploy := func() error {
				// TODO: Implement deploy logic
				return nil
			}

			// Start watching
			pterm.Info.Printf("👀 Watching for changes in %s...\n", configFile)

			// Debounce timer to handle multiple rapid changes
			var debounceTimer *time.Timer
			const debounceDelay = 500 * time.Millisecond

			// Watch for changes
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return fmt.Errorf("watcher channel closed unexpectedly")
					}

					if event.Op&fsnotify.Write == fsnotify.Write {
						// Reset or start debounce timer
						if debounceTimer != nil {
							debounceTimer.Stop()
						}
						debounceTimer = time.AfterFunc(debounceDelay, func() {
							// Validate changes
							pterm.Info.Println("🔍 Validating changes...")
							yamlData, err := os.ReadFile(absPath)
							if err != nil {
								pterm.Error.Printf("Failed to read config: %v\n", err)
								return
							}

							var config schema.NexlayerYAML
							if err := yaml.Unmarshal(yamlData, &config); err != nil {
								pterm.Error.Printf("Failed to parse config: %v\n", err)
								return
							}

							validationErrors := validator.ValidateYAML(&config)
							if len(validationErrors) > 0 {
								pterm.Warning.Println("⚠️  Validation found issues:")
								for _, issue := range validationErrors {
									fmt.Printf("  - %s\n", issue)
								}
								return
							}
							pterm.Success.Println("✅ Configuration is valid")

							// Run sync if enabled
							if !noSync {
								pterm.Info.Println("🔄 Running sync...")
								if err := sync(); err != nil {
									pterm.Error.Printf("Sync failed: %v\n", err)
									return
								}
								pterm.Success.Println("✅ Sync completed")
							}

							// Run deploy if enabled
							if !noDeploy {
								pterm.Info.Println("🚀 Triggering deployment...")
								if err := deploy(); err != nil {
									pterm.Error.Printf("Deployment failed: %v\n", err)
									return
								}
								pterm.Success.Println("✅ Deployment started")
							}
						})
					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return fmt.Errorf("watcher error channel closed unexpectedly")
					}
					pterm.Error.Printf("Watch error: %v\n", err)
				}
			}
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&configFile, "file", "f", "nexlayer.yaml", "Path to nexlayer.yaml")
	cmd.Flags().BoolVar(&noSync, "no-sync", false, "Disable auto-sync")
	cmd.Flags().BoolVar(&noDeploy, "no-deploy", false, "Disable auto-deploy")

	return cmd
}
