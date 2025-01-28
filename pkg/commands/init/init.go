package init

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	// Stack types - Traditional
	StackMERN = "mern"  // MongoDB, Express, React, Node.js
	StackMEAN = "mean"  // MongoDB, Express, Angular, Node.js
	StackMEVN = "mevn"  // MongoDB, Express, Vue.js, Node.js
	StackPERN = "pern"  // PostgreSQL, Express, React, Node.js
	StackMNFA = "mnfa"  // MongoDB, Neo4j, FastAPI, Angular
	StackPDN  = "pdn"   // PostgreSQL, Django, Node.js

	// Stack types - AI/LLM
	StackLangChainJS = "langchain-js"    // LangChain.js, Next.js, MongoDB
	StackLangChainPy = "langchain-py"    // LangChain Python, FastAPI, PostgreSQL
	StackOpenAINode  = "openai-node"     // OpenAI Node.js SDK, Express, React
	StackOpenAIPy    = "openai-py"       // OpenAI Python SDK, FastAPI, Vue
	StackLlamaNode   = "llama-node"      // Llama.cpp Node.js, Next.js, PostgreSQL
	StackLlamaPy     = "llama-py"        // Llama.cpp Python, FastAPI, MongoDB
	StackVertexAI    = "vertex-ai"       // Google Vertex AI, Flask, React
	StackHuggingface = "huggingface"     // Hugging Face Transformers, FastAPI, React
	StackAnthropicPy = "anthropic-py"    // Anthropic Claude, FastAPI, Svelte
	StackAnthropicJS = "anthropic-js"    // Anthropic Claude, Next.js, MongoDB
)

