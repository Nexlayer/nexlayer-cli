package initcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/template"
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

const (
	// Stack types for modern web applications
	StackMERN          = "mern"           // MongoDB, Express, React, Node.js
	StackMEAN          = "mean"           // MongoDB, Express, Angular, Node.js
	StackMEVN          = "mevn"           // MongoDB, Express, Vue.js, Node.js
	StackPERN          = "pern"           // PostgreSQL, Express, React, Node.js
	StackNextJS        = "nextjs"         // Next.js
	
	// Stack types for ML/AI applications
	StackKubeflow      = "kubeflow"       // Kubeflow ML platform
	StackMLflow        = "mlflow"         // MLflow ML platform
	StackHuggingface   = "huggingface"    // Hugging Face
	StackMNFA          = "mnfa"           // MongoDB, Node.js, FastAPI
	StackPDN           = "pdn"            // PostgreSQL, Django, Node.js
	
	// Stack types for LLM applications
	StackLlamaNode     = "llama-node"     // Llama, Node.js
	StackLlamaPy       = "llama-py"       // Llama, Python
	StackLlamaJS       = "llama-js"       // Llama, Next.js
	StackAnthropicNode = "anthropic-node" // Anthropic Claude, Node.js
	StackAnthropicJS   = "anthropic-js"   // Anthropic Claude, Next.js
	StackAnthropicPy   = "anthropic-py"   // Anthropic Claude, Python
	StackLangChainJS   = "langchain-js"   // LangChain, Node.js
	StackLangChainPy   = "langchain-py"   // LangChain, Python
	StackOpenAINode    = "openai-node"    // OpenAI, Node.js
	StackOpenAIPy      = "openai-py"      // OpenAI, Python
	StackVertexAI      = "vertex-ai"      // Google Vertex AI
)

// Use type aliases to ensure consistency
type Config = types.Config
type PodConfig = types.PodConfig
type VarPair = types.VarPair
type BuildConfig = types.BuildConfig
type RegistryAuth = types.RegistryAuth

type Dependency struct {
	Name     string
	Type     string
	Required bool
}

type ServiceDependency struct {
	Name     string
	Type     string
	Image    string
	Required bool
}

func addTemplateConfig(config *Config, templateName string, pods []PodConfig) {
	config.Application.Template.Name = templateName

	for _, pod := range pods {
		addPod(config, pod.Type, pod.Name, pod.ExposeHttp, pod.Vars)
	}
}

func addPod(config *Config, podType string, name string, exposeHttp bool, vars []VarPair) {
	config.Application.Template.Pods = append(config.Application.Template.Pods, types.PodConfig{
		Type:       podType,
		Name:       name,
		ExposeHttp: exposeHttp,
		Vars:       vars,
	})
}

func getPodTypeFromService(name string, image string) string {
	// Database types
	if strings.Contains(image, "postgres") || strings.Contains(name, "postgres") {
		return "postgres"
	}
	if strings.Contains(image, "mysql") || strings.Contains(name, "mysql") {
		return "mysql"
	}
	if strings.Contains(image, "neo4j") || strings.Contains(name, "neo4j") {
		return "neo4j"
	}
	if strings.Contains(image, "redis") || strings.Contains(name, "redis") {
		return "redis"
	}
	if strings.Contains(image, "mongo") || strings.Contains(name, "mongo") {
		return "mongo"
	}
	if strings.Contains(image, "pinecone") || strings.Contains(name, "pinecone") {
		return "pinecone"
	}

	// Frontend types
	if strings.Contains(image, "react") || strings.Contains(name, "react") {
		return "react"
	}
	if strings.Contains(image, "angular") || strings.Contains(name, "angular") {
		return "angular"
	}
	if strings.Contains(image, "vue") || strings.Contains(name, "vue") {
		return "vue"
	}

	// Backend types
	if strings.Contains(image, "django") || strings.Contains(name, "django") {
		return "django"
	}
	if strings.Contains(image, "fastapi") || strings.Contains(name, "fastapi") {
		return "fastapi"
	}
	if strings.Contains(image, "express") || strings.Contains(name, "express") {
		return "express"
	}

	// Other types
	if strings.Contains(image, "nginx") || strings.Contains(name, "nginx") {
		return "nginx"
	}
	if strings.Contains(image, "llm") || strings.Contains(name, "llm") {
		return "llm"
	}

	return ""
}

