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

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/fsnotify/fsnotify"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"gopkg.in/yaml.v3"
)

var (
	// Styles for different types of output
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

	// Enhanced diff styles
	diffStyleAdded = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")). // Bright green for additions
			MarginLeft(2)

	diffStyleRemoved = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ff0000")). // Bright red for deletions
				MarginLeft(2)

	diffStyleInfo = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")). // Gray for context
			MarginLeft(2)

	diffStyleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00ffff")). // Cyan for headers
			MarginLeft(2).
			MarginBottom(1)
)

// NewCommand creates a new watch command
func NewCommand() *cobra.Command {
	var previewMode bool

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch for project changes and update configuration",
		Long: `Watch for changes in your project and automatically update the nexlayer.yaml configuration.
This is useful during development when your project structure or dependencies change.

In preview mode (--preview), changes will be shown but not applied automatically.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWatchCommand(cmd, previewMode)
		},
	}

	cmd.Flags().BoolVar(&previewMode, "preview", false, "Show changes without applying them")
	return cmd
}

func runWatchCommand(cmd *cobra.Command, previewMode bool) error {
	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Get current directory
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Add directories to watch
	if err := addDirsToWatch(watcher, dir); err != nil {
		return fmt.Errorf("failed to add directories to watch: %w", err)
	}

	// Create a debouncer to avoid too frequent updates
	var lastUpdate time.Time
	debounceInterval := 2 * time.Second

	// Create context for cancellation
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	// Print initial message
	mode := "auto-update"
	if previewMode {
		mode = "preview"
	}
	fmt.Println(titleStyle.Render(fmt.Sprintf("👀 Watching for changes (%s mode)...", mode)))

	for {
		select {
		case <-ctx.Done():
			return nil

		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Skip temporary files and directories
			if shouldSkipFile(event.Name) {
				continue
			}

			// Debounce updates
			if time.Since(lastUpdate) < debounceInterval {
				continue
			}
			lastUpdate = time.Now()

			// Show progress
			spinner := ui.NewSpinner("Analyzing changes...")
			spinner.Start()

			// Generate new configuration
			newConfig, err := generateConfiguration(dir)
			if err != nil {
				spinner.Stop()
				fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Error generating configuration: %v", err)))
				continue
			}

			spinner.Stop()

			// In preview mode, show changes but don't apply them
			if previewMode {
				if err := showConfigurationDiff(newConfig); err != nil {
					fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Error showing changes: %v", err)))
					continue
				}
				fmt.Println(infoStyle.Render("\nℹ️  Run without --preview to apply these changes"))
				continue
			}

			// Apply configuration changes
			if err := applyConfiguration(newConfig); err != nil {
				fmt.Println(errorStyle.Render(fmt.Sprintf("❌ Error applying configuration: %v", err)))
				continue
			}

			fmt.Println(infoStyle.Render(fmt.Sprintf("✨ Configuration updated at %s", time.Now().Format("15:04:05"))))

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Println(warningStyle.Render(fmt.Sprintf("⚠️  Watch error: %v", err)))
		}
	}
}

// showConfigurationDiff displays the differences between current and new configuration
func showConfigurationDiff(newConfig []byte) error {
	// Read current configuration
	currentConfig, err := os.ReadFile("nexlayer.yaml")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read current configuration: %w", err)
	}

	// If no current config exists, show the entire new config
	if len(currentConfig) == 0 {
		fmt.Println(titleStyle.Render("📝 New configuration to be created:"))
		rendered, err := renderMarkdown(string(newConfig))
		if err != nil {
			return err
		}
		fmt.Println(rendered)
		return nil
	}

	// Show diff
	fmt.Println(titleStyle.Render("📝 Configuration changes:"))
	diff := generateEnhancedDiff(string(currentConfig), string(newConfig))
	fmt.Println(diff)
	return nil
}

// generateEnhancedDiff creates a detailed and colorized diff output
func generateEnhancedDiff(current, new string) string {
	dmp := diffmatchpatch.New()

	// Convert both texts to a common line ending format
	current = strings.ReplaceAll(current, "\r\n", "\n")
	new = strings.ReplaceAll(new, "\r\n", "\n")

	// Generate diff
	diffs := dmp.DiffMain(current, new, true)
	diffs = dmp.DiffCleanupSemantic(diffs)

	// Build formatted output
	var output strings.Builder
	var stats struct {
		additions int
		deletions int
		changes   int
	}

	// Context settings
	const (
		contextLines = 3                  // Number of context lines before and after changes
		separator    = "..."              // Separator for non-adjacent changes
		minDistance  = contextLines*2 + 1 // Minimum lines between changes to show separator
	)

	// Convert diffs to lines for context processing
	var lines []struct {
		content string
		typ     diffmatchpatch.Operation
		lineNum int
	}

	lineNum := 1
	for _, d := range diffs {
		diffLines := strings.Split(d.Text, "\n")
		for _, line := range diffLines {
			if line != "" {
				lines = append(lines, struct {
					content string
					typ     diffmatchpatch.Operation
					lineNum int
				}{line, d.Type, lineNum})
				lineNum++
			}
		}
	}

	// Process lines with context
	var lastPrintedLine int
	inChange := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// Check if this line is part of a change
		isChange := line.typ != diffmatchpatch.DiffEqual

		if isChange {
			// If we're not already in a change block, print preceding context
			if !inChange {
				// Add separator if needed
				if lastPrintedLine > 0 && i-lastPrintedLine > minDistance {
					output.WriteString(diffStyleInfo.Render(fmt.Sprintf("  %s\n", separator)))
				}

				// Print preceding context
				start := max(0, i-contextLines)
				for j := start; j < i; j++ {
					output.WriteString(diffStyleInfo.Render(fmt.Sprintf("  %s\n", lines[j].content)))
				}
			}

			// Print the change
			switch line.typ {
			case diffmatchpatch.DiffInsert:
				stats.additions++
				output.WriteString(diffStyleAdded.Render(fmt.Sprintf("+ %s\n", line.content)))
			case diffmatchpatch.DiffDelete:
				stats.deletions++
				output.WriteString(diffStyleRemoved.Render(fmt.Sprintf("- %s\n", line.content)))
			}

			inChange = true
			lastPrintedLine = i
		} else {
			// If we were in a change block, print following context
			if inChange {
				end := min(len(lines), i+contextLines)
				for j := i; j < end; j++ {
					output.WriteString(diffStyleInfo.Render(fmt.Sprintf("  %s\n", lines[j].content)))
				}
				inChange = false
				lastPrintedLine = i + contextLines - 1
				i = end - 1 // Skip the context lines we just printed
			}
		}
	}

	// Add diff summary header
	summary := fmt.Sprintf("Changes: %d additions(+), %d deletions(-)", stats.additions, stats.deletions)
	header := diffStyleHeader.Render(summary)

	return header + "\n" + output.String()
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// renderMarkdown renders YAML as styled markdown
func renderMarkdown(content string) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return "", err
	}

	return r.Render(fmt.Sprintf("```yaml\n%s\n```", content))
}

// generateConfiguration creates a new configuration based on current project state
func generateConfiguration(dir string) ([]byte, error) {
	// Create detector registry
	registry := detection.NewDetectorRegistry()

	// Detect project type and info
	info, err := registry.DetectProject(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to detect project: %w", err)
	}

	// Create template generator
	generator := template.NewGenerator()

	// Generate template
	tmpl, err := generator.GenerateFromProjectInfo(info.Name, string(info.Type), info.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to generate template: %w", err)
	}

	// Add database if needed
	if hasDatabase(info) {
		if err := generator.AddPod(tmpl, template.PodTypePostgres, 0); err != nil {
			return nil, fmt.Errorf("failed to add database: %w", err)
		}
	}

	// Marshal template to YAML using yaml.v3
	return yaml.Marshal(tmpl)
}

// applyConfiguration writes the new configuration to disk
func applyConfiguration(config []byte) error {
	return os.WriteFile("nexlayer.yaml", config, 0644)
}

func addDirsToWatch(watcher *fsnotify.Watcher, root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden directories and node_modules
		if info.IsDir() {
			if shouldSkipDir(path) {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})
}

func shouldSkipDir(path string) bool {
	base := filepath.Base(path)
	return base[0] == '.' || // Hidden directories
		base == "node_modules" ||
		base == "vendor" ||
		base == "dist" ||
		base == "build" ||
		base == "__pycache__"
}

func shouldSkipFile(path string) bool {
	// Skip temporary files and certain extensions
	base := filepath.Base(path)
	return base[0] == '.' || // Hidden files
		filepath.Ext(path) == ".swp" ||
		filepath.Ext(path) == ".swx" ||
		filepath.Ext(path) == ".tmp"
}

// hasDatabase checks if the project uses a database based on its dependencies
func hasDatabase(info *types.ProjectInfo) bool {
	// Check for common database dependencies
	dbDeps := []string{
		"pg",             // PostgreSQL
		"mysql",          // MySQL
		"mysql2",         // MySQL2
		"sequelize",      // SQL ORM
		"mongoose",       // MongoDB
		"mongodb",        // MongoDB
		"typeorm",        // TypeORM
		"prisma",         // Prisma
		"sqlite3",        // SQLite
		"better-sqlite3", // Better SQLite3
		"redis",          // Redis
		"ioredis",        // Redis
	}

	for _, dep := range dbDeps {
		if _, ok := info.Dependencies[dep]; ok {
			return true
		}
	}

	return false
}
