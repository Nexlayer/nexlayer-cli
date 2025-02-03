package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	nexerrors "github.com/Nexlayer/nexlayer-cli/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// NexlayerYAML represents the structure of a Nexlayer deployment template.
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

// yamlPrompt is the prompt template for AI assistants.
const yamlPrompt = `As a Nexlayer AI assistant, generate a deployment template YAML that follows these requirements:

1. Structure Requirements:
   - Must include application.template.name, deploymentName, and registryLogin.
   - Must define a pods array with application components.

2. Pod Configuration:
   - Each pod needs: type, name, tag, privateTag, and a vars array.
   - Supported pod types:
     * Frontend: react, angular, vue
     * Backend: express, django, fastapi
     * Database: postgres, mongodb, redis, neo4j
     * Others: nginx, llm

3. Environment Variables:
   - Database pods: Include DATABASE_CONNECTION_STRING.
   - Frontend/Backend: Include appropriate connection URLs.
   - LLM pods: Include model-specific variables.

4. HTTP Exposure:
   - Set exposeHttp: true for web-accessible components.

Application Details:
- Name: %s
- Stack Type: %s
- Components: %s

Output the YAML content only, no explanations.`

// GenerateYAML generates a Nexlayer deployment template using available AI assistants.
func GenerateYAML(appName string, stackType string, components []string) (string, error) {
	provider := GetPreferredProvider(context.Background(), CapCodeGeneration)
	// Format the prompt with application details.
	prompt := fmt.Sprintf(yamlPrompt, appName, stackType, strings.Join(components, ", "))
	if provider == nil {
		// Fallback: Use a default template if no provider is active.
		return generateDefaultTemplate(appName, stackType, components)
	}

	// In a production scenario, an API call would be made to provider.Endpoint here.
	// For now, simulate the response with a mock.
	rawYAML := mockGenerateYAML(appName, stackType, components, prompt)

	// Validate and fix the generated YAML.
	return validateAndFixYAML(rawYAML)
}

// TemplateRequest represents a request to generate a deployment template
type TemplateRequest struct {
	ProjectName    string                 `json:"projectName"`
	TemplateType   string                 `json:"templateType"`
	RequiredFields map[string]interface{} `json:"requiredFields"`
}

// GenerateTemplate generates a deployment template using AI assistance
func GenerateTemplate(ctx context.Context, req TemplateRequest) (string, error) {
	provider := GetPreferredProvider(ctx, CapCodeGeneration)
	if provider == nil {
		return "", fmt.Errorf("no AI provider available for template generation")
	}

	// For now, use the mock implementation
	yamlStr := mockGenerateYAML(req.ProjectName, req.TemplateType, []string{}, "")
	return yamlStr, nil
}

// generateDefaultTemplate creates a basic template based on detected components.
func generateDefaultTemplate(appName string, stackType string, components []string) (string, error) {
	template := NexlayerYAML{}
	template.Application.Template.Name = appName
	template.Application.Template.DeploymentName = appName
	template.Application.Template.RegistryLogin.Registry = "ghcr.io"
	template.Application.Template.RegistryLogin.Username = "<your-username>"
	template.Application.Template.RegistryLogin.PersonalAccessToken = "<your-pat>"

	// For each detected component, add a pod with default configuration.
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

		// For database-type pods, add a default DATABASE_CONNECTION_STRING.
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

	// Set the build command based on the stack type.
	template.Application.Template.Build.Command = getBuildCommand(stackType)
	template.Application.Template.Build.Output = "dist"

	// Marshal the structure to YAML.
	out, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to generate template: %w", err)
	}
	return string(out), nil
}

// Helper functions for template generation

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
	case "react", "angular", "vue", "express":
		return "node:18"
	case "django", "fastapi":
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

// NewCommand creates the "ai" command with its subcommands.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-powered features for Nexlayer",
		Long: `AI-powered features for Nexlayer CLI.
Provides intelligent assistance for template generation, debugging, and optimization.`,
	}

	// Add subcommands
	cmd.AddCommand(
		newGenerateCommand(),
		newDetectCommand(),
	)

	return cmd
}

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [app-name]",
		Short: "Generate deployment template using AI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			// Detect stack type and components in the current directory
			stackType, components := detectStack(".")
			yamlOut, err := GenerateYAML(appName, stackType, components)
			if err != nil {
				return err
			}
			// Print the generated YAML template
			fmt.Println(yamlOut)
			return nil
		},
	}
}

func newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect available AI assistants",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider := GetPreferredProvider(cmd.Context(), CapCodeGeneration)
			if provider == nil {
				fmt.Println("❌ No AI assistants detected")
				return nil
			}
			fmt.Printf("✅ Detected AI assistant: %s\n", provider.Name)
			fmt.Printf("   Description: %s\n", provider.Description)
			return nil
		},
	}
}

