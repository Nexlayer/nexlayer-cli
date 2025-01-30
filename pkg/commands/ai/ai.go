package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// AIProvider represents an AI code assistant provider
type AIProvider struct {
	Name        string
	EnvVarKey   string
	Description string
	Endpoint    string
}

// Supported AI providers
var (
	GitHubCopilot = AIProvider{
		Name:        "GitHub Copilot",
		EnvVarKey:   "GITHUB_COPILOT_TOKEN",
		Description: "GitHub's AI pair programmer",
		Endpoint:    "https://api.github.com/copilot",
	}
	CursorAI = AIProvider{
		Name:        "Cursor AI",
		EnvVarKey:   "CURSOR_AI_KEY",
		Description: "AI-powered code editor",
		Endpoint:    "https://api.cursor.sh",
	}
	WindsurfEditor = AIProvider{
		Name:        "Windsurf Editor by Codeium",
		EnvVarKey:   "WINDSURF_API_KEY",
		Description: "World's first agentic IDE",
		Endpoint:    "https://api.codeium.com/windsurf",
	}
	Tabnine = AIProvider{
		Name:        "Tabnine",
		EnvVarKey:   "TABNINE_API_KEY",
		Description: "AI code completion assistant",
		Endpoint:    "https://api.tabnine.com",
	}
	JetBrainsAI = AIProvider{
		Name:        "JetBrains AI",
		EnvVarKey:   "JETBRAINS_AI_KEY",
		Description: "AI-powered development in JetBrains IDEs",
		Endpoint:    "https://api.jetbrains.com/ai",
	}
	IntelliCode = AIProvider{
		Name:        "Microsoft IntelliCode",
		EnvVarKey:   "INTELLICODE_KEY",
		Description: "AI-assisted development in Visual Studio",
		Endpoint:    "https://api.intellicode.microsoft.com",
	}
	CodeWhisperer = AIProvider{
		Name:        "Amazon CodeWhisperer",
		EnvVarKey:   "CODEWHISPERER_KEY",
		Description: "Amazon's AI code companion",
		Endpoint:    "https://api.aws.amazon.com/codewhisperer",
	}
	ClaudeSonnet = AIProvider{
		Name:        "Claude Sonnet 3.5",
		EnvVarKey:   "CLAUDE_API_KEY",
		Description: "Anthropic's advanced code assistant",
		Endpoint:    "https://api.anthropic.com/v1",
	}
	ChatGPTCode = AIProvider{
		Name:        "ChatGPT Code Interpreter",
		EnvVarKey:   "OPENAI_API_KEY",
		Description: "OpenAI's code interpreter",
		Endpoint:    "https://api.openai.com/v1",
	}
)

// All supported AI providers
var AllProviders = []AIProvider{
	GitHubCopilot,
	CursorAI,
	WindsurfEditor,
	Tabnine,
	JetBrainsAI,
	IntelliCode,
	CodeWhisperer,
	ClaudeSonnet,
	ChatGPTCode,
}

const yamlPrompt = `You are an AI assistant integrated into the Nexlayer CLI. Your job is to generate YAML deployment templates that strictly follow the Nexlayer template standards and specifications.

### Nexlayer Template Requirements:
1. **General Structure**:
   - The YAML must include application.template.name, application.template.deploymentName, and application.template.registryLogin fields.
   - The pods array must define the application's components, such as frontend, backend, database, or others (e.g., nginx, llm).

2. **Pod Configuration**:
   - Each pod must include:
     - type: The type of component (frontend, backend, database, nginx, llm, etc.).
     - name: A specific, descriptive name for the pod (e.g., backend-api, frontend-react, custom-llm-service).
     - tag: A valid Docker image tag (e.g., node:18, redis:7, nginx:1.23).
     - vars: Environment variables for the pod in the format [{ "key": "<VAR_NAME>", "value": "<VALUE>" }].

3. **Supported Pod Types**:
   - Frontend: react, angular, vue
   - Backend: express, django, fastapi
   - Database: mongodb, postgres, redis, neo4j, pinecone
   - Others:
     - nginx: Used for load balancing or serving static assets.
     - llm: Custom pods for large language models or custom workloads (naming allowed).

4. **Expose HTTP**:
   - Use exposeHttp: true for components that need to be accessible via HTTP.

5. **Predefined Environment Variables**:
   - Include the following mandatory variables based on the stack:
     - DATABASE_CONNECTION_STRING for database connectivity.
     - FRONTEND_CONNECTION_URL and BACKEND_CONNECTION_URL for connecting frontend and backend services.
     - PINECONE_API_KEY, PINECONE_ENVIRONMENT, PINECONE_INDEX for Pinecone vector database.

6. **Validation**:
   - Ensure the generated YAML conforms to Nexlayer's standards and avoid hallucinating unsupported configurations.

Based on the following inputs, generate a YAML deployment file:
- Application Name: %s
- Detected Stack: %s
- Required Components: %s

The YAML should be valid and follow all the requirements above. Only output the YAML content, nothing else.`

