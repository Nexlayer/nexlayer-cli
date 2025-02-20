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
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/analysis"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/knowledge"
	"github.com/Nexlayer/nexlayer-cli/pkg/template"
)

// llmYamlPrompt provides structured instructions for AI-generated Nexlayer templates.
const llmYamlPrompt = `You are an expert in cloud automation for Nexlayer AI Cloud Platform.
Generate a deployment template YAML that seamlessly integrates into Nexlayer Cloud following the Nexlayer YAML Schema Template Documentation (v1.0).`

// TemplateRequest represents a request to generate a Nexlayer deployment template.
type TemplateRequest struct {
	ProjectName string
	ProjectDir  string
}

// result types for parallel processing
type aiResult struct {
	response string
	err      error
}

type detectionResult struct {
	info *detection.ProjectInfo
	err  error
}

type analysisResult struct {
	analysis *types.ProjectAnalysis
	err      error
}

type graphResult struct {
	graph *knowledge.Graph
	err   error
}

// processResults handles waiting for and collecting all parallel processing results
func processResults(ctx context.Context, detectChan chan detectionResult, analyzeChan chan analysisResult, graphChan chan graphResult, aiChan chan aiResult) (*detection.ProjectInfo, *types.ProjectAnalysis, *knowledge.Graph, string, error) {
	select {
	case <-ctx.Done():
		return nil, nil, nil, "", ctx.Err()
	case result := <-detectChan:
		if result.err != nil {
			return nil, nil, nil, "", fmt.Errorf("failed to detect project info: %w", result.err)
		}
		info := result.info

		select {
		case <-ctx.Done():
			return nil, nil, nil, "", ctx.Err()
		case result := <-analyzeChan:
			if result.err != nil {
				return nil, nil, nil, "", fmt.Errorf("failed to analyze project: %w", result.err)
			}
			analysis := result.analysis

			select {
			case <-ctx.Done():
				return nil, nil, nil, "", ctx.Err()
			case result := <-graphChan:
				if result.err != nil {
					return nil, nil, nil, "", fmt.Errorf("failed to build knowledge graph: %w", result.err)
				}
				graph := result.graph

				var aiResponse string
				select {
				case <-ctx.Done():
					return nil, nil, nil, "", ctx.Err()
				case result := <-aiChan:
					aiResponse = result.response
				case <-time.After(5 * time.Second):
					// Timeout waiting for AI response is acceptable
				}

				return info, analysis, graph, aiResponse, nil
			}
		}
	}
}

// startAnalysis initiates the project analysis
func startAnalysis(ctx context.Context, projectDir string) (*types.ProjectAnalysis, error) {
	parser := analysis.NewParser()
	return parser.AnalyzeProject(ctx, projectDir)
}

// buildKnowledgeGraph constructs the knowledge graph from analysis results
func buildKnowledgeGraph(ctx context.Context, analysis *types.ProjectAnalysis, projectDir string) (*knowledge.Graph, error) {
	graph := knowledge.NewGraph()
	if err := graph.BuildFromAnalysis(ctx, analysis); err != nil {
		return nil, fmt.Errorf("failed to build graph: %w", err)
	}

	// Run go-callvis in the background
	if err := addCallGraphData(graph, projectDir); err != nil {
		fmt.Printf("Warning: Failed to add call graph data: %v\n", err)
	}

	// Setup file watcher
	if err := setupGraphWatcher(ctx, graph, projectDir); err != nil {
		fmt.Printf("Warning: Failed to setup graph watcher: %v\n", err)
	}

	return graph, nil
}

// addCallGraphData adds call graph information to the knowledge graph
func addCallGraphData(graph *knowledge.Graph, projectDir string) error {
	cmd := exec.Command("go-callvis", "-focus", projectDir, "-group", "pkg,type", "-nostd", "-format", "json")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return graph.AddCallGraphData(output)
}

// setupGraphWatcher initializes and starts the file watcher for the graph
func setupGraphWatcher(ctx context.Context, graph *knowledge.Graph, projectDir string) error {
	watcher, err := knowledge.NewWatcher(graph, projectDir)
	if err != nil {
		return err
	}
	return watcher.Start(ctx)
}

