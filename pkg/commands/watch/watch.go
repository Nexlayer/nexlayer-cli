// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package watch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

	// lastKnownModTime tracks when the watch command last modified the config file
	lastKnownModTime time.Time
)

// NewCommand creates a new watch command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Monitor project changes and update configuration",
		Long: `Watch the project for changes and automatically update nexlayer.yaml configuration.
When changes are detected (new dependencies, frameworks, services, Docker images, etc.),
the configuration will be updated to match the current project state.

The command runs in the foreground. Press Ctrl+C to stop watching.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find nexlayer.yaml in current directory
			configFile, err := findConfigFile()
			if err != nil {
				return fmt.Errorf("nexlayer.yaml not found in current directory: %w", err)
			}

			return runWatch(cmd, configFile)
		},
	}

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

// runWatch starts watching for project changes and updates configuration
func runWatch(cmd *cobra.Command, configFile string) error {
	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Initialize lastKnownModTime
	if fileInfo, err := os.Stat(configFile); err == nil {
		lastKnownModTime = fileInfo.ModTime()
	}

	// Add current directory to watch
	if err := addDirsToWatch(watcher, "."); err != nil {
		return fmt.Errorf("failed to watch current directory: %w", err)
	}

	// Create context with cancellation
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nStopping watch mode...")
		cancel()
	}()

	// Load initial configuration
	currentConfig, err := loadCurrentConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load current configuration: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Watching for project changes...\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Configuration will be updated when new components are detected.\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Manual edits to nexlayer.yaml will be preserved unless you choose to overwrite them.\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Press Ctrl+C to stop watching.\n\n")

	// Use a fixed debounce time
	debounceTime := 2 * time.Second
	debounceCh := make(chan struct{})
	var timer *time.Timer

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

			// Skip the nexlayer.yaml file itself to avoid loops
			if filepath.Base(event.Name) == filepath.Base(configFile) {
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
			// Analyze project for changes
			fmt.Fprintf(cmd.OutOrStdout(), "Analyzing project changes...\n")

			// Create project detector
			registry := detection.NewDetectorRegistry()
			detectedInfo, err := registry.DetectProject(".")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error detecting project changes: %v\n", err)
				continue
			}

			// Convert to schema.ProjectInfo
			projectInfo := &schema.ProjectInfo{
				Type:         schema.ProjectType(detectedInfo.Type),
				Name:         detectedInfo.Name,
				Version:      detectedInfo.Version,
				Dependencies: detectedInfo.Dependencies,
				Scripts:      detectedInfo.Scripts,
				Port:         detectedInfo.Port,
				HasDocker:    detectedInfo.HasDocker,
				LLMProvider:  detectedInfo.LLMProvider,
				LLMModel:     detectedInfo.LLMModel,
				ImageTag:     detectedInfo.ImageTag,
			}

			// Generate new configuration
			generator := schema.NewGenerator()
			newConfig, err := generator.GenerateFromProjectInfo(projectInfo.Name, string(projectInfo.Type), projectInfo.Port)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error generating configuration: %v\n", err)
				continue
			}

			// Add database if needed
			if hasDatabase(projectInfo) {
				if err := generator.AddPod(newConfig, "postgres", 0); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Error adding database configuration: %v\n", err)
				}
			}

			// Add AI-specific configurations if detected
			if projectInfo.LLMProvider != "" {
				generator.AddAIConfigurations(newConfig, projectInfo.LLMProvider)
			}

			// Check for Docker images
			dockerImages, err := findDockerImages(".")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error scanning for Docker images: %v\n", err)
			} else {
				for _, image := range dockerImages {
					if err := generator.AddPod(newConfig, "docker", 0); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Error adding Docker pod: %v\n", err)
						continue
					}
					pod := &newConfig.Application.Pods[len(newConfig.Application.Pods)-1]
					pod.Name = fmt.Sprintf("docker-%s", strings.Split(image, ":")[0])
					pod.Image = image
				}
			}

			// Compare configurations
			if configsEqual(currentConfig, newConfig) {
				fmt.Fprintf(cmd.OutOrStdout(), "No configuration changes needed.\n")
				continue
			}

			// Show changes
			fmt.Fprintf(cmd.OutOrStdout(), "\nConfiguration changes detected:\n")
			showConfigurationDiff(currentConfig, newConfig)

			// Check for manual edits
			fileInfo, err := os.Stat(configFile)
			if err == nil && fileInfo.ModTime().After(lastKnownModTime) {
				// File has been modified externally
				fmt.Fprintf(cmd.OutOrStdout(), warningStyle.Render("\n⚠️  The configuration file has been manually edited since the last automatic update.\n"))
				if !promptYesNo("Overwrite manual changes with new configuration?") {
					fmt.Fprintf(cmd.OutOrStdout(), "Skipping configuration update to preserve manual changes.\n")

					// Reload the current config to include manual changes
					if newCurrent, err := loadCurrentConfig(configFile); err == nil {
						currentConfig = newCurrent
						lastKnownModTime = fileInfo.ModTime()
					}
					continue
				}
			}

			// Apply changes
			if err := writeYAMLToFile(configFile, newConfig); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error writing configuration: %v\n", err)
				continue
			}

			// Update last known modification time
			if fileInfo, err := os.Stat(configFile); err == nil {
				lastKnownModTime = fileInfo.ModTime()
			}

			currentConfig = newConfig
			fmt.Fprintf(cmd.OutOrStdout(), "Configuration updated successfully.\n")

		case <-ctx.Done():
			fmt.Fprintf(cmd.OutOrStdout(), "Watch mode stopped.\n")
			return nil
		}
	}
}

// showConfigurationDiff displays the differences between current and new configuration
func showConfigurationDiff(current, new *schema.NexlayerYAML) {
	// Convert configs to YAML for comparison
	currentYAML, _ := yaml.Marshal(current)
	newYAML, _ := yaml.Marshal(new)

	// Show diff
	fmt.Println(titleStyle.Render("📝 Configuration changes:"))
	diff := generateColorCodedDiff(string(currentYAML), string(newYAML))
	fmt.Println(diff)
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

// findDockerImages scans the project for Dockerfile and docker-compose.yml files
func findDockerImages(dir string) ([]string, error) {
	var images []string

	// Look for Dockerfile
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		// Parse Dockerfile for base image
		content, err := os.ReadFile(filepath.Join(dir, "Dockerfile"))
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "FROM ") {
					image := strings.TrimSpace(strings.TrimPrefix(line, "FROM "))
					images = append(images, image)
				}
			}
		}
	}

	// Look for docker-compose.yml
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err == nil {
		// Parse docker-compose.yml for images
		content, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
		if err == nil {
			var compose struct {
				Services map[string]struct {
					Image string `yaml:"image"`
				} `yaml:"services"`
			}
			if err := yaml.Unmarshal(content, &compose); err == nil {
				for _, service := range compose.Services {
					if service.Image != "" {
						images = append(images, service.Image)
					}
				}
			}
		}
	}

	return images, nil
}

// loadCurrentConfig loads the current nexlayer.yaml configuration
func loadCurrentConfig(configFile string) (*schema.NexlayerYAML, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &schema.NexlayerYAML{}, nil
		}
		return nil, err
	}

	var config schema.NexlayerYAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// configsEqual compares two configurations for equality
func configsEqual(a, b *schema.NexlayerYAML) bool {
	aData, err := yaml.Marshal(a)
	if err != nil {
		return false
	}
	bData, err := yaml.Marshal(b)
	if err != nil {
		return false
	}
	return string(aData) == string(bData)
}

// writeYAMLToFile writes the configuration to a file
func writeYAMLToFile(configFile string, config *schema.NexlayerYAML) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

// hasDatabase checks if the project has a database
func hasDatabase(projectInfo *schema.ProjectInfo) bool {
	// Check dependencies for database-related packages
	for name := range projectInfo.Dependencies {
		switch name {
		case "pg", "postgres", "postgresql", "sequelize", "typeorm", "prisma",
			"mongoose", "mongodb", "mysql", "mysql2", "sqlite3", "redis":
			return true
		}
	}
	return false
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