// GenerateYAML generates a YAML template using AI
func GenerateYAML(appName string, stackType string, components []string) (string, error) {
	// Get preferred AI provider from environment
	provider := getPreferredProvider()
	if provider == nil {
		return "", fmt.Errorf("no AI provider configured. Set one of the following environment variables: %s", getSupportedEnvVars())
	}

	// For now, return a mock response
	// TODO: Make actual API call to the provider using the prompt
	return mockGenerateYAML(appName, stackType, components), nil
}

// getPreferredProvider returns the first configured AI provider
func getPreferredProvider() *AIProvider {
	for _, provider := range AllProviders {
		if os.Getenv(provider.EnvVarKey) != "" {
			return &provider
		}
	}
	return nil
}

// getSupportedEnvVars returns a list of supported environment variables
func getSupportedEnvVars() string {
	var vars []string
	for _, provider := range AllProviders {
		vars = append(vars, provider.EnvVarKey)
	}
	return strings.Join(vars, ", ")
}

// mockGenerateYAML generates a mock YAML response
func mockGenerateYAML(appName string, stackType string, components []string) string {
	// Base template following Nexlayer standards
	yaml := fmt.Sprintf(`application:
  template:
    name: "%s"
    deploymentName: "%s"
    registryLogin:
      registry: ghcr.io
      username: <Github username>
      personalAccessToken: <Github Packages Read-Only PAT>
    pods:`, appName, appName)

	// Add pods based on stack type and components
	if strings.Contains(stackType, "node") {
		// Node.js backend
		yaml += `
      - type: backend
        name: express
        tag: "node:18"
        vars:
          - key: PORT
            value: "3000"
          - key: NODE_ENV
            value: "development"
          - key: BACKEND_CONNECTION_URL
            value: "http://express:3000"
        exposeHttp: true`
	}

	if strings.Contains(stackType, "react") {
		// React frontend
		yaml += `
      - type: frontend
        name: react
        tag: "node:18"
        vars:
          - key: NODE_ENV
            value: "development"
          - key: PORT
            value: "3000"
          - key: FRONTEND_CONNECTION_URL
            value: "http://react:3000"
          - key: BACKEND_CONNECTION_URL
            value: "http://express:3000"
        exposeHttp: true`
	}

	// Add database pods based on components
	for _, comp := range components {
		switch {
		case strings.Contains(comp, "redis"):
			yaml += `
      - type: database
        name: redis
        tag: "redis:7"
        vars:
          - key: REDIS_MAX_MEMORY
            value: "256mb"
          - key: DATABASE_CONNECTION_STRING
            value: "redis://redis:6379"
        exposeHttp: false`
		case strings.Contains(comp, "mongo"):
			yaml += `
      - type: database
        name: mongodb
        tag: "mongo:6"
        vars:
          - key: DATABASE_CONNECTION_STRING
            value: "mongodb://mongodb:27017"
        exposeHttp: false`
		case strings.Contains(comp, "postgres"):
			yaml += `
      - type: database
        name: postgres
        tag: "postgres:15"
        vars:
          - key: POSTGRES_DB
            value: "app"
          - key: POSTGRES_USER
            value: "postgres"
          - key: POSTGRES_PASSWORD
            value: "<your-postgres-password>"
          - key: DATABASE_CONNECTION_STRING
            value: "postgresql://postgres:password@postgres:5432/app"
        exposeHttp: false`
		case strings.Contains(comp, "pinecone"):
			yaml += `
      - type: database
        name: pinecone
        tag: "pinecone/pinecone-client:latest"
        vars:
          - key: PINECONE_API_KEY
            value: "<your-pinecone-api-key>"
          - key: PINECONE_ENVIRONMENT
            value: "<your-pinecone-environment>"
          - key: PINECONE_INDEX
            value: "<your-pinecone-index>"
        exposeHttp: false`
		case strings.Contains(comp, "llm"):
			yaml += `
      - type: llm
        name: llm-service
        tag: "custom-llm:latest"
        vars:
          - key: LLM_CONNECTION_URL
            value: "http://llm-service:8000"
          - key: MODEL_PATH
            value: "/models"
        exposeHttp: true`
		case strings.Contains(comp, "nginx"):
			yaml += `
      - type: nginx
        name: nginx
        tag: "nginx:1.25"
        vars:
          - key: NGINX_PORT
            value: "80"
        exposeHttp: true`
		}
	}

	// Add build configuration
	yaml += `
    build:
      command: "npm install && npm run build"
      output: "build"`

	return yaml
}

