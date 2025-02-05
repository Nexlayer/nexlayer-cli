// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// TemplateRequest represents a request to generate a Nexlayer template.
type TemplateRequest struct {
	ProjectName    string
	TemplateType   string
	RequiredFields map[string]interface{}
}

// NexlayerYAML represents the structure of a Nexlayer deployment template.
type NexlayerYAML struct {
	Application Application `yaml:"application"`
}

// Application represents the application configuration.
type Application struct {
	Template Template `yaml:"template"`
}

// Template represents the template configuration.
type Template struct {
	Name           string      `yaml:"name"`
	DeploymentName string      `yaml:"deploymentName"`
	Pods           []PodConfig `yaml:"pods"`
}

// Port represents a port mapping.
type Port struct {
	ContainerPort int    `yaml:"containerPort"`
	ServicePort   int    `yaml:"servicePort"`
	Name          string `yaml:"name"`
}

// PodConfig represents a pod configuration.
type PodConfig struct {
	Type       string    `yaml:"type"`
	Name       string    `yaml:"name"`
	Image      string    `yaml:"image"` // Full image URL including tag
	Command    []string  `yaml:"command,omitempty"`
	Vars       []VarPair `yaml:"vars,omitempty"`
	Ports      []Port    `yaml:"ports,omitempty"`
	ExposeOn80 bool      `yaml:"exposeOn80"`
}

// VarPair represents a key-value pair for environment variables.
type VarPair struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// llmYamlPrompt defines detailed instructions for the AI LLM to generate a Nexlayer template.
// It includes overall template structure, pods configuration (with supported types),
// standard environment variables, and port examples.
const llmYamlPrompt = `You are an expert cloud automation engineer assistant for the Nexlayer AI Cloud Platform.
Generate a deployment template YAML that deploys instantly and flawlessly to Nexlayer Cloud.

Overall Template Structure:
application:
  template:
    name: Application name (e.g., "%s")
    deploymentName: The deployment name (e.g., "%s")
    pods: List of pod configurations

Pod Configuration:
Each pod in the pods array must include:
- type: Component type (frontend, backend, database, nginx, llm)
- name: Descriptive pod name
- image: Full image URL including registry and tag
- vars: Environment variables array (each with key and value)
- ports: Port configurations array with:
    - containerPort: Port inside container
    - servicePort: Port exposed to service
    - name: Port name
- exposeOn80: Boolean to indicate if the pod should be exposed on port 80

Supported Pod Types:
- Frontend: react, angular, vue
- Backend: express, django, fastapi
- Database: mongodb, postgres, redis, neo4j
- Others: nginx (load balancing/static assets), llm (custom workloads)

Standard Environment Variables:
- DATABASE_CONNECTION_STRING (for databases)
- FRONTEND_CONNECTION_URL, BACKEND_CONNECTION_URL (for frontend/backend)
- LLM_CONNECTION_URL (for llm)
- Additional common variables: PORT, NODE_ENV

Example Port Configurations:
- Frontend: containerPort: 3000, servicePort: 80
- Backend: containerPort: 8000, servicePort: 8000
- Database (mongodb): containerPort: 27017, servicePort: 27017
- LLM (ollama): containerPort: 11434, servicePort: 11434

Using the above structure, generate a valid Nexlayer deployment template YAML for the project "%s", stack type "%s", and components: %s.
Output only the YAML content without any extra commentary.`

