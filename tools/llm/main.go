package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// LLMMetadata represents the top-level structure for LLM-optimized metadata
type LLMMetadata struct {
	// Structured in natural language format for LLMs
	Purpose string `json:"purpose"`
	Version string `json:"version"`
	
	// Core capabilities and concepts
	Capabilities []Capability `json:"capabilities"`
	
	// Deployment patterns with examples
	DeploymentPatterns []DeploymentPattern `json:"deployment_patterns"`
	
	// API endpoints with natural language descriptions
	APIEndpoints []APIEndpoint `json:"api_endpoints"`
	
	// Common user intents and how to handle them
	UserIntents []UserIntent `json:"user_intents"`
}

type Capability struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Examples    []string `json:"examples"`
	Keywords    []string `json:"keywords"` // For semantic search
}

type DeploymentPattern struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Template    string `json:"template"`
	// Natural language explanation of the pattern
	Explanation string   `json:"explanation"`
	UseCase     string   `json:"use_case"`
	Keywords    []string `json:"keywords"`
}

type APIEndpoint struct {
	Path        string `json:"path"`
	Method      string `json:"method"`
	Description string `json:"description"`
	// Natural language examples of when to use this endpoint
	UsageExamples []string `json:"usage_examples"`
	// Common patterns in natural language
	CommonPatterns []string `json:"common_patterns"`
}

type UserIntent struct {
	Intent      string   `json:"intent"`
	Keywords    []string `json:"keywords"`
	Actions     []string `json:"actions"`
	Examples    []string `json:"examples"`
	Suggestions []string `json:"suggestions"`
}

func main() {
	metadata := LLMMetadata{
		Purpose: "This metadata helps AI agents understand how to deploy applications to Nexlayer Cloud, either through the CLI or direct API calls",
		Version: "2.0.0",
		Capabilities: []Capability{
			{
				Name:        "Container Deployment",
				Description: "Deploy containerized applications to Nexlayer Cloud",
				Examples: []string{
					"Deploy a Node.js application with MongoDB",
					"Deploy a React frontend with a Go backend",
				},
				Keywords: []string{"deploy", "container", "docker", "kubernetes", "pod"},
			},
			// Add more capabilities...
		},
		DeploymentPatterns: []DeploymentPattern{
			{
				Name:        "Frontend with Backend",
				Description: "Deploy a web application with separate frontend and backend services",
				Template: `application:
  name: web-app
  pods:
    - name: frontend
      image: nginx
    - name: backend
      image: node`,
				Explanation: "This pattern creates two pods: one for the frontend (typically a web server) and one for the backend (application server). They can communicate internally while the frontend is exposed to the internet.",
				UseCase:     "Modern web applications that separate presentation from business logic",
				Keywords:    []string{"frontend", "backend", "web", "microservices"},
			},
			// Add more patterns...
		},
		APIEndpoints: []APIEndpoint{
			{
				Path:        "/startUserDeployment",
				Method:     "POST",
				Description: "Start a new deployment from a YAML configuration",
				UsageExamples: []string{
					"When a user wants to deploy their application",
					"When updating an existing deployment with new configuration",
				},
				CommonPatterns: []string{
					"First validate the YAML, then call this endpoint",
					"Use this when deploying from CI/CD pipelines",
				},
			},
			// Add more endpoints...
		},
		UserIntents: []UserIntent{
			{
				Intent:   "Deploy a web application",
				Keywords: []string{"deploy", "web", "app", "website"},
				Actions: []string{
					"1. Validate the deployment configuration",
					"2. Check for existing deployments",
					"3. Start new deployment",
				},
				Examples: []string{
					"I want to deploy my React app",
					"How do I deploy a Node.js backend?",
				},
				Suggestions: []string{
					"Consider using the frontend-backend pattern",
					"Make sure to configure environment variables",
				},
			},
			// Add more intents...
		},
	}

	// Create AI training metadata directory if it doesn't exist
	aiDir := filepath.Join("ai_training", "metadata")
	if err := os.MkdirAll(aiDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create AI training directory: %v\n", err)
		os.Exit(1)
	}

	// Write LLM-optimized metadata
	llmFile := filepath.Join(aiDir, "llm_metadata.json")
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal metadata: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(llmFile, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write metadata file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Generated LLM-optimized metadata in build/llm_metadata.json")
}