// GenerateTemplate uses AI assistance to generate a valid Nexlayer deployment template.
func GenerateTemplate(ctx context.Context, req TemplateRequest) (string, error) {
	// Create channels for parallel processing
	aiChan := make(chan aiResult, 1)
	detectChan := make(chan detectionResult, 1)
	analyzeChan := make(chan analysisResult, 1)
	graphChan := make(chan graphResult, 1)

	// Start stack detection
	go func() {
		info, err := DetectStack(req.ProjectDir)
		detectChan <- detectionResult{info, err}
	}()

	// Start tree-sitter analysis
	go func() {
		analysis, err := startAnalysis(ctx, req.ProjectDir)
		analyzeChan <- analysisResult{analysis, err}
	}()

	// Start knowledge graph construction
	go func() {
		analysisRes := <-analyzeChan
		if analysisRes.err != nil {
			graphChan <- graphResult{nil, fmt.Errorf("analysis failed: %w", analysisRes.err)}
			return
		}

		graph, err := buildKnowledgeGraph(ctx, analysisRes.analysis, req.ProjectDir)
		if err != nil {
			graphChan <- graphResult{nil, err}
			return
		}

		// Create LLM enricher and generate enhanced prompt
		enricher := knowledge.NewLLMEnricher(graph)
		if err := enricher.LoadMetadata("tools"); err == nil {
			if enhancedPrompt, err := enricher.GeneratePrompt(ctx, llmYamlPrompt); err == nil {
				// Get the AI provider
				provider := NewDefaultProvider()
				if response, err := provider.GenerateText(ctx, enhancedPrompt); err == nil {
					aiChan <- aiResult{response, nil}
				}
			}
		}

		graphChan <- graphResult{graph, nil}
	}()

	// Wait for all results
	info, analysis, graph, aiResponse, err := processResults(ctx, detectChan, analyzeChan, graphChan, aiChan)
	if err != nil {
		return "", err
	}

	// Try to use AI response first if available
	if aiResponse != "" {
		var aiTemplate template.NexlayerYAML
		if err := yaml.Unmarshal([]byte(aiResponse), &aiTemplate); err == nil {
			return aiResponse, nil
		}
	}

	// Convert detection.ProjectInfo to our local ProjectInfo
	projectInfo := &ProjectInfo{
		Name: info.Name,
		Type: string(info.Type),
		Port: info.Port,
	}

	// Convert analysis.ProjectAnalysis to our AnalysisResult
	analysisResult := &AnalysisResult{
		Components: make([]*Component, 0),
	}

	// Add components from analysis
	for _, functions := range analysis.Functions {
		for _, fn := range functions {
			// Create a component for each detected function
			component := &Component{
				Name:  fn.Name,
				Type:  "function",
				Image: fmt.Sprintf("ghcr.io/nexlayer/%s:latest", req.ProjectName),
				Ports: []int{projectInfo.Port},
			}
			analysisResult.Components = append(analysisResult.Components, component)
		}
	}

	// Convert knowledge.Graph to our GraphResult
	graphResult := &GraphResult{
		Nodes: make([]*GraphNode, 0),
	}

	if graph != nil {
		for _, node := range graph.Nodes {
			graphNode := &GraphNode{
				Name:        node.Name,
				EnvVars:     make(map[string]string),
				Annotations: make(map[string]string),
			}
			// Add any environment variables and annotations from node properties
			if props := node.Properties; props != nil {
				if envVars, ok := props["env"].(map[string]string); ok {
					graphNode.EnvVars = envVars
				}
				if annotations, ok := props["annotations"].(map[string]string); ok {
					graphNode.Annotations = annotations
				}
			}
			graphResult.Nodes = append(graphResult.Nodes, graphNode)
		}
	}

	// Fall back to standard template generation
	tmpl, err := createTemplate(projectInfo, analysisResult, graphResult)
	if err != nil {
		return "", fmt.Errorf("failed to create template: %w", err)
	}

	// Marshal final template
	data, err := yaml.Marshal(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to marshal template: %w", err)
	}

	return string(data), nil
}