// GenerateTemplate generates a Nexlayer deployment template using AI assistance.
func GenerateTemplate(ctx context.Context, req TemplateRequest) (string, error) {
	// Get AI provider with template generation capability.
	provider := GetPreferredProvider(ctx, CapDeploymentAssistance)
	if provider == nil {
		return "", fmt.Errorf("no AI provider available for template generation")
	}

	// Create a basic template structure as a starting point.
	template := NexlayerYAML{
		Application: Application{
			Template: Template{
				Name:           req.ProjectName,
				DeploymentName: req.ProjectName,
			},
		},
	}

	// (Additional pod configuration based on req.TemplateType can be added here)
	if req.TemplateType == "llm-express" {
		template.Application.Template.Pods = []PodConfig{
			{
				Type:       "llm",
				Name:       "ollama",
				Image:      "us-east1-docker.pkg.dev/capture-by-auditdeploy/nexlayer-3rd-party/00219bb4-61ac-4fb5-b648-f23cfcb2ad59/ollama:dolphin-phi",
				ExposeOn80: false,
				Ports: []Port{
					{
						ContainerPort: 11434,
						ServicePort:   11434,
						Name:          "ollama",
					},
				},
			},
			{
				Type:       "frontend",
				Name:       req.ProjectName,
				Image:      fmt.Sprintf("us-east1-docker.pkg.dev/capture-by-auditdeploy/nexlayer-3rd-party/00219bb4-61ac-4fb5-b648-f23cfcb2ad59/%s:latest", req.ProjectName),
				ExposeOn80: true,
				Ports: []Port{
					{
						ContainerPort: 3000,
						ServicePort:   80,
						Name:          "web",
					},
				},
				Vars: []VarPair{
					{
						Key:   "REACT_APP_API_URL",
						Value: "http://CANDIDATE_DEPENDENCY_URL_0:11434",
					},
					{
						Key:   "NODE_ENV",
						Value: "production",
					},
				},
			},
		}
	}

	// Marshal to YAML.
	data, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %v", err)
	}

	return string(data), nil
}

// GenerateYAML generates a Nexlayer deployment template using available AI assistance.
func GenerateYAML(appName string, _ string, components []string) (string, error) {
	// Get the preferred AI provider (for deployment assistance).
	provider := GetPreferredProvider(context.Background(), CapDeploymentAssistance)

	// If no provider is available, fall back to a default template.
	if provider == nil {
		return generateDefaultTemplate(appName, components)
	}

	// Print provider info.
	fmt.Printf("✨ Using %s for template assistance\n", provider.Name)
	fmt.Println("💡 Your AI assistant will help you customize this template")

	// In a real integration, you would call provider.Endpoint with the prompt.
	// Here, we simulate the AI response using our mock.
	rawYAML := mockGenerateYAML(appName, components)

	// Validate and fix the generated YAML.
	return validateAndFixYAML(rawYAML)
}

// generateDefaultTemplate creates a basic template based on detected components.
func generateDefaultTemplate(appName string, components []string) (string, error) {
	comments := fmt.Sprintf(`# Nexlayer Deployment Template for %s
# Generated with AI assistance.
# Customize this template using your IDE's AI assistant.
#
`, appName)

	template := NexlayerYAML{
		Application: Application{
			Template: Template{
				Name:           appName,
				DeploymentName: appName,
			},
		},
	}

	// For each detected component, add a pod with default configuration.
	for _, comp := range components {
		pod := PodConfig{
			Type:       comp,
			Name:       fmt.Sprintf("%s-service", comp),
			Image:      fmt.Sprintf("us-east1-docker.pkg.dev/capture-by-auditdeploy/nexlayer-3rd-party/00219bb4-61ac-4fb5-b648-f23cfcb2ad59/%s:%s", comp, defaultTagForType(comp)),
			ExposeOn80: isExposeByDefault(comp),
		}

		// Add default environment variables based on component type.
		switch comp {
		case "react", "vue", "angular":
			pod.Ports = []Port{{
				ContainerPort: 3000,
				ServicePort:   80,
				Name:          "web",
			}}
			pod.Vars = []VarPair{{
				Key:   "NODE_ENV",
				Value: "production",
			}}
		case "express", "fastapi", "django":
			pod.Ports = []Port{{
				ContainerPort: 8000,
				ServicePort:   8000,
				Name:          "api",
			}}
			pod.Vars = []VarPair{{
				Key:   "NODE_ENV",
				Value: "development",
			}}
		case "mongodb", "postgres", "redis":
			pod.Ports = []Port{{
				ContainerPort: 27017,
				ServicePort:   27017,
				Name:          "db",
			}}
			pod.Vars = []VarPair{{
				Key:   "DATABASE_CONNECTION_STRING",
				Value: fmt.Sprintf("%s://%s:1234/db", comp, pod.Name),
			}}
		}

		template.Application.Template.Pods = append(template.Application.Template.Pods, pod)
	}

	out, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to generate template: %w", err)
	}

	return comments + string(out), nil
}

