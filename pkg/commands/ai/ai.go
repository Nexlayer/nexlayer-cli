// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/analysis"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/knowledge"
	tmpl "github.com/Nexlayer/nexlayer-cli/pkg/template"
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
	analysis *analysis.ProjectAnalysis
	err      error
}

type graphResult struct {
	graph *knowledge.Graph
	err   error
}

// processResults handles waiting for and collecting all parallel processing results
func processResults(ctx context.Context, detectChan chan detectionResult, analyzeChan chan analysisResult, graphChan chan graphResult, aiChan chan aiResult) (*detection.ProjectInfo, *analysis.ProjectAnalysis, *knowledge.Graph, string, error) {
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
func startAnalysis(ctx context.Context, projectDir string) (*analysis.ProjectAnalysis, error) {
	parser := analysis.NewParser()
	return parser.AnalyzeProject(ctx, projectDir)
}

// buildKnowledgeGraph constructs the knowledge graph from analysis results
func buildKnowledgeGraph(ctx context.Context, analysis *analysis.ProjectAnalysis, projectDir string) (*knowledge.Graph, error) {
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
		var aiTemplate tmpl.NexlayerYAML
		if err := yaml.Unmarshal([]byte(aiResponse), &aiTemplate); err == nil {
			return aiResponse, nil
		}
	}

	// Fall back to standard template generation
	yamlTemplate := createTemplate(info, analysis, graph, req)

	// Marshal final template
	data, err := yaml.Marshal(&yamlTemplate)
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
	// For now, use a simple detection based on package.json or go.mod
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		// Node.js project
		return &ProjectInfo{
			Name: filepath.Base(dir),
			Type: "node",
			Port: 3000,
		}, nil
	}

	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		// Go project
		return &ProjectInfo{
			Name: filepath.Base(dir),
			Type: "go",
			Port: 8080,
		}, nil
	}

	// Default to unknown
	return &ProjectInfo{
		Name: filepath.Base(dir),
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
			projectInfo, err := detectProject(cwd)
			if err != nil {
				return fmt.Errorf("failed to detect project: %w", err)
			}

			// Analyze project
			analysis := &AnalysisResult{
				Components: []*Component{
					{
						Name:  "app",
						Type:  projectInfo.Type,
						Image: fmt.Sprintf("ghcr.io/nexlayer/%s:latest", projectInfo.Name),
						Ports: []int{projectInfo.Port},
					},
				},
			}

			// Generate template
			tmpl, err := createTemplate(projectInfo, analysis, nil)
			if err != nil {
				return fmt.Errorf("failed to create template: %w", err)
			}

			// Write template to file
			data, err := yaml.Marshal(tmpl)
			if err != nil {
				return fmt.Errorf("failed to marshal template: %w", err)
			}

			outputFile := "nexlayer.yaml"
			if err := os.WriteFile(outputFile, data, 0644); err != nil {
				return fmt.Errorf("failed to write template: %w", err)
			}

			fmt.Printf("Generated deployment template: %s\n", outputFile)
			return nil
		},
	}

	return cmd
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
