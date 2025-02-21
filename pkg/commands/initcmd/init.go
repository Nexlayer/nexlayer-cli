// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package initcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

const (
	cacheDir  = ".nexlayer"
	cacheFile = "detection-cache.json"
)

var (
	// Styles for different types of output
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ffff"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff0000"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffff00"))
)

// detectionCache represents cached project detection results
type detectionCache struct {
	ProjectInfo *types.ProjectInfo `json:"project_info"`
	Timestamp   time.Time          `json:"timestamp"`
}

// NewCommand initializes a new Nexlayer project
func NewCommand() *cobra.Command {
	var interactive bool
	var force bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Nexlayer project",
		Long: `Initialize a new Nexlayer project by creating a nexlayer.yaml file in the current directory.
The command will automatically detect your project type and configure it appropriately.

Examples:
  # Auto-detect and initialize
  nexlayer init

  # Interactive mode
  nexlayer init --interactive

  # Force re-detection (ignore cache)
  nexlayer init --force`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInitCommand(cmd, interactive, force)
		},
	}

	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Enable interactive mode")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-detection (ignore cache)")

	return cmd
}

// runInitCommand handles the execution of the init command
func runInitCommand(cmd *cobra.Command, interactive, force bool) error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Show welcome message
	fmt.Println(infoStyle.Render("🚀 Initializing Nexlayer project..."))

	// Try to load from cache first
	var info *types.ProjectInfo
	if !force {
		info = loadFromCache(cwd)
	}

	// If not in cache or force flag is set, detect project
	if info == nil {
		var err error
		info, err = detectProjectParallel(cwd)
		if err != nil && interactive {
			// If detection fails in interactive mode, prompt user
			info, err = promptForProjectType()
		}
		if err != nil {
			return fmt.Errorf("failed to detect project type: %w", err)
		}

		// Save to cache
		if err := saveToCache(cwd, info); err != nil {
			fmt.Println(warningStyle.Render("⚠️  Warning: Failed to cache detection results"))
		}
	}

	// Show progress spinner
	p := tea.NewProgram(newSpinnerModel("Generating configuration..."))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to show progress: %w", err)
	}

	// Generate template
	generator := template.NewGenerator()
	tmpl, err := generator.GenerateFromProjectInfo(info.Name, string(info.Type), info.Port)
	if err != nil {
		return fmt.Errorf("failed to generate template: %w\nTry running with --interactive flag", err)
	}

	// Add database if needed
	if hasDatabase(info) {
		if err := generator.AddPod(tmpl, template.PodTypePostgres, 0); err != nil {
			return fmt.Errorf("failed to add database: %w", err)
		}
	}

	// Add AI-specific configurations if AI IDE is detected
	if info.LLMProvider != "" {
		addAIConfigurations(tmpl, info)
	}

	// Write configuration
	if err := writeYAMLToFile("nexlayer.yaml", tmpl); err != nil {
		return fmt.Errorf("failed to write configuration: %w\nCheck file permissions and disk space", err)
	}

	// Print success message with detected info
	printSuccessMessage(info, tmpl)

	return nil
}

// detectProjectParallel runs project detection in parallel
func detectProjectParallel(dir string) (*types.ProjectInfo, error) {
	registry := detection.NewDetectorRegistry()
	detectors := registry.GetDetectors()

	// Create channels for results and errors
	resultCh := make(chan *types.ProjectInfo, len(detectors))
	errCh := make(chan error, len(detectors))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run detectors in parallel
	var wg sync.WaitGroup
	for _, d := range detectors {
		wg.Add(1)
		go func(det detection.ProjectDetector) {
			defer wg.Done()
			if info, err := det.Detect(dir); err == nil && info != nil {
				select {
				case resultCh <- info:
				case <-ctx.Done():
				}
			}
		}(d)
	}

	// Wait for first successful result or all failures
	go func() {
		wg.Wait()
		close(resultCh)
		close(errCh)
	}()

	select {
	case info := <-resultCh:
		return info, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("detection timed out")
	case err := <-errCh:
		return nil, err
	}
}

