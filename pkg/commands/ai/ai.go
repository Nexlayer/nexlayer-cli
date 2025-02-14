// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/analysis"
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
Generate a deployment template YAML that seamlessly integrates into Nexlayer Cloud following the Nexlayer YAML Schema Template Documentation (v1.0).

The template must follow this structure:

application:
  name: "<deployment-name>"          # The deployment name (must be lowercase, alphanumeric, '-', '.')
  url: "<permanent-domain>"            # Permanent domain URL (optional; omit if not needed)
  registryLogin:
    registry: "<registry-url>"         # The registry where private images are stored (if required)
    username: "<registry-username>"    # Registry username (if required)
    personalAccessToken: "<registry-access-token>"  # Read-only registry PAT (if required)
  pods:
    - name: "<pod-name>"               # Pod name (must start with a lowercase letter and include only alphanumeric characters, '-', '.')
      path: "<public-route-path>"      # Path to render the pod at (e.g., "/" for front-facing pods; optional for internal services)
      image: "<image-name>:<tag>"        # Docker image for the pod.
                                       # For private images, use the format: '<% REGISTRY %>/some/path/image:tag'
                                       # Do not prefix with "docker.io"—use a simple organization/repository format.
      volumes:
        # Array of volumes to be mounted for this pod.
        - name: "<volume-name>"        # Volume name (lowercase, alphanumeric, '-')
          size: "<volume-size>"        # Required: Volume size (e.g., "1Gi", "500Mi")
          mountPath: "<mount-directory>"  # Required: Must start with '/'
      secrets:
        # Array of secret files for this pod.
        - name: "<secret-name>"        # Secret name (lowercase, alphanumeric, '-')
          data: "<base64-encoded-data>"# Raw or Base64-encoded secret data (e.g., JSON files should be encoded)
          mountPath: "<secret-directory>" # Required: Must start with '/'
          fileName: "<secret-file-name>"  # For example, "secret-file.txt"
      vars:
        # Array of environment variables for this pod.
        - key: "<env-var-key>"         # Environment variable name
          value: "<env-var-value>"     # Its value. Can reference other pods dynamically (e.g., "http://<pod-name>.pod:<port>")
                                      # or use <% URL %> to reference the deployment's base URL.
      servicePorts:
        # Array of ports to expose for this pod.
        - <port-number>               # Exposing a port (shorthand notation is acceptable)
  entrypoint: "<custom-entrypoint>"    # Optional: Custom container entrypoint
  command: "<custom-command>"          # Optional: Custom container command
`

// CallGraph represents the structure of the call graph JSON output
// This is a placeholder structure and should be adjusted based on the actual JSON format
// Example structure:
type CallGraph struct {
	Nodes []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"nodes"`
	Edges []struct {
		Source string `json:"source"`
		Target string `json:"target"`
	} `json:"edges"`
}

// parseCallGraph parses the JSON call graph data into a CallGraph struct
func parseCallGraph(data string) (*CallGraph, error) {
	var graph CallGraph
	err := json.Unmarshal([]byte(data), &graph)
	if err != nil {
		return nil, fmt.Errorf("failed to parse call graph JSON: %v", err)
	}
	return &graph, nil
}

