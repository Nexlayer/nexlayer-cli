package initcmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/commands/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/templates"
	"github.com/Nexlayer/nexlayer-cli/pkg/types"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

const (
	// Frontend pod types
	PodTypeReact   = "react"
	PodTypeAngular = "angular"
	PodTypeVue     = "vue"

	// Backend pod types
	PodTypeExpress = "express"
	PodTypeDjango  = "django"
	PodTypeFastAPI = "fastapi"

	// Database pod types
	PodTypeMongoDB  = "mongodb"
	PodTypePostgres = "postgres"
	PodTypeMySQL    = "mysql"
	PodTypeNeo4j    = "neo4j"
	PodTypeRedis    = "redis"
	PodTypePinecone = "pinecone"

	// Other pod types
	PodTypeNginx = "nginx"
	PodTypeLLM   = "llm"

	// Stack types
	StackMERN        = "mern"        // MongoDB, Express, React, Node.js
	StackMEAN        = "mean"        // MongoDB, Express, Angular, Node.js
	StackMEVN        = "mevn"        // MongoDB, Express, Vue.js, Node.js
	StackPERN        = "pern"        // PostgreSQL, Express, React, Node.js
	StackMNFA        = "mnfa"        // MongoDB, Neo4j, FastAPI, Angular
	StackPDN         = "pdn"         // PostgreSQL, Django, Node.js

	// Stack types - ML
	StackKubeflow = "kubeflow" // Kubeflow ML Pipeline
	StackMLflow   = "mlflow"   // MLflow with tracking server

	// Stack types - AI/LLM
	StackLangChainJS = "langchain-nextjs"  // LangChain.js, Next.js, MongoDB
	StackLangChainPy = "langchain-fastapi" // LangChain Python, FastAPI, PostgreSQL
	StackOpenAINode  = "openai-node"       // OpenAI Node.js SDK, Express, React
	StackOpenAIPy    = "openai-py"         // OpenAI Python SDK, FastAPI, Vue
	StackLlamaNode   = "llama-node"        // Llama.cpp Node.js, Next.js, PostgreSQL
	StackLlamaPy     = "llama-py"          // Llama.cpp Python, FastAPI, MongoDB
	StackVertexAI    = "vertex-ai"         // Google Vertex AI, Flask, React
	StackHuggingface = "huggingface"       // Hugging Face Transformers, FastAPI, React
	StackAnthropicPy = "anthropic-py"      // Anthropic Claude, FastAPI, Svelte
	StackAnthropicJS = "anthropic-js"      // Anthropic Claude, Next.js, MongoDB
)

// Config is an alias for types.Config
type Config = types.Config

// RegistryAuth is an alias for types.RegistryAuth
type RegistryAuth = types.RegistryAuth

// PodConfig is an alias for types.PodConfig
type PodConfig = types.PodConfig

// VarPair is an alias for types.VarPair
type VarPair = types.VarPair

type DockerCompose struct {
	Services map[string]struct {
		Image       string            `yaml:"image"`
		Build       string            `yaml:"build"`
		Environment []string          `yaml:"environment"`
		Env         map[string]string `yaml:"env"`
	} `yaml:"services"`
}

func addTemplateConfig(config *Config, templateName string, pods []PodConfig) {
	config.Application.Template.Name = templateName

	for _, pod := range pods {
		addPod(config, pod.Type, pod.Name, pod.ExposeHttp, pod.Vars)
	}
}

func addPod(config *Config, podType string, name string, exposeHttp bool, vars []VarPair) {
	pod := PodConfig{
		Type:       podType,
		Name:       name,
		Tag:        fmt.Sprintf("ghcr.io/your-username/%s:latest", name),
		ExposeHttp: exposeHttp,
		Vars:       vars,
	}

	config.Application.Template.Pods = append(config.Application.Template.Pods, pod)
}