func getPodTypeFromDependency(name string) string {
	switch name {
	// Database dependencies
	case "pg", "postgres", "postgresql":
		return "postgres"
	case "mysql", "mysql2":
		return "mysql"
	case "neo4j-driver":
		return "neo4j"
	case "redis":
		return "redis"
	case "mongodb", "mongoose":
		return "mongo"
	case "pinecone-client":
		return "pinecone"

	// Frontend dependencies
	case "react", "react-dom":
		return "react"
	case "@angular/core":
		return "angular"
	case "vue":
		return "vue"

	// Backend dependencies
	case "django":
		return "django"
	case "fastapi":
		return "fastapi"
	case "express":
		return "express"

	// Other dependencies
	case "nginx":
		return "nginx"
	case "@langchain/core", "langchain":
		return "llm"
	}

	return ""
}

func getDependencyType(name string) string {
	return getPodTypeFromDependency(name)
}

func detectServiceDependencies(dir string) []Dependency {
	var deps []Dependency

	// Check package.json for Node.js dependencies
	if pkgJson, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(pkgJson, &pkg); err == nil {
			for name := range pkg.Dependencies {
				if depType := getDependencyType(name); depType != "" {
					deps = append(deps, Dependency{
						Name:     name,
						Type:     depType,
						Required: true,
					})
				}
			}
		}
	}

	// Check requirements.txt for Python dependencies
	if reqsTxt, err := os.ReadFile(filepath.Join(dir, "requirements.txt")); err == nil {
		for _, line := range strings.Split(string(reqsTxt), "\n") {
			name := strings.Split(strings.TrimSpace(line), "==")[0]
			if depType := getDependencyType(name); depType != "" {
				deps = append(deps, Dependency{
					Name:     name,
					Type:     depType,
					Required: true,
				})
			}
		}
	}

	return deps
}

func createDefaultConfig(projectName, stackType string, deps []Dependency) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = projectName
	config.Application.Template.RegistryLogin = defaultRegistryLogin()

	// Add pods based on detected dependencies
	for _, dep := range deps {
		switch dep.Type {
		case "frontend":
			addPod(&config, "react", "frontend", true, []types.VarPair{
				{Key: "NODE_ENV", Value: "development"},
				{Key: "PORT", Value: "3000"},
			})
		case "backend":
			addPod(&config, "express", "backend", true, []types.VarPair{
				{Key: "NODE_ENV", Value: "development"},
				{Key: "PORT", Value: "5000"},
			})
		case "database":
			if strings.Contains(dep.Name, "mongodb") {
				addPod(&config, "mongodb", "mongodb", false, []types.VarPair{
					{Key: "MONGODB_URI", Value: "mongodb://mongodb:27017/myapp"},
				})
			} else if strings.Contains(dep.Name, "postgres") {
				addPod(&config, "postgres", "postgres", false, []types.VarPair{
					{Key: "POSTGRES_DB", Value: "myapp"},
					{Key: "POSTGRES_USER", Value: "postgres"},
					{Key: "POSTGRES_PASSWORD", Value: "postgres"},
				})
			}
		}
	}

	return config
}

func createConfig(projectName, templateName string) (Config, error) {
	cfg, exists := template.GetTemplate(templateName)
	if !exists {
		return Config{}, fmt.Errorf("unknown template: %s", templateName)
	}

	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = projectName
	config.Application.Template.RegistryLogin = defaultRegistryLogin()

	// Convert template.VarPair to types.VarPair
	for _, pod := range cfg.DefaultPods {
		vars := make([]types.VarPair, len(pod.Vars))
		for i, v := range pod.Vars {
			vars[i] = types.VarPair{
				Key:   v.Key,
				Value: v.Value,
			}
		}
		addPod(&config, pod.Type, pod.Name, pod.ExposeHttp, vars)
	}

	return config, nil
}

func createLlamaPyConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackLlamaPy)
}

func createHuggingFaceConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackHuggingface)
}

func defaultRegistryLogin() RegistryAuth {
	return RegistryAuth{
		Registry:            "ghcr.io",
		Username:           "<Github username>",
		PersonalAccessToken: "<Github Packages Read-Only PAT>",
	}
}

func createKubeflowConfig(projectName string) (Config, error) {
	cfg := Config{}
	cfg.Application.Template.Name = projectName
	cfg.Application.Template.DeploymentName = projectName
	cfg.Application.Template.RegistryLogin = defaultRegistryLogin()

	// Add Kubeflow-specific pods
	addPod(&cfg, "kubeflow", "notebook", true, []types.VarPair{
		{Key: "JUPYTER_ENABLE_LAB", Value: "yes"},
		{Key: "NB_PREFIX", Value: "/"},
	})

	return cfg, nil
}

func createLlamaNodeConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackLlamaNode)
}

func createLlamaJSConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackLlamaJS)
}

func createAnthropicNodeConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackAnthropicNode)
}

func createAnthropicJSConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackAnthropicJS)
}

func createMernConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackMERN)
}

func createNextjsConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackNextJS)
}

func createMLConfig(projectName string, stackType string) (Config, error) {
	var config Config
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = projectName

	switch stackType {
	case StackKubeflow:
		config.Application.Template.Pods = []PodConfig{
			{
				Type: "llm",
				Name: "ml-pipeline",
				Tag:  "python:3.11-slim",
				Vars: []VarPair{
					{Key: "PIPELINE_ROOT", Value: "/tmp/pipeline"},
					{Key: "DATA_PATH", Value: "/tmp/data"},
					{Key: "MODEL_PATH", Value: "/tmp/model"},
					{Key: "KUBEFLOW_URL", Value: "http://localhost:8080"},
				},
				ExposeHttp: true,
			},
		}
	case StackMLflow:
		config.Application.Template.Pods = []PodConfig{
			{
				Type: "llm",
				Name: "mlflow-server",
				Tag:  "mlflow:latest",
				Vars: []types.VarPair{
					{Key: "MLFLOW_TRACKING_URI", Value: "http://localhost:5000"},
					{Key: "MLFLOW_S3_ENDPOINT_URL", Value: "http://minio:9000"},
					{Key: "AWS_ACCESS_KEY_ID", Value: "minioadmin"},
					{Key: "AWS_SECRET_ACCESS_KEY", Value: "minioadmin"},
				},
				ExposeHttp: true,
			},
			{
				Type: "database",
				Name: "minio",
				Tag:  "minio/minio:latest",
				Vars: []types.VarPair{
					{Key: "MINIO_ROOT_USER", Value: "minioadmin"},
					{Key: "MINIO_ROOT_PASSWORD", Value: "minioadmin"},
				},
			},
		}
	}

	return config, nil
}

func createMEANConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackMEAN)
}

func createMEVNConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackMEVN)
}

func createPERNConfig(projectName string) (Config, error) {
	return createConfig(projectName, StackPERN)
}

func setBuildConfig(config *Config, stackType string) {
	// Set build configuration based on stack type
	switch {
	case strings.Contains(stackType, "node"):
		config.Application.Template.Build.Command = "npm install && npm run build"
		config.Application.Template.Build.Output = "build"
	case strings.Contains(stackType, "py"):
		config.Application.Template.Build.Command = "pip install -r requirements.txt"
		config.Application.Template.Build.Output = "dist"
	default:
		config.Application.Template.Build.Command = "npm install && npm run build"
		config.Application.Template.Build.Output = "build"
	}
}