// promptForProjectType asks the user to select a project type
func promptForProjectType() (*types.ProjectInfo, error) {
	prompt := promptui.Select{
		Label: "Select your project type",
		Items: []string{
			"Next.js",
			"React",
			"Node.js",
			"Python",
			"Go",
			"Docker",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// Convert selection to ProjectType
	var projectType types.ProjectType
	switch result {
	case "Next.js":
		projectType = types.TypeNextjs
	case "React":
		projectType = types.TypeReact
	case "Node.js":
		projectType = types.TypeNode
	case "Python":
		projectType = types.TypePython
	case "Go":
		projectType = types.TypeGo
	case "Docker":
		projectType = types.TypeDockerRaw
	}

	return &types.ProjectInfo{
		Type: projectType,
		Name: filepath.Base(filepath.Dir("")),
	}, nil
}

// loadFromCache attempts to load project info from cache
func loadFromCache(dir string) *types.ProjectInfo {
	cachePath := filepath.Join(dir, cacheDir, cacheFile)
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil
	}

	var cache detectionCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil
	}

	// Check if cache is still valid (24 hours)
	if time.Since(cache.Timestamp) > 24*time.Hour {
		return nil
	}

	return cache.ProjectInfo
}

// saveToCache saves project info to cache
func saveToCache(dir string, info *types.ProjectInfo) error {
	cachePath := filepath.Join(dir, cacheDir)
	if err := os.MkdirAll(cachePath, 0755); err != nil {
		return err
	}

	cache := detectionCache{
		ProjectInfo: info,
		Timestamp:   time.Now(),
	}

	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(cachePath, cacheFile), data, 0644)
}

// writeYAMLToFile writes the template to a YAML file
func writeYAMLToFile(filename string, tmpl *template.NexlayerYAML) error {
	// Create backup if file exists
	if _, err := os.Stat(filename); err == nil {
		backupFile := filename + ".backup"
		if err := os.Rename(filename, backupFile); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Printf("📦 Backed up existing %s to %s\n", filename, backupFile)
	}

	// Write new file
	data, err := yaml.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// addAIConfigurations adds AI-specific settings to the template
func addAIConfigurations(tmpl *template.NexlayerYAML, info *types.ProjectInfo) {
	// Add AI-specific annotations
	for i := range tmpl.Application.Pods {
		if tmpl.Application.Pods[i].Annotations == nil {
			tmpl.Application.Pods[i].Annotations = make(map[string]string)
		}
		tmpl.Application.Pods[i].Annotations["ai.nexlayer.io/enabled"] = "true"
		tmpl.Application.Pods[i].Annotations["ai.nexlayer.io/provider"] = info.LLMProvider
		tmpl.Application.Pods[i].Annotations["ai.nexlayer.io/model"] = info.LLMModel
	}
}

// printSuccessMessage prints a detailed success message
func printSuccessMessage(info *types.ProjectInfo, tmpl *template.NexlayerYAML) {
	fmt.Println(successStyle.Render("\n✨ Project initialized successfully!"))
	fmt.Println(infoStyle.Render("\nDetected Configuration:"))
	fmt.Printf("• Project Type: %s\n", info.Type)
	if info.Version != "" {
		fmt.Printf("• Version: %s\n", info.Version)
	}
	if info.LLMProvider != "" {
		fmt.Printf("• AI Integration: %s (%s)\n", info.LLMProvider, info.LLMModel)
	}
	fmt.Printf("• Port: %d\n", info.Port)

	fmt.Println(infoStyle.Render("\nNext Steps:"))
	fmt.Println("1. Review nexlayer.yaml")
	fmt.Println("2. Run 'nexlayer deploy' to deploy your application")
	fmt.Println("3. Run 'nexlayer help' for more commands")
}

// spinnerModel represents the progress spinner
type spinnerModel struct {
	spinner  spinner.Model
	message  string
	quitting bool
}

func newSpinnerModel(message string) spinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return spinnerModel{spinner: s, message: message}
}

func (m spinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m spinnerModel) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

// hasDatabase checks if the project needs a database
func hasDatabase(info *types.ProjectInfo) bool {
	// Check dependencies for database-related packages
	for name := range info.Dependencies {
		switch name {
		case "pg", "postgres", "postgresql", "sequelize", "typeorm", "prisma",
			"mongoose", "mongodb", "mysql", "mysql2", "sqlite3", "redis":
			return true
		}
	}
	return false
}
