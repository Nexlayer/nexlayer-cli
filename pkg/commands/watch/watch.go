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
	"github.com/Nexlayer/nexlayer-cli/pkg/core/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
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
)

// NewCommand creates a new watch command
func NewCommand(apiClient api.APIClient) *cobra.Command {
	var previewMode bool

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Monitor & auto-update configuration",
		Long: `Watch the project for changes and automatically update nexlayer.yaml configuration.
When changes are detected (new dependencies, frameworks, services, Docker images, etc.),
the configuration will be updated to match the current project state.

The command will automatically detect nexlayer.yaml in the current directory.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find nexlayer.yaml in current directory
			configFile, err := findConfigFile()
			if err != nil {
				return fmt.Errorf("nexlayer.yaml not found in current directory: %w", err)
			}

			return runWatch(cmd, configFile, previewMode)
		},
	}

	cmd.Flags().BoolVar(&previewMode, "preview", false, "(Future) Show changes without applying them")

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
func runWatch(cmd *cobra.Command, configFile string, previewMode bool) error {
	// Create new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add current directory to watch
	if err := addDirsToWatch(watcher, "."); err != nil {
		return fmt.Errorf("failed to watch current directory: %w", err)
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	// Load initial configuration
	currentConfig, err := loadCurrentConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load current configuration: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Watching for project changes...\n")
	fmt.Fprintf(cmd.OutOrStdout(), "Configuration will be updated when new components are detected.\n\n")

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
			projectInfo, err := registry.DetectProject(".")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error detecting project changes: %v\n", err)
				continue
			}

			// Generate new configuration
			generator := template.NewGenerator()
			newConfig, err := generator.GenerateFromProjectInfo(projectInfo.Name, string(projectInfo.Type), projectInfo.Port)
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error generating configuration: %v\n", err)
				continue
			}

			// Add database if needed
			if hasDatabase(projectInfo) {
				if err := generator.AddPod(newConfig, template.PodTypePostgres, 0); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Error adding database configuration: %v\n", err)
				}
			}

			// Add AI-specific configurations if detected
			if projectInfo.LLMProvider != "" {
				addAIConfigurations(newConfig, projectInfo)
			}

			// Check for Docker images
			dockerImages, err := findDockerImages(".")
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error scanning for Docker images: %v\n", err)
			} else {
				for _, image := range dockerImages {
					pod := template.Pod{
						Name:  fmt.Sprintf("docker-%s", strings.Split(image, ":")[0]),
						Type:  "docker",
						Image: image,
						ServicePorts: []template.ServicePort{
							{Name: "http", Port: 80, TargetPort: 80},
						},
					}
					newConfig.Application.Pods = append(newConfig.Application.Pods, pod)
				}
			}

			// Compare configurations
			if configsEqual(currentConfig, newConfig) {
				fmt.Fprintf(cmd.OutOrStdout(), "No configuration changes needed.\n")
				continue
			}

			if previewMode {
				// Show changes
				fmt.Fprintf(cmd.OutOrStdout(), "\nConfiguration changes detected:\n")
				showConfigurationDiff(currentConfig, newConfig)

				if promptYesNo("Apply these changes?") {
					if err := writeYAMLToFile(configFile, newConfig); err != nil {
						fmt.Fprintf(cmd.ErrOrStderr(), "Error writing configuration: %v\n", err)
						continue
					}
					currentConfig = newConfig
					fmt.Fprintf(cmd.OutOrStdout(), "Configuration updated successfully.\n")
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Changes not applied.\n")
				}
			} else {
				// Apply changes directly
				if err := writeYAMLToFile(configFile, newConfig); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "Error writing configuration: %v\n", err)
					continue
				}
				currentConfig = newConfig
				fmt.Fprintf(cmd.OutOrStdout(), "Configuration updated successfully.\n")
			}

		case <-ctx.Done():
			return nil
		}
	}
}

// showConfigurationDiff displays the differences between current and new configuration
func showConfigurationDiff(current, new *template.NexlayerYAML) {
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
func loadCurrentConfig(configFile string) (*template.NexlayerYAML, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &template.NexlayerYAML{}, nil
		}
		return nil, err
	}

	var config template.NexlayerYAML
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// configsEqual compares two configurations for equality
func configsEqual(a, b *template.NexlayerYAML) bool {
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
func writeYAMLToFile(configFile string, config *template.NexlayerYAML) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

// hasDatabase checks if the project has a database
func hasDatabase(projectInfo *types.ProjectInfo) bool {
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

// addAIConfigurations adds AI-specific configurations to the configuration
func addAIConfigurations(config *template.NexlayerYAML, projectInfo *types.ProjectInfo) {
	// Add AI-specific annotations to all pods
	for i := range config.Application.Pods {
		if config.Application.Pods[i].Annotations == nil {
			config.Application.Pods[i].Annotations = make(map[string]string)
		}
		config.Application.Pods[i].Annotations["ai.nexlayer.io/provider"] = projectInfo.LLMProvider
		config.Application.Pods[i].Annotations["ai.nexlayer.io/enabled"] = "true"
	}
}
