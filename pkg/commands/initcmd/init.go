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
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	cacheDir         = ".nexlayer"
	cacheFile        = "detection-cache.json"
	podRefPattern    = `([a-z][a-z0-9-]*).pod`
	urlRefPattern    = `<% URL %>`
	envVarRefPattern = `<%\s*([A-Z_][A-Z0-9_]*)\s*%>`
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

// NewCommand creates a new init command
func NewCommand() *cobra.Command {
	var (
		interactive bool
		force       bool
		appName     string
		podName     string
		podImage    string
		podPort     int
		podPath     string
	)

	cmd := &cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize a new Nexlayer project",
		Long: `Initialize a new Nexlayer project by creating a nexlayer.yaml file.
The command will auto-detect your project type and configure it appropriately.

Examples:
  # Auto-detect and initialize in current directory
  nexlayer init

  # Initialize with a custom name
  nexlayer init --name my-app

  # Interactive mode
  nexlayer init --interactive

  # Force re-detection (ignore cache)
  nexlayer init --force

Required Fields in nexlayer.yaml:
  - application.name: The name of the application
  - pods[].name: The pod name (e.g., "web" or "api")
  - pods[].image: The container image (e.g., "nginx:latest")
  - pods[].servicePorts: List of ports to expose (e.g., [3000])
  - pods[].path: Only for forward-facing pods (e.g., "/")

Optional Fields (included when needed):
  - volumes: For database pods (mountPath, size)
  - vars: For environment variables (AI, database configs)
  - registryLogin: For private images (registry, username, password)`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get target directory
			dir := "."
			if len(args) > 0 {
				dir = args[0]
			}

			// Create InitOptions
			opts := &InitOptions{
				Directory:   dir,
				Interactive: interactive,
				Force:       force,
				AppName:     appName,
				PodName:     podName,
				PodImage:    podImage,
				PodPort:     podPort,
				PodPath:     podPath,
			}

			return runInitCommand(cmd, opts)
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Enable interactive mode")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force re-detection (ignore cache)")
	cmd.Flags().StringVar(&appName, "name", "", "Application name (default: directory name)")
	cmd.Flags().StringVar(&podName, "pod-name", "", "Main pod name (default: based on project type)")
	cmd.Flags().StringVar(&podImage, "pod-image", "", "Main pod image (default: based on project type)")
	cmd.Flags().IntVar(&podPort, "pod-port", 0, "Main pod port (default: based on project type)")
	cmd.Flags().StringVar(&podPath, "pod-path", "", "Main pod path (default: / for web/api pods)")

	return cmd
}

// InitOptions holds configuration for the init command
type InitOptions struct {
	Directory   string
	Interactive bool
	Force       bool
	AppName     string
	PodName     string
	PodImage    string
	PodPort     int
	PodPath     string
}

