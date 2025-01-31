package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	nexerrors "github.com/Nexlayer/nexlayer-cli/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// AIProvider represents an AI code assistant provider
type AIProvider struct {
	Name        string
	EnvVarKey   string
	Description string
	Endpoint    string
	Priority    int // Higher number = higher priority
}

// Supported AI providers
var (
	WindsurfEditor = AIProvider{
		Name:        "Windsurf Editor by Codeium",
		EnvVarKey:   "WINDSURF_API_KEY",
		Description: "World's first agentic IDE",
		Endpoint:    "https://api.codeium.com/windsurf",
		Priority:    100,
	}
	GitHubCopilot = AIProvider{
		Name:        "GitHub Copilot",
		EnvVarKey:   "GITHUB_COPILOT_TOKEN",
		Description: "GitHub's AI pair programmer",
		Endpoint:    "https://api.github.com/copilot",
		Priority:    90,
	}
	CursorAI = AIProvider{
		Name:        "Cursor AI",
		EnvVarKey:   "CURSOR_AI_KEY",
		Description: "AI-powered code editor",
		Endpoint:    "https://api.cursor.sh",
		Priority:    80,
	}
	JetBrainsAI = AIProvider{
		Name:        "JetBrains AI",
		EnvVarKey:   "JETBRAINS_AI_KEY",
		Description: "AI-powered development in JetBrains IDEs",
		Endpoint:    "https://api.jetbrains.com/ai",
		Priority:    70,
	}
	VSCodeAI = AIProvider{
		Name:        "VS Code AI",
		EnvVarKey:   "VSCODE_AI_KEY",
		Description: "AI assistance in VS Code",
		Endpoint:    "https://api.vscode.dev/ai",
		Priority:    60,
	}
)

// AllProviders holds all recognized AI providers, sorted by priority
var AllProviders = []AIProvider{
	WindsurfEditor,
	GitHubCopilot,
	CursorAI,
	JetBrainsAI,
	VSCodeAI,
}

// NexlayerYAML represents the structure of a Nexlayer deployment template
type NexlayerYAML struct {
	Application struct {
		Template struct {
			Name           string `yaml:"name"`
			DeploymentName string `yaml:"deploymentName"`
			RegistryLogin  struct {
				Registry            string `yaml:"registry"`
				Username            string `yaml:"username"`
				PersonalAccessToken string `yaml:"personalAccessToken"`
			} `yaml:"registryLogin"`
			Pods []struct {
				Type       string `yaml:"type"`
				Name       string `yaml:"name"`
				Tag        string `yaml:"tag"`
				PrivateTag bool   `yaml:"privateTag"`
				ExposeHTTP bool   `yaml:"exposeHttp"`
				Vars       []struct {
					Key   string `yaml:"key"`
					Value string `yaml:"value"`
				} `yaml:"vars"`
			} `yaml:"pods"`
			Build struct {
				Command string `yaml:"command"`
				Output  string `yaml:"output"`
			} `yaml:"build"`
		} `yaml:"template"`
	} `yaml:"application"`
}

// Template generation prompt for AI assistants
const yamlPrompt = `As a Nexlayer AI assistant, generate a deployment template YAML that follows these requirements:

1. Structure Requirements:
   - Must include application.template.name, deploymentName, and registryLogin
   - Must define pods array with application components

2. Pod Configuration:
   - Each pod needs: type, name, tag, privateTag, vars array
   - Supported pod types:
     * Frontend: react, angular, vue
     * Backend: express, django, fastapi
     * Database: postgres, mongodb, redis, neo4j
     * Others: nginx, llm

3. Environment Variables:
   - Database pods: Include DATABASE_CONNECTION_STRING
   - Frontend/Backend: Include appropriate connection URLs
   - LLM pods: Include model-specific variables

4. HTTP Exposure:
   - Set exposeHttp: true for web-accessible components

Application Details:
- Name: %s
- Stack Type: %s
- Components: %s

Output the YAML content only, no explanations.`

// GenerateYAML generates a Nexlayer deployment template using available AI assistants
func GenerateYAML(appName string, stackType string, components []string) (string, error) {
	// Try to get an AI provider from the IDE
	provider := getPreferredProvider()
	if provider == nil {
		// No AI provider available, use default template
		return generateDefaultTemplate(appName, stackType, components)
	}

	// Format the prompt for the AI
	prompt := fmt.Sprintf(yamlPrompt, appName, stackType, strings.Join(components, ", "))

	// In a real implementation, we would:
	// 1. Call the AI provider's API with the prompt
	// 2. Process the response
	// For now, use a mock response with the formatted prompt
	rawYAML := mockGenerateYAML(appName, stackType, components, prompt)

	// Validate and fix the generated YAML
	return validateAndFixYAML(rawYAML)
}