// GenerateTemplate uses AI assistance to generate a valid Nexlayer deployment template.
func GenerateTemplate(ctx context.Context, req TemplateRequest) (string, error) {
	// Create channels for parallel processing
	type aiResult struct {
		response string
		err      error
	}
	type detectionResult struct {
		info *detection.ProjectInfo
		err  error
	}
	type analysisResult struct {
		analysis *analysis.ProjectAnalysis
		err      error
	}

	aiChan := make(chan aiResult, 1)
	detectChan := make(chan detectionResult, 1)
	analyzeChan := make(chan analysisResult, 1)

	// Create a channel for call graph results
	type callGraphResult struct {
		graphData string
		err       error
	}
	callGraphChan := make(chan callGraphResult, 1)

	// Start stack detection (always needed)
	go func() {
		info, err := DetectStack(req.ProjectDir)
		detectChan <- detectionResult{info, err}
	}()

	// Start tree-sitter analysis
	go func() {
		parser := analysis.NewParser()
		analysis, err := parser.AnalyzeProject(ctx, req.ProjectDir)
		analyzeChan <- analysisResult{analysis, err}
	}()

	// Start call graph generation
	go func() {
		// Example command to generate call graph
		cmd := exec.Command("go-callvis", "-focus", req.ProjectDir, "-group", "pkg,type", "-nostd", "-format", "svg", "-o", "callgraph.svg")
		output, err := cmd.CombinedOutput()
		if err != nil {
			callGraphChan <- callGraphResult{"", fmt.Errorf("failed to generate call graph: %v", err)}
			return
		}
		callGraphChan <- callGraphResult{string(output), nil}
	}()

	// Wait for detection, analysis, and call graph
	info, analysis, callGraph, err := func() (*detection.ProjectInfo, *analysis.ProjectAnalysis, string, error) {
		var info *detection.ProjectInfo
		var analysis *analysis.ProjectAnalysis
		var callGraph string
		var err error

		// Wait for detection
		select {
		case <-ctx.Done():
			return nil, nil, "", ctx.Err()
		case result := <-detectChan:
			info = result.info
			err = result.err
		}
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to detect project info: %v", err)
		}

		// Wait for analysis
		select {
		case <-ctx.Done():
			return nil, nil, "", ctx.Err()
		case result := <-analyzeChan:
			analysis = result.analysis
			err = result.err
		}
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to analyze project: %v", err)
		}

		// Wait for call graph
		select {
		case <-ctx.Done():
			return nil, nil, "", ctx.Err()
		case result := <-callGraphChan:
			callGraph = result.graphData
			err = result.err
		}
		if err != nil {
			return nil, nil, "", fmt.Errorf("failed to generate call graph: %v", err)
		}

		return info, analysis, callGraph, nil
	}()
	if err != nil {
		return "", err
	}

	// Example: Parse callGraph data and use it to refine the template
	if callGraph != "" {
		fmt.Println("Call graph data available for template refinement.")
		graph, err := parseCallGraph(callGraph)
		if err != nil {
			return "", fmt.Errorf("failed to parse call graph: %v", err)
		}
		// Example: Use graph data to enhance the template
		for _, node := range graph.Nodes {
			fmt.Printf("Node ID: %s, Name: %s\n", node.ID, node.Name)
			// TODO: Implement logic to use node data in the template
		}
	}

	// Create template based on detected stack and analysis
	var yamlTemplate tmpl.NexlayerYAML

	// Map detected stack to pod type
	var podType tmpl.PodType
	switch info.Type {
	case detection.TypeReact:
		podType = tmpl.React
	case detection.TypePython:
		// Use analysis to determine framework
		if containsAny(analysis.Frameworks, "django") {
			podType = tmpl.Django
		} else if containsAny(analysis.Frameworks, "fastapi") {
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

			// Enhance prompt with analysis results
			enhancedPrompt := enhancePromptWithAnalysis(llmYamlPrompt, analysis)
			response, err := provider.GenerateText(ctx, enhancedPrompt)
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
			// Use detected ports from analysis if available
			if len(analysis.APIEndpoints) > 0 {
				ports = make([]tmpl.Port, 0)
				for _, endpoint := range analysis.APIEndpoints {
					// Extract port from endpoint if possible
					if port := extractPortFromEndpoint(endpoint.Path); port > 0 {
						ports = append(ports, tmpl.Port{
							ContainerPort: port,
							ServicePort:   port,
							Name:          "api",
						})
					}
				}
			}
			// Fallback to info.Port if no ports detected
			if len(ports) == 0 && info.Port > 0 {
				ports = []tmpl.Port{{ContainerPort: info.Port, ServicePort: info.Port, Name: "app"}}
			}
		}

		// Get default environment variables
		vars := tmpl.DefaultEnvVars[podType]

		// Add database environment variables if needed
		if len(analysis.DatabaseTypes) > 0 {
			vars = append(vars, generateDatabaseEnvVars(analysis.DatabaseTypes)...)
		}

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
	if info.LLMProvider != "" {
		// TODO: Update this once we add AI environment to template package
	}

	// Marshal final template
	data, err := yaml.Marshal(&yamlTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %v", err)
	}

	return string(data), nil
}

// Helper functions

func enhancePromptWithAnalysis(basePrompt string, analysis *analysis.ProjectAnalysis) string {
	var sb strings.Builder
	sb.WriteString(basePrompt)
	sb.WriteString("\n\nProject Analysis:\n")

	// Add detected frameworks
	if len(analysis.Frameworks) > 0 {
		sb.WriteString("\nFrameworks:\n")
		for _, fw := range analysis.Frameworks {
			sb.WriteString(fmt.Sprintf("- %s\n", fw))
		}
	}

	// Add detected API endpoints
	if len(analysis.APIEndpoints) > 0 {
		sb.WriteString("\nAPI Endpoints:\n")
		for _, ep := range analysis.APIEndpoints {
			sb.WriteString(fmt.Sprintf("- %s %s\n", ep.Method, ep.Path))
		}
	}

	// Add detected database types
	if len(analysis.DatabaseTypes) > 0 {
		sb.WriteString("\nDatabases:\n")
		for _, db := range analysis.DatabaseTypes {
			sb.WriteString(fmt.Sprintf("- %s\n", db))
		}
	}

	return sb.String()
}

func generateDatabaseEnvVars(dbTypes []string) []tmpl.EnvVar {
	var vars []tmpl.EnvVar
	for _, dbType := range dbTypes {
		switch strings.ToLower(dbType) {
		case "postgres", "postgresql":
			vars = append(vars, []tmpl.EnvVar{
				{Key: "DB_HOST", Value: "localhost"},
				{Key: "DB_PORT", Value: "5432"},
				{Key: "DB_NAME", Value: "app"},
				{Key: "DB_USER", Value: "postgres"},
				{Key: "DB_PASSWORD", Value: "<% DB_PASSWORD %>"},
			}...)
		case "mysql", "mariadb":
			vars = append(vars, []tmpl.EnvVar{
				{Key: "DB_HOST", Value: "localhost"},
				{Key: "DB_PORT", Value: "3306"},
				{Key: "DB_NAME", Value: "app"},
				{Key: "DB_USER", Value: "root"},
				{Key: "DB_PASSWORD", Value: "<% DB_PASSWORD %>"},
			}...)
		case "mongodb":
			vars = append(vars, []tmpl.EnvVar{
				{Key: "MONGODB_URI", Value: "mongodb://localhost:27017/app"},
			}...)
		}
	}
	return vars
}

func extractPortFromEndpoint(path string) int {
	// Simple port extraction from URL-like paths
	// e.g., "http://localhost:3000" -> 3000
	parts := strings.Split(path, ":")
	if len(parts) > 1 {
		lastPart := parts[len(parts)-1]
		// Remove any trailing path
		if idx := strings.Index(lastPart, "/"); idx != -1 {
			lastPart = lastPart[:idx]
		}
		var port int
		if _, err := fmt.Sscanf(lastPart, "%d", &port); err == nil {
			return port
		}
	}
	return 0
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
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			workDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %v", err)
			}

			// Create request with minimal info
			req := TemplateRequest{
				ProjectName: appName,
				ProjectDir:  workDir,
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