// runInitCommand handles the execution of the init command
func runInitCommand(cmd *cobra.Command, opts *InitOptions) error {
	// Show welcome message
	fmt.Println(infoStyle.Render("🚀 Initializing Nexlayer project..."))

	// Try to load from cache first
	var info *types.ProjectInfo
	if !opts.Force {
		info = loadFromCache(opts.Directory)
	}

	// If not in cache or force flag is set, detect project
	if info == nil {
		var err error
		info, err = detectProjectParallel(opts.Directory)
		if err != nil && opts.Interactive {
			// If detection fails in interactive mode, prompt user
			info, err = promptForProjectType(opts.Directory)
		}
		if err != nil {
			return fmt.Errorf("failed to detect project type: %w", err)
		}

		// Save to cache
		if err := saveToCache(opts.Directory, info); err != nil {
			fmt.Println(warningStyle.Render("⚠️  Warning: Failed to cache detection results"))
		}
	}

	// Apply user overrides
	if err := applyUserOverrides(info, opts); err != nil {
		return fmt.Errorf("failed to apply overrides: %w", err)
	}

	// Generate configuration
	config, err := generateConfiguration(info, opts)
	if err != nil {
		return fmt.Errorf("failed to generate configuration: %w", err)
	}

	// Validate configuration
	if err := validateConfiguration(config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	// Write configuration
	if err := writeYAMLToFile(filepath.Join(opts.Directory, "nexlayer.yaml"), config); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	// Print success message
	printSuccessMessage(info, config)

	return nil
}

// applyUserOverrides applies user-provided overrides to the project info
func applyUserOverrides(info *types.ProjectInfo, opts *InitOptions) error {
	if opts.AppName != "" {
		info.Name = opts.AppName
	}

	// If in interactive mode, prompt for confirmation/changes
	if opts.Interactive {
		if err := promptForOverrides(info, opts); err != nil {
			return err
		}
	}

	return nil
}

// generateConfiguration creates a minimal but complete nexlayer.yaml configuration
func generateConfiguration(info *types.ProjectInfo, opts *InitOptions) (*template.NexlayerYAML, error) {
	// Create base configuration
	config := &template.NexlayerYAML{
		Application: template.Application{
			Name: info.Name,
			Pods: []template.Pod{},
		},
	}

	// Add main pod
	mainPod := generateMainPod(info, opts)
	config.Application.Pods = append(config.Application.Pods, mainPod)

	// Add database if needed
	if hasDatabase(info) {
		dbPod := generateDatabasePod(info)
		config.Application.Pods = append(config.Application.Pods, dbPod)
	}

	// Add AI configurations if detected
	if info.LLMProvider != "" {
		addAIConfigurations(config, info)
	}

	return config, nil
}

// generateMainPod creates the main pod configuration based on project type
func generateMainPod(info *types.ProjectInfo, opts *InitOptions) template.Pod {
	pod := template.Pod{
		Name: opts.PodName,
		Type: string(info.Type),
	}

	// Set defaults based on project type if not overridden
	if pod.Name == "" {
		switch info.Type {
		case types.TypeNextjs, types.TypeReact:
			pod.Name = "web"
		case types.TypeNode, types.TypePython, types.TypeGo:
			pod.Name = "api"
		default:
			pod.Name = "app"
		}
	}

	// Set image based on project type if not overridden
	if opts.PodImage != "" {
		pod.Image = opts.PodImage
	} else {
		pod.Image = getDefaultImage(info.Type)
	}

	// Set port based on project type if not overridden
	port := opts.PodPort
	if port == 0 {
		port = info.Port
	}
	pod.ServicePorts = []template.ServicePort{
		{Name: "http", Port: port, TargetPort: port},
	}

	// Set path for web/api pods
	if opts.PodPath != "" {
		pod.Path = opts.PodPath
	} else if isWebOrAPI(info.Type) {
		pod.Path = "/"
	}

	// Add environment variables for service dependencies
	pod.Vars = generateEnvironmentVars(info)

	return pod
}

// generateEnvironmentVars creates environment variables with pod references
func generateEnvironmentVars(info *types.ProjectInfo) []template.EnvVar {
	var vars []template.EnvVar

	// Add base URL if needed
	if isWebOrAPI(info.Type) {
		vars = append(vars, template.EnvVar{
			Key:   "BASE_URL",
			Value: "<% URL %>",
		})
	}

	// Add service URLs based on dependencies
	for name := range info.Dependencies {
		switch {
		case strings.Contains(name, "postgres"):
			vars = append(vars, template.EnvVar{
				Key:   "DATABASE_URL",
				Value: "postgresql://postgres:<% DB_PASSWORD %>@postgres.pod:5432/app",
			})
		case strings.Contains(name, "mongodb"):
			vars = append(vars, template.EnvVar{
				Key:   "MONGODB_URI",
				Value: "mongodb://root:<% MONGO_ROOT_PASSWORD %>@mongodb.pod:27017/app",
			})
		case strings.Contains(name, "mysql"):
			vars = append(vars, template.EnvVar{
				Key:   "MYSQL_URL",
				Value: "mysql://root:<% MYSQL_ROOT_PASSWORD %>@mysql.pod:3306/app",
			})
		case strings.Contains(name, "redis"):
			vars = append(vars, template.EnvVar{
				Key:   "REDIS_URL",
				Value: "redis://:<% REDIS_PASSWORD %>@redis.pod:6379",
			})
		case strings.Contains(name, "ai-model"):
			vars = append(vars, template.EnvVar{
				Key:   "AI_MODEL_URL",
				Value: "http://ai-model.pod:5000",
			})
		case strings.Contains(name, "vector-db"):
			vars = append(vars, template.EnvVar{
				Key:   "VECTOR_DB_URL",
				Value: "http://vector-db.pod:8080",
			})
		case strings.Contains(name, "minio"):
			vars = append(vars, []template.EnvVar{
				{Key: "MINIO_ENDPOINT", Value: "minio.pod:9000"},
				{Key: "MINIO_ACCESS_KEY", Value: "<% MINIO_ACCESS_KEY %>"},
				{Key: "MINIO_SECRET_KEY", Value: "<% MINIO_SECRET_KEY %>"},
			}...)
		}
	}

	// Add AI-specific environment variables if needed
	if info.LLMProvider != "" {
		vars = append(vars, []template.EnvVar{
			{Key: "LLM_PROVIDER", Value: info.LLMProvider},
			{Key: "LLM_MODEL", Value: info.LLMModel},
			{Key: "LLM_API_KEY", Value: "<% LLM_API_KEY %>"},
		}...)
	}

	return vars
}

// generateDatabasePod creates a database pod configuration
func generateDatabasePod(info *types.ProjectInfo) template.Pod {
	dbType := detectDatabaseType(info)
	pod := template.Pod{
		Name:  fmt.Sprintf("db-%s", dbType),
		Type:  dbType,
		Image: fmt.Sprintf("%s:latest", dbType),
		ServicePorts: []template.ServicePort{
			{Name: "db", Port: getDefaultDBPort(dbType), TargetPort: getDefaultDBPort(dbType)},
		},
		Volumes: []template.Volume{
			{
				Name: fmt.Sprintf("%s-data", dbType),
				Path: getDefaultDBPath(dbType),
				Size: "5Gi",
			},
		},
	}

	// Add default environment variables
	pod.Vars = getDefaultDBVars(dbType)

	// Add health check environment variables
	pod.Vars = append(pod.Vars, template.EnvVar{
		Key:   "POD_NAME",
		Value: fmt.Sprintf("%s.pod", pod.Name),
	})

	return pod
}

// validateConfiguration ensures the configuration is valid
func validateConfiguration(config *template.NexlayerYAML) error {
	if config.Application.Name == "" {
		return fmt.Errorf("application name is required")
	}

	if len(config.Application.Pods) == 0 {
		return fmt.Errorf("at least one pod is required")
	}

	// Validate individual pods
	for i, pod := range config.Application.Pods {
		if err := validatePod(pod, i); err != nil {
			return err
		}
	}

	// Validate pod references in environment variables
	if errors := validatePodReferences(config); len(errors) > 0 {
		var errMsg strings.Builder
		errMsg.WriteString("Invalid pod references found:\n")
		for _, err := range errors {
			errMsg.WriteString(fmt.Sprintf("- %s: %s\n", err.Field, err.Message))
			for _, suggestion := range err.Suggestions {
				errMsg.WriteString(fmt.Sprintf("  %s\n", suggestion))
			}
		}
		return fmt.Errorf(errMsg.String())
	}

	return nil
}

// validatePod validates a single pod configuration
func validatePod(pod template.Pod, index int) error {
	if pod.Name == "" {
		return fmt.Errorf("pod[%d]: name is required", index)
	}

	if !isValidPodName(pod.Name) {
		return fmt.Errorf("pod[%d]: invalid name '%s' (must start with lowercase letter, contain only alphanumeric characters, '-', or '.')", index, pod.Name)
	}

	if pod.Image == "" {
		return fmt.Errorf("pod[%d]: image is required", index)
	}

	if len(pod.ServicePorts) == 0 {
		return fmt.Errorf("pod[%d]: at least one service port is required", index)
	}

	for _, port := range pod.ServicePorts {
		if port.Port < 1 || port.Port > 65535 {
			return fmt.Errorf("pod[%d]: invalid port %d (must be between 1 and 65535)", index, port.Port)
		}
	}

	for _, volume := range pod.Volumes {
		if !strings.HasPrefix(volume.Path, "/") {
			return fmt.Errorf("pod[%d]: volume path '%s' must start with '/'", index, volume.Path)
		}
	}

	return nil
}

// Helper functions for default values and validation

func getDefaultImage(projectType types.ProjectType) string {
	switch projectType {
	case types.TypeNextjs:
		return "node:18-alpine"
	case types.TypeReact:
		return "nginx:alpine"
	case types.TypeNode:
		return "node:18-alpine"
	case types.TypePython:
		return "python:3.9-slim"
	case types.TypeGo:
		return "golang:1.23-alpine"
	default:
		return "alpine:latest"
	}
}

func isWebOrAPI(projectType types.ProjectType) bool {
	switch projectType {
	case types.TypeNextjs, types.TypeReact, types.TypeNode, types.TypePython, types.TypeGo:
		return true
	default:
		return false
	}
}

func detectDatabaseType(info *types.ProjectInfo) string {
	for name := range info.Dependencies {
		switch {
		case strings.Contains(name, "pg"), strings.Contains(name, "postgres"):
			return "postgres"
		case strings.Contains(name, "mongodb"), strings.Contains(name, "mongoose"):
			return "mongodb"
		case strings.Contains(name, "mysql"):
			return "mysql"
		case strings.Contains(name, "redis"):
			return "redis"
		}
	}
	return "postgres" // Default to PostgreSQL
}

func getDefaultDBPort(dbType string) int {
	switch dbType {
	case "postgres":
		return 5432
	case "mongodb":
		return 27017
	case "mysql":
		return 3306
	case "redis":
		return 6379
	default:
		return 5432
	}
}

func getDefaultDBPath(dbType string) string {
	switch dbType {
	case "postgres":
		return "/var/lib/postgresql/data"
	case "mongodb":
		return "/data/db"
	case "mysql":
		return "/var/lib/mysql"
	case "redis":
		return "/data"
	default:
		return "/data"
	}
}

func getDefaultDBVars(dbType string) []template.EnvVar {
	switch dbType {
	case "postgres":
		return []template.EnvVar{
			{Key: "POSTGRES_USER", Value: "postgres"},
			{Key: "POSTGRES_PASSWORD", Value: "<% DB_PASSWORD %>"},
			{Key: "POSTGRES_DB", Value: "app"},
		}
	case "mongodb":
		return []template.EnvVar{
			{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
			{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "<% MONGO_ROOT_PASSWORD %>"},
		}
	case "mysql":
		return []template.EnvVar{
			{Key: "MYSQL_ROOT_PASSWORD", Value: "<% MYSQL_ROOT_PASSWORD %>"},
			{Key: "MYSQL_DATABASE", Value: "app"},
		}
	case "redis":
		return []template.EnvVar{
			{Key: "REDIS_PASSWORD", Value: "<% REDIS_PASSWORD %>"},
		}
	default:
		return nil
	}
}

func isValidPodName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if name[0] < 'a' || name[0] > 'z' {
		return false
	}
	for _, c := range name {
		if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' || c == '.') {
			return false
		}
	}
	return true
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
func promptForProjectType(dir string) (*types.ProjectInfo, error) {
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
		Name: filepath.Base(dir),
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

// promptForOverrides prompts the user to confirm or modify detected settings
func promptForOverrides(info *types.ProjectInfo, opts *InitOptions) error {
	// Confirm application name
	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Application name [%s]", info.Name),
		Default:   info.Name,
		AllowEdit: true,
	}
	if result, err := prompt.Run(); err != nil {
		if err != promptui.ErrInterrupt {
			return fmt.Errorf("prompt failed: %w", err)
		}
	} else if result != "" && result != info.Name {
		info.Name = result
	}

	// Confirm project type
	typePrompt := promptui.Select{
		Label: "Project type",
		Items: []string{
			"Next.js",
			"React",
			"Node.js",
			"Python",
			"Go",
			"Docker",
		},
	}
	if _, result, err := typePrompt.Run(); err != nil {
		if err != promptui.ErrInterrupt {
			return fmt.Errorf("prompt failed: %w", err)
		}
	} else {
		switch result {
		case "Next.js":
			info.Type = types.TypeNextjs
		case "React":
			info.Type = types.TypeReact
		case "Node.js":
			info.Type = types.TypeNode
		case "Python":
			info.Type = types.TypePython
		case "Go":
			info.Type = types.TypeGo
		case "Docker":
			info.Type = types.TypeDockerRaw
		}
	}

	// Confirm port
	portPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Port [%d]", info.Port),
		Default:   fmt.Sprintf("%d", info.Port),
		AllowEdit: true,
		Validate: func(input string) error {
			if input == "" {
				return nil
			}
			port, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("port must be a number")
			}
			if port < 1 || port > 65535 {
				return fmt.Errorf("port must be between 1 and 65535")
			}
			return nil
		},
	}
	if result, err := portPrompt.Run(); err != nil {
		if err != promptui.ErrInterrupt {
			return fmt.Errorf("prompt failed: %w", err)
		}
	} else if result != "" {
		if port, err := strconv.Atoi(result); err == nil {
			info.Port = port
		}
	}

	// If database dependencies are detected, confirm database type
	if hasDatabase(info) {
		dbPrompt := promptui.Select{
			Label: "Database type",
			Items: []string{
				"PostgreSQL",
				"MongoDB",
				"MySQL",
				"Redis",
			},
		}
		if _, result, err := dbPrompt.Run(); err != nil {
			if err != promptui.ErrInterrupt {
				return fmt.Errorf("prompt failed: %w", err)
			}
		} else {
			// Store the selected database type for later use
			info.Dependencies[strings.ToLower(result)] = "latest"
		}
	}

	// Prompt for environment variables if needed
	if hasEnvironmentVars(info) {
		fmt.Println(infoStyle.Render("\nEnvironment Variables:"))
		fmt.Println("Available pod references:", strings.Join(getDefaultPodNames(info), ", "))
		fmt.Println("Use <pod-name>.pod to reference other pods (e.g., postgres.pod:5432)")
		fmt.Println("Use <% URL %> to reference the deployment's base URL")

		for name, value := range info.Dependencies {
			if isServiceDependency(name) {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("%s URL", name),
					Default:   getDefaultServiceURL(name, value),
					AllowEdit: true,
				}
				if result, err := prompt.Run(); err != nil {
					if err != promptui.ErrInterrupt {
						return fmt.Errorf("prompt failed: %w", err)
					}
				} else if result != "" {
					info.Dependencies[name] = result
				}
			}
		}
	}

	return nil
}