type Config struct {
	Name        string            `yaml:"name"`
	Type        string            `yaml:"type"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Build       struct {
		Command string   `yaml:"command,omitempty"`
		Output  string   `yaml:"output,omitempty"`
		Env     []string `yaml:"env,omitempty"`
	} `yaml:"build,omitempty"`
	Deploy struct {
		Resources struct {
			CPU    string `yaml:"cpu,omitempty"`
			Memory string `yaml:"memory,omitempty"`
		} `yaml:"resources,omitempty"`
		Port int `yaml:"port,omitempty"`
	} `yaml:"deploy,omitempty"`
	Application struct {
		Template struct {
			Name            string `yaml:"name"`
			DeploymentName  string `yaml:"deploymentName"`
			RegistryLogin   struct {
				Registry            string `yaml:"registry"`
				Username           string `yaml:"username"`
				PersonalAccessToken string `yaml:"personalAccessToken"`
			} `yaml:"registryLogin"`
			Pods []struct {
				Type       string `yaml:"type"`
				ExposeHttp bool   `yaml:"exposeHttp"`
				Name       string `yaml:"name"`
				Tag        string `yaml:"tag"`
				PrivateTag bool   `yaml:"privateTag"`
				Vars       []VarPair `yaml:"vars"`
				GPU       bool `yaml:"gpu"`
				Resources struct {
					Limits   map[string]string `yaml:"limits"`
					Requests map[string]string `yaml:"requests"`
				} `yaml:"resources"`
			} `yaml:"pods"`
		} `yaml:"template"`
	} `yaml:"application"`
}

type PackageJSON struct {
	Name string `json:"name"`
}

type PyProjectTOML struct {
	Project struct {
		Name string `toml:"name"`
	} `toml:"project"`
}

type PodConfig struct {
	Type       string
	Name       string
	ExposeHttp bool
	Vars       []VarPair
}

// VarPair represents a key-value pair for environment variables
type VarPair struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

func addTemplateConfig(config *Config, templateName string, pods []PodConfig) {
	config.Application.Template.Name = templateName
	
	for _, pod := range pods {
		addPod(config, pod.Type, pod.Name, pod.ExposeHttp, pod.Vars)
	}
}

func addPod(config *Config, podType string, name string, exposeHttp bool, vars []VarPair) {
	pod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       podType,
		ExposeHttp: exposeHttp,
		Name:       name,
		Tag:        fmt.Sprintf("ghcr.io/your-username/%s:latest", name),
		PrivateTag: false,
		Vars:       vars,
	}

	config.Application.Template.Pods = append(config.Application.Template.Pods, pod)
}

func addGPUResources(pod *struct {
	Type       string    `yaml:"type"`
	ExposeHttp bool      `yaml:"exposeHttp"`
	Name       string    `yaml:"name"`
	Tag        string    `yaml:"tag"`
	PrivateTag bool      `yaml:"privateTag"`
	Vars       []VarPair `yaml:"vars"`
	GPU        bool      `yaml:"gpu"`
	Resources  struct {
		Limits   map[string]string `yaml:"limits"`
		Requests map[string]string `yaml:"requests"`
	} `yaml:"resources"`
}) {
	pod.GPU = true
	pod.Resources = struct {
		Limits   map[string]string `yaml:"limits"`
		Requests map[string]string `yaml:"requests"`
	}{
		Limits: map[string]string{
			"nvidia.com/gpu": "1",
		},
		Requests: map[string]string{
			"nvidia.com/gpu": "1",
		},
	}
}

func createDefaultConfig(projectName, stackType string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper(stackType))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	switch stackType {
	case StackLangChainJS:
		addTemplateConfig(&config, "langchain-nextjs-mongodb", []PodConfig{
			{Type: "database", Name: "mongodb", ExposeHttp: false, Vars: []VarPair{
				{"MONGO_INITDB_ROOT_USERNAME", "mongo"},
				{"MONGO_INITDB_ROOT_PASSWORD", "passw0rd"},
				{"MONGO_INITDB_DATABASE", "langchain"},
			}},
			{Type: "nextjs", Name: "app", ExposeHttp: true, Vars: []VarPair{
				{"MONGODB_URL", "DATABASE_CONNECTION_STRING"},
				{"OPENAI_API_KEY", "your-openai-api-key"},
				{"LANGCHAIN_TRACING_V2", "true"},
				{"LANGCHAIN_ENDPOINT", "https://api.smith.langchain.com"},
				{"LANGCHAIN_API_KEY", "your-langchain-api-key"},
				{"LANGCHAIN_PROJECT", projectName},
			}},
		})

	case StackLangChainPy:
		addTemplateConfig(&config, "langchain-fastapi-postgres", []PodConfig{
			{Type: "database", Name: "postgres", ExposeHttp: false, Vars: []VarPair{
				{"POSTGRES_USER", "postgres"},
				{"POSTGRES_PASSWORD", "passw0rd"},
				{"POSTGRES_DB", "langchain"},
			}},
			{Type: "fastapi", Name: "app", ExposeHttp: true, Vars: []VarPair{
				{"DATABASE_URL", "postgresql://postgres:passw0rd@postgres:5432/langchain"},
				{"OPENAI_API_KEY", "your-openai-api-key"},
				{"LANGCHAIN_TRACING_V2", "true"},
				{"LANGCHAIN_ENDPOINT", "https://api.smith.langchain.com"},
				{"LANGCHAIN_API_KEY", "your-langchain-api-key"},
				{"LANGCHAIN_PROJECT", projectName},
			}},
		})

	case StackOpenAINode:
		addTemplateConfig(&config, "openai-express-react", []PodConfig{
			{Type: "express", Name: "backend", ExposeHttp: false, Vars: []VarPair{
				{"OPENAI_API_KEY", "your-openai-api-key"},
				{"OPENAI_ORG_ID", "your-org-id"},
				{"MODEL", "gpt-4-turbo-preview"},
				{"MAX_TOKENS", "2048"},
				{"TEMPERATURE", "0.7"},
			}},
			{Type: "nginx", Name: "frontend", ExposeHttp: true, Vars: []VarPair{
				{"BACKEND_URL", "BACKEND_CONNECTION_URL"},
				{"VITE_API_URL", "/api"},
			}},
		})

	case StackOpenAIPy:
		addTemplateConfig(&config, "openai-fastapi-vue", []PodConfig{
			{Type: "fastapi", Name: "backend", ExposeHttp: false, Vars: []VarPair{
				{"OPENAI_API_KEY", "your-openai-api-key"},
				{"OPENAI_ORG_ID", "your-org-id"},
				{"MODEL", "gpt-4-turbo-preview"},
				{"MAX_TOKENS", "2048"},
				{"TEMPERATURE", "0.7"},
			}},
			{Type: "nginx", Name: "frontend", ExposeHttp: true, Vars: []VarPair{
				{"BACKEND_URL", "BACKEND_CONNECTION_URL"},
				{"VITE_API_URL", "/api"},
			}},
		})

	case StackLlamaNode:
		config = createLlamaNodeConfig(projectName)
	case StackLlamaPy:
		config = createLlamaPyConfig(projectName)
	case StackVertexAI:
		addTemplateConfig(&config, "vertex-ai-flask-react", []PodConfig{
			{Type: "flask", Name: "backend", ExposeHttp: false, Vars: []VarPair{
				{"GOOGLE_CLOUD_PROJECT", "your-project-id"},
				{"GOOGLE_APPLICATION_CREDENTIALS", "/secrets/credentials.json"},
				{"VERTEX_LOCATION", "us-central1"},
				{"MODEL_NAME", "text-bison@002"},
				{"MAX_OUTPUT_TOKENS", "1024"},
				{"TEMPERATURE", "0.7"},
			}},
			{Type: "nginx", Name: "frontend", ExposeHttp: true, Vars: []VarPair{
				{"BACKEND_URL", "BACKEND_CONNECTION_URL"},
				{"VITE_API_URL", "/api"},
			}},
		})

	case StackHuggingface:
		config = createHuggingFaceConfig(projectName)
	case StackAnthropicPy:
		addTemplateConfig(&config, "anthropic-fastapi-svelte", []PodConfig{
			{Type: "fastapi", Name: "backend", ExposeHttp: false, Vars: []VarPair{
				{"ANTHROPIC_API_KEY", "your-anthropic-api-key"},
				{"MODEL", "claude-3-opus-20240229"},
				{"MAX_TOKENS", "4096"},
				{"TEMPERATURE", "0.7"},
			}},
			{Type: "nginx", Name: "frontend", ExposeHttp: true, Vars: []VarPair{
				{"BACKEND_URL", "BACKEND_CONNECTION_URL"},
				{"VITE_API_URL", "/api"},
			}},
		})

	case StackAnthropicJS:
		addTemplateConfig(&config, "anthropic-nextjs-mongodb", []PodConfig{
			{Type: "database", Name: "mongodb", ExposeHttp: false, Vars: []VarPair{
				{"MONGO_INITDB_ROOT_USERNAME", "mongo"},
				{"MONGO_INITDB_ROOT_PASSWORD", "passw0rd"},
				{"MONGO_INITDB_DATABASE", "anthropic"},
			}},
			{Type: "nextjs", Name: "app", ExposeHttp: true, Vars: []VarPair{
				{"MONGODB_URL", "DATABASE_CONNECTION_STRING"},
				{"ANTHROPIC_API_KEY", "your-anthropic-api-key"},
				{"MODEL", "claude-3-opus-20240229"},
				{"MAX_TOKENS", "4096"},
				{"TEMPERATURE", "0.7"},
			}},
		})
	}

	// Set default resource limits
	config.Deploy.Resources.CPU = "1000m"    // 1 CPU core
	config.Deploy.Resources.Memory = "2048Mi" // 2GB RAM
	config.Deploy.Port = detectPort(stackType)

	// Set build configuration
	setBuildConfig(&config, stackType)

	return config
}

func createLlamaNodeConfig(projectName string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper("llama-node"))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	dbPod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       "database",
		ExposeHttp: false,
		Name:       "postgres",
		Tag:        "postgres:latest",
		PrivateTag: false,
		Vars: []VarPair{
			{"POSTGRES_USER", "postgres"},
			{"POSTGRES_PASSWORD", "passw0rd"},
			{"POSTGRES_DB", "llama"},
		},
	}

	appPod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       "nextjs",
		ExposeHttp: true,
		Name:       "app",
		Tag:        "ghcr.io/your-username/llama-app:latest",
		PrivateTag: false,
		Vars: []VarPair{
			{"DATABASE_URL", "postgresql://postgres:passw0rd@postgres:5432/llama"},
			{"MODEL_PATH", "/models/llama-2-70b-chat.Q4_K_M.gguf"},
			{"NUM_GPU_LAYERS", "35"},
			{"CONTEXT_SIZE", "4096"},
			{"NUM_THREADS", "4"},
			{"GPU_LAYERS", "all"},
		},
	}
	addGPUResources(&appPod)

	config.Application.Template.Pods = []struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{dbPod, appPod}
	config.Deploy.Resources.CPU = "2000m"
	config.Deploy.Resources.Memory = "16384Mi"
	config.Deploy.Port = 3000

	return config
}

func createLlamaPyConfig(projectName string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper("llama-py"))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	dbPod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       "database",
		ExposeHttp: false,
		Name:       "mongodb",
		Tag:        "mongo:latest",
		PrivateTag: false,
		Vars: []VarPair{
			{"MONGO_INITDB_ROOT_USERNAME", "mongo"},
			{"MONGO_INITDB_ROOT_PASSWORD", "passw0rd"},
			{"MONGO_INITDB_DATABASE", "llama"},
		},
	}

	appPod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       "fastapi",
		ExposeHttp: true,
		Name:       "app",
		Tag:        "ghcr.io/your-username/llama-app:latest",
		PrivateTag: false,
		Vars: []VarPair{
			{"MONGODB_URL", "DATABASE_CONNECTION_STRING"},
			{"MODEL_PATH", "/models/llama-2-70b-chat.Q4_K_M.gguf"},
			{"NUM_GPU_LAYERS", "35"},
			{"CONTEXT_SIZE", "4096"},
			{"NUM_THREADS", "4"},
			{"USE_MLOCK", "true"},
			{"GPU_LAYERS", "all"},
		},
	}
	addGPUResources(&appPod)

	config.Application.Template.Pods = []struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{dbPod, appPod}
	config.Deploy.Resources.CPU = "2000m"
	config.Deploy.Resources.Memory = "16384Mi"
	config.Deploy.Port = 8000

	return config
}

func createHuggingFaceConfig(projectName string) Config {
	config := Config{}
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = fmt.Sprintf("My %s App", strings.ToUpper("huggingface"))
	config.Application.Template.RegistryLogin.Registry = "ghcr.io"

	backendPod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       "fastapi",
		ExposeHttp: false,
		Name:       "backend",
		Tag:        "ghcr.io/your-username/hf-app:latest",
		PrivateTag: false,
		Vars: []VarPair{
			{"HF_API_KEY", "your-huggingface-api-key"},
			{"MODEL_ID", "mistralai/Mixtral-8x7B-Instruct-v0.1"},
			{"CUDA_VISIBLE_DEVICES", "0"},
			{"MAX_LENGTH", "2048"},
			{"TOP_K", "50"},
			{"TOP_P", "0.9"},
		},
	}
	addGPUResources(&backendPod)

	frontendPod := struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{
		Type:       "nginx",
		ExposeHttp: true,
		Name:       "frontend",
		Tag:        "ghcr.io/your-username/hf-frontend:latest",
		PrivateTag: false,
		Vars: []VarPair{
			{"BACKEND_URL", "BACKEND_CONNECTION_URL"},
			{"VITE_API_URL", "/api"},
		},
	}

	config.Application.Template.Pods = []struct {
		Type       string    `yaml:"type"`
		ExposeHttp bool      `yaml:"exposeHttp"`
		Name       string    `yaml:"name"`
		Tag        string    `yaml:"tag"`
		PrivateTag bool      `yaml:"privateTag"`
		Vars       []VarPair `yaml:"vars"`
		GPU        bool      `yaml:"gpu"`
		Resources  struct {
			Limits   map[string]string `yaml:"limits"`
			Requests map[string]string `yaml:"requests"`
		} `yaml:"resources"`
	}{backendPod, frontendPod}
	config.Deploy.Resources.CPU = "2000m"
	config.Deploy.Resources.Memory = "8192Mi"
	config.Deploy.Port = 8000

	return config
}

// NewCommand creates a new init command
func NewCommand(client api.APIClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new project",
		Long:  "Initialize a new project with a template configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			uiManager := ui.NewManager()
			progress := uiManager.StartProgress("Initializing project")
			defer progress.Complete()

			projectName, err := cmd.Flags().GetString("name")
			if err != nil {
				return fmt.Errorf("failed to get project name: %w", err)
			}
			progress.Update(20.0, "Got project name")

			stackType, err := cmd.Flags().GetString("template")
			if err != nil {
				return fmt.Errorf("failed to get template type: %w", err)
			}
			progress.Update(40.0, fmt.Sprintf("Using template: %s", stackType))

			config := createDefaultConfig(projectName, stackType)
			progress.Update(60.0, "Created configuration")

			yamlData, err := yaml.Marshal(&config)
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}
			progress.Update(80.0, "Generated YAML configuration")

			err = os.WriteFile("nexlayer.yaml", yamlData, 0644)
			if err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}
			progress.Update(100.0, "Wrote configuration file")

			fmt.Printf("\nSuccessfully created nexlayer.yaml with %s template!\n", stackType)
			fmt.Println("To deploy your application, run: nexlayer deploy")
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "Project name")
	cmd.MarkFlagRequired("name")
	
	cmd.Flags().StringP("template", "t", "", "Template type (e.g., langchain-nextjs, llama-fastapi)")
	cmd.MarkFlagRequired("template")
	
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

func detectProjectType(dir string) string {
	// Check for package.json (Node.js)
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		var packageJSON struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		data, err := os.ReadFile(filepath.Join(dir, "package.json"))
		if err == nil {
			if err := json.Unmarshal(data, &packageJSON); err == nil {
				// Check for specific frameworks
				if _, hasNext := packageJSON.Dependencies["next"]; hasNext {
					return StackLangChainJS
				}
				if _, hasExpress := packageJSON.Dependencies["express"]; hasExpress {
					return StackOpenAINode
				}
				if _, hasLlama := packageJSON.Dependencies["llama-node"]; hasLlama {
					return StackLlamaNode
				}
			}
		}
		return StackOpenAINode
	}

	// Check for requirements.txt (Python)
	if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		data, err := os.ReadFile(filepath.Join(dir, "requirements.txt"))
		if err == nil {
			content := string(data)
			if strings.Contains(content, "langchain") {
				return StackLangChainPy
			}
			if strings.Contains(content, "llama-cpp-python") {
				return StackLlamaPy
			}
			if strings.Contains(content, "fastapi") {
				return StackOpenAIPy
			}
		}
		return StackOpenAIPy
	}

	// Check for pyproject.toml (Python with poetry)
	if _, err := os.Stat(filepath.Join(dir, "pyproject.toml")); err == nil {
		data, err := os.ReadFile(filepath.Join(dir, "pyproject.toml"))
		if err == nil {
			content := string(data)
			if strings.Contains(content, "langchain") {
				return StackLangChainPy
			}
			if strings.Contains(content, "llama-cpp-python") {
				return StackLlamaPy
			}
			if strings.Contains(content, "fastapi") {
				return StackOpenAIPy
			}
		}
		return StackOpenAIPy
	}

	// Default to OpenAI Node.js if no specific stack is detected
	return StackOpenAINode
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

func setBuildConfig(config *Config, projectType string) {
	switch projectType {
	case "nodejs":
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case "python":
		config.Build.Command = "pip install -r requirements.txt"
	case "golang":
		config.Build.Command = "go build -o app"
		config.Build.Output = "app"
	case "static":
		// No build needed for static sites
	case StackMERN:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackMEAN:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackMEVN:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackPERN:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackMNFA:
		config.Build.Command = "pip install -r requirements.txt"
	case StackPDN:
		config.Build.Command = "pip install -r requirements.txt"
	case StackLangChainJS:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackLangChainPy:
		config.Build.Command = "pip install -r requirements.txt"
	case StackOpenAINode:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackOpenAIPy:
		config.Build.Command = "pip install -r requirements.txt"
	case StackLlamaNode:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	case StackLlamaPy:
		config.Build.Command = "pip install -r requirements.txt"
	case StackVertexAI:
		config.Build.Command = "pip install -r requirements.txt"
	case StackHuggingface:
		config.Build.Command = "pip install -r requirements.txt"
	case StackAnthropicPy:
		config.Build.Command = "pip install -r requirements.txt"
	case StackAnthropicJS:
		config.Build.Command = "npm install && npm run build"
		config.Build.Output = "build"
	}
}

func writeConfig(file string, config Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(file, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