// detectStack inspects the given directory to determine the project's stack type and components.
func detectStack(dir string) (string, []string) {
	var components []string

	// Check for package.json (Node.js projects)
	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var pkg struct {
			Dependencies    map[string]string `json:"dependencies"`
			DevDependencies map[string]string `json:"devDependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			// Detect frontend frameworks
			if _, hasReact := pkg.Dependencies["react"]; hasReact {
				components = append(components, "react")
			}
			if _, hasVue := pkg.Dependencies["vue"]; hasVue {
				components = append(components, "vue")
			}
			if _, hasAngular := pkg.Dependencies["@angular/core"]; hasAngular {
				components = append(components, "angular")
			}
			// Detect backend frameworks
			if _, hasExpress := pkg.Dependencies["express"]; hasExpress {
				components = append(components, "express")
			}
			// Detect LLM frameworks
			if _, hasLangchain := pkg.Dependencies["langchain"]; hasLangchain {
				components = append(components, "llm")
			}
			if _, hasOpenAI := pkg.Dependencies["openai"]; hasOpenAI {
				components = append(components, "llm")
			}
		}
	}

	// Check for requirements.txt (Python projects)
	if data, err := os.ReadFile(filepath.Join(dir, "requirements.txt")); err == nil {
		content := string(data)
		// Detect Python frameworks
		if strings.Contains(content, "fastapi") {
			components = append(components, "fastapi")
		}
		if strings.Contains(content, "django") {
			components = append(components, "django")
		}
		// Detect LLM frameworks
		if strings.Contains(content, "langchain") {
			components = append(components, "llm")
		}
		if strings.Contains(content, "openai") {
			components = append(components, "llm")
		}
		if strings.Contains(content, "anthropic") {
			components = append(components, "llm")
		}
	}

	// Determine stack type based on components
	stackType := "unknown"
	if containsAny(components, "react", "vue", "angular") {
		if containsAny(components, "express") {
			stackType = "node"
		} else if containsAny(components, "fastapi", "django") {
			stackType = "python"
		}
	}
	if containsAny(components, "llm") {
		if containsAny(components, "langchain") {
			stackType = "langchain"
		} else if containsAny(components, "openai") {
			stackType = "openai"
		} else if containsAny(components, "anthropic") {
			stackType = "anthropic"
		}
	}

	return stackType, components
}

// containsAny returns true if any of the values are present in the slice.
func containsAny(slice []string, values ...string) bool {
	for _, v := range values {
		for _, s := range slice {
			if strings.ToLower(s) == strings.ToLower(v) {
				return true
			}
		}
	}
	return false
}

// getLLMVarsForTemplate returns environment variables based on the LLM template type
func getLLMVarsForTemplate(templateType string) []struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
} {
	switch templateType {
	case "langchain-nextjs", "langchain-fastapi":
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "OPENAI_API_KEY", Value: "${OPENAI_API_KEY}"},
			{Key: "LANGCHAIN_TRACING_V2", Value: "true"},
			{Key: "LANGCHAIN_ENDPOINT", Value: "https://api.smith.langchain.com"},
		}
	case "openai-node", "openai-py":
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "OPENAI_API_KEY", Value: "${OPENAI_API_KEY}"},
			{Key: "OPENAI_ORG_ID", Value: "${OPENAI_ORG_ID}"},
		}
	case "llama-node", "llama-py":
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "MODEL_PATH", Value: "/models"},
			{Key: "CUDA_VISIBLE_DEVICES", Value: "0"},
		}
	case "vertex-ai":
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "GOOGLE_APPLICATION_CREDENTIALS", Value: "${GOOGLE_APPLICATION_CREDENTIALS}"},
			{Key: "PROJECT_ID", Value: "${GCP_PROJECT_ID}"},
		}
	case "huggingface":
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "HUGGINGFACE_API_KEY", Value: "${HUGGINGFACE_API_KEY}"},
			{Key: "MODEL_ID", Value: "gpt2"},
		}
	case "anthropic-py", "anthropic-js":
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "ANTHROPIC_API_KEY", Value: "${ANTHROPIC_API_KEY}"},
			{Key: "ANTHROPIC_MODEL", Value: "claude-2"},
		}
	default:
		return []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		}{
			{Key: "LLM_CONNECTION_URL", Value: "http://localhost:8000"},
		}
	}
}

// mockGenerateYAML simulates an AI response for YAML generation.
func mockGenerateYAML(appName string, stackType string, components []string, prompt string) string {
	// Use prompt to customize the template if provided
	customVars := make(map[string]string)
	if prompt != "" {
		// Parse the prompt for any custom configuration
		// This is a simple example - in a real implementation, we'd use NLP
		if strings.Contains(prompt, "production") {
			customVars["NODE_ENV"] = "production"
		}
		if strings.Contains(prompt, "development") {
			customVars["NODE_ENV"] = "development"
		}
	}

	// Build a base YAML using Nexlayer standards.
	yamlStr := fmt.Sprintf(`application:
  template:
    name: "%s"
    deploymentName: "%s"
    registryLogin:
      registry: ghcr.io
      username: ${GITHUB_USERNAME}
      personalAccessToken: ${GITHUB_PAT}
    pods:`, appName, appName)

	// Handle LLM-specific templates
	if strings.Contains(stackType, "langchain") || 
	   strings.Contains(stackType, "openai") || 
	   strings.Contains(stackType, "llama") || 
	   strings.Contains(stackType, "vertex-ai") || 
	   strings.Contains(stackType, "huggingface") || 
	   strings.Contains(stackType, "anthropic") {
		
		// Add frontend for templates that need it
		if strings.Contains(stackType, "-nextjs") || strings.Contains(stackType, "-node") {
			yamlStr += `
      - type: frontend
        name: next-app
        tag: "node:18"
        exposeHttp: true
        vars:`
			for _, v := range getLLMVarsForTemplate(stackType) {
				yamlStr += fmt.Sprintf(`
          - key: %s
            value: "%s"`, v.Key, v.Value)
			}
		}

		// Add backend/LLM service
		yamlStr += `
      - type: llm
        name: llm-service
        tag: "python:3.9"
        exposeHttp: true
        vars:`
		for _, v := range getLLMVarsForTemplate(stackType) {
			yamlStr += fmt.Sprintf(`
          - key: %s
            value: "%s"`, v.Key, v.Value)
		}

		// Add vector database for RAG applications
		if strings.Contains(stackType, "langchain") {
			yamlStr += `
      - type: database
        name: postgres-vector
        tag: "ankane/pgvector:latest"
        vars:
          - key: POSTGRES_DB
            value: "vectorstore"
          - key: POSTGRES_USER
            value: "postgres"
          - key: POSTGRES_PASSWORD
            value: "${POSTGRES_PASSWORD}"
          - key: DATABASE_CONNECTION_STRING
            value: "postgresql://postgres:${POSTGRES_PASSWORD}@postgres-vector:5432/vectorstore"`
		}
	}

	// For each detected component, add a pod with default configuration.
	for _, comp := range components {
		if !strings.Contains(stackType, "langchain") && !strings.Contains(stackType, "openai") && !strings.Contains(stackType, "llama") && !strings.Contains(stackType, "vertex-ai") && !strings.Contains(stackType, "huggingface") && !strings.Contains(stackType, "anthropic") {
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

			yamlStr += fmt.Sprintf(`
      - type: %s
        name: %s
        tag: "%s"
        privateTag: %v
        exposeHttp: %v
        vars:`, pod.Type, pod.Name, pod.Tag, pod.PrivateTag, pod.ExposeHTTP)
			for _, v := range pod.Vars {
				yamlStr += fmt.Sprintf(`
          - key: %s
            value: "%s"`, v.Key, v.Value)
			}
		}
	}

	// Append build configuration.
	yamlStr += `
    build:
      command: "npm install && npm run build"
      output: "build"`

	return yamlStr
}

// validateAndFixYAML validates the generated YAML against Nexlayer standards and fixes common issues.
func validateAndFixYAML(yamlString string) (string, error) {
	var template NexlayerYAML
	if err := yaml.Unmarshal([]byte(yamlString), &template); err != nil {
		return "", fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	// Validate template section.
	if err := validateTemplate(&template); err != nil {
		return "", err
	}
	// Validate pods.
	if err := validatePods(&template); err != nil {
		return "", err
	}
	// Validate build configuration.
	if err := validateBuild(&template); err != nil {
		return "", err
	}

	// Marshal the fixed template back to YAML.
	out, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	return string(out), nil
}

// Validation helper functions

func validateTemplate(template *NexlayerYAML) error {
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
			},
		)
	}
	if template.Application.Template.DeploymentName == "" {
		template.Application.Template.DeploymentName = template.Application.Template.Name
	}

	// Set default registry login values if missing.
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
			},
		)
	}

	for _, pod := range template.Application.Template.Pods {
		if err := validatePodType(pod.Type); err != nil {
			return err
		}

		if pod.Name == "" {
			return nexerrors.NewValidationError(
				fmt.Sprintf("pod of type %s has no name", pod.Type),
				&nexerrors.ValidationContext{
					Field:       "application.template.pods[].name",
					ActualValue: "",
					ResolutionHints: []string{
						"Add a name field to the pod configuration",
						"Ensure the name is a valid identifier for the pod",
					},
				},
			)
		}

		if pod.Tag == "" {
			return nexerrors.NewValidationError(
				fmt.Sprintf("pod %s has no tag", pod.Name),
				&nexerrors.ValidationContext{
					Field:       "application.template.pods[].tag",
					ActualValue: "",
					ResolutionHints: []string{
						"Add a tag field with a valid Docker image tag",
						"Ensure the tag matches the pod's type and version",
					},
				},
			)
		}
	}

	return nil
}

func validateBuild(template *NexlayerYAML) error {
	build := &template.Application.Template.Build
	if build.Command == "" {
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
				Field:         "application.template.pods[].type",
				ActualValue:   podType,
				AllowedValues: []string{"react", "angular", "vue", "express", "django", "fastapi", "postgres", "mongodb", "redis", "neo4j", "nginx", "llm"},
				ResolutionHints: []string{
					fmt.Sprintf("Replace '%s' with a supported type", podType),
					"Check the pod's type and ensure it matches the application's requirements",
				},
			},
		)
	}
	return nil
}