// hasEnvironmentVars checks if the project needs environment variables
func hasEnvironmentVars(info *types.ProjectInfo) bool {
	for name := range info.Dependencies {
		if isServiceDependency(name) {
			return true
		}
	}
	return false
}

// isServiceDependency checks if a dependency requires service connection
func isServiceDependency(name string) bool {
	services := []string{
		"postgres", "mongodb", "mysql", "redis",
		"ai-model", "vector-db", "minio",
	}
	for _, service := range services {
		if strings.Contains(name, service) {
			return true
		}
	}
	return false
}

// getDefaultServiceURL returns a default URL for a service
func getDefaultServiceURL(name, version string) string {
	switch {
	case strings.Contains(name, "postgres"):
		return "postgres.pod:5432"
	case strings.Contains(name, "mongodb"):
		return "mongodb.pod:27017"
	case strings.Contains(name, "mysql"):
		return "mysql.pod:3306"
	case strings.Contains(name, "redis"):
		return "redis.pod:6379"
	case strings.Contains(name, "ai-model"):
		return "ai-model.pod:5000"
	case strings.Contains(name, "vector-db"):
		return "vector-db.pod:8080"
	case strings.Contains(name, "minio"):
		return "minio.pod:9000"
	default:
		return fmt.Sprintf("%s.pod", name)
	}
}