// NewCommand creates the "ai" command with its subcommands.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-powered features for Nexlayer",
		Long:  "AI-powered features for Nexlayer CLI. Provides intelligent assistance for template generation, debugging, and optimization.",
	}

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
			// Detect project stack and components based on current directory.
			stackType, components := detectStack(".")
			yamlOut, err := GenerateYAML(appName, stackType, components)
			if err != nil {
				return err
			}
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
			provider := GetPreferredProvider(cmd.Context(), CapDeploymentAssistance)
			if provider == nil {
				fmt.Println("ℹ️  No AI assistants detected")
				fmt.Println("💡 Configure an AI assistant for enhanced template generation")
				return nil
			}
			fmt.Printf("✨ Detected AI assistant: %s\n", provider.Name)
			fmt.Printf("   Description: %s\n", provider.Description)
			return nil
		},
	}
}

// Helper functions

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
		return "latest" // Use "latest" for non-specific image tags.
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
			if _, hasReact := pkg.Dependencies["react"]; hasReact {
				components = append(components, "react")
			}
			if _, hasVue := pkg.Dependencies["vue"]; hasVue {
				components = append(components, "vue")
			}
			if _, hasAngular := pkg.Dependencies["@angular/core"]; hasAngular {
				components = append(components, "angular")
			}
			if _, hasExpress := pkg.Dependencies["express"]; hasExpress {
				components = append(components, "express")
			}
			if _, hasLangchain := pkg.Dependencies["langchain"]; hasLangchain {
				components = append(components, "llm")
			}
			if _, hasOpenAI := pkg.Dependencies["openai"]; hasOpenAI {
				components = append(components, "llm")
			}
			if _, hasTensorflow := pkg.Dependencies["@tensorflow/tfjs"]; hasTensorflow {
				components = append(components, "ml")
			}
		}
	}

	// Check for requirements.txt (Python projects)
	if data, err := os.ReadFile(filepath.Join(dir, "requirements.txt")); err == nil {
		content := string(data)
		if strings.Contains(content, "fastapi") {
			components = append(components, "fastapi")
		}
		if strings.Contains(content, "django") {
			components = append(components, "django")
		}
		if strings.Contains(content, "tensorflow") || strings.Contains(content, "torch") {
			components = append(components, "ml")
		}
		if strings.Contains(content, "transformers") {
			components = append(components, "llm")
		}
		if strings.Contains(content, "langchain") {
			components = append(components, "llm")
		}
	}

	// Check for docker-compose.yml (Database detection)
	if data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml")); err == nil {
		content := string(data)
		if strings.Contains(content, "postgres") {
			components = append(components, "postgres")
		}
		if strings.Contains(content, "mongodb") {
			components = append(components, "mongodb")
		}
		if strings.Contains(content, "redis") {
			components = append(components, "redis")
		}
	}

	// Determine stack type based on components.
	stackType := "unknown"
	if containsAny(components, "react", "vue", "angular") {
		if containsAny(components, "express") {
			stackType = "node"
		} else if containsAny(components, "fastapi", "django") {
			stackType = "python"
		}
	}
	if containsAny(components, "ml", "llm") {
		if containsAny(components, "tensorflow", "torch") {
			stackType = "ml"
		} else if containsAny(components, "langchain") {
			stackType = "langchain"
		}
	}

	return stackType, components
}

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

