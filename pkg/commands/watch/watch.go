// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package watch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Styles for different types of output
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ff00")).
			MarginBottom(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffff"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00"))

	diffStyleAdded = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00"))

	diffStyleRemoved = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0000"))

	diffStyleUnchanged = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888"))
)

// NewCommand creates a new watch command
func NewCommand(apiClient api.APIClient) *cobra.Command {
	var configFile string
	var debounceTime time.Duration
	var previewMode bool
	var watchDirs []string

	cmd := &cobra.Command{
		Use:   "watch [applicationID]",
		Short: "Watch for changes and redeploy",
		Long: `Watch the specified directories for changes and automatically redeploy.
When files change, the application will be redeployed using the specified configuration.

Example:
  nexlayer watch myapp --file deployment.yaml --watch-dirs ./src`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get app ID if provided
			var appID string
			if len(args) > 0 {
				appID = args[0]
			}

			// If no config file specified, try to find one
			if configFile == "" {
				file, err := findConfigFile()
				if err != nil {
					return err
				}
				configFile = file
				cmd.Printf("Using config file: %s\n", configFile)
			}

			return runWatch(cmd, apiClient, appID, configFile, debounceTime, previewMode, watchDirs)
		},
	}

	cmd.Flags().StringVarP(&configFile, "file", "f", "", "Path to deployment YAML file")
	cmd.Flags().DurationVarP(&debounceTime, "debounce", "d", 2*time.Second, "Debounce time between deployments")
	cmd.Flags().BoolVar(&previewMode, "preview", false, "Show changes without applying them")
	cmd.Flags().StringSliceVar(&watchDirs, "watch-dirs", []string{"."}, "Directories to watch for changes")

	return cmd
}

// findConfigFile looks for a deployment configuration file
func findConfigFile() (string, error) {
	possibleFiles := []string{
		"deployment.yaml",
		"deployment.yml",
		"nexlayer.yaml",
		"nexlayer.yml",
	}

	for _, file := range possibleFiles {
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}

	return "", fmt.Errorf("no deployment file found in current directory. Expected one of: %v", possibleFiles)
}

// runWatch starts watching for file changes and triggers redeployment
func runWatch(cmd *cobra.Command, client api.APIClient, appID, configFile string, debounceTime time.Duration, previewMode bool, watchDirs []string) error {
	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add directories to watch
	for _, dir := range watchDirs {
		if err := addDirsToWatch(watcher, dir); err != nil {
			return fmt.Errorf("failed to add directory to watch: %w", err)
		}
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	// Channel for debounced events
	debounceCh := make(chan struct{})
	var timer *time.Timer

	// Watch for changes
	fmt.Fprintf(cmd.OutOrStdout(), "Watching for changes...\n")
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Skip temporary files and hidden directories
			if shouldIgnoreFile(event.Name) {
				continue
			}

			// Reset timer for debouncing
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(debounceTime, func() {
				debounceCh <- struct{}{}
			})

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)

		case <-debounceCh:
			// Trigger deployment
			fmt.Fprintf(cmd.OutOrStdout(), "Changes detected, preparing to deploy...\n")
			if previewMode {
				if err := showConfigurationDiff(configFile); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Error showing changes: %v\n", err)
					continue
				}
				if promptYesNo("Apply these changes?") {
					resp, err := client.StartDeployment(ctx, appID, configFile)
					if err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Deployment failed: %v\n", err)
						continue
					}
					fmt.Fprintf(cmd.OutOrStdout(), "Deployment started: %s\n", resp.Data.URL)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Changes not applied\n")
				}
			} else {
				resp, err := client.StartDeployment(ctx, appID, configFile)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Deployment failed: %v\n", err)
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "Deployment started: %s\n", resp.Data.URL)
			}

		case <-ctx.Done():
			return nil
		}
	}
}

// showConfigurationDiff displays the differences between current and new configuration
func showConfigurationDiff(configFile string) error {
	// Read current configuration
	currentConfig, err := os.ReadFile(configFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read current configuration: %w", err)
	}

	// Simulate new configuration (replace with actual new config generation logic)
	newConfig := []byte("new config content") // Placeholder

	// If no current config exists, show the entire new config
	if len(currentConfig) == 0 {
		fmt.Println(titleStyle.Render("📝 New configuration to be created:"))
		fmt.Println(string(newConfig))
		return nil
	}

	// Show diff
	fmt.Println(titleStyle.Render("📝 Configuration changes:"))
	diff := generateColorCodedDiff(string(currentConfig), string(newConfig))
	fmt.Println(diff)
	return nil
}

// generateColorCodedDiff creates a color-coded diff output
func generateColorCodedDiff(current, new string) string {
	currentLines := strings.Split(current, "\n")
	newLines := strings.Split(new, "\n")

	var diff strings.Builder
	for i := 0; i < len(newLines); i++ {
		if i >= len(currentLines) {
			// New lines added
			diff.WriteString(diffStyleAdded.Render(fmt.Sprintf("+ %s", newLines[i])) + "\n")
			continue
		}
		if currentLines[i] != newLines[i] {
			diff.WriteString(diffStyleRemoved.Render(fmt.Sprintf("- %s", currentLines[i])) + "\n")
			diff.WriteString(diffStyleAdded.Render(fmt.Sprintf("+ %s", newLines[i])) + "\n")
		} else {
			diff.WriteString(diffStyleUnchanged.Render(fmt.Sprintf("  %s", newLines[i])) + "\n")
		}
	}

	return diff.String()
}

// addDirsToWatch recursively adds directories to the watcher
func addDirsToWatch(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and irrelevant ones
		if info.IsDir() && shouldIgnoreDir(path) {
			return filepath.SkipDir
		}

		// Add directory to watcher
		if info.IsDir() {
			return watcher.Add(path)
		}

		return nil
	})
}

// shouldIgnoreFile checks if a file should be ignored
func shouldIgnoreFile(path string) bool {
	base := filepath.Base(path)
	return base[0] == '.' || strings.HasSuffix(path, ".swp") || strings.HasSuffix(path, ".swx") || strings.HasSuffix(path, ".tmp")
}

// shouldIgnoreDir checks if a directory should be ignored
func shouldIgnoreDir(path string) bool {
	base := filepath.Base(path)
	return base[0] == '.' || base == "node_modules" || base == "vendor" || base == "dist" || base == "build" || base == "__pycache__" || base == ".git"
}

// promptYesNo asks the user to confirm an action
func promptYesNo(label string) bool {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	result, err := prompt.Run()
	if err != nil {
		return false
	}
	return strings.ToLower(result) == "y"
}