func getPodTypeFromService(name string, image string) string {
	// Database types
	if strings.Contains(image, "postgres") || strings.Contains(name, "postgres") {
		return PodTypePostgres
	}
	if strings.Contains(image, "mysql") || strings.Contains(name, "mysql") {
		return PodTypeMySQL
	}
	if strings.Contains(image, "neo4j") || strings.Contains(name, "neo4j") {
		return PodTypeNeo4j
	}
	if strings.Contains(image, "redis") || strings.Contains(name, "redis") {
		return PodTypeRedis
	}
	if strings.Contains(image, "mongo") || strings.Contains(name, "mongo") {
		return PodTypeMongoDB
	}
	if strings.Contains(image, "pinecone") || strings.Contains(name, "pinecone") {
		return PodTypePinecone
	}

	// Frontend types
	if strings.Contains(image, "react") || strings.Contains(name, "react") {
		return PodTypeReact
	}
	if strings.Contains(image, "angular") || strings.Contains(name, "angular") {
		return PodTypeAngular
	}
	if strings.Contains(image, "vue") || strings.Contains(name, "vue") {
		return PodTypeVue
	}

	// Backend types
	if strings.Contains(image, "django") || strings.Contains(name, "django") {
		return PodTypeDjango
	}
	if strings.Contains(image, "fastapi") || strings.Contains(name, "fastapi") {
		return PodTypeFastAPI
	}
	if strings.Contains(image, "express") || strings.Contains(name, "express") {
		return PodTypeExpress
	}

	// Other types
	if strings.Contains(image, "nginx") || strings.Contains(name, "nginx") {
		return PodTypeNginx
	}
	if strings.Contains(image, "llm") || strings.Contains(name, "llm") {
		return PodTypeLLM
	}

	return ""
}

func getPodTypeFromDependency(name string) string {
	switch name {
	// Database dependencies
	case "pg", "postgres", "postgresql":
		return PodTypePostgres
	case "mysql", "mysql2":
		return PodTypeMySQL
	case "neo4j-driver":
		return PodTypeNeo4j
	case "redis":
		return PodTypeRedis
	case "mongodb", "mongoose":
		return PodTypeMongoDB
	case "pinecone-client":
		return PodTypePinecone

	// Frontend dependencies
	case "react", "react-dom":
		return PodTypeReact
	case "@angular/core":
		return PodTypeAngular
	case "vue":
		return PodTypeVue

	// Backend dependencies
	case "django":
		return PodTypeDjango
	case "fastapi":
		return PodTypeFastAPI
	case "express":
		return PodTypeExpress

	// Other dependencies
	case "nginx":
		return PodTypeNginx
	case "@langchain/core", "langchain":
		return PodTypeLLM
	}

	return ""
}

func detectServiceDependencies(dockerComposePath string) []ServiceDependency {
	deps := []ServiceDependency{}
	seenServices := make(map[string]bool)

	// Read docker-compose.yml if it exists
	if data, err := os.ReadFile(dockerComposePath); err == nil {
		var compose DockerCompose
		if err := yaml.Unmarshal(data, &compose); err == nil {
			for name, service := range compose.Services {
				if !seenServices[name] {
					seenServices[name] = true

					// Determine pod type and image
					podType := ""
					image := service.Image
					switch {
					case name == "frontend" || strings.Contains(name, "react"):
						podType = "frontend"
						if image == "" {
							image = "node:18"
						}
					case strings.Contains(name, "redis") || strings.Contains(image, "redis"):
						podType = "database"
						if image == "" {
							image = "redis:7"
						}
					case strings.Contains(name, "postgres") || strings.Contains(image, "postgres"):
						podType = "database"
						if image == "" {
							image = "postgres:latest"
						}
					case strings.Contains(name, "mysql") || strings.Contains(image, "mysql"):
						podType = "database"
						if image == "" {
							image = "mysql:latest"
						}
					case strings.Contains(name, "mongodb") || strings.Contains(image, "mongo"):
						podType = "database"
						if image == "" {
							image = "mongo:latest"
						}
					case strings.Contains(name, "neo4j") || strings.Contains(image, "neo4j"):
						podType = "database"
						if image == "" {
							image = "neo4j:latest"
						}
					case strings.Contains(name, "nginx") || strings.Contains(image, "nginx"):
						podType = "nginx"
						if image == "" {
							image = "nginx:latest"
						}
					case strings.Contains(name, "llm") || strings.Contains(image, "llm"):
						podType = "llm"
					case strings.Contains(name, "pinecone") || strings.Contains(image, "pinecone"):
						podType = "pinecone"
					}

					if podType != "" {
						deps = append(deps, ServiceDependency{
							Type:     podType,
							Name:     name,
							Image:    image,
							Required: true,
						})
					}
				}
			}
		}
	}

	return deps
}