// generateDefaultTemplate creates a basic template based on detected components
func generateDefaultTemplate(appName string, stackType string, components []string) (string, error) {
	template := NexlayerYAML{}
	template.Application.Template.Name = appName
	template.Application.Template.DeploymentName = appName
	template.Application.Template.RegistryLogin.Registry = "ghcr.io"
	template.Application.Template.RegistryLogin.Username = "<your-username>"
	template.Application.Template.RegistryLogin.PersonalAccessToken = "<your-pat>"

	// Add pods based on components
	for _, comp := range components {
		pod := struct {
			Type       string `yaml:"type"`
			Name       string `yaml:"name"`
			Tag        string `yaml:"tag"`
			PrivateTag bool   `yaml:"privateTag"`
			ExposeHTTP bool   `yaml:"exposeHttp"`
			Vars       []struct {
				Key   string `yaml:"key"`
				Value string `yaml:"value"`
			} `yaml:"vars"`
		}{
			Type:       comp,
			Name:       fmt.Sprintf("%s-service", comp),
			Tag:        defaultTagForType(comp),
			PrivateTag: false,
			ExposeHTTP: isExposeByDefault(comp),
		}

		// Add default environment variables
		if isDatabaseType(comp) {
			pod.Vars = append(pod.Vars, struct {
				Key   string `yaml:"key"`
				Value string `yaml:"value"`
			}{
				Key:   "DATABASE_CONNECTION_STRING",
				Value: fmt.Sprintf("%s://%s:1234/db", comp, pod.Name),
			})
		}

		template.Application.Template.Pods = append(template.Application.Template.Pods, pod)
	}

	// Set build configuration
	template.Application.Template.Build.Command = getBuildCommand(stackType)
	template.Application.Template.Build.Output = "dist"

	// Convert to YAML
	out, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to generate template: %w", err)
	}

	return string(out), nil
}

func getBuildCommand(stackType string) string {
	switch strings.ToLower(stackType) {
	case "python":
		return "pip install -r requirements.txt"
	case "node":
		return "npm install && npm run build"
	default:
		return "npm install && npm run build"
	}
}

func isExposeByDefault(podType string) bool {
	switch strings.ToLower(podType) {
	case "react", "angular", "vue", "express", "django", "fastapi", "nginx":
		return true
	default:
		return false
	}
}

func defaultTagForType(podType string) string {
	switch strings.ToLower(podType) {
	case "react":
		return "node:18"
	case "angular":
		return "node:18"
	case "vue":
		return "node:18"
	case "express":
		return "node:18"
	case "django":
		return "python:3.10"
	case "fastapi":
		return "python:3.10"
	case "postgres":
		return "postgres:15"
	case "mongodb":
		return "mongo:6"
	case "redis":
		return "redis:7"
	case "nginx":
		return "nginx:1.25"
	default:
		return "latest"
	}
}