// ProjectInfo represents basic project information
type ProjectInfo struct {
	Name string
	Type string
	Port int
}

// detectProject detects project information from a directory
func detectProject(dir string) (*ProjectInfo, error) {
	// Get project name from directory
	projectName := filepath.Base(dir)
	// Clean the project name
	projectName = strings.ToLower(projectName)
	projectName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, projectName)

	// Check for package.json (Next.js/Node.js)
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		// Read package.json to determine if it's Next.js
		data, err := os.ReadFile(filepath.Join(dir, "package.json"))
		if err == nil {
			var pkg struct {
				Dependencies map[string]string `json:"dependencies"`
			}
			if err := json.Unmarshal(data, &pkg); err == nil {
				if _, hasNext := pkg.Dependencies["next"]; hasNext {
					return &ProjectInfo{
						Name: projectName,
						Type: "nextjs",
						Port: 3000,
					}, nil
				}
			}
		}

		// Default to Node.js
		return &ProjectInfo{
			Name: projectName,
			Type: "node",
			Port: 3000,
		}, nil
	}

	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		// Go project
		return &ProjectInfo{
			Name: projectName,
			Type: "go",
			Port: 8080,
		}, nil
	}

	// Default to unknown
	return &ProjectInfo{
		Name: projectName,
		Type: "unknown",
	}, nil
}

// NewCommand creates the "ai" command with its subcommands.
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-powered features for Nexlayer",
		Long:  "AI-powered features for generating deployment templates and detecting project types",
	}

	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newDetectCommand())

	return cmd
}

// newGenerateCommand creates a new generate command
func newGenerateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a deployment template using AI",
		Long:  "Generate a deployment template by analyzing your project using AI",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current directory
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			// Detect project type
			info, err := DetectStack(cwd)
			if err != nil {
				return fmt.Errorf("failed to detect project: %w", err)
			}

			// Create template generator
			generator := template.NewGenerator()

			// Generate template
			yamlTemplate, err := generator.GenerateFromProjectInfo(info.Name, string(info.Type), info.Port)
			if err != nil {
				return fmt.Errorf("failed to generate template: %w", err)
			}

			// Add database if needed
			if hasDatabase(info) {
				if err := generator.AddPod(yamlTemplate, template.PodTypePostgres, 0); err != nil {
					return fmt.Errorf("failed to add database: %w", err)
				}
			}

			// Marshal template to YAML
			yamlData, err := yaml.Marshal(yamlTemplate)
			if err != nil {
				return fmt.Errorf("failed to marshal template: %w", err)
			}

			// Write template to file
			outputFile := "nexlayer.yaml"
			if err := os.WriteFile(outputFile, yamlData, 0644); err != nil {
				return fmt.Errorf("failed to write template: %w", err)
			}

			fmt.Printf("Generated deployment template: %s\n", outputFile)
			return nil
		},
	}

	return cmd
}

// hasDatabase checks if the project needs a database
func hasDatabase(info *detection.ProjectInfo) bool {
	// Check dependencies for database-related packages
	for _, dep := range info.Dependencies {
		switch dep {
		case "pg", "postgres", "postgresql", "sequelize", "typeorm", "prisma",
			"mongoose", "mongodb", "mysql", "mysql2", "sqlite3", "redis":
			return true
		}
	}
	return false
}

// newDetectCommand creates a new detect command
func newDetectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detect",
		Short: "Detect project type using AI",
		Long:  "Detect your project type and configuration using AI analysis",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get current directory
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			// Detect project type
			projectInfo, err := detectProject(cwd)
			if err != nil {
				return fmt.Errorf("failed to detect project: %w", err)
			}

			fmt.Printf("Detected project type: %s\n", projectInfo.Type)
			fmt.Printf("Project name: %s\n", projectInfo.Name)
			if projectInfo.Port > 0 {
				fmt.Printf("Default port: %d\n", projectInfo.Port)
			}

			return nil
		},
	}

	return cmd
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