// NewCommand creates a new init command
func NewCommand() *cobra.Command {
	var templateFlag string

	cmd := &cobra.Command{
		Use:   "init [project-name]",
		Short: "Initialize a new Nexlayer project",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) > 0 {
				projectName = args[0]
			} else {
				// Use current directory name as project name
				dir, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get current directory: %w", err)
				}
				projectName = filepath.Base(dir)
			}

			// If no template specified, show interactive prompt
			if templateFlag == "" {
				selectedTemplate, _ := pterm.DefaultInteractiveSelect.
					WithOptions([]string{
						// Web Templates
						"MERN - MongoDB, Express, React, Node.js",
						"MEAN - MongoDB, Express, Angular, Node.js",
						"MEVN - MongoDB, Express, Vue.js, Node.js",
						"PERN - PostgreSQL, Express, React, Node.js",
						// ML Templates
						"Kubeflow - ML Pipeline with Kubeflow",
						"MLflow - MLflow with tracking server",
						// AI Templates
						"OpenAI Node.js - OpenAI with Express and React",
						"OpenAI Python - OpenAI with FastAPI and Vue",
					}).
					WithDefaultText("Select a template:").
					Show()

				// Extract template ID from selection
				templateFlag = strings.Split(selectedTemplate, " - ")[0]
				templateFlag = strings.ToLower(templateFlag)
			}

			// Create progress bar
			progress, _ := pterm.DefaultProgressbar.WithTotal(100).Start()
			progress.Title = "Initializing project"

			var config Config

			// Create config based on template type
			switch templateFlag {
			case "mern":
				config, _ = createMernConfig(projectName)
			case "mean":
				config, _ = createMEANConfig(projectName)
			case "mevn":
				config, _ = createMEVNConfig(projectName)
			case "pern":
				config, _ = createPERNConfig(projectName)
			case "kubeflow":
				config, _ = createKubeflowConfig(projectName)
			case "mlflow":
				config, _ = createMLConfig(projectName, templateFlag)
			case "openai-node":
				config, _ = createLlamaNodeConfig(projectName)
			case "openai-py":
				config, _ = createLlamaPyConfig(projectName)
			case "huggingface":
				config, _ = createHuggingFaceConfig(projectName)
			default:
				// For unknown templates, detect project type and dependencies
				deps := detectServiceDependencies(".")
				stackType := detectStackType(deps)
				config = createDefaultConfig(projectName, stackType, deps)
			}

			progress.Add(60)

			// Write config file
			yamlFileName := config.Application.Template.Name + ".yaml"
			err := writeConfig(config, yamlFileName)
			if err != nil {
				cmd.SilenceUsage = true
				return fmt.Errorf("failed to write config file: %w", err)
			}

			progress.Add(40)

			fmt.Printf("\nSuccessfully created %s with %s template!\n\n", yamlFileName, templateFlag)
			fmt.Println("To deploy your application, run: nexlayer deploy")

			return nil
		},
	}

	cmd.Flags().StringVarP(&templateFlag, "template", "t", "", "Template to use (e.g., mern, kubeflow)")

	return cmd
}

func detectProjectName(dir string, projectType string) (string, error) {
	switch projectType {
	case "nodejs":
		// Try to get name from package.json
		if _, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
			var packageJSON struct {
				Name string `json:"name"`
			}
			data, err := os.ReadFile(filepath.Join(dir, "package.json"))
			if err == nil {
				if err := json.Unmarshal(data, &packageJSON); err == nil && packageJSON.Name != "" {
					return packageJSON.Name, nil
				}
			}
		}
	case "python":
		// Try to get name from pyproject.toml
		if _, err := os.ReadFile(filepath.Join(dir, "pyproject.toml")); err == nil {
			data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml"))
			if err == nil {
				content := string(data)
				for _, line := range strings.Split(content, "\n") {
					if strings.HasPrefix(line, "name") {
						parts := strings.Split(line, "=")
						if len(parts) == 2 {
							name := strings.TrimSpace(parts[1])
							name = strings.Trim(name, "\"'")
							if name != "" {
								return name, nil
							}
						}
					}
				}
			}
		}
	}

	// If no name found, use directory name
	return filepath.Base(dir), nil
}