func createDefaultConfig(projectName string, stackType string, deps []ServiceDependency) Config {
	config := Config{
		Application: types.Application{
			Template: types.Template{
				Name:           projectName,
				DeploymentName: projectName,
				RegistryLogin: types.RegistryAuth{
					Registry:            "ghcr.io",
					Username:           "<Github username>",
					PersonalAccessToken: "<Github Packages Read-Only PAT>",
				},
				Build: struct {
					Command string `yaml:"command"`
					Output  string `yaml:"output"`
				}{
					Command: "npm install && npm run build",
					Output:  "build",
				},
			},
		},
	}

	// Add detected service dependencies
	var pods []types.PodConfig
	seenServices := make(map[string]bool)

	for _, dep := range deps {
		if seenServices[dep.Name] {
			continue
		}
		seenServices[dep.Name] = true

		switch dep.Type {
		case "frontend":
			pods = append(pods, types.PodConfig{
				Type: PodTypeReact,
				Name: "frontend",
				Tag:  dep.Image,
				Vars: []types.VarPair{
					{Key: "NODE_ENV", Value: "development"},
					{Key: "PORT", Value: "3000"},
				},
				ExposeHttp: true,
			})
		case "backend":
			pods = append(pods, types.PodConfig{
				Type: PodTypeExpress,
				Name: "backend",
				Tag:  dep.Image,
				Vars: []types.VarPair{
					{Key: "NODE_ENV", Value: "development"},
					{Key: "PORT", Value: "3001"},
				},
				ExposeHttp: true,
			})
		case "database":
			if strings.Contains(dep.Name, "redis") || strings.Contains(dep.Image, "redis") {
				pods = append(pods, types.PodConfig{
					Type: PodTypeRedis,
					Name: "redis",
					Tag:  dep.Image,
					Vars: []types.VarPair{
						{Key: "REDIS_MAX_MEMORY", Value: "256mb"},
					},
				})
			} else if strings.Contains(dep.Name, "mongodb") || strings.Contains(dep.Image, "mongo") {
				pods = append(pods, types.PodConfig{
					Type: PodTypeMongoDB,
					Name: "mongodb",
					Tag:  dep.Image,
					Vars: []types.VarPair{
						{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
						{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "<your-mongodb-password>"},
						{Key: "MONGO_INITDB_DATABASE", Value: projectName},
					},
				})
			}
		}
	}

	config.Application.Template.Pods = pods
	return config
}

func createMERNConfig(projectName string) Config {
	config := Config{
		Application: types.Application{
			Template: types.Template{
				Name:           projectName,
				DeploymentName: projectName,
				RegistryLogin: types.RegistryAuth{
					Registry:            "ghcr.io",
					Username:           "<Github username>",
					PersonalAccessToken: "<Github Packages Read-Only PAT>",
				},
				Pods: []types.PodConfig{
					{
						Type: PodTypeMongoDB,
						Name: "mongodb",
						Tag:  "mongo:6",
						Vars: []types.VarPair{
							{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
							{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "<your-mongodb-password>"},
							{Key: "MONGO_INITDB_DATABASE", Value: projectName},
						},
					},
					{
						Type: PodTypeExpress,
						Name: "express",
						Tag:  "node:18",
						Vars: []types.VarPair{
							{Key: "PORT", Value: "3000"},
							{Key: "NODE_ENV", Value: "development"},
							{Key: "MONGODB_URI", Value: "mongodb://root:<your-mongodb-password>@mongodb:27017/" + projectName + "?authSource=admin"},
						},
						ExposeHttp: true,
					},
					{
						Type: PodTypeReact,
						Name: "react",
						Tag:  "node:18",
						Vars: []types.VarPair{
							{Key: "PORT", Value: "3001"},
							{Key: "REACT_APP_API_URL", Value: "http://localhost:3000"},
						},
						ExposeHttp: true,
					},
				},
				Build: struct {
					Command string `yaml:"command"`
					Output  string `yaml:"output"`
				}{
					Command: "npm install && npm run build",
					Output:  "build",
				},
			},
		},
	}

	return config
}

func createLlamaNodeConfig(projectName string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper("llama-node"))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	dbPod := PodConfig{
		Type:       PodTypePostgres,
		Name:       "postgres",
		Tag:        "postgres:latest",
		ExposeHttp: false,
		Vars: []VarPair{
			{Key: "POSTGRES_USER", Value: "postgres"},
			{Key: "POSTGRES_PASSWORD", Value: "passw0rd"},
			{Key: "POSTGRES_DB", Value: "llama"},
		},
	}

	appPod := PodConfig{
		Type:       PodTypeExpress,
		Name:       "app",
		Tag:        "ghcr.io/your-username/llama-app:latest",
		ExposeHttp: true,
		Vars: []VarPair{
			{Key: "DATABASE_URL", Value: "postgresql://postgres:passw0rd@postgres:5432/llama"},
			{Key: "MODEL_PATH", Value: "/models/llama-2-70b-chat.Q4_K_M.gguf"},
			{Key: "NUM_GPU_LAYERS", Value: "35"},
			{Key: "CONTEXT_SIZE", Value: "4096"},
			{Key: "NUM_THREADS", Value: "4"},
			{Key: "GPU_LAYERS", Value: "all"},
		},
	}

	config.Application.Template.Pods = []PodConfig{dbPod, appPod}

	return config
}

func createLlamaPyConfig(projectName string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper("llama-py"))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	dbPod := PodConfig{
		Type:       PodTypeMongoDB,
		Name:       "mongodb",
		Tag:        "mongo:latest",
		ExposeHttp: false,
		Vars: []VarPair{
			{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "mongo"},
			{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "passw0rd"},
			{Key: "MONGO_INITDB_DATABASE", Value: "llama"},
		},
	}

	appPod := PodConfig{
		Type:       PodTypeFastAPI,
		Name:       "app",
		Tag:        "ghcr.io/your-username/llama-app:latest",
		ExposeHttp: true,
		Vars: []VarPair{
			{Key: "MONGODB_URL", Value: "DATABASE_CONNECTION_STRING"},
			{Key: "MODEL_PATH", Value: "/models/llama-2-70b-chat.Q4_K_M.gguf"},
			{Key: "NUM_GPU_LAYERS", Value: "35"},
			{Key: "CONTEXT_SIZE", Value: "4096"},
			{Key: "NUM_THREADS", Value: "4"},
			{Key: "USE_MLOCK", Value: "true"},
			{Key: "GPU_LAYERS", Value: "all"},
		},
	}

	config.Application.Template.Pods = []PodConfig{dbPod, appPod}

	return config
}