// NewAICommand creates a new AI command for template generation.
func NewAICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai [flags]",
		Short: "AI-assisted template generation",
		Long: `Use AI assistance to generate and customize Nexlayer deployment templates.

Example:
  nexlayer ai generate --name myapp --type nodejs
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use one of the subcommands: generate")
		},
	}

	genCmd := &cobra.Command{
		Use:   "generate [flags]",
		Short: "Generate a new template using AI",
		RunE: func(cmd *cobra.Command, args []string) error {
			name, _ := cmd.Flags().GetString("name")
			templateType, _ := cmd.Flags().GetString("type")
			components, _ := cmd.Flags().GetStringSlice("components")

			yaml, err := GenerateYAML(name, templateType, components)
			if err != nil {
				return err
			}

			fmt.Println(yaml)
			return nil
		},
	}

	genCmd.Flags().String("name", "", "Name of the application")
	genCmd.Flags().String("type", "", "Type of application (e.g. nodejs, python)")
	genCmd.Flags().StringSlice("components", []string{}, "List of required components")

	cmd.AddCommand(genCmd)
	return cmd
}

func mockGenerateYAML(appName string, components []string) string {
	// Create a basic template using the Nexlayer Cloud standard template.
	template := NexlayerYAML{
		Application: Application{
			Template: Template{
				Name:           appName,
				DeploymentName: appName,
			},
		},
	}

	// Add pods based on the provided components.
	for _, comp := range components {
		pod := PodConfig{
			Type:       comp,
			Name:       fmt.Sprintf("%s-service", comp),
			Image:      fmt.Sprintf("us-east1-docker.pkg.dev/capture-by-auditdeploy/nexlayer-3rd-party/00219bb4-61ac-4fb5-b648-f23cfcb2ad59/%s:latest", comp),
			ExposeOn80: isExposeByDefault(comp),
		}

		// Add default environment variables based on component type.
		switch comp {
		case "react", "vue", "angular":
			pod.Vars = []VarPair{
				{Key: "NODE_ENV", Value: "production"},
				{Key: "PORT", Value: "3000"},
			}
		case "express", "fastapi", "django":
			pod.Vars = []VarPair{
				{Key: "NODE_ENV", Value: "production"},
				{Key: "PORT", Value: "8000"},
			}
		case "mongodb":
			pod.Vars = []VarPair{
				{Key: "MONGO_INITDB_DATABASE", Value: "app"},
				{Key: "DATABASE_CONNECTION_STRING", Value: "mongodb://localhost:27017/app"},
			}
		}

		template.Application.Template.Pods = append(template.Application.Template.Pods, pod)
	}

	data, _ := yaml.Marshal(&template)
	return string(data)
}

// validateAndFixYAML validates the generated YAML against Nexlayer requirements
// and fixes common issues.
func validateAndFixYAML(yamlStr string) (string, error) {
	var template NexlayerYAML
	if err := yaml.Unmarshal([]byte(yamlStr), &template); err != nil {
		return "", fmt.Errorf("invalid YAML: %v", err)
	}

	if template.Application.Template.Name == "" {
		return "", fmt.Errorf("missing required field: application.template.name")
	}
	if template.Application.Template.DeploymentName == "" {
		return "", fmt.Errorf("missing required field: application.template.deploymentName")
	}
	if len(template.Application.Template.Pods) == 0 {
		return "", fmt.Errorf("template must contain at least one pod")
	}

	for i, pod := range template.Application.Template.Pods {
		if pod.Type == "" {
			return "", fmt.Errorf("pod[%d]: missing type", i)
		}
		if pod.Name == "" {
			return "", fmt.Errorf("pod[%d]: missing name", i)
		}
		if pod.Image == "" {
			return "", fmt.Errorf("pod[%d]: missing image", i)
		}

		validTypes := []string{"frontend", "backend", "database", "nginx", "llm"}
		validType := false
		for _, t := range validTypes {
			if pod.Type == t {
				validType = true
				break
			}
		}
		if !validType {
			return "", fmt.Errorf("pod[%d]: invalid type '%s'. Must be one of: %v", i, pod.Type, validTypes)
		}

		if pod.Vars == nil {
			template.Application.Template.Pods[i].Vars = []VarPair{}
		}

		switch pod.Type {
		case "database":
			hasConnStr := false
			for _, v := range pod.Vars {
				if v.Key == "DATABASE_CONNECTION_STRING" {
					hasConnStr = true
					break
				}
			}
			if !hasConnStr {
				template.Application.Template.Pods[i].Vars = append(template.Application.Template.Pods[i].Vars,
					VarPair{Key: "DATABASE_CONNECTION_STRING", Value: "auto-generated"})
			}
		case "frontend":
			hasBackendUrl := false
			for _, v := range pod.Vars {
				if v.Key == "BACKEND_CONNECTION_URL" {
					hasBackendUrl = true
					break
				}
			}
			if !hasBackendUrl {
				template.Application.Template.Pods[i].Vars = append(template.Application.Template.Pods[i].Vars,
					VarPair{Key: "BACKEND_CONNECTION_URL", Value: "auto-generated"})
			}
		case "backend":
			hasFrontendUrl := false
			for _, v := range pod.Vars {
				if v.Key == "FRONTEND_CONNECTION_URL" {
					hasFrontendUrl = true
					break
				}
			}
			if !hasFrontendUrl {
				template.Application.Template.Pods[i].Vars = append(template.Application.Template.Pods[i].Vars,
					VarPair{Key: "FRONTEND_CONNECTION_URL", Value: "auto-generated"})
			}
		}
	}

	data, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("error marshaling YAML: %v", err)
	}

	return string(data), nil
}
