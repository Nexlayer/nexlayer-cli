// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	tmpl "github.com/Nexlayer/nexlayer-cli/pkg/template"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

// TemplateRequest represents a request to generate a Nexlayer deployment template.
type TemplateRequest struct {
	ProjectName string
	ProjectDir  string
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
	// Create channels for parallel processing
	type aiResult struct {
		response string
		err     error
	}
	type detectionResult struct {
		info *detection.ProjectInfo
		err  error
	}
	aiChan := make(chan aiResult, 1)
	detectChan := make(chan detectionResult, 1)

	// Start stack detection (always needed)
	go func() {
		info, err := DetectStack(req.ProjectDir)
		detectChan <- detectionResult{info, err}
	}()

	// Wait for detection first
	info, err := func() (*detection.ProjectInfo, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-detectChan:
			return result.info, result.err
		}
	}()
	if err != nil {
		return "", fmt.Errorf("failed to detect project info: %v", err)
	}

	// Create template based on detected stack
	var yamlTemplate tmpl.NexlayerYAML

	// Map detected stack to pod type
	var podType tmpl.PodType
	switch info.Type {
	case detection.TypeReact:
		podType = tmpl.React
	case detection.TypePython:
		// Assume Django/FastAPI based on dependencies
		if _, hasDjango := info.Dependencies["django"]; hasDjango {
			podType = tmpl.Django
		} else if _, hasFastAPI := info.Dependencies["fastapi"]; hasFastAPI {
			podType = tmpl.FastAPI
		} else {
			podType = tmpl.Backend
		}
	default:
		// Unknown stack type, use AI for suggestions
		go func() {
			provider := GetPreferredProvider(ctx, CapDeploymentAssistance)
			if provider == nil {
				aiChan <- aiResult{"", fmt.Errorf("no AI provider configured")}
				return
			}
			response, err := provider.GenerateText(ctx, llmYamlPrompt)
			aiChan <- aiResult{response, err}
		}()

		// Wait for AI response
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case result := <-aiChan:
			if result.err != nil {
				// Use generic template if AI fails
				podType = tmpl.Backend
			} else {
				// Parse AI response
				if err := yaml.Unmarshal([]byte(result.response), &yamlTemplate); err != nil {
					// Use generic template if parsing fails
					podType = tmpl.Backend
				}
			}
		}
	}

	// If we're using a predefined template
	if podType != "" {
		// Get default ports for the pod type
		ports := tmpl.DefaultPorts[podType]
		if len(ports) == 0 {
			ports = []tmpl.Port{{ContainerPort: info.Port, ServicePort: info.Port, Name: "app"}}
		}

		// Get default environment variables
		vars := tmpl.DefaultEnvVars[podType]

		// Create template
		yamlTemplate = tmpl.NexlayerYAML{
			Application: tmpl.Application{
				Name: req.ProjectName,
				Pods: []tmpl.Pod{
					{
						Name:  "app",
						Type:  podType,
						Image: fmt.Sprintf("ghcr.io/nexlayer/%s:latest", req.ProjectName),
						Ports: ports,
						Vars:  vars,
					},
				},
			},
		}
	}

	// Add AI environment if detected
	if info.LLMProvider != "" || info.LLMModel != "" {
		// TODO: Update this once we add AI environment to template package
		// template.Application.Environment = &template.Environment{...}
	}

	// Marshal final template
	data, err := yaml.Marshal(&yamlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %v", err)
	}

	return string(data), nil
}

// NewCommand creates the "ai" command with its subcommands.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai [subcommand]",
		Short: "AI-powered features for Nexlayer",
		Long: `AI-powered features for Nexlayer.

Subcommands:
  generate        Generate AI-powered deployment template
  detect          Detect AI assistants & project type

Examples:
  nexlayer ai generate myapp
  nexlayer ai detect`,
	}

	cmd.AddCommand(
		newGenerateCommand(),
		newDetectCommand(),
	)

	return cmd
}

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate <app-name>",
		Short: "Generate AI-powered deployment template",
		Long: `Generate an AI-powered deployment template for your application.

Arguments:
  app-name        Name of your application

Example:
  nexlayer ai generate myapp`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %v", err)
			}

			// Create request with minimal info
			req := TemplateRequest{
				ProjectName: appName,
				ProjectDir: workDir,
			}

			// Generate template (detection and AI will run in parallel)
			yamlOut, err := GenerateTemplate(cmd.Context(), req)
			if err != nil {
				return err
			}

			// Write to nexlayer.yaml
			if err := os.WriteFile("nexlayer.yaml", []byte(yamlOut), 0644); err != nil {
				return fmt.Errorf("failed to write nexlayer.yaml: %v", err)
			}

			fmt.Println("Successfully generated nexlayer.yaml")
			return nil
		},
	}
}

func newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect AI assistants & project type",
		Long: `Detect AI assistants and project type in the current directory.

Example:
  nexlayer ai detect`,
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

// DetectStack analyzes a directory to determine its stack type, components, and AI environment.
func DetectStack(dir string) (*detection.ProjectInfo, error) {
	// Create detector registry
	registry := detection.NewDetectorRegistry()

	// Detect project type and info
	info, err := registry.DetectProject(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to detect project: %v", err)
	}

	// Ensure we have a valid project info
	if info == nil {
		info = &detection.ProjectInfo{
			Type: detection.TypeUnknown,
			Name: filepath.Base(dir),
		}
	}

	// Detect AI IDE and LLM if not already set
	if info.LLMProvider == "" {
		info.LLMProvider = detection.DetectAIIDE()
	}
	if info.LLMModel == "" {
		info.LLMModel = detection.DetectLLMModel()
	}

	return info, nil
}

func containsAny(slice []string, values ...string) bool {
	for _, v := range values {
		for _, s := range slice {
			if strings.EqualFold(s, v) {
				return true
			}
		}
	}
	return false
}