// detectStack detects the project's stack and components
func detectStack(dir string) (string, []string) {
	var components []string

	// Check for package.json
	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			// Check frontend frameworks
			if _, hasReact := pkg.Dependencies["react"]; hasReact {
				components = append(components, "react")
			}
			if _, hasVue := pkg.Dependencies["vue"]; hasVue {
				components = append(components, "vue")
			}
			if _, hasAngular := pkg.Dependencies["@angular/core"]; hasAngular {
				components = append(components, "angular")
			}

			// Check backend frameworks
			if _, hasExpress := pkg.Dependencies["express"]; hasExpress {
				components = append(components, "express")
			}

			// Check databases
			if _, hasMongo := pkg.Dependencies["mongodb"]; hasMongo {
				components = append(components, "mongodb")
			}
			if _, hasRedis := pkg.Dependencies["redis"]; hasRedis {
				components = append(components, "redis")
			}
			if _, hasPg := pkg.Dependencies["pg"]; hasPg {
				components = append(components, "postgres")
			}
			if _, hasPinecone := pkg.Dependencies["@pinecone-database/pinecone"]; hasPinecone {
				components = append(components, "pinecone")
			}

			// Check AI/LLM
			if _, hasLangchain := pkg.Dependencies["langchain"]; hasLangchain {
				components = append(components, "llm")
			}
		}
	}

	// Check for requirements.txt
	if data, err := os.ReadFile(filepath.Join(dir, "requirements.txt")); err == nil {
		reqs := strings.Split(string(data), "\n")
		for _, req := range reqs {
			req = strings.TrimSpace(req)
			switch {
			case strings.HasPrefix(req, "fastapi"):
				components = append(components, "fastapi")
			case strings.HasPrefix(req, "django"):
				components = append(components, "django")
			case strings.HasPrefix(req, "pymongo"):
				components = append(components, "mongodb")
			case strings.HasPrefix(req, "redis"):
				components = append(components, "redis")
			case strings.HasPrefix(req, "psycopg"):
				components = append(components, "postgres")
			case strings.HasPrefix(req, "pinecone-client"):
				components = append(components, "pinecone")
			case strings.HasPrefix(req, "langchain"):
				components = append(components, "llm")
			}
		}
	}

	// Check for docker-compose.yml
	if data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml")); err == nil {
		// Simple string matching for now
		content := string(data)
		if strings.Contains(content, "nginx") {
			components = append(components, "nginx")
		}
		if strings.Contains(content, "redis") {
			components = append(components, "redis")
		}
		if strings.Contains(content, "mongo") {
			components = append(components, "mongodb")
		}
		if strings.Contains(content, "postgres") {
			components = append(components, "postgres")
		}
		if strings.Contains(content, "pinecone") {
			components = append(components, "pinecone")
		}
	}

	// Determine stack type
	stackType := "unknown"
	if containsAny(components, "react", "vue", "angular") {
		if containsAny(components, "express") {
			stackType = "node"
		} else if containsAny(components, "fastapi", "django") {
			stackType = "python"
		}
	} else if containsAny(components, "express") {
		stackType = "node"
	} else if containsAny(components, "fastapi", "django") {
		stackType = "python"
	}

	return stackType, components
}

// containsAny checks if slice contains any of the values
func containsAny(slice []string, values ...string) bool {
	for _, v := range values {
		for _, s := range slice {
			if s == v {
				return true
			}
		}
	}
	return false
}

// Command represents the ai command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-powered features",
		Long:  "AI-powered features for generating and optimizing deployment templates",
	}

	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newDetectCommand())

	return cmd
}

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [app-name]",
		Short: "Generate deployment template using AI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			stackType, components := detectStack(".")
			yaml, err := GenerateYAML(appName, stackType, components)
			if err != nil {
				return err
			}
			fmt.Println(yaml)
			return nil
		},
	}
}

func newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect available AI providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Available AI Providers:")
			for _, provider := range AllProviders {
				configured := ""
				if os.Getenv(provider.EnvVarKey) != "" {
					configured = " (configured)"
				}
				fmt.Printf("- %s: %s%s\n", provider.Name, provider.Description, configured)
			}
			return nil
		},
	}
}
