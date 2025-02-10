// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

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

// TemplateRequest represents a request to generate a Nexlayer deployment template.
type TemplateRequest struct {
	ProjectName    string
	TemplateType   string
	RequiredFields map[string]interface{}
}

// NexlayerYAML represents the structure of a Nexlayer deployment template.
type NexlayerYAML struct {
	Application Application `yaml:"application"`
}

// Application represents the application-level configuration.
type Application struct {
	Name          string         `yaml:"name"`
	URL           string         `yaml:"url,omitempty"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
	Pods          []PodConfig    `yaml:"pods"`
}

// RegistryLogin contains authentication details for private registries.
type RegistryLogin struct {
	Registry            string `yaml:"registry"`
	Username            string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// PodConfig represents a pod configuration.
type PodConfig struct {
	Name         string    `yaml:"name"`
	Path         string    `yaml:"path,omitempty"`
	Image        string    `yaml:"image"`
	Volumes      []Volume  `yaml:"volumes,omitempty"`
	Secrets      []Secret  `yaml:"secrets,omitempty"`
	Vars         []VarPair `yaml:"vars,omitempty"`
	ServicePorts []int     `yaml:"servicePorts,omitempty"`
}

// Volume represents a storage volume configuration.
type Volume struct {
	Name      string `yaml:"name"`
	Size      string `yaml:"size"`
	MountPath string `yaml:"mountPath"`
}

// Secret represents a secret file configuration.
type Secret struct {
	Name      string `yaml:"name"`
	Data      string `yaml:"data"`
	MountPath string `yaml:"mountPath"`
	FileName  string `yaml:"fileName"`
}

// VarPair represents an environment variable key-value pair.
type VarPair struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// llmYamlPrompt provides structured instructions for AI-generated Nexlayer templates.
const llmYamlPrompt = `You are an expert in cloud automation for Nexlayer AI Cloud Platform.
Generate a deployment template YAML that seamlessly integrates into Nexlayer Cloud.

Overall Template Structure:
application:
  name: "<deployment-name>"
  url: "<permanent-domain>"
  registryLogin:
    registry: "<docker-registry-url>"
    username: "<registry-username>"
    personalAccessToken: "<registry-access-token>"
  pods:
    - name: "<pod-name>"
      path: "<public-route-path>"
      image: "<docker-image>"
      volumes:
        - name: "<volume-name>"
          size: "<volume-size>"
          mountPath: "<mount-directory>"
      secrets:
        - name: "<secret-name>"
          data: "<base64-encoded-data>"
          mountPath: "<secret-directory>"
          fileName: "<secret-file-name>"
      vars:
        - key: "<env-var-key>"
          value: "<env-var-value>"
      servicePorts: ["<port-number>"]

Supported Pod Types:
- **Frontend**: react, angular, vue
- **Backend**: express, django, fastapi
- **Database**: mongodb, postgres, redis
- **Other Services**: nginx (proxy/load balancer), llm (AI workloads)

Example:
application:
  name: "my-ai-app"
  url: "https://my-ai-app.nexlayer.ai"
  registryLogin:
    registry: "ghcr.io/nexlayer"
    username: "nexlayer-user"
    personalAccessToken: "ghp_xxx"
  pods:
    - name: "backend"
      path: "/"
      image: "ghcr.io/nexlayer/backend-app:latest"
      servicePorts: [3000]
`

// GenerateTemplate uses AI assistance to generate a valid Nexlayer deployment template.
func GenerateTemplate(ctx context.Context, req TemplateRequest) (string, error) {
	provider := GetPreferredProvider(ctx, CapDeploymentAssistance)
	if provider == nil {
		return "", fmt.Errorf("no AI provider available for template generation")
	}

	// Create base template with project name
	template := NexlayerYAML{
		Application: Application{
			Name: req.ProjectName,
			Pods: []PodConfig{},
		},
	}

	// Get AI assistance for template generation
	aiResponse, err := provider.GenerateText(ctx, llmYamlPrompt)
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %v", err)
	}

	// Parse AI response and merge with base template
	var aiTemplate NexlayerYAML
	if err := yaml.Unmarshal([]byte(aiResponse), &aiTemplate); err != nil {
		return "", fmt.Errorf("failed to parse AI response: %v", err)
	}

	// Merge AI suggestions with base template
	template.Application.Pods = aiTemplate.Application.Pods
	if aiTemplate.Application.URL != "" {
		template.Application.URL = aiTemplate.Application.URL
	}
	if aiTemplate.Application.RegistryLogin != nil {
		template.Application.RegistryLogin = aiTemplate.Application.RegistryLogin
	}

	// Marshal final template
	data, err := yaml.Marshal(&template)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %v", err)
	}

	return string(data), nil
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
			stackType, components := DetectStack(".")
			
			req := TemplateRequest{
				ProjectName: appName,
				TemplateType: stackType,
				RequiredFields: map[string]interface{}{
					"components": components,
				},
			}
			
			yamlOut, err := GenerateTemplate(cmd.Context(), req)
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

// DetectStack analyzes a directory to determine its stack type and components.
func DetectStack(dir string) (string, []string) {
	var components []string

	// Check for package.json (Node.js projects)
	if data, err := os.ReadFile(filepath.Join(dir, "package.json")); err == nil {
		var pkg struct {
			Dependencies map[string]string `json:"dependencies"`
		}
		if err := json.Unmarshal(data, &pkg); err == nil {
			// Frontend frameworks
			for _, fw := range []string{"react", "vue", "angular"} {
				if _, has := pkg.Dependencies[fw]; has {
					components = append(components, fw)
				}
			}
			// Backend frameworks
			for _, fw := range []string{"express", "nest", "koa"} {
				if _, has := pkg.Dependencies[fw]; has {
					components = append(components, fw)
				}
			}
			// Database clients
			for _, db := range []string{"mongodb", "pg", "mysql2", "redis"} {
				if _, has := pkg.Dependencies[db]; has {
					components = append(components, strings.TrimSuffix(db, "2"))
				}
			}
		}
	}

	// Check for requirements.txt (Python projects)
	if data, err := os.ReadFile(filepath.Join(dir, "requirements.txt")); err == nil {
		content := string(data)
		// Backend frameworks
		for _, fw := range []string{"fastapi", "django", "flask"} {
			if strings.Contains(content, fw) {
				components = append(components, fw)
			}
		}
		// Database clients
		for _, db := range []string{"psycopg2", "pymongo", "redis", "mysql-connector"} {
			if strings.Contains(content, db) {
				switch db {
				case "psycopg2":
					components = append(components, "postgres")
				case "mysql-connector":
					components = append(components, "mysql")
				default:
					components = append(components, strings.Split(db, "-")[0])
				}
			}
		}
		// AI/ML dependencies
		for _, ml := range []string{"tensorflow", "pytorch", "transformers"} {
			if strings.Contains(content, ml) {
				components = append(components, "llm")
				break
			}
		}
	}

	// Check for go.mod (Go projects)
	if data, err := os.ReadFile(filepath.Join(dir, "go.mod")); err == nil {
		content := string(data)
		// Backend frameworks
		for _, fw := range []string{"gin-gonic/gin", "gorilla/mux", "labstack/echo"} {
			if strings.Contains(content, fw) {
				components = append(components, strings.Split(fw, "/")[1])
			}
		}
	}

	// Check for pom.xml or build.gradle (Java projects)
	for _, f := range []string{"pom.xml", "build.gradle"} {
		if data, err := os.ReadFile(filepath.Join(dir, f)); err == nil {
			content := string(data)
			// Spring Boot
			if strings.Contains(content, "spring-boot") {
				components = append(components, "spring")
			}
			break
		}
	}

	// Determine stack type based on detected components
	stackType := "unknown"
	if containsAny(components, "react", "vue", "angular") {
		stackType = "frontend"
	} else if containsAny(components, "express", "fastapi", "django", "gin", "echo", "spring") {
		stackType = "backend"
	} else if containsAny(components, "postgres", "mongodb", "redis", "mysql") {
		stackType = "database"
	} else if containsAny(components, "llm") {
		stackType = "ai"
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