func detectPort(projectType string) int {
	switch projectType {
	case "nodejs":
		return 3000
	case "python":
		return 8000
	case "golang":
		return 8080
	case "static":
		return 80
	case StackMERN:
		return 3000
	case StackMEAN:
		return 3000
	case StackMEVN:
		return 3000
	case StackPERN:
		return 3000
	case StackMNFA:
		return 8000
	case StackPDN:
		return 8000
	case StackLangChainJS:
		return 3000
	case StackLangChainPy:
		return 8000
	case StackOpenAINode:
		return 3000
	case StackOpenAIPy:
		return 8000
	case StackVertexAI:
		return 8080
	case StackHuggingface:
		return 8000
	case StackAnthropicPy:
		return 8000
	case StackAnthropicJS:
		return 3000
	default:
		return 8080
	}
}

func writeConfig(config Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func detectProjectType(dir string) (string, []ServiceDependency) {
	// Detect service dependencies from docker-compose.yml
	deps := detectServiceDependencies(filepath.Join(dir, "docker-compose.yml"))
	
	// Convert Dependency to ServiceDependency
	serviceDeps := make([]ServiceDependency, len(deps))
	for i, dep := range deps {
		serviceDeps[i] = ServiceDependency{
			Name:     dep.Name,
			Type:     dep.Type,
			Required: dep.Required,
		}
	}

	// Check for package.json
	if _, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		return StackOpenAINode, serviceDeps
	}

	// Check for pyproject.toml
	if _, err := os.ReadFile(filepath.Join(dir, "pyproject.toml")); err == nil {
		return StackOpenAIPy, serviceDeps
	}

	// Default to Node.js if no specific markers found
	return StackOpenAINode, serviceDeps
}

type PackageJSON struct {
	Name string `json:"name"`
}

type PyProject struct {
	Project struct {
		Name string `toml:"name"`
	} `toml:"project"`
}

type DockerCompose struct {
	Services map[string]struct {
		Image       string            `yaml:"image"`
		Build       string            `yaml:"build"`
		Environment []string          `yaml:"environment"`
		Env         map[string]string `yaml:"env"`
	} `yaml:"services"`
}

func detectStackType(deps []Dependency) string {
	hasReact := false
	hasExpress := false
	hasMongo := false
	hasPostgres := false
	hasAngular := false
	hasVue := false
	hasNextJS := false
	hasLlama := false
	hasAnthropic := false

	for _, dep := range deps {
		switch dep.Type {
		case "react", "react-dom":
			hasReact = true
		case "express":
			hasExpress = true
		case "mongodb", "mongoose":
			hasMongo = true
		case "pg", "postgres":
			hasPostgres = true
		case "@angular/core":
			hasAngular = true
		case "vue":
			hasVue = true
		case "next":
			hasNextJS = true
		case "llama":
			hasLlama = true
		case "@anthropic-ai/sdk":
			hasAnthropic = true
		}
	}

	// Determine stack type based on combinations
	switch {
	case hasReact && hasExpress && hasMongo:
		return StackMERN
	case hasAngular && hasExpress && hasMongo:
		return StackMEAN
	case hasVue && hasExpress && hasMongo:
		return StackMEVN
	case hasReact && hasExpress && hasPostgres:
		return StackPERN
	case hasNextJS:
		return StackNextJS
	case hasLlama && hasNextJS:
		return StackLlamaJS
	case hasLlama:
		return StackLlamaNode
	case hasAnthropic && hasNextJS:
		return StackAnthropicJS
	case hasAnthropic:
		return StackAnthropicNode
	default:
		return ""
	}
}
