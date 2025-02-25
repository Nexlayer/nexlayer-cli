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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
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

			return runInitCommand(opts)
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
func runInitCommand(opts *InitOptions) error {
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
	// Convert template.NexlayerYAML to schema.NexlayerYAML
	schemaConfig := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: config.Application.Name,
			URL:  config.Application.URL,
		},
	}

	// Convert registry login if present
	if config.Application.RegistryLogin != nil {
		schemaConfig.Application.RegistryLogin = &schema.RegistryLogin{
			Registry:            config.Application.RegistryLogin.Registry,
			Username:            config.Application.RegistryLogin.Username,
			PersonalAccessToken: config.Application.RegistryLogin.PersonalAccessToken,
		}
	}

	// Convert pods
	for _, pod := range config.Application.Pods {
		schemaPod := schema.Pod{
			Name:        pod.Name,
			Type:        pod.Type,
			Path:        pod.Path,
			Image:       pod.Image,
			Command:     pod.Command,
			Entrypoint:  pod.Entrypoint,
			Annotations: pod.Annotations,
		}

		// Convert service ports
		for _, port := range pod.ServicePorts {
			// Ensure port name is set (required in schema.ServicePort)
			portName := port.Name
			if portName == "" {
				portName = "http" // Default name if not specified
			}

			schemaPod.ServicePorts = append(schemaPod.ServicePorts, schema.ServicePort{
				Name:       portName,
				Port:       port.Port,
				TargetPort: port.TargetPort,
				Protocol:   port.Protocol,
			})
		}

		// Convert environment variables
		for _, envVar := range pod.Vars {
			schemaPod.Vars = append(schemaPod.Vars, schema.EnvVar{
				Key:   envVar.Key,
				Value: envVar.Value,
			})
		}

		// Convert volumes
		for _, volume := range pod.Volumes {
			schemaPod.Volumes = append(schemaPod.Volumes, schema.Volume{
				Name:     volume.Name,
				Path:     volume.Path,
				Size:     volume.Size,
				Type:     volume.Type,
				ReadOnly: volume.ReadOnly,
			})
		}

		// Convert secrets
		for _, secret := range pod.Secrets {
			schemaPod.Secrets = append(schemaPod.Secrets, schema.Secret{
				Name:     secret.Name,
				Data:     secret.Data,
				Path:     secret.Path,
				FileName: secret.FileName,
			})
		}

		schemaConfig.Application.Pods = append(schemaConfig.Application.Pods, schemaPod)
	}

	// Validate using schema validator
	validator := schema.NewValidator(true)
	errors := validator.ValidateYAML(schemaConfig)
	if len(errors) > 0 {
		report := schema.NewValidationReport()
		report.AddErrors(errors)
		return fmt.Errorf("validation failed:\n%s", report.String())
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
func printSuccessMessage(info *types.ProjectInfo, _ *template.NexlayerYAML) {
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
func getDefaultServiceURL(name, _ string) string {
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