func isDatabaseType(podType string) bool {
	switch strings.ToLower(podType) {
	case "postgres", "mongodb", "redis":
		return true
	default:
		return false
	}
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

// mockGenerateYAML generates a mock YAML response
func mockGenerateYAML(appName string, stackType string, components []string, prompt string) string {
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

	// Add prompt to the YAML
	yaml += `
    prompt: "%s"`

	return fmt.Sprintf(yaml, prompt)
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

func validateAndFixYAML(yamlString string) (string, error) {
	// First, unify exposeOn80 -> exposeHttp (legacy support)
	yamlString = strings.ReplaceAll(yamlString, "exposeOn80:", "exposeHttp:")

	var template NexlayerYAML
	if err := yaml.Unmarshal([]byte(yamlString), &template); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// 1. Validate and fix template section
	if err := validateTemplate(&template); err != nil {
		return "", err
	}

	// 2. Validate and fix pods
	if err := validatePods(&template); err != nil {
		return "", err
	}

	// 3. Validate and fix build section
	if err := validateBuild(&template); err != nil {
		return "", err
	}

	// Marshal back to YAML
	out, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	return string(out), nil
}

func validateTemplate(template *NexlayerYAML) error {
	// Check required fields
	if template.Application.Template.Name == "" {
		return nexerrors.NewValidationError(
			"template name is required",
			&nexerrors.ValidationContext{
				Field:         "application.template.name",
				ExpectedType:  "string",
				ActualValue:   "",
				ResolutionHints: []string{
					"Add 'name' field under application.template",
					"The name should be a valid identifier for your application",
				},
				Example: `application:
  template:
    name: my-app`,
			},
		)
	}
	if template.Application.Template.DeploymentName == "" {
		template.Application.Template.DeploymentName = template.Application.Template.Name
	}

	// Validate registry login
	reg := &template.Application.Template.RegistryLogin
	if reg.Registry == "" {
		reg.Registry = "ghcr.io"
	}
	if reg.Username == "" {
		reg.Username = "<your-username>"
	}
	if reg.PersonalAccessToken == "" {
		reg.PersonalAccessToken = "<your-pat>"
	}

	return nil
}

func validatePods(template *NexlayerYAML) error {
	if len(template.Application.Template.Pods) == 0 {
		return nexerrors.NewValidationError(
			"no pods defined in template",
			&nexerrors.ValidationContext{
				Field:         "application.template.pods",
				ExpectedType:  "array",
				ActualValue:   "[]",
				ResolutionHints: []string{
					"Add at least one pod configuration",
					"Each pod must have a type, name, and tag",
				},
				Example: `application:
  template:
    pods:
      - type: backend
        name: api
        tag: node:18`,
			},
		)
	}

	// Required environment variables by pod type
	requiredVars := map[string][]string{
		"database": {
			"DATABASE_CONNECTION_STRING",
			"DATABASE_HOST",
		},
		"neo4j": {
			"NEO4J_URI",
		},
		"frontend": {
			"PROXY_URL",
			"PROXY_DOMAIN",
			"BACKEND_CONNECTION_URL",
			"BACKEND_CONNECTION_DOMAIN",
		},
		"backend": {
			"PROXY_URL",
			"PROXY_DOMAIN",
			"FRONTEND_CONNECTION_URL",
			"FRONTEND_CONNECTION_DOMAIN",
			"DATABASE_CONNECTION_STRING",
		},
		"llm": {
			"PROXY_URL",
			"PROXY_DOMAIN",
			"LLM_CONNECTION_URL",
			"LLM_CONNECTION_DOMAIN",
		},
	}

	for _, pod := range template.Application.Template.Pods {
		if err := validatePodType(pod.Type); err != nil {
			return err
		}

		if pod.Name == "" {
			return nexerrors.NewValidationError(
				fmt.Sprintf("pod of type %s has no name", pod.Type),
				&nexerrors.ValidationContext{
					Field:       fmt.Sprintf("application.template.pods[].name"),
					ActualValue: "",
					ResolutionHints: []string{
						"Add a name field to the pod configuration",
						"The name should be a valid identifier for the pod",
					},
				},
			)
		}

		if pod.Tag == "" {
			return nexerrors.NewValidationError(
				fmt.Sprintf("pod %s has no tag", pod.Name),
				&nexerrors.ValidationContext{
					Field:       fmt.Sprintf("application.template.pods[].tag"),
					ActualValue: "",
					ResolutionHints: []string{
						"Add a tag field with a valid Docker image tag",
						"The tag should match the pod's type and version",
					},
				},
			)
		}

		// Validate required environment variables
		if vars, ok := requiredVars[pod.Type]; ok {
			for _, required := range vars {
				found := false
				for _, v := range pod.Vars {
					if v.Key == required {
						found = true
						break
					}
				}
				if !found {
					return nexerrors.NewValidationError(
						fmt.Sprintf("pod %s is missing required environment variable %s", pod.Name, required),
						&nexerrors.ValidationContext{
							Field:       fmt.Sprintf("application.template.pods[].vars"),
							MissingVar:  required,
							ResolutionHints: []string{
								fmt.Sprintf("Add %s to the pod's vars section", required),
								fmt.Sprintf("Ensure the value for %s is correctly set", required),
							},
						},
					)
				}
			}
		}
	}

	return nil
}

func validateBuild(template *NexlayerYAML) error {
	build := &template.Application.Template.Build
	if build.Command == "" {
		// Default to npm since it's most common
		build.Command = "npm install && npm run build"
	}
	if build.Output == "" {
		build.Output = "dist"
	}
	return nil
}

func validatePodType(podType string) error {
	validTypes := map[string]bool{
		"react":    true,
		"angular":  true,
		"vue":      true,
		"express":  true,
		"django":   true,
		"fastapi":  true,
		"postgres": true,
		"mongodb":  true,
		"redis":    true,
		"neo4j":    true,
		"nginx":    true,
		"llm":      true,
	}
	if !validTypes[strings.ToLower(podType)] {
		return nexerrors.NewValidationError(
			fmt.Sprintf("invalid pod type: %s", podType),
			&nexerrors.ValidationContext{
				Field:       "application.template.pods[].type",
				ActualValue: podType,
				AllowedValues: []string{
					"react",
					"angular",
					"vue",
					"express",
					"django",
					"fastapi",
					"postgres",
					"mongodb",
					"redis",
					"neo4j",
					"nginx",
					"llm",
				},
				ResolutionHints: []string{
					fmt.Sprintf("Replace '%s' with a supported type", podType),
					"Check the pod's type and ensure it matches the application's requirements",
				},
			},
		)
	}
	return nil
}