func createHuggingFaceConfig(projectName string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper("huggingface"))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	backendPod := PodConfig{
		Type:       PodTypeFastAPI,
		Name:       "backend",
		Tag:        "ghcr.io/your-username/hf-app:latest",
		ExposeHttp: false,
		Vars: []VarPair{
			{Key: "HF_API_KEY", Value: "your-huggingface-api-key"},
			{Key: "MODEL_ID", Value: "mistralai/Mixtral-8x7B-Instruct-v0.1"},
			{Key: "CUDA_VISIBLE_DEVICES", Value: "0"},
			{Key: "MAX_LENGTH", Value: "2048"},
			{Key: "TOP_K", Value: "50"},
			{Key: "TOP_P", Value: "0.9"},
		},
	}

	frontendPod := PodConfig{
		Type:       PodTypeNginx,
		Name:       "frontend",
		Tag:        "ghcr.io/your-username/hf-frontend:latest",
		ExposeHttp: true,
		Vars: []VarPair{
			{Key: "BACKEND_URL", Value: "BACKEND_CONNECTION_URL"},
			{Key: "VITE_API_URL", Value: "/api"},
		},
	}

	config.Application.Template.Pods = []PodConfig{backendPod, frontendPod}

	return config
}

func createMLConfig(projectName string, stackType string) Config {
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
				Vars: []VarPair{
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
				Vars: []VarPair{
					{Key: "MINIO_ROOT_USER", Value: "minioadmin"},
					{Key: "MINIO_ROOT_PASSWORD", Value: "minioadmin"},
				},
			},
		}
	}

	return config
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
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName := args[0]

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
				config = createMERNConfig(projectName)
			case "mean":
				config = createMEANConfig(projectName)
			case "mevn":
				config = createMEVNConfig(projectName)
			case "pern":
				config = createPERNConfig(projectName)
			case "kubeflow":
				config = templates.CreateKubeflowConfig(projectName)
			case "mlflow":
				config = createMLConfig(projectName, templateFlag)
			case "openai-node":
				config = createLlamaNodeConfig(projectName)
			case "openai-py":
				config = createLlamaPyConfig(projectName)
			case "huggingface":
				config = createHuggingFaceConfig(projectName)
			default:
				// For unknown templates, detect project type and dependencies
				stackType, deps := detectProjectType(".")
				config = createDefaultConfig(projectName, stackType, deps)
			}

			progress.Add(60)

			// Write config file
			yamlFileName := config.Application.Template.DeploymentName + ".yaml"
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
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
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
		if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err == nil {
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
	case StackLlamaNode:
		return 3000
	case StackLlamaPy:
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

	// Check for package.json
	if _, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		return StackOpenAINode, deps
	}

	// Check for pyproject.toml
	if _, err := os.ReadFile(filepath.Join(dir, "pyproject.toml")); err == nil {
		return StackOpenAIPy, deps
	}

	// Default to Node.js if no specific markers found
	return StackOpenAINode, deps
}

type PackageJSON struct {
	Name string `json:"name"`
}

type PyProject struct {
	Project struct {
		Name string `toml:"name"`
	} `toml:"project"`
}

type ServiceDependency struct {
	Type     string
	Name     string
	Image    string
	Required bool
}