// getDefaultPodNames returns a list of default pod names based on project type
func getDefaultPodNames(info *types.ProjectInfo) []string {
	pods := []string{"web", "api"}

	// Add database pods if needed
	if hasDatabase(info) {
		dbType := detectDatabaseType(info)
		pods = append(pods, fmt.Sprintf("db-%s", dbType))
	}

	// Add AI-specific pods if needed
	if info.LLMProvider != "" {
		pods = append(pods, "ai-model", "vector-db")
	}

	return pods
}

// validatePodReferences checks if all referenced pods exist
func validatePodReferences(config *template.NexlayerYAML) []ValidationError {
	var errors []ValidationError
	podNames := make(map[string]bool)

	// Build map of existing pod names
	for _, pod := range config.Application.Pods {
		podNames[pod.Name] = true
	}

	// Check each pod's environment variables
	for podIndex, pod := range config.Application.Pods {
		for _, envVar := range pod.Vars {
			refs := extractPodReferences(envVar.Value)
			for _, ref := range refs {
				if !podNames[ref] {
					suggestion := findClosestPodName(ref, podNames)
					err := ValidationError{
						Field:   fmt.Sprintf("pods[%d].vars[%s]", podIndex, envVar.Key),
						Message: fmt.Sprintf("referenced pod '%s' not found", ref),
					}
					if suggestion != "" {
						err.Suggestions = []string{
							fmt.Sprintf("Did you mean '%s'?", suggestion),
							fmt.Sprintf("Available pods: %s", strings.Join(getAvailablePods(podNames), ", ")),
						}
					}
					errors = append(errors, err)
				}
			}
		}
	}

	return errors
}

// extractPodReferences finds all pod references in a string
func extractPodReferences(value string) []string {
	re := regexp.MustCompile(podRefPattern)
	matches := re.FindAllStringSubmatch(value, -1)
	refs := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			refs = append(refs, match[1])
		}
	}
	return refs
}

// findClosestPodName finds the most similar pod name using Levenshtein distance
func findClosestPodName(ref string, podNames map[string]bool) string {
	minDist := 1000
	var closest string
	for name := range podNames {
		dist := levenshteinDistance(ref, name)
		if dist < minDist {
			minDist = dist
			closest = name
		}
	}
	if minDist <= len(ref)/2 {
		return closest
	}
	return ""
}

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// getAvailablePods returns a sorted list of pod names
func getAvailablePods(podNames map[string]bool) []string {
	pods := make([]string, 0, len(podNames))
	for name := range podNames {
		pods = append(pods, name)
	}
	sort.Strings(pods)
	return pods
}

// ValidationError type
type ValidationError struct {
	Field       string
	Message     string
	Suggestions []string
}
